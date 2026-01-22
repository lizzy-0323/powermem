package deepseek

import (
	"context"
	"errors"

	"github.com/oceanbase/powermem-go/pkg/llm"
	openai "github.com/sashabaranov/go-openai"
)

// Client DeepSeek LLM 客户端
// 实现了 llm.Provider 接口，提供基于 DeepSeek API 的文本生成功能
// DeepSeek 使用 OpenAI 兼容的 API 格式，因此可以复用 OpenAI SDK
type Client struct {
	client *openai.Client
	model  string
}

// Config DeepSeek LLM 配置
// APIKey: DeepSeek API 密钥（必需）
// Model: 使用的模型名称，默认为 "deepseek-chat"
// BaseURL: API 基础 URL，默认为 "https://api.deepseek.com"
type Config struct {
	APIKey  string
	Model   string
	BaseURL string
}

// NewClient 创建新的 DeepSeek LLM 客户端
// 参数:
//   - cfg: DeepSeek 配置，包含 APIKey、Model 和 BaseURL
//
// 返回:
//   - *Client: DeepSeek 客户端实例
//   - error: 如果配置无效或初始化失败则返回错误
func NewClient(cfg *Config) (*Client, error) {
	config := openai.DefaultConfig(cfg.APIKey)

	// DeepSeek 使用 OpenAI 兼容的 API，但 base URL 不同
	if cfg.BaseURL != "" {
		config.BaseURL = cfg.BaseURL
	} else {
		// 默认 DeepSeek API base URL
		config.BaseURL = "https://api.deepseek.com"
	}

	client := openai.NewClientWithConfig(config)

	return &Client{
		client: client,
		model:  cfg.Model,
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
// 参数:
//   - ctx: 上下文，用于控制请求生命周期
//   - messages: 消息历史列表，每个消息包含 role 和 content
//   - opts: 可选的生成参数（temperature, max_tokens, top_p 等）
//
// 返回:
//   - string: 生成的文本内容
//   - error: 如果生成失败则返回错误
func (c *Client) GenerateWithMessages(ctx context.Context, messages []llm.Message, opts ...llm.GenerateOption) (string, error) {
	options := llm.ApplyGenerateOptions(opts)

	// 转换消息格式
	chatMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		chatMessages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	req := openai.ChatCompletionRequest{
		Model:       c.model,
		Messages:    chatMessages,
		Temperature: float32(options.Temperature),
		MaxTokens:   options.MaxTokens,
		TopP:        float32(options.TopP),
		Stop:        options.Stop,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("llm generation failed: no choices returned from DeepSeek API")
	}

	return resp.Choices[0].Message.Content, nil
}

// Close 关闭客户端连接
// DeepSeek 客户端（基于 OpenAI SDK）不需要显式关闭，此方法为接口兼容性保留
// 返回:
//   - error: 始终返回 nil
func (c *Client) Close() error {
	return nil
}
