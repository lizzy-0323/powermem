package core_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	powermem "github.com/oceanbase/powermem-go/pkg/core"
)

func TestLoadConfigFromEnv(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
	}{
		{
			name: "valid config with SQLite",
			envVars: map[string]string{
				"VECTOR_STORE_PROVIDER": "sqlite",
				"VECTOR_STORE_DB_PATH":  "./test.db",
				"LLM_PROVIDER":          "openai",
				"LLM_API_KEY":           "test-key",
				"LLM_MODEL":             "gpt-4",
				"EMBEDDING_PROVIDER":    "openai",
				"EMBEDDING_API_KEY":     "test-key",
				"EMBEDDING_MODEL":       "text-embedding-3-small",
			},
			wantErr: false,
		},
		{
			name: "valid config with Qwen",
			envVars: map[string]string{
				"VECTOR_STORE_PROVIDER": "sqlite",
				"VECTOR_STORE_DB_PATH":  "./test.db",
				"LLM_PROVIDER":          "qwen",
				"LLM_API_KEY":           "test-key",
				"LLM_MODEL":             "qwen-plus",
				"EMBEDDING_PROVIDER":    "qwen",
				"EMBEDDING_API_KEY":     "test-key",
				"EMBEDDING_MODEL":       "text-embedding-v4",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置环境变量
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			}()

			config, err := powermem.LoadConfigFromEnv()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, config)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, config)
				assert.Equal(t, tt.envVars["VECTOR_STORE_PROVIDER"], config.VectorStore.Provider)
				assert.Equal(t, tt.envVars["LLM_PROVIDER"], config.LLM.Provider)
				assert.Equal(t, tt.envVars["EMBEDDING_PROVIDER"], config.Embedder.Provider)
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *powermem.Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &powermem.Config{
				LLM: powermem.LLMConfig{
					Provider: "openai",
					APIKey:   "test-key",
					Model:    "gpt-4",
				},
				Embedder: powermem.EmbedderConfig{
					Provider: "openai",
					APIKey:   "test-key",
					Model:    "text-embedding-3-small",
				},
				VectorStore: powermem.VectorStoreConfig{
					Provider: "sqlite",
					Config: map[string]interface{}{
						"db_path": "./test.db",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing LLM provider",
			config: &powermem.Config{
				LLM: powermem.LLMConfig{
					Provider: "",
				},
				Embedder: powermem.EmbedderConfig{
					Provider: "openai",
				},
				VectorStore: powermem.VectorStoreConfig{
					Provider: "sqlite",
				},
			},
			wantErr: true,
		},
		{
			name: "missing Embedder provider",
			config: &powermem.Config{
				LLM: powermem.LLMConfig{
					Provider: "openai",
				},
				Embedder: powermem.EmbedderConfig{
					Provider: "",
				},
				VectorStore: powermem.VectorStoreConfig{
					Provider: "sqlite",
				},
			},
			wantErr: true,
		},
		{
			name: "missing VectorStore provider",
			config: &powermem.Config{
				LLM: powermem.LLMConfig{
					Provider: "openai",
				},
				Embedder: powermem.EmbedderConfig{
					Provider: "openai",
				},
				VectorStore: powermem.VectorStoreConfig{
					Provider: "",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFindEnvFile(t *testing.T) {
	// 测试查找 .env 文件
	envPath, found := powermem.FindEnvFile()
	
	// 这个测试取决于实际的文件系统状态
	// 我们只验证函数不会 panic
	assert.NotNil(t, envPath)
	// found 可能是 true 或 false，取决于是否存在 .env 文件
	_ = found
}

func TestDefaultConfig(t *testing.T) {
	// 测试默认配置值
	config := &powermem.Config{
		LLM: powermem.LLMConfig{
			Provider: "openai",
			Model:    "gpt-4",
		},
		Embedder: powermem.EmbedderConfig{
			Provider: "openai",
			Model:    "text-embedding-3-small",
		},
		VectorStore: powermem.VectorStoreConfig{
			Provider: "sqlite",
			Config: map[string]interface{}{
				"db_path": "./test.db",
			},
		},
	}

	err := config.Validate()
	require.NoError(t, err)
	
	assert.Equal(t, "openai", config.LLM.Provider)
	assert.Equal(t, "openai", config.Embedder.Provider)
	assert.Equal(t, "sqlite", config.VectorStore.Provider)
}
