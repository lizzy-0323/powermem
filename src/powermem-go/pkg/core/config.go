package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

// Config PowerMem 主配置
type Config struct {
	LLM          LLMConfig           `json:"llm"`
	Embedder     EmbedderConfig      `json:"embedder"`
	VectorStore  VectorStoreConfig   `json:"vector_store"`
	Intelligence *IntelligenceConfig `json:"intelligence,omitempty"`
	AgentMemory  *AgentMemoryConfig  `json:"agent_memory,omitempty"`
}

// LLMConfig LLM 配置
type LLMConfig struct {
	Provider   string                 `json:"provider"` // openai, qwen, anthropic, gemini, ollama
	APIKey     string                 `json:"api_key"`
	Model      string                 `json:"model"`
	BaseURL    string                 `json:"base_url,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// EmbedderConfig Embedder 配置
type EmbedderConfig struct {
	Provider   string                 `json:"provider"` // openai, qwen, huggingface, ollama
	APIKey     string                 `json:"api_key"`
	Model      string                 `json:"model"`
	BaseURL    string                 `json:"base_url,omitempty"`
	Dimensions int                    `json:"dimensions,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// VectorStoreConfig 向量存储配置
type VectorStoreConfig struct {
	Provider string                 `json:"provider"` // oceanbase, sqlite, postgres
	Config   map[string]interface{} `json:"config"`
}

// IntelligenceConfig 智能记忆配置
type IntelligenceConfig struct {
	Enabled             bool    `json:"enabled"`
	DecayRate           float64 `json:"decay_rate"`           // 遗忘曲线衰减率
	ReinforcementFactor float64 `json:"reinforcement_factor"` // 强化因子
	DuplicateThreshold  float64 `json:"duplicate_threshold"`  // 去重相似度阈值
}

// AgentMemoryConfig 多代理记忆配置
type AgentMemoryConfig struct {
	DefaultScope          MemoryScope `json:"default_scope"`
	AllowCrossAgentAccess bool        `json:"allow_cross_agent_access"`
	CollaborationLevel    string      `json:"collaboration_level"` // none, read_only, full
}

// LoadConfigFromEnv 从环境变量加载配置
// 支持从 .env 文件或 .env.example 文件加载
// 会自动向上查找 .env 文件（最多5层目录）
func LoadConfigFromEnv() (*Config, error) {
	// 使用 FindEnvFile 查找 .env 文件（支持向上查找）
	envPath, found := FindEnvFile()
	if found {
		// 如果找到 .env 文件，加载它
		_ = godotenv.Load(envPath)
	} else {
		// 如果找不到，尝试从当前目录加载（godotenv 的默认行为）
		_ = godotenv.Load()
	}

	// 获取向量存储提供者
	provider := getEnvOrDefault("VECTOR_STORE_PROVIDER", "sqlite")

	// 根据提供者构建不同的配置
	vectorStoreConfig := make(map[string]interface{})

	switch provider {
	case "oceanbase":
		port, _ := strconv.Atoi(getEnvOrDefault("VECTOR_STORE_PORT", "2881"))
		dims, _ := strconv.Atoi(getEnvOrDefault("VECTOR_STORE_EMBEDDING_MODEL_DIMS", "1536"))
		vectorStoreConfig = map[string]interface{}{
			"host":                 getEnvOrDefault("VECTOR_STORE_HOST", "127.0.0.1"),
			"port":                 port,
			"user":                 getEnvOrDefault("VECTOR_STORE_USER", "root@sys"),
			"password":             os.Getenv("VECTOR_STORE_PASSWORD"),
			"db_name":              getEnvOrDefault("VECTOR_STORE_DB", "powermem"),
			"collection_name":      getEnvOrDefault("VECTOR_STORE_COLLECTION", "memories"),
			"embedding_model_dims": dims,
		}
	case "sqlite":
		dims, _ := strconv.Atoi(getEnvOrDefault("VECTOR_STORE_EMBEDDING_MODEL_DIMS", "1536"))
		vectorStoreConfig = map[string]interface{}{
			"db_path":              getEnvOrDefault("VECTOR_STORE_DB_PATH", "./powermem.db"),
			"collection_name":      getEnvOrDefault("VECTOR_STORE_COLLECTION", "memories"),
			"embedding_model_dims": dims,
		}
	case "postgres":
		port, _ := strconv.Atoi(getEnvOrDefault("VECTOR_STORE_PORT", "5432"))
		dims, _ := strconv.Atoi(getEnvOrDefault("VECTOR_STORE_EMBEDDING_MODEL_DIMS", "1536"))
		vectorStoreConfig = map[string]interface{}{
			"host":                 getEnvOrDefault("VECTOR_STORE_HOST", "localhost"),
			"port":                 port,
			"user":                 getEnvOrDefault("VECTOR_STORE_USER", "postgres"),
			"password":             os.Getenv("VECTOR_STORE_PASSWORD"),
			"db_name":              getEnvOrDefault("VECTOR_STORE_DB", "powermem"),
			"collection_name":      getEnvOrDefault("VECTOR_STORE_COLLECTION", "memories"),
			"embedding_model_dims": dims,
			"ssl_mode":             getEnvOrDefault("VECTOR_STORE_SSL_MODE", "disable"),
		}
	}

	// 获取 LLM provider 以确定使用哪个 base URL 环境变量和默认模型
	llmProvider := getEnvOrDefault("LLM_PROVIDER", "openai")
	var llmBaseURL string
	var defaultModel string

	switch llmProvider {
	case "deepseek":
		llmBaseURL = os.Getenv("DEEPSEEK_LLM_BASE_URL")
		if llmBaseURL == "" {
			llmBaseURL = "https://api.deepseek.com"
		}
		defaultModel = "deepseek-chat"
	case "qwen":
		defaultModel = "qwen-plus"
	case "ollama":
		llmBaseURL = os.Getenv("OLLAMA_LLM_BASE_URL")
		if llmBaseURL == "" {
			llmBaseURL = "http://localhost:11434"
		}
		defaultModel = "llama3.1:70b"
	case "anthropic":
		llmBaseURL = os.Getenv("ANTHROPIC_LLM_BASE_URL")
		if llmBaseURL == "" {
			llmBaseURL = "https://api.anthropic.com"
		}
		defaultModel = "claude-3-5-sonnet-20240620"
	default:
		llmBaseURL = os.Getenv("LLM_BASE_URL")
		defaultModel = "gpt-4"
	}

	// 使用 Python SDK 风格的环境变量命名：EMBEDDING_*
	embedderProvider := getEnvOrDefault("EMBEDDING_PROVIDER", "qwen")
	embedderAPIKey := os.Getenv("EMBEDDING_API_KEY")
	embedderModel := os.Getenv("EMBEDDING_MODEL")

	// 根据 provider 设置默认 base URL
	var embedderFinalBaseURL string
	switch embedderProvider {
	case "qwen":
		embedderFinalBaseURL = os.Getenv("QWEN_EMBEDDING_BASE_URL")
		if embedderFinalBaseURL == "" {
			embedderFinalBaseURL = "https://dashscope.aliyuncs.com/api/v1"
		}
		if embedderModel == "" {
			embedderModel = "text-embedding-v4"
		}
	case "openai":
		embedderFinalBaseURL = os.Getenv("OPENAI_EMBEDDING_BASE_URL")
		if embedderFinalBaseURL == "" {
			embedderFinalBaseURL = "https://api.openai.com/v1"
		}
		if embedderModel == "" {
			embedderModel = "text-embedding-3-small"
		}
	default:
		embedderFinalBaseURL = os.Getenv("EMBEDDING_BASE_URL")
		if embedderModel == "" {
			embedderModel = "text-embedding-3-small"
		}
	}

	config := &Config{
		LLM: LLMConfig{
			Provider: llmProvider,
			APIKey:   os.Getenv("LLM_API_KEY"),
			Model:    getEnvOrDefault("LLM_MODEL", defaultModel),
			BaseURL:  llmBaseURL,
		},
		Embedder: EmbedderConfig{
			Provider: embedderProvider,
			APIKey:   embedderAPIKey,
			Model:    embedderModel,
			BaseURL:  embedderFinalBaseURL,
		},
		VectorStore: VectorStoreConfig{
			Provider: provider,
			Config:   vectorStoreConfig,
		},
	}

	// 智能记忆配置（可选）
	if os.Getenv("INTELLIGENCE_ENABLED") == "true" {
		config.Intelligence = &IntelligenceConfig{
			Enabled:             true,
			DecayRate:           0.1,
			ReinforcementFactor: 0.3,
			DuplicateThreshold:  0.95,
		}
	}

	return config, nil
}

// LoadConfigFromEnvFile 从指定的 .env 文件加载配置
func LoadConfigFromEnvFile(envPath string) (*Config, error) {
	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}
	return LoadConfigFromEnv()
}

