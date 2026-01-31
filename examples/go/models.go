package main

import (
	"strings"
	"time"
)

type MemoryID string

func (m *MemoryID) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	*m = MemoryID(str)
	return nil
}

func (m MemoryID) String() string {
	return string(m)
}

type Memory struct {
	MemoryID  MemoryID               `json:"memory_id,omitempty"`
	Content   string                 `json:"content"`
	UserID    string                 `json:"user_id,omitempty"`
	AgentID   string                 `json:"agent_id,omitempty"`
	RunID     string                 `json:"run_id,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Filters   map[string]interface{} `json:"filters,omitempty"`
	Scope     string                 `json:"scope,omitempty"`
	Type      string                 `json:"memory_type,omitempty"`
	CreatedAt string                 `json:"created_at,omitempty"`
	UpdatedAt string                 `json:"updated_at,omitempty"`
}

type CreateMemoryRequest struct {
	Content    string                 `json:"content"`
	UserID     string                 `json:"user_id,omitempty"`
	AgentID    string                 `json:"agent_id,omitempty"`
	RunID      string                 `json:"run_id,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Filters    map[string]interface{} `json:"filters,omitempty"`
	Scope      string                 `json:"scope,omitempty"`
	MemoryType string                 `json:"memory_type,omitempty"`
	Infer      bool                   `json:"infer,omitempty"`
}

type UpdateMemoryRequest struct {
	Content  string                 `json:"content,omitempty"`
	UserID   string                 `json:"user_id,omitempty"`
	AgentID  string                 `json:"agent_id,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type SearchMemoryRequest struct {
	Query   string                 `json:"query"`
	UserID  string                 `json:"user_id,omitempty"`
	AgentID string                 `json:"agent_id,omitempty"`
	RunID   string                 `json:"run_id,omitempty"`
	Filters map[string]interface{} `json:"filters,omitempty"`
	Limit   int                    `json:"limit,omitempty"`
}

type ListMemoriesRequest struct {
	UserID  string `json:"user_id,omitempty"`
	AgentID string `json:"agent_id,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Offset  int    `json:"offset,omitempty"`
	SortBy  string `json:"sort_by,omitempty"`
	Order   string `json:"order,omitempty"`
}

type SearchResult struct {
	MemoryID  MemoryID               `json:"memory_id"`
	Content   string                 `json:"content"`
	Score     float64                `json:"score"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt string                 `json:"created_at,omitempty"`
	UpdatedAt string                 `json:"updated_at,omitempty"`
}

type APIResponse struct {
	Success   bool                   `json:"success"`
	Data      interface{}            `json:"data,omitempty"`
	Message   string                 `json:"message,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Error     map[string]interface{} `json:"error,omitempty"`
}

type MemoryListData struct {
	Memories []Memory `json:"memories"`
	Total    int      `json:"total"`
	Limit    int      `json:"limit"`
	Offset   int      `json:"offset"`
}

type SearchData struct {
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
	Query   string         `json:"query"`
}
