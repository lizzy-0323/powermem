package core

import (
	"time"
)

// Memory 代表一条记忆
type Memory struct {
	ID                int64                  `json:"id"`
	UserID            string                 `json:"user_id"`
	AgentID           string                 `json:"agent_id,omitempty"`
	Content           string                 `json:"content"`
	Embedding         []float64              `json:"embedding,omitempty"`
	SparseEmbedding   map[int]float64        `json:"sparse_embedding,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	RetentionStrength float64                `json:"retention_strength"`
	LastAccessedAt    *time.Time             `json:"last_accessed_at,omitempty"`
	Score             float64                `json:"score,omitempty"` // 搜索相关性分数
}

// MemoryScope 定义记忆的作用域
type MemoryScope string

const (
	ScopePrivate    MemoryScope = "private"     // 私有（单个 agent）
	ScopeAgentGroup MemoryScope = "agent_group" // 代理组共享
	ScopeGlobal     MemoryScope = "global"      // 全局共享
)

// MetricType 定义向量距离度量类型
type MetricType string

const (
	MetricCosine MetricType = "cosine"
	MetricL2     MetricType = "l2"
	MetricIP     MetricType = "ip" // 内积
)

// VectorIndexType 定义向量索引类型
type VectorIndexType string

const (
	IndexTypeHNSW    VectorIndexType = "HNSW"
	IndexTypeIVFFlat VectorIndexType = "IVF_FLAT"
	IndexTypeIVFPQ   VectorIndexType = "IVF_PQ"
)

// HNSWParams HNSW 索引参数
type HNSWParams struct {
	M              int // 每个节点的最大连接数
	EfConstruction int // 构建时的搜索深度
	EfSearch       int // 搜索时的深度
}

// IVFParams IVF 索引参数
type IVFParams struct {
	Nlist  int // 聚类中心数量
	Nprobe int // 搜索时探测的聚类数量
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

// SearchResult 搜索结果
type SearchResult struct {
	Memories   []*Memory
	TotalCount int
}
