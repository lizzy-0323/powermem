package storage

import (
	"context"
	"time"
)

// Memory 记忆类型（避免循环依赖）
type Memory struct {
	ID                int64
	UserID            string
	AgentID           string
	Content           string
	Embedding         []float64
	SparseEmbedding   map[int]float64
	Metadata          map[string]interface{}
	CreatedAt         time.Time
	UpdatedAt         time.Time
	RetentionStrength float64
	LastAccessedAt    *time.Time
	Score             float64
}

// VectorIndexType 向量索引类型
type VectorIndexType string

const (
	IndexTypeHNSW    VectorIndexType = "HNSW"
	IndexTypeIVFFlat VectorIndexType = "IVF_FLAT"
	IndexTypeIVFPQ   VectorIndexType = "IVF_PQ"
)

// MetricType 距离度量类型
type MetricType string

const (
	MetricCosine MetricType = "cosine"
	MetricL2     MetricType = "l2"
	MetricIP     MetricType = "ip"
)

// HNSWParams HNSW 索引参数
type HNSWParams struct {
	M              int
	EfConstruction int
	EfSearch       int
}

// IVFParams IVF 索引参数
type IVFParams struct {
	Nlist  int
	Nprobe int
}

// VectorIndexConfig 向量索引配置
type VectorIndexConfig struct {
	IndexName   string
	TableName   string
	VectorField string
	IndexType   VectorIndexType
	MetricType  MetricType
	HNSWParams  *HNSWParams
	IVFParams   *IVFParams
}

// VectorStore 向量存储接口
type VectorStore interface {
	// Insert 插入一条记忆
	Insert(ctx context.Context, memory *Memory) error

	// Search 向量搜索
	Search(ctx context.Context, embedding []float64, opts *SearchOptions) ([]*Memory, error)

	// Get 根据 ID 获取记忆
	Get(ctx context.Context, id int64) (*Memory, error)

	// Update 更新记忆
	Update(ctx context.Context, id int64, content string, embedding []float64) (*Memory, error)

	// Delete 删除记忆
	Delete(ctx context.Context, id int64) error

	// GetAll 获取所有记忆
	GetAll(ctx context.Context, opts *GetAllOptions) ([]*Memory, error)

	// DeleteAll 删除所有记忆
	DeleteAll(ctx context.Context, opts *DeleteAllOptions) error

	// Close 关闭连接
	Close() error

	// CreateIndex 创建向量索引
	CreateIndex(ctx context.Context, config *VectorIndexConfig) error
}

// SearchOptions 搜索选项
type SearchOptions struct {
	UserID   string
	AgentID  string
	Limit    int
	MinScore float64
	Filters  map[string]interface{}
}

// GetAllOptions 获取所有记忆的选项
type GetAllOptions struct {
	UserID  string
	AgentID string
	Limit   int
	Offset  int
}

// DeleteAllOptions 删除所有记忆的选项
type DeleteAllOptions struct {
	UserID  string
	AgentID string
}
