# PowerMem Go SDK

<div align="center">

[English](./README.md) | [简体中文](./README_CN.md)

[![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GitHub Issues](https://img.shields.io/github/issues/oceanbase/powermem)](https://github.com/oceanbase/powermem/issues/143)

</div>

PowerMem Go SDK 是 [PowerMem](https://github.com/oceanbase/powermem) 的 Go 语言实现，为 Go 开发者提供原生的智能记忆管理能力。

## ✨ 特性

- 🚀 **高性能**: 专为 Go 的高并发场景优化
- 🔐 **并发安全**: 所有操作都支持 context 和 goroutine 安全
- 🎯 **完整功能**: 与 Python SDK 功能对等
- 🔌 **易于集成**: 简洁的 API 设计，易于集成到现有项目
- 🧠 **智能去重**: 自动检测和合并相似记忆
- 📊 **多种存储**: 支持 OceanBase、SQLite、PostgreSQL
- 🤖 **多代理支持**: 完善的多代理记忆管理
- 📈 **Ebbinghaus 算法**: 内置遗忘曲线算法

## 📦 安装

```bash
go get github.com/oceanbase/powermem-go
```

## 🚀 快速开始

### 基础使用

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/oceanbase/powermem-go/pkg/powermem"
)

func main() {
    // 1. 配置
    config := &powermem.Config{
        LLM: powermem.LLMConfig{
            Provider: "openai",
            APIKey:   "your-api-key",
            Model:    "gpt-4",
        },
        Embedder: powermem.EmbedderConfig{
            Provider: "openai",
            APIKey:   "your-api-key",
            Model:    "text-embedding-3-small",
        },
        VectorStore: powermem.VectorStoreConfig{
            Provider: "oceanbase",
            Config: map[string]interface{}{
                "host":                 "127.0.0.1",
                "port":                 2881,
                "user":                 "root@sys",
                "password":             "password",
                "db_name":              "powermem",
                "collection_name":      "memories",
                "embedding_model_dims": 1536,
            },
        },
    }
    
    // 2. 创建客户端
    client, err := powermem.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    ctx := context.Background()
    
    // 3. 添加记忆
    memory, err := client.Add(ctx, "User likes Python programming",
        powermem.WithUserID("user123"),
        powermem.WithMetadata(map[string]interface{}{
            "category": "preference",
        }),
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Added memory: %d\n", memory.ID)
    
    // 4. 搜索记忆
    results, err := client.Search(ctx, "user preferences",
        powermem.WithUserID("user123"),
        powermem.WithLimit(5),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    for _, mem := range results {
        fmt.Printf("- %s (score: %.3f)\n", mem.Content, mem.Score)
    }
}
```

### 使用环境变量配置

创建 `.env` 文件：

```bash
# .env
LLM_PROVIDER=openai
LLM_API_KEY=your-api-key
LLM_MODEL=gpt-4

EMBEDDING_PROVIDER=openai
EMBEDDING_API_KEY=your-api-key
EMBEDDING_MODEL=text-embedding-3-small

VECTOR_STORE_PROVIDER=oceanbase
VECTOR_STORE_HOST=127.0.0.1
VECTOR_STORE_PORT=2881
VECTOR_STORE_USER=root@sys
VECTOR_STORE_PASSWORD=password
VECTOR_STORE_DB=powermem
VECTOR_STORE_COLLECTION=memories
```

然后在代码中：

```go
config, err := powermem.LoadConfigFromEnv()
if err != nil {
    log.Fatal(err)
}

client, err := powermem.NewClient(config)
```

## 📚 核心功能

### 1. 记忆管理

```go
// 添加记忆
memory, err := client.Add(ctx, "content",
    powermem.WithUserID("user123"),
    powermem.WithAgentID("agent1"),
    powermem.WithMetadata(map[string]interface{}{"key": "value"}),
)

// 搜索记忆
results, err := client.Search(ctx, "query",
    powermem.WithUserID("user123"),
    powermem.WithLimit(10),
    powermem.WithMinScore(0.7),
)

// 更新记忆
updated, err := client.Update(ctx, memoryID, "new content")

// 删除记忆
err := client.Delete(ctx, memoryID)

// 获取所有记忆
memories, err := client.GetAll(ctx,
    powermem.WithUserID("user123"),
    powermem.WithLimit(100),
)

// 删除所有记忆
err := client.DeleteAll(ctx, powermem.WithUserID("user123"))
```

### 2. 智能去重

```go
config.Intelligence = &powermem.IntelligenceConfig{
    Enabled:             true,
    DecayRate:           0.1,
    ReinforcementFactor: 0.3,
    DuplicateThreshold:  0.95,
}

// 添加记忆时启用智能去重
memory, err := client.Add(ctx, "content",
    powermem.WithUserID("user123"),
    powermem.WithInfer(true), // 启用去重
)
```

### 3. 多代理支持

```go
// Agent1 添加私有记忆
_, err := client.Add(ctx, "Agent1's private data",
    powermem.WithAgentID("agent1"),
    powermem.WithUserID("user123"),
    powermem.WithScope(powermem.ScopePrivate),
)

// Agent2 添加共享记忆
_, err := client.Add(ctx, "Shared knowledge",
    powermem.WithAgentID("agent2"),
    powermem.WithUserID("user123"),
    powermem.WithScope(powermem.ScopeAgentGroup),
)

// Agent1 搜索（只能看到自己的私有记忆 + 共享记忆）
results, _ := client.Search(ctx, "query",
    powermem.WithAgentID("agent1"),
    powermem.WithUserID("user123"),
)
```

### 4. 高级搜索

```go
// 带过滤器的搜索
results, err := client.Search(ctx, "query",
    powermem.WithUserID("user123"),
    powermem.WithFilters(map[string]interface{}{
        "category": "important",
    }),
    powermem.WithMinScore(0.8),
    powermem.WithLimit(10),
)
```

## 🏗️ 架构

```
powermem-go/
├── pkg/
│   ├── powermem/          # 核心客户端
│   │   ├── memory.go      # Memory 客户端
│   │   ├── config.go      # 配置管理
│   │   ├── types.go       # 类型定义
│   │   ├── options.go     # 选项模式
│   │   └── errors.go      # 错误定义
│   ├── storage/           # 存储层
│   │   ├── base.go        # 存储接口
│   │   └── oceanbase/     # OceanBase 实现
│   ├── llm/               # LLM 提供商
│   │   ├── base.go        # LLM 接口
│   │   └── openai/        # OpenAI 实现
│   ├── embedder/          # Embedder 提供商
│   │   ├── base.go        # Embedder 接口
│   │   └── openai/        # OpenAI 实现
│   └── intelligence/      # 智能功能
│       ├── dedup.go       # 去重逻辑
│       └── ebbinghaus.go  # 遗忘曲线
└── examples/              # 示例代码
    ├── basic/             # 基础示例
    ├── advanced/          # 高级示例
    └── multi_agent/       # 多代理示例
```

## 📖 示例

查看 `examples/` 目录获取完整示例：

- **基础使用**: `examples/basic/main.go`
- **高级功能**: `examples/advanced/main.go`
- **多代理协作**: `examples/multi_agent/main.go`

运行示例：

```bash
cd examples/basic
go run main.go
```

## 🔧 配置选项

### LLM 提供商

目前支持：

- OpenAI (`openai`)

即将支持：

- Qwen (`qwen`)
- Anthropic (`anthropic`)
- Gemini (`gemini`)
- Ollama (`ollama`)

### Embedder 提供商

目前支持：

- OpenAI (`openai`)

即将支持：

- Qwen (`qwen`)
- HuggingFace (`huggingface`)
- Ollama (`ollama`)

### 向量存储

目前支持：

- OceanBase (`oceanbase`)

即将支持：

- SQLite (`sqlite`)
- PostgreSQL (`postgres`)

## 🧪 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./pkg/powermem

# 运行测试并显示覆盖率
go test -cover ./...
```

## 🤝 贡献

欢迎贡献！请查看 [CONTRIBUTING.md](../../CONTRIBUTING.md) 了解详情。

## 📄 许可证

Apache License 2.0 - 详见 [LICENSE](../../LICENSE)

## 🔗 相关链接

- [PowerMem Python SDK](https://github.com/oceanbase/powermem)
- [Issue #143](https://github.com/oceanbase/powermem/issues/143)
- [文档](https://powermem.oceanbase.com)
- [OceanBase](https://github.com/oceanbase/oceanbase)

## 💬 联系我们

- 提交 Issue: [GitHub Issues](https://github.com/oceanbase/powermem/issues)
- 加入讨论: [GitHub Discussions](https://github.com/oceanbase/powermem/discussions)

---

Made with ❤️ by the OceanBase team
