package core

import (
	"github.com/oceanbase/powermem-go/pkg/intelligence"
	"github.com/oceanbase/powermem-go/pkg/storage"
)

// toStorageMemory 将 powermem.Memory 转换为 storage.Memory
func toStorageMemory(m *Memory) *storage.Memory {
	return &storage.Memory{
		ID:                m.ID,
		UserID:            m.UserID,
		AgentID:           m.AgentID,
		Content:           m.Content,
		Embedding:         m.Embedding,
		SparseEmbedding:   m.SparseEmbedding,
		Metadata:          m.Metadata,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
		RetentionStrength: m.RetentionStrength,
		LastAccessedAt:    m.LastAccessedAt,
		Score:             m.Score,
	}
}

// fromStorageMemory 将 storage.Memory 转换为 powermem.Memory
func fromStorageMemory(m *storage.Memory) *Memory {
	return &Memory{
		ID:                m.ID,
		UserID:            m.UserID,
		AgentID:           m.AgentID,
		Content:           m.Content,
		Embedding:         m.Embedding,
		SparseEmbedding:   m.SparseEmbedding,
		Metadata:          m.Metadata,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
		RetentionStrength: m.RetentionStrength,
		LastAccessedAt:    m.LastAccessedAt,
		Score:             m.Score,
	}
}

// fromStorageMemories 批量转换 storage.Memory 到 powermem.Memory
func fromStorageMemories(memories []*storage.Memory) []*Memory {
	result := make([]*Memory, len(memories))
	for i, m := range memories {
		result[i] = fromStorageMemory(m)
	}
	return result
}

// fromIntelligenceMemory 将 intelligence.Memory 转换为 powermem.Memory
func fromIntelligenceMemory(m *intelligence.Memory) *Memory {
	return &Memory{
		ID:                m.ID,
		UserID:            m.UserID,
		AgentID:           m.AgentID,
		Content:           m.Content,
		Embedding:         m.Embedding,
		Metadata:          m.Metadata,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
		RetentionStrength: m.RetentionStrength,
		LastAccessedAt:    m.LastAccessedAt,
		Score:             m.Score,
	}
}
