package embedder

import "context"

// Provider Embedder 提供商接口
type Provider interface {
	// Embed 将文本转换为向量
	Embed(ctx context.Context, text string) ([]float64, error)

	// EmbedBatch 批量转换文本为向量
	EmbedBatch(ctx context.Context, texts []string) ([][]float64, error)

	// Dimensions 返回向量维度
	Dimensions() int

	// Close 关闭连接
	Close() error
}
