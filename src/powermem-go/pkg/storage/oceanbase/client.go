package oceanbase

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/oceanbase/powermem-go/pkg/storage"
)

// Client OceanBase 客户端
type Client struct {
	db             *sql.DB
	config         *Config
	collectionName string
}

// Config OceanBase 配置
type Config struct {
	Host               string
	Port               int
	User               string
	Password           string
	DBName             string
	CollectionName     string
	EmbeddingModelDims int
}

// NewClient 创建新的 OceanBase 客户端
func NewClient(cfg *Config) (*Client, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("NewOceanBaseClient: %w", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("NewOceanBaseClient: %w", err)
	}

	client := &Client{
		db:             db,
		config:         cfg,
		collectionName: cfg.CollectionName,
	}

	// 初始化表结构
	if err := client.initTables(context.Background()); err != nil {
		return nil, err
	}

	return client, nil
}

// initTables 初始化数据库表
func (c *Client) initTables(ctx context.Context) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id BIGINT PRIMARY KEY,
			user_id VARCHAR(255) NOT NULL,
			agent_id VARCHAR(255),
			content LONGTEXT NOT NULL,
			embedding VECTOR(%d) NOT NULL,
			metadata JSON,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			retention_strength FLOAT DEFAULT 1.0,
			last_accessed_at DATETIME,
			INDEX idx_user_agent (user_id, agent_id)
		)
	`, c.collectionName, c.config.EmbeddingModelDims)

	_, err := c.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("initTables: %w", err)
	}

	return nil
}

// Insert 插入记忆
func (c *Client) Insert(ctx context.Context, memory *storage.Memory) error {
	query := fmt.Sprintf(`
		INSERT INTO %s 
		(id, user_id, agent_id, content, embedding, metadata, created_at, retention_strength)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, c.collectionName)

	vectorStr := vectorToString(memory.Embedding)

	metadataJSON, err := json.Marshal(memory.Metadata)
	if err != nil {
		return fmt.Errorf("Insert: %w", err)
	}

	_, err = c.db.ExecContext(ctx, query,
		memory.ID,
		memory.UserID,
		memory.AgentID,
		memory.Content,
		vectorStr,
		metadataJSON,
		time.Now(),
		memory.RetentionStrength,
	)

	if err != nil {
		return fmt.Errorf("Insert: %w", err)
	}

	return nil
}

// Search 向量搜索
func (c *Client) Search(ctx context.Context, embedding []float64, opts *storage.SearchOptions) ([]*storage.Memory, error) {
	queryVectorStr := vectorToString(embedding)

	whereClause, args := buildWhereClause(opts.UserID, opts.AgentID, opts.Filters)

	query := fmt.Sprintf(`
		SELECT 
			id, user_id, agent_id, content, embedding, metadata,
			created_at, updated_at, retention_strength, last_accessed_at,
			cosine_distance(embedding, ?) as distance
		FROM %s
		%s
		ORDER BY distance ASC
		LIMIT ?
	`, c.collectionName, whereClause)

	// 将查询向量添加到参数列表的开头
	allArgs := append([]interface{}{queryVectorStr}, args...)
	allArgs = append(allArgs, opts.Limit)

	rows, err := c.db.QueryContext(ctx, query, allArgs...)
	if err != nil {
		return nil, fmt.Errorf("Search: %w", err)
	}
	defer rows.Close()

	return c.scanMemories(rows, true)
}

// Get 根据 ID 获取记忆
func (c *Client) Get(ctx context.Context, id int64) (*storage.Memory, error) {
	query := fmt.Sprintf(`
		SELECT id, user_id, agent_id, content, embedding, metadata,
		       created_at, updated_at, retention_strength, last_accessed_at
		FROM %s
		WHERE id = ?
	`, c.collectionName)

	row := c.db.QueryRowContext(ctx, query, id)

	memory, err := c.scanMemory(row)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("Get: not found")
	}
	if err != nil {
		return nil, fmt.Errorf("Get: %w", err)
	}

	return memory, nil
}

// Update 更新记忆
func (c *Client) Update(ctx context.Context, id int64, content string, embedding []float64) (*storage.Memory, error) {
	vectorStr := vectorToString(embedding)

	query := fmt.Sprintf(`
		UPDATE %s
		SET content = ?, embedding = ?, updated_at = ?
		WHERE id = ?
	`, c.collectionName)

	result, err := c.db.ExecContext(ctx, query, content, vectorStr, time.Now(), id)
	if err != nil {
		return nil, fmt.Errorf("Update: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("Update: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("Update: not found")
	}

	// 返回更新后的记忆
	return c.Get(ctx, id)
}

// Delete 删除记忆
func (c *Client) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", c.collectionName)

	result, err := c.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("Delete: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Delete: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Delete: not found")
	}

	return nil
}

