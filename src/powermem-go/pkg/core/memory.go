package core

import (
	"context"
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/oceanbase/powermem-go/pkg/embedder"
	openaiEmbedder "github.com/oceanbase/powermem-go/pkg/embedder/openai"
	qwenEmbedder "github.com/oceanbase/powermem-go/pkg/embedder/qwen"
	"github.com/oceanbase/powermem-go/pkg/intelligence"
	"github.com/oceanbase/powermem-go/pkg/llm"
	anthropicLLM "github.com/oceanbase/powermem-go/pkg/llm/anthropic"
	deepseekLLM "github.com/oceanbase/powermem-go/pkg/llm/deepseek"
	ollamaLLM "github.com/oceanbase/powermem-go/pkg/llm/ollama"
	openaiLLM "github.com/oceanbase/powermem-go/pkg/llm/openai"
	qwenLLM "github.com/oceanbase/powermem-go/pkg/llm/qwen"
	"github.com/oceanbase/powermem-go/pkg/storage"
	"github.com/oceanbase/powermem-go/pkg/storage/oceanbase"
	postgresStore "github.com/oceanbase/powermem-go/pkg/storage/postgres"
	sqliteStore "github.com/oceanbase/powermem-go/pkg/storage/sqlite"
)

// Client PowerMem 客户端
type Client struct {
	config            *Config
	storage           storage.VectorStore
	llm               llm.Provider
	embedder          embedder.Provider
	dedupManager      *intelligence.DedupManager
	ebbinghausManager *intelligence.EbbinghausManager
	snowflakeNode     *snowflake.Node
	mu                sync.RWMutex
}

// NewClient 创建新的 PowerMem 客户端
func NewClient(cfg *Config) (*Client, error) {
	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// 初始化存储
	store, err := initStorage(cfg.VectorStore)
	if err != nil {
		return nil, err
	}

	// 初始化 LLM
	llmProvider, err := initLLM(cfg.LLM)
	if err != nil {
		return nil, err
	}

	// 初始化 Embedder
	embedderProvider, err := initEmbedder(cfg.Embedder)
	if err != nil {
		return nil, err
	}

	// 初始化 Snowflake ID 生成器
	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, NewMemoryError("NewClient", err)
	}

	client := &Client{
		config:        cfg,
		storage:       store,
		llm:           llmProvider,
		embedder:      embedderProvider,
		snowflakeNode: node,
	}

	// 初始化智能功能（如果启用）
	if cfg.Intelligence != nil && cfg.Intelligence.Enabled {
		client.dedupManager = intelligence.NewDedupManager(
			store,
			cfg.Intelligence.DuplicateThreshold,
		)
		client.ebbinghausManager = intelligence.NewEbbinghausManager(
			cfg.Intelligence.DecayRate,
			cfg.Intelligence.ReinforcementFactor,
		)
	}

	return client, nil
}

// Add 添加记忆
func (c *Client) Add(ctx context.Context, content string, opts ...AddOption) (*Memory, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 应用选项
	addOpts := applyAddOptions(opts)

	// 检查 context
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// 1. 生成 embedding
	embedding, err := c.embedder.Embed(ctx, content)
	if err != nil {
		return nil, NewMemoryError("Add", err)
	}

	// 2. 智能去重（如果启用）
	if addOpts.Infer && c.dedupManager != nil {
		isDup, existingID, err := c.dedupManager.CheckDuplicate(ctx, embedding, addOpts.UserID, addOpts.AgentID)
		if err != nil {
			return nil, NewMemoryError("Add", err)
		}
		if isDup {
			// 合并记忆
			merged, err := c.dedupManager.MergeMemories(ctx, existingID, content, embedding)
			if err != nil {
				return nil, NewMemoryError("Add", err)
			}
			// 转换回powermem.Memory类型
			return fromIntelligenceMemory(merged), nil
		}
	}

	// 3. 插入存储
	memory := &Memory{
		ID:                c.snowflakeNode.Generate().Int64(),
		UserID:            addOpts.UserID,
		AgentID:           addOpts.AgentID,
		Content:           content,
		Embedding:         embedding,
		Metadata:          addOpts.Metadata,
		RetentionStrength: 1.0, // 初始强度为 1.0
	}

	if err := c.storage.Insert(ctx, toStorageMemory(memory)); err != nil {
		return nil, NewMemoryError("Add", err)
	}

	return memory, nil
}

// Search 搜索记忆
func (c *Client) Search(ctx context.Context, query string, opts ...SearchOption) ([]*Memory, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 应用选项
	searchOpts := applySearchOptions(opts)

	// 1. 生成查询 embedding
	queryEmbedding, err := c.embedder.Embed(ctx, query)
	if err != nil {
		return nil, NewMemoryError("Search", err)
	}

	// 2. 执行向量搜索
	storageOpts := &storage.SearchOptions{
		UserID:   searchOpts.UserID,
		AgentID:  searchOpts.AgentID,
		Limit:    searchOpts.Limit,
		MinScore: searchOpts.MinScore,
		Filters:  searchOpts.Filters,
	}

	memories, err := c.storage.Search(ctx, queryEmbedding, storageOpts)
	if err != nil {
		return nil, NewMemoryError("Search", err)
	}

	return fromStorageMemories(memories), nil
}

// Get 根据 ID 获取记忆
func (c *Client) Get(ctx context.Context, id int64) (*Memory, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	memory, err := c.storage.Get(ctx, id)
	if err != nil {
		return nil, NewMemoryError("Get", err)
	}

	return fromStorageMemory(memory), nil
}

