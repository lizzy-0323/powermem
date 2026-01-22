package openai

import (
	"context"
	"errors"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

// Client OpenAI Embedder 客户端
// 实现了 embedder.Provider 接口，提供基于 OpenAI Embeddings API 的文本向量化功能
type Client struct {
	client     *openai.Client
	model      openai.EmbeddingModel
	dimensions int
}

// Config OpenAI Embedder 配置
// APIKey: OpenAI API 密钥（必需）
// Model: 使用的模型名称，当前固定使用 AdaEmbeddingV2
// BaseURL: API 基础 URL，默认为 OpenAI 官方地址
// Dimensions: 向量维度，默认为 1536（AdaEmbeddingV2 的默认维度）
type Config struct {
	APIKey     string
	Model      string
	BaseURL    string
	Dimensions int
}

// NewClient 创建新的 OpenAI Embedder 客户端
// 参数:
//   - cfg: OpenAI Embedder 配置，包含 APIKey、BaseURL、Dimensions 等
// 返回:
//   - *Client: OpenAI Embedder 客户端实例
//   - error: 如果配置无效或初始化失败则返回错误
func NewClient(cfg *Config) (*Client, error) {
	config := openai.DefaultConfig(cfg.APIKey)
	if cfg.BaseURL != "" {
		config.BaseURL = cfg.BaseURL
	}

	client := openai.NewClientWithConfig(config)

	// 默认使用 Ada v2 模型
	model := openai.AdaEmbeddingV2

	dimensions := cfg.Dimensions
	if dimensions == 0 {
		dimensions = 1536 // AdaEmbeddingV2 的默认维度
	}

	return &Client{
		client:     client,
		model:      model,
		dimensions: dimensions,
	}, nil
}

// Embed 将单个文本转换为向量
// 参数:
//   - ctx: 上下文，用于控制请求生命周期
//   - text: 要向量化的文本内容
// 返回:
//   - []float64: 文本的向量表示（维度由配置决定）
//   - error: 如果向量化失败则返回错误
func (c *Client) Embed(ctx context.Context, text string) ([]float64, error) {
	resp, err := c.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: []string{text},
		Model: c.model,
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, errors.New("embedding generation failed: no data returned from OpenAI API")
	}

	// 转换 float32 到 float64
	embedding32 := resp.Data[0].Embedding
	embedding64 := make([]float64, len(embedding32))
	for i, v := range embedding32 {
		embedding64[i] = float64(v)
	}

	return embedding64, nil
}

// EmbedBatch 批量将多个文本转换为向量
// 参数:
//   - ctx: 上下文，用于控制请求生命周期
//   - texts: 要向量化的文本列表
// 返回:
//   - [][]float64: 每个文本对应的向量表示（顺序与输入文本一致）
//   - error: 如果向量化失败或返回结果数量不匹配则返回错误
func (c *Client) EmbedBatch(ctx context.Context, texts []string) ([][]float64, error) {
	resp, err := c.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: texts,
		Model: c.model,
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Data) != len(texts) {
		return nil, fmt.Errorf("embedding generation failed: unexpected number of results from OpenAI API (got %d, expected %d)", len(resp.Data), len(texts))
	}

	embeddings := make([][]float64, len(texts))
	for i, data := range resp.Data {
		embedding32 := data.Embedding
		embedding64 := make([]float64, len(embedding32))
		for j, v := range embedding32 {
			embedding64[j] = float64(v)
		}
		embeddings[i] = embedding64
	}

	return embeddings, nil
}

// Dimensions 返回向量维度
// 返回:
//   - int: 向量维度数
func (c *Client) Dimensions() int {
	return c.dimensions
}

// Close 关闭客户端连接
// OpenAI SDK 的客户端不需要显式关闭，此方法为接口兼容性保留
// 返回:
//   - error: 始终返回 nil
func (c *Client) Close() error {
	return nil
}