// GetAll 获取所有记忆
func (c *Client) GetAll(ctx context.Context, opts *storage.GetAllOptions) ([]*storage.Memory, error) {
	whereClause, args := buildWhereClause(opts.UserID, opts.AgentID, nil)

	query := fmt.Sprintf(`
		SELECT id, user_id, agent_id, content, embedding, metadata,
		       created_at, updated_at, retention_strength, last_accessed_at
		FROM %s
		%s
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, c.collectionName, whereClause)

	args = append(args, opts.Limit, opts.Offset)

	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetAll: %w", err)
	}
	defer rows.Close()

	return c.scanMemories(rows, false)
}

// DeleteAll 删除所有记忆
func (c *Client) DeleteAll(ctx context.Context, opts *storage.DeleteAllOptions) error {
	whereClause, args := buildWhereClause(opts.UserID, opts.AgentID, nil)

	query := fmt.Sprintf("DELETE FROM %s %s", c.collectionName, whereClause)

	_, err := c.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("DeleteAll: %w", err)
	}

	return nil
}

// Close 关闭数据库连接
func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// CreateIndex 创建向量索引
func (c *Client) CreateIndex(ctx context.Context, config *storage.VectorIndexConfig) error {
	var query string

	switch config.IndexType {
	case storage.IndexTypeHNSW:
		query = fmt.Sprintf(`
			CREATE VECTOR INDEX %s ON %s (%s) WITH (
				index_type = HNSW,
				M = %d,
				efConstruction = %d,
				metric_type = %s
			)`,
			config.IndexName, config.TableName, config.VectorField,
			config.HNSWParams.M,
			config.HNSWParams.EfConstruction,
			config.MetricType,
		)
	case storage.IndexTypeIVFFlat:
		query = fmt.Sprintf(`
			CREATE VECTOR INDEX %s ON %s (%s) WITH (
				index_type = IVF_FLAT,
				nlist = %d,
				metric_type = %s
			)`,
			config.IndexName, config.TableName, config.VectorField,
			config.IVFParams.Nlist,
			config.MetricType,
		)
	default:
		return fmt.Errorf("CreateIndex: invalid index type")
	}

	_, err := c.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("CreateIndex: %w", err)
	}

	return nil
}

// scanMemory 扫描单条记忆
func (c *Client) scanMemory(row *sql.Row) (*storage.Memory, error) {
	var memory storage.Memory
	var embeddingStr string
	var metadataJSON []byte
	var lastAccessedAt sql.NullTime

	err := row.Scan(
		&memory.ID,
		&memory.UserID,
		&memory.AgentID,
		&memory.Content,
		&embeddingStr,
		&metadataJSON,
		&memory.CreatedAt,
		&memory.UpdatedAt,
		&memory.RetentionStrength,
		&lastAccessedAt,
	)
	if err != nil {
		return nil, err
	}

	// 解析 embedding
	if embeddingStr != "" {
		embedding, err := stringToVector(embeddingStr)
		if err != nil {
			return nil, err
		}
		memory.Embedding = embedding
	}

	// 解析 metadata
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &memory.Metadata); err != nil {
			return nil, err
		}
	}

	// 处理 last_accessed_at
	if lastAccessedAt.Valid {
		memory.LastAccessedAt = &lastAccessedAt.Time
	}

	return &memory, nil
}

// scanMemories 扫描多条记忆
func (c *Client) scanMemories(rows *sql.Rows, hasScore bool) ([]*storage.Memory, error) {
	var memories []*storage.Memory

	for rows.Next() {
		var memory storage.Memory
		var embeddingStr string
		var metadataJSON []byte
		var lastAccessedAt sql.NullTime
		var distance float64

		if hasScore {
			err := rows.Scan(
				&memory.ID,
				&memory.UserID,
				&memory.AgentID,
				&memory.Content,
				&embeddingStr,
				&metadataJSON,
				&memory.CreatedAt,
				&memory.UpdatedAt,
				&memory.RetentionStrength,
				&lastAccessedAt,
				&distance,
			)
			if err != nil {
				return nil, err
			}
			// 将距离转换为相似度分数（1 - distance）
			memory.Score = 1.0 - distance
		} else {
			err := rows.Scan(
				&memory.ID,
				&memory.UserID,
				&memory.AgentID,
				&memory.Content,
				&embeddingStr,
				&metadataJSON,
				&memory.CreatedAt,
				&memory.UpdatedAt,
				&memory.RetentionStrength,
				&lastAccessedAt,
			)
			if err != nil {
				return nil, err
			}
		}

		// 解析 embedding
		if embeddingStr != "" {
			embedding, err := stringToVector(embeddingStr)
			if err != nil {
				return nil, err
			}
			memory.Embedding = embedding
		}

		// 解析 metadata
		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &memory.Metadata); err != nil {
				return nil, err
			}
		}

		// 处理 last_accessed_at
		if lastAccessedAt.Valid {
			memory.LastAccessedAt = &lastAccessedAt.Time
		}

		memories = append(memories, &memory)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return memories, nil
}