// Update 更新记忆
func (c *Client) Update(ctx context.Context, id int64, content string) (*Memory, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 生成新的 embedding
	embedding, err := c.embedder.Embed(ctx, content)
	if err != nil {
		return nil, NewMemoryError("Update", err)
	}

	// 更新存储
	memory, err := c.storage.Update(ctx, id, content, embedding)
	if err != nil {
		return nil, NewMemoryError("Update", err)
	}

	return fromStorageMemory(memory), nil
}

// Delete 删除记忆
func (c *Client) Delete(ctx context.Context, id int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.storage.Delete(ctx, id); err != nil {
		return NewMemoryError("Delete", err)
	}

	return nil
}

// GetAll 获取所有记忆
func (c *Client) GetAll(ctx context.Context, opts ...GetAllOption) ([]*Memory, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	getAllOpts := applyGetAllOptions(opts)

	storageOpts := &storage.GetAllOptions{
		UserID:  getAllOpts.UserID,
		AgentID: getAllOpts.AgentID,
		Limit:   getAllOpts.Limit,
		Offset:  getAllOpts.Offset,
	}

	memories, err := c.storage.GetAll(ctx, storageOpts)
	if err != nil {
		return nil, NewMemoryError("GetAll", err)
	}

	return fromStorageMemories(memories), nil
}

// DeleteAll 删除所有记忆
func (c *Client) DeleteAll(ctx context.Context, opts ...DeleteAllOption) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	deleteAllOpts := applyDeleteAllOptions(opts)

	storageOpts := &storage.DeleteAllOptions{
		UserID:  deleteAllOpts.UserID,
		AgentID: deleteAllOpts.AgentID,
	}

	if err := c.storage.DeleteAll(ctx, storageOpts); err != nil {
		return NewMemoryError("DeleteAll", err)
	}

	return nil
}

// Close 关闭客户端
func (c *Client) Close() error {
	var errs []error

	if c.storage != nil {
		if err := c.storage.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if c.llm != nil {
		if err := c.llm.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if c.embedder != nil {
		if err := c.embedder.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs[0] // 返回第一个错误
	}

	return nil
}

// initStorage 初始化存储
func initStorage(cfg VectorStoreConfig) (storage.VectorStore, error) {
	switch cfg.Provider {
	case "oceanbase":
		return oceanbase.NewClient(&oceanbase.Config{
			Host:               cfg.Config["host"].(string),
			Port:               cfg.Config["port"].(int),
			User:               cfg.Config["user"].(string),
			Password:           cfg.Config["password"].(string),
			DBName:             cfg.Config["db_name"].(string),
			CollectionName:     cfg.Config["collection_name"].(string),
			EmbeddingModelDims: cfg.Config["embedding_model_dims"].(int),
		})
	case "sqlite":
		return sqliteStore.NewClient(&sqliteStore.Config{
			DBPath:             cfg.Config["db_path"].(string),
			CollectionName:     cfg.Config["collection_name"].(string),
			EmbeddingModelDims: cfg.Config["embedding_model_dims"].(int),
		})
	case "postgres":
		sslMode := "disable"
		if s, ok := cfg.Config["ssl_mode"].(string); ok {
			sslMode = s
		}
		return postgresStore.NewClient(&postgresStore.Config{
			Host:               cfg.Config["host"].(string),
			Port:               cfg.Config["port"].(int),
			User:               cfg.Config["user"].(string),
			Password:           cfg.Config["password"].(string),
			DBName:             cfg.Config["db_name"].(string),
			CollectionName:     cfg.Config["collection_name"].(string),
			EmbeddingModelDims: cfg.Config["embedding_model_dims"].(int),
			SSLMode:            sslMode,
		})
	default:
		return nil, NewMemoryError("initStorage", ErrInvalidConfig)
	}
}

// initLLM 初始化 LLM
func initLLM(cfg LLMConfig) (llm.Provider, error) {
	switch cfg.Provider {
	case "openai":
		return openaiLLM.NewClient(&openaiLLM.Config{
			APIKey:  cfg.APIKey,
			Model:   cfg.Model,
			BaseURL: cfg.BaseURL,
		})
	case "qwen":
		return qwenLLM.NewClient(&qwenLLM.Config{
			APIKey:  cfg.APIKey,
			Model:   cfg.Model,
			BaseURL: cfg.BaseURL,
		})
	case "deepseek":
		return deepseekLLM.NewClient(&deepseekLLM.Config{
			APIKey:  cfg.APIKey,
			Model:   cfg.Model,
			BaseURL: cfg.BaseURL,
		})
	case "ollama":
		return ollamaLLM.NewClient(&ollamaLLM.Config{
			APIKey:  cfg.APIKey,
			Model:   cfg.Model,
			BaseURL: cfg.BaseURL,
		})
	case "anthropic":
		return anthropicLLM.NewClient(&anthropicLLM.Config{
			APIKey:  cfg.APIKey,
			Model:   cfg.Model,
			BaseURL: cfg.BaseURL,
		})
	default:
		return nil, NewMemoryError("initLLM", ErrInvalidConfig)
	}
}

// initEmbedder 初始化 Embedder
func initEmbedder(cfg EmbedderConfig) (embedder.Provider, error) {
	switch cfg.Provider {
	case "openai":
		return openaiEmbedder.NewClient(&openaiEmbedder.Config{
			APIKey:     cfg.APIKey,
			Model:      cfg.Model,
			BaseURL:    cfg.BaseURL,
			Dimensions: cfg.Dimensions,
		})
	case "qwen":
		return qwenEmbedder.NewClient(&qwenEmbedder.Config{
			APIKey:     cfg.APIKey,
			Model:      cfg.Model,
			BaseURL:    cfg.BaseURL,
			Dimensions: cfg.Dimensions,
		})
	default:
		return nil, NewMemoryError("initEmbedder", ErrInvalidConfig)
	}
}
