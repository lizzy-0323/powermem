package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/oceanbase/powermem-go/pkg/llm"
)

// Client Anthropic LLM 客户端
// 实现了 llm.Provider 接口，提供基于 Anthropic Claude API 的文本生成功能
// 支持 system 消息分离，符合 Anthropic Messages API 规范
type Client struct {
	client  *http.Client
	apiKey  string
	model   string
	baseURL string
}

// Config Anthropic LLM 配置
// APIKey: Anthropic API 密钥（必需）
// Model: 使用的模型名称，默认为 "claude-3-5-sonnet-20240620"
// BaseURL: API 基础 URL，默认为 "https://api.anthropic.com"
// HTTPClient: 自定义 HTTP 客户端，如果为 nil 则使用默认客户端（超时 120 秒）
type Config struct {
	APIKey     string
	Model      string
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient 创建新的 Anthropic LLM 客户端
// 参数:
//   - cfg: Anthropic 配置，包含 APIKey、Model、BaseURL 等
//
// 返回:
//   - *Client: Anthropic 客户端实例
//   - error: 如果配置无效（如缺少 APIKey）或初始化失败则返回错误
func NewClient(cfg *Config) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, errors.New("API key is required")
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.anthropic.com"
	}

	model := cfg.Model
	if model == "" {
		model = "claude-3-5-sonnet-20240620"
	}

	client := cfg.HTTPClient
	if client == nil {
		client = &http.Client{
			Timeout: 120 * time.Second,
		}
	}

	return &Client{
		client:  client,
		apiKey:  cfg.APIKey,
		model:   model,
		baseURL: baseURL,
	}, nil
}

// Generate 根据提示词生成文本
// 参数:
//   - ctx: 上下文，用于控制请求生命周期
//   - prompt: 用户输入的提示词
//   - opts: 可选的生成参数（temperature, max_tokens, top_p 等）
//
// 返回:
//   - string: 生成的文本内容
//   - error: 如果生成失败则返回错误
func (c *Client) Generate(ctx context.Context, prompt string, opts ...llm.GenerateOption) (string, error) {
	messages := []llm.Message{
		{Role: "user", Content: prompt},
	}
	return c.GenerateWithMessages(ctx, messages, opts...)
}

// GenerateWithMessages 使用消息历史生成文本
// 支持多轮对话，可以传入完整的消息历史（包括 system、user、assistant 消息）
// 注意: Anthropic API 要求将 system 消息单独传递，而不是放在 messages 数组中
// 参数:
//   - ctx: 上下文，用于控制请求生命周期
//   - messages: 消息历史列表，每个消息包含 role 和 content（system 消息会被自动分离）
//   - opts: 可选的生成参数（temperature, max_tokens, top_p 等）
//
// 返回:
//   - string: 生成的文本内容
//   - error: 如果生成失败则返回错误
func (c *Client) GenerateWithMessages(ctx context.Context, messages []llm.Message, opts ...llm.GenerateOption) (string, error) {
	options := llm.ApplyGenerateOptions(opts)

	// 分离 system 消息和其他消息
	var systemMessage string
	var filteredMessages []map[string]string

	for _, msg := range messages {
		if msg.Role == "system" {
			systemMessage = msg.Content
		} else {
			filteredMessages = append(filteredMessages, map[string]string{
				"role":    msg.Role,
				"content": msg.Content,
			})
		}
	}

	// 构建请求体
	reqBody := map[string]interface{}{
		"model":       c.model,
		"max_tokens":  options.MaxTokens,
		"temperature": options.Temperature,
		"top_p":       options.TopP,
		"messages":    filteredMessages,
	}

	if systemMessage != "" {
		reqBody["system"] = systemMessage
	}

	if len(options.Stop) > 0 {
		reqBody["stop_sequences"] = options.Stop
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	// 创建 HTTP 请求
	url := fmt.Sprintf("%s/v1/messages", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var response struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(response.Content) == 0 {
		return "", errors.New("llm generation failed: no content returned from Anthropic API")
	}

	return response.Content[0].Text, nil
}

// Close 关闭客户端连接
// HTTP 客户端不需要显式关闭，此方法为接口兼容性保留
// 返回:
//   - error: 始终返回 nil
func (c *Client) Close() error {
	return nil
}
