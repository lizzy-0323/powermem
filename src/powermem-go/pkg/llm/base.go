package llm

import "context"

// Provider LLM 提供商接口
type Provider interface {
	// Generate 生成文本
	Generate(ctx context.Context, prompt string, opts ...GenerateOption) (string, error)

	// GenerateWithMessages 使用消息历史生成
	GenerateWithMessages(ctx context.Context, messages []Message, opts ...GenerateOption) (string, error)

	// Close 关闭连接
	Close() error
}

// Message 聊天消息
type Message struct {
	Role    string `json:"role"` // system, user, assistant
	Content string `json:"content"`
}

// GenerateOptions 生成选项
type GenerateOptions struct {
	Temperature float64
	MaxTokens   int
	TopP        float64
	Stop        []string
}

// GenerateOption 生成选项函数
type GenerateOption func(*GenerateOptions)

// WithTemperature 设置温度
func WithTemperature(temp float64) GenerateOption {
	return func(opts *GenerateOptions) {
		opts.Temperature = temp
	}
}

// WithMaxTokens 设置最大 token 数
func WithMaxTokens(max int) GenerateOption {
	return func(opts *GenerateOptions) {
		opts.MaxTokens = max
	}
}

// WithTopP 设置 top_p
func WithTopP(topP float64) GenerateOption {
	return func(opts *GenerateOptions) {
		opts.TopP = topP
	}
}

// ApplyGenerateOptions 应用生成选项
func ApplyGenerateOptions(opts []GenerateOption) *GenerateOptions {
	options := &GenerateOptions{
		Temperature: 0.7,
		MaxTokens:   1000,
		TopP:        1.0,
	}
	for _, opt := range opts {
		opt(options)
	}
	return options
}
