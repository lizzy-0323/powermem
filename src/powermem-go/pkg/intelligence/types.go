package intelligence

import "time"

// Memory 记忆类型（避免循环依赖）
type Memory struct {
	ID                int64
	UserID            string
	AgentID           string
	Content           string
	Embedding         []float64
	Metadata          map[string]interface{}
	CreatedAt         time.Time
	UpdatedAt         time.Time
	RetentionStrength float64
	LastAccessedAt    *time.Time
	Score             float64
}
