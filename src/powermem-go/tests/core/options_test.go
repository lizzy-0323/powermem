package core_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	powermem "github.com/oceanbase/powermem-go/pkg/core"
)

func TestAddOptions(t *testing.T) {
	// 测试 WithUserID
	opt := powermem.WithUserID("user123")
	assert.NotNil(t, opt)
	
	// 测试 WithAgentID
	opt = powermem.WithAgentID("agent456")
	assert.NotNil(t, opt)
	
	// 测试 WithMetadata
	metadata := map[string]interface{}{"key": "value"}
	opt = powermem.WithMetadata(metadata)
	assert.NotNil(t, opt)
}

func TestSearchOptions(t *testing.T) {
	// 测试 WithUserIDForSearch
	opt := powermem.WithUserIDForSearch("user123")
	assert.NotNil(t, opt)
	
	// 测试 WithAgentIDForSearch
	opt = powermem.WithAgentIDForSearch("agent456")
	assert.NotNil(t, opt)
	
	// 测试 WithLimit
	opt = powermem.WithLimit(10)
	assert.NotNil(t, opt)
	
	// 测试 WithFilters（SearchOptions 使用 Filters 而不是 MetadataFilter）
	filter := map[string]interface{}{"key": "value"}
	opt = powermem.WithFilters(filter)
	assert.NotNil(t, opt)
}

func TestGetAllOptions(t *testing.T) {
	// 测试 WithUserIDForGetAll
	opt := powermem.WithUserIDForGetAll("user123")
	assert.NotNil(t, opt)
	
	// 测试 WithAgentIDForGetAll
	opt = powermem.WithAgentIDForGetAll("agent456")
	assert.NotNil(t, opt)
	
	// 测试 WithLimitForGetAll
	opt = powermem.WithLimitForGetAll(10)
	assert.NotNil(t, opt)
	
	// 测试 WithOffset（GetAll 操作）
	opt = powermem.WithOffset(5)
	assert.NotNil(t, opt)
}

func TestDeleteAllOptions(t *testing.T) {
	// 测试 WithUserIDForDeleteAll
	opt := powermem.WithUserIDForDeleteAll("user123")
	assert.NotNil(t, opt)
	
	// 测试 WithAgentIDForDeleteAll
	opt = powermem.WithAgentIDForDeleteAll("agent456")
	assert.NotNil(t, opt)
}
