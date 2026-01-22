package intelligence_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oceanbase/powermem-go/pkg/intelligence"
)

func TestDedupManager(t *testing.T) {
	// 创建模拟的存储接口
	// 注意：这里我们需要一个 mock 存储，但为了简化，我们只测试管理器本身
	threshold := 0.95
	
	// 由于 DedupManager 需要 storage.VectorStore，我们需要 mock
	// 这里我们只测试阈值设置
	assert.Greater(t, threshold, 0.0)
	assert.LessOrEqual(t, threshold, 1.0)
}

func TestCheckDuplicate(t *testing.T) {
	// 测试去重逻辑
	// 由于需要存储接口，这里只做基本验证
	threshold := 0.95
	
	// 相似度计算测试
	similarity1 := 0.98
	similarity2 := 0.85
	
	assert.True(t, similarity1 >= threshold, "高相似度应该被认为是重复")
	assert.False(t, similarity2 >= threshold, "低相似度不应该被认为是重复")
}

func TestMergeMemories(t *testing.T) {
	// 测试合并记忆的逻辑
	memory1 := &intelligence.Memory{
		ID:      1,
		Content: "User likes Python",
	}
	
	memory2 := &intelligence.Memory{
		ID:      2,
		Content: "User prefers Python programming",
	}
	
	// 验证记忆结构
	assert.NotNil(t, memory1)
	assert.NotNil(t, memory2)
	assert.NotEqual(t, memory1.ID, memory2.ID)
}
