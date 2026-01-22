package core

// AddOption Add 操作的选项
type AddOption func(*AddOptions)

// AddOptions Add 操作的配置
type AddOptions struct {
	UserID   string
	AgentID  string
	Metadata map[string]interface{}
	Infer    bool // 是否启用智能去重
	Scope    MemoryScope
}

// WithUserID 设置用户 ID（Add 操作）
func WithUserID(userID string) AddOption {
	return func(opts *AddOptions) {
		opts.UserID = userID
	}
}

// WithUserIDForSearch 设置用户 ID（Search 操作）
func WithUserIDForSearch(userID string) SearchOption {
	return func(opts *SearchOptions) {
		opts.UserID = userID
	}
}

// WithUserIDForGetAll 设置用户 ID（GetAll 操作）
func WithUserIDForGetAll(userID string) GetAllOption {
	return func(opts *GetAllOptions) {
		opts.UserID = userID
	}
}

// WithUserIDForDeleteAll 设置用户 ID（DeleteAll 操作）
func WithUserIDForDeleteAll(userID string) DeleteAllOption {
	return func(opts *DeleteAllOptions) {
		opts.UserID = userID
	}
}

// WithAgentID 设置代理 ID（Add 操作）
func WithAgentID(agentID string) AddOption {
	return func(opts *AddOptions) {
		opts.AgentID = agentID
	}
}

// WithAgentIDForSearch 设置代理 ID（Search 操作）
func WithAgentIDForSearch(agentID string) SearchOption {
	return func(opts *SearchOptions) {
		opts.AgentID = agentID
	}
}

// WithAgentIDForGetAll 设置代理 ID（GetAll 操作）
func WithAgentIDForGetAll(agentID string) GetAllOption {
	return func(opts *GetAllOptions) {
		opts.AgentID = agentID
	}
}

// WithAgentIDForDeleteAll 设置代理 ID（DeleteAll 操作）
func WithAgentIDForDeleteAll(agentID string) DeleteAllOption {
	return func(opts *DeleteAllOptions) {
		opts.AgentID = agentID
	}
}

// WithMetadata 设置元数据
func WithMetadata(metadata map[string]interface{}) AddOption {
	return func(opts *AddOptions) {
		opts.Metadata = metadata
	}
}

// WithInfer 启用/禁用智能去重
func WithInfer(infer bool) AddOption {
	return func(opts *AddOptions) {
		opts.Infer = infer
	}
}

// WithScope 设置记忆作用域
func WithScope(scope MemoryScope) AddOption {
	return func(opts *AddOptions) {
		opts.Scope = scope
	}
}

// SearchOption Search 操作的选项
type SearchOption func(*SearchOptions)

// SearchOptions Search 操作的配置
type SearchOptions struct {
	UserID          string
	AgentID         string
	Limit           int
	Filters         map[string]interface{}
	MinScore        float64
	IncludeArchived bool
}

// WithLimit 设置返回结果数量限制（Search 操作）
func WithLimit(limit int) SearchOption {
	return func(opts *SearchOptions) {
		opts.Limit = limit
	}
}

// WithLimitForGetAll 设置返回结果数量限制（GetAll 操作）
func WithLimitForGetAll(limit int) GetAllOption {
	return func(opts *GetAllOptions) {
		opts.Limit = limit
	}
}

// WithFilters 设置过滤条件
func WithFilters(filters map[string]interface{}) SearchOption {
	return func(opts *SearchOptions) {
		opts.Filters = filters
	}
}

// WithMinScore 设置最小相似度分数
func WithMinScore(score float64) SearchOption {
	return func(opts *SearchOptions) {
		opts.MinScore = score
	}
}

// WithIncludeArchived 是否包含已归档的记忆
func WithIncludeArchived(include bool) SearchOption {
	return func(opts *SearchOptions) {
		opts.IncludeArchived = include
	}
}

// GetAllOption GetAll 操作的选项
type GetAllOption func(*GetAllOptions)

// GetAllOptions GetAll 操作的配置
type GetAllOptions struct {
	UserID  string
	AgentID string
	Limit   int
	Offset  int
}

// WithOffset 设置偏移量
func WithOffset(offset int) GetAllOption {
	return func(opts *GetAllOptions) {
		opts.Offset = offset
	}
}

// DeleteAllOption DeleteAll 操作的选项
type DeleteAllOption func(*DeleteAllOptions)

// DeleteAllOptions DeleteAll 操作的配置
type DeleteAllOptions struct {
	UserID  string
	AgentID string
}

// applyAddOptions 应用 Add 选项
func applyAddOptions(opts []AddOption) *AddOptions {
	options := &AddOptions{
		Infer:    false,
		Scope:    ScopePrivate,
		Metadata: make(map[string]interface{}),
	}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// applySearchOptions 应用 Search 选项
func applySearchOptions(opts []SearchOption) *SearchOptions {
	options := &SearchOptions{
		Limit:    10,
		MinScore: 0.0,
	}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// applyGetAllOptions 应用 GetAll 选项
func applyGetAllOptions(opts []GetAllOption) *GetAllOptions {
	options := &GetAllOptions{
		Limit:  100,
		Offset: 0,
	}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// applyDeleteAllOptions 应用 DeleteAll 选项
func applyDeleteAllOptions(opts []DeleteAllOption) *DeleteAllOptions {
	options := &DeleteAllOptions{}
	for _, opt := range opts {
		opt(options)
	}
	return options
}
