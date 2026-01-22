package storage_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oceanbase/powermem-go/pkg/storage"
	sqliteStore "github.com/oceanbase/powermem-go/pkg/storage/sqlite"
)

func setupSQLiteTest(t *testing.T) (storage.VectorStore, func()) {
	testDBPath := "./test_powermem.db"
	
	// 清理可能存在的测试数据库
	os.Remove(testDBPath)
	
	config := &sqliteStore.Config{
		DBPath:             testDBPath,
		CollectionName:     "memories",
		EmbeddingModelDims: 1536,
	}
	
	store, err := sqliteStore.NewClient(config)
	require.NoError(t, err)
	require.NotNil(t, store)
	
	cleanup := func() {
		store.Close()
		os.Remove(testDBPath)
	}
	
	return store, cleanup
}

func TestSQLiteClient_Insert(t *testing.T) {
	store, cleanup := setupSQLiteTest(t)
	defer cleanup()
	
	ctx := context.Background()
	
	memory := &storage.Memory{
		ID:        100, // SQLite 需要手动设置 ID
		UserID:    "test_user",
		AgentID:   "test_agent",
		Content:   "Test memory content",
		Embedding: []float64{0.1, 0.2, 0.3, 0.4, 0.5},
		Metadata:  map[string]interface{}{"key": "value"},
	}
	
	err := store.Insert(ctx, memory)
	assert.NoError(t, err)
	assert.Equal(t, int64(100), memory.ID)
}

func TestSQLiteClient_Get(t *testing.T) {
	store, cleanup := setupSQLiteTest(t)
	defer cleanup()
	
	ctx := context.Background()
	
	// 先插入一条记忆
	memory := &storage.Memory{
		ID:        1, // SQLite 需要手动设置 ID 或使用自增
		UserID:    "test_user",
		Content:   "Test memory content",
		Embedding: []float64{0.1, 0.2, 0.3},
	}
	
	err := store.Insert(ctx, memory)
	require.NoError(t, err)
	id := memory.ID
	
	// 获取记忆
	retrieved, err := store.Get(ctx, id)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, id, retrieved.ID)
	assert.Equal(t, "test_user", retrieved.UserID)
	assert.Equal(t, "Test memory content", retrieved.Content)
}

func TestSQLiteClient_Update(t *testing.T) {
	store, cleanup := setupSQLiteTest(t)
	defer cleanup()
	
	ctx := context.Background()
	
	// 先插入一条记忆
	memory := &storage.Memory{
		ID:        2,
		UserID:    "test_user",
		Content:   "Original content",
		Embedding: []float64{0.1, 0.2, 0.3},
	}
	
	err := store.Insert(ctx, memory)
	require.NoError(t, err)
	id := memory.ID
	
	// 更新记忆
	updatedContent := "Updated content"
	updatedEmbedding := []float64{0.2, 0.3, 0.4}
	
	updated, err := store.Update(ctx, id, updatedContent, updatedEmbedding)
	assert.NoError(t, err)
	assert.Equal(t, updatedContent, updated.Content)
	
	// 验证更新（已在上面验证）
	_ = updated
}

func TestSQLiteClient_Delete(t *testing.T) {
	store, cleanup := setupSQLiteTest(t)
	defer cleanup()
	
	ctx := context.Background()
	
	// 先插入一条记忆
	memory := &storage.Memory{
		ID:        3,
		UserID:    "test_user",
		Content:   "Test memory content",
		Embedding: []float64{0.1, 0.2, 0.3},
	}
	
	err := store.Insert(ctx, memory)
	require.NoError(t, err)
	id := memory.ID
	
	// 删除记忆
	err = store.Delete(ctx, id)
	assert.NoError(t, err)
	
	// 验证删除
	_, err = store.Get(ctx, id)
	assert.Error(t, err)
}

func TestSQLiteClient_Search(t *testing.T) {
	store, cleanup := setupSQLiteTest(t)
	defer cleanup()
	
	ctx := context.Background()
	
	// 插入几条记忆
	memories := []*storage.Memory{
		{
			ID:        4,
			UserID:    "test_user",
			Content:   "Python programming",
			Embedding: []float64{0.1, 0.2, 0.3, 0.4, 0.5},
		},
		{
			ID:        5,
			UserID:    "test_user",
			Content:   "Go programming",
			Embedding: []float64{0.2, 0.3, 0.4, 0.5, 0.6},
		},
		{
			ID:        6,
			UserID:    "test_user",
			Content:   "JavaScript programming",
			Embedding: []float64{0.3, 0.4, 0.5, 0.6, 0.7},
		},
	}
	
	for _, mem := range memories {
		err := store.Insert(ctx, mem)
		require.NoError(t, err)
	}
	
	// 搜索
	queryVector := []float64{0.15, 0.25, 0.35, 0.45, 0.55}
	options := &storage.SearchOptions{
		Limit: 2,
	}
	
	results, err := store.Search(ctx, queryVector, options)
	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.LessOrEqual(t, len(results), 2)
}

func TestSQLiteClient_GetAll(t *testing.T) {
	store, cleanup := setupSQLiteTest(t)
	defer cleanup()
	
	ctx := context.Background()
	
	// 插入几条记忆
	for i := 0; i < 3; i++ {
		memory := &storage.Memory{
			ID:        int64(10 + i),
			UserID:    "test_user",
			Content:   "Test memory",
			Embedding: []float64{0.1, 0.2, 0.3},
		}
		err := store.Insert(ctx, memory)
		require.NoError(t, err)
	}
	
	// 获取所有记忆
	options := &storage.GetAllOptions{
		Limit: 10,
	}
	
	results, err := store.GetAll(ctx, options)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 3)
}

func TestSQLiteClient_DeleteAll(t *testing.T) {
	store, cleanup := setupSQLiteTest(t)
	defer cleanup()
	
	ctx := context.Background()
	
	// 插入几条记忆
	for i := 0; i < 3; i++ {
		memory := &storage.Memory{
			ID:        int64(20 + i),
			UserID:    "test_user",
			Content:   "Test memory",
			Embedding: []float64{0.1, 0.2, 0.3},
		}
		err := store.Insert(ctx, memory)
		require.NoError(t, err)
	}
	
	// 删除所有记忆
	options := &storage.DeleteAllOptions{}
	err := store.DeleteAll(ctx, options)
	assert.NoError(t, err)
	
	// 验证删除
	getOptions := &storage.GetAllOptions{Limit: 10}
	results, err := store.GetAll(ctx, getOptions)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(results))
}
