package intelligence

import (
	"context"
	"math"

	"github.com/oceanbase/powermem-go/pkg/storage"
)

// DedupManager 去重管理器
type DedupManager struct {
	store     storage.VectorStore
	threshold float64 // 相似度阈值
}

// NewDedupManager 创建去重管理器
func NewDedupManager(store storage.VectorStore, threshold float64) *DedupManager {
	if threshold == 0 {
		threshold = 0.95 // 默认阈值
	}
	return &DedupManager{
		store:     store,
		threshold: threshold,
	}
}

// CheckDuplicate 检查是否存在重复记忆
// 返回: (是否重复, 重复的记忆ID, 错误)
func (m *DedupManager) CheckDuplicate(ctx context.Context, embedding []float64, userID, agentID string) (bool, int64, error) {
	// 搜索相似记忆
	opts := &storage.SearchOptions{
		UserID:  userID,
		AgentID: agentID,
		Limit:   5, // 只查找前5个最相似的
	}

	memories, err := m.store.Search(ctx, embedding, opts)
	if err != nil {
		return false, 0, err
	}

	// 检查是否有相似度超过阈值的记忆
	for _, mem := range memories {
		if mem.Score >= m.threshold {
			return true, mem.ID, nil
		}
	}

	return false, 0, nil
}

// MergeMemories 合并两条记忆
func (m *DedupManager) MergeMemories(ctx context.Context, existingID int64, newContent string, newEmbedding []float64) (*Memory, error) {
	// 获取现有记忆
	existing, err := m.store.Get(ctx, existingID)
	if err != nil {
		return nil, err
	}

	// 简单合并策略：将新内容附加到现有内容
	// 更复杂的策略可以使用 LLM 来智能合并
	mergedContent := existing.Content + " " + newContent

	// 计算新的 embedding（取平均）
	mergedEmbedding := averageEmbeddings(existing.Embedding, newEmbedding)

	// 更新记忆
	updated, err := m.store.Update(ctx, existingID, mergedContent, mergedEmbedding)
	if err != nil {
		return nil, err
	}

	// 转换类型
	return &Memory{
		ID:                updated.ID,
		UserID:            updated.UserID,
		AgentID:           updated.AgentID,
		Content:           updated.Content,
		Embedding:         updated.Embedding,
		Metadata:          updated.Metadata,
		CreatedAt:         updated.CreatedAt,
		UpdatedAt:         updated.UpdatedAt,
		RetentionStrength: updated.RetentionStrength,
		LastAccessedAt:    updated.LastAccessedAt,
		Score:             updated.Score,
	}, nil
}

// averageEmbeddings 计算两个向量的平均值
func averageEmbeddings(e1, e2 []float64) []float64 {
	if len(e1) != len(e2) {
		return e1 // 如果维度不匹配，返回第一个
	}

	result := make([]float64, len(e1))
	for i := range e1 {
		result[i] = (e1[i] + e2[i]) / 2.0
	}

	// 归一化
	return normalizeVector(result)
}

// normalizeVector 归一化向量
func normalizeVector(v []float64) []float64 {
	var sum float64
	for _, val := range v {
		sum += val * val
	}
	norm := math.Sqrt(sum)

	if norm == 0 {
		return v
	}

	result := make([]float64, len(v))
	for i, val := range v {
		result[i] = val / norm
	}

	return result
}

// CosineSimilarity 计算余弦相似度
func CosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