// LoadConfigFromJSON 从 JSON 文件加载配置
func LoadConfigFromJSON(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, NewMemoryError("LoadConfigFromJSON", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, NewMemoryError("LoadConfigFromJSON", err)
	}

	return &config, nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.LLM.Provider == "" {
		return NewMemoryError("Validate", ErrInvalidConfig)
	}
	if c.Embedder.Provider == "" {
		return NewMemoryError("Validate", ErrInvalidConfig)
	}
	if c.VectorStore.Provider == "" {
		return NewMemoryError("Validate", ErrInvalidConfig)
	}
	return nil
}

// getEnvOrDefault 获取环境变量或返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// FindEnvFile 查找 .env 或 .env.example 文件
func FindEnvFile() (string, bool) {
	// 首先检查当前目录
	if _, err := os.Stat(".env"); err == nil {
		return ".env", true
	}
	if _, err := os.Stat(".env.example"); err == nil {
		return ".env.example", true
	}

	// 检查项目根目录（向上查找）
	dir, _ := os.Getwd()
	for i := 0; i < 5; i++ {
		envPath := filepath.Join(dir, ".env")
		envExamplePath := filepath.Join(dir, ".env.example")

		if _, err := os.Stat(envPath); err == nil {
			return envPath, true
		}
		if _, err := os.Stat(envExamplePath); err == nil {
			return envExamplePath, true
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", false
}
