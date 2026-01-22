package qwen

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client Qwen Embedder 客户端
// 实现了 embedder.Provider 接口，提供基于阿里云 DashScope Text Embedding API 的文本向量化功能
type Client struct {
	client     *http.Client
	apiKey     string
	model      string
	baseURL    string
	dimensions int
}

// Config Qwen Embedder 配置
// APIKey: DashScope API 密钥（必需）
// Model: 使用的模型名称，默认为 "text-embedding-v4"
// BaseURL: API 基础 URL，默认为 DashScope 官方地址
// Dimensions: 向量维度，默认为 1536（text-embedding-v4 的默认维度）
// HTTPClient: 自定义 HTTP 客户端，如果为 nil 则使用默认客户端
type Config struct {
	APIKey     string
	Model      string
	BaseURL    string
	Dimensions int
	HTTPClient *http.Client
}

// NewClient 创建新的 Qwen Embedder 客户端
// 参数:
//   - cfg: Qwen Embedder 配置，包含 APIKey、Model、BaseURL、Dimensions 等
// 返回:
//   - *Client: Qwen Embedder 客户端实例
//   - error: 如果配置无效（如缺少 APIKey）或初始化失败则返回错误
func NewClient(cfg *Config) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, errors.New("API key is required")
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://dashscope.aliyuncs.com/api/v1"
	}

	model := cfg.Model
	if model == "" {
		model = "text-embedding-v4"
	}

	dimensions := cfg.Dimensions
	if dimensions == 0 {
		dimensions = 1536 // text-embedding-v4 默认维度
	}

	client := cfg.HTTPClient
	if client == nil {
		client = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	return &Client{
		client:     client,
		apiKey:     cfg.APIKey,
		model:      model,
		baseURL:    baseURL,
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
	// 构建请求
	reqBody := map[string]interface{}{
		"model": c.model,
		"input": map[string]interface{}{
			"texts": []string{text},
		},
	}

	// 添加维度参数
	if c.dimensions > 0 {
		reqBody["parameters"] = map[string]interface{}{
			"dimension": c.dimensions,
		}
	}

	// 默认使用 document 类型
	reqBody["text_type"] = "document"

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// 创建 HTTP 请求
	url := fmt.Sprintf("%s/services/embeddings/text-embedding/text-embedding", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var response struct {
		Output struct {
			Embeddings []struct {
				Embedding []float64 `json:"embedding"`
			} `json:"embeddings"`
		} `json:"output"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(response.Output.Embeddings) == 0 {
		return nil, errors.New("embedding generation failed: no embeddings returned from Qwen API")
	}

	return response.Output.Embeddings[0].Embedding, nil
}

// EmbedBatch 批量将多个文本转换为向量
// 参数:
//   - ctx: 上下文，用于控制请求生命周期
//   - texts: 要向量化的文本列表
// 返回:
//   - [][]float64: 每个文本对应的向量表示（顺序与输入文本一致）
//   - error: 如果向量化失败或返回结果数量不匹配则返回错误
func (c *Client) EmbedBatch(ctx context.Context, texts []string) ([][]float64, error) {
	// 构建请求
	reqBody := map[string]interface{}{
		"model": c.model,
		"input": map[string]interface{}{
			"texts": texts,
		},
	}

	// 添加维度参数
	if c.dimensions > 0 {
		reqBody["parameters"] = map[string]interface{}{
			"dimension": c.dimensions,
		}
	}

	// 默认使用 document 类型
	reqBody["text_type"] = "document"

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// 创建 HTTP 请求
	url := fmt.Sprintf("%s/services/embeddings/text-embedding/text-embedding", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var response struct {
		Output struct {
			Embeddings []struct {
				Embedding []float64 `json:"embedding"`
			} `json:"embeddings"`
		} `json:"output"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(response.Output.Embeddings) != len(texts) {
		return nil, fmt.Errorf("embedding generation failed: unexpected number of results from Qwen API (got %d, expected %d)", len(response.Output.Embeddings), len(texts))
	}

	embeddings := make([][]float64, len(texts))
	for i, emb := range response.Output.Embeddings {
		embeddings[i] = emb.Embedding
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
// HTTP 客户端不需要显式关闭，此方法为接口兼容性保留
// 返回:
//   - error: 始终返回 nil
func (c *Client) Close() error {
	return nil
}
