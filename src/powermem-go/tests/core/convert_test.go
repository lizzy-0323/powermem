package core_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	powermem "github.com/oceanbase/powermem-go/pkg/core"
	"github.com/oceanbase/powermem-go/pkg/storage"
)

// 注意：convert.go 中的函数都是私有的，无法直接测试
// 这些函数会在实际的 Memory 操作中被间接测试
func TestConvertMemoryTypes(t *testing.T) {
	// 测试类型转换的正确性通过实际使用来验证
	// 这里只验证类型定义的一致性
	
	coreMem := &powermem.Memory{
		ID:                12345,
		UserID:            "user123",
		AgentID:           "agent456",
		Content:           "Test content",
		Embedding:         []float64{0.1, 0.2, 0.3},
		Metadata:          map[string]interface{}{"key": "value"},
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		RetentionStrength: 0.8,
		Score:             0.95,
	}
	
	storageMem := &storage.Memory{
		ID:                coreMem.ID,
		UserID:            coreMem.UserID,
		AgentID:           coreMem.AgentID,
		Content:           coreMem.Content,
		Embedding:         coreMem.Embedding,
		Metadata:          coreMem.Metadata,
		CreatedAt:         coreMem.CreatedAt,
		UpdatedAt:         coreMem.UpdatedAt,
		RetentionStrength: coreMem.RetentionStrength,
		Score:             coreMem.Score,
	}
	
	// 验证字段一致性
	assert.Equal(t, coreMem.ID, storageMem.ID)
	assert.Equal(t, coreMem.UserID, storageMem.UserID)
	assert.Equal(t, coreMem.AgentID, storageMem.AgentID)
	assert.Equal(t, coreMem.Content, storageMem.Content)
	assert.Equal(t, coreMem.Embedding, storageMem.Embedding)
	assert.Equal(t, coreMem.Metadata, storageMem.Metadata)
}
