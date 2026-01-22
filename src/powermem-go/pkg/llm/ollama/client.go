package ollama

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

// Client Ollama LLM 客户端
// 实现了 llm.Provider 接口，提供基于 Ollama 本地/远程服务的文本生成功能
// Ollama 是一个本地运行大语言模型的工具，支持本地部署和远程访问
type Client struct {
	client  *http.Client
	apiKey  string
	model   string
	baseURL string
}

// Config Ollama LLM 配置
// APIKey: Ollama API 密钥（可选，本地部署通常不需要）
// Model: 使用的模型名称，默认为 "llama3.1:70b"
// BaseURL: Ollama 服务地址，默认为 "http://localhost:11434"
// HTTPClient: 自定义 HTTP 客户端，如果为 nil 则使用默认客户端（超时 120 秒）
type Config struct {
	APIKey     string
	Model      string
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient 创建新的 Ollama LLM 客户端
// 参数:
//   - cfg: Ollama 配置，包含 Model、BaseURL 等（APIKey 可选）
// 返回:
//   - *Client: Ollama 客户端实例
//   - error: 如果初始化失败则返回错误
func NewClient(cfg *Config) (*Client, error) {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	model := cfg.Model
	if model == "" {
		model = "llama3.1:70b"
	}

	client := cfg.HTTPClient
	if client == nil {
		// Ollama 可能需要更长的超时时间，特别是对于大模型
		client = &http.Client{
			Timeout: 120 * time.Second,
		}
	}

	return &Client{
		client:  client,
		apiKey:  cfg.APIKey, // Ollama 本地部署通常不需要 API key，但保留以支持需要认证的远程部署
		model:   model,
		baseURL: baseURL,
	}, nil
}

// Generate 根据提示词生成文本
// 参数:
//   - ctx: 上下文，用于控制请求生命周期
//   - prompt: 用户输入的提示词
//   - opts: 可选的生成参数（temperature, max_tokens, top_p 等）
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
// 注意: Ollama 使用不同的参数名（num_predict 而不是 max_tokens）
// 参数:
//   - ctx: 上下文，用于控制请求生命周期
//   - messages: 消息历史列表，每个消息包含 role 和 content
//   - opts: 可选的生成参数（temperature, max_tokens, top_p 等）
// 返回:
//   - string: 生成的文本内容
//   - error: 如果生成失败则返回错误
func (c *Client) GenerateWithMessages(ctx context.Context, messages []llm.Message, opts ...llm.GenerateOption) (string, error) {
	options := llm.ApplyGenerateOptions(opts)

	// 转换消息格式
	chatMessages := make([]map[string]string, len(messages))
	for i, msg := range messages {
		chatMessages[i] = map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		}
	}

	// 构建请求体
	reqBody := map[string]interface{}{
		"model":    c.model,
		"messages": chatMessages,
		"options": map[string]interface{}{
			"temperature": options.Temperature,
			"num_predict": options.MaxTokens,
			"top_p":       options.TopP,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	// 创建 HTTP 请求
	url := fmt.Sprintf("%s/api/chat", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}

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
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if response.Message.Content == "" {
		return "", errors.New("llm generation failed: empty response from Ollama API")
	}

	return response.Message.Content, nil
}

// Close 关闭客户端连接
// HTTP 客户端不需要显式关闭，此方法为接口兼容性保留
// 返回:
//   - error: 始终返回 nil
func (c *Client) Close() error {
	return nil
}
