# PowerMem Go SDK

<div align="center">

[English](./README.md) | [简体中文](./README_CN.md)

[![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GitHub Issues](https://img.shields.io/github/issues/oceanbase/powermem)](https://github.com/oceanbase/powermem/issues/143)

</div>

PowerMem Go SDK 是 [PowerMem](https://github.com/oceanbase/powermem) 的 Go 语言实现，为 Go 开发者提供原生的智能记忆管理能力。

## ✨ 特性

- 🚀 **高性能**: 专为 Go 的高并发场景优化，适合微服务架构
- 🔐 **并发安全**: 所有操作都支持 context 和 goroutine 安全
- 🎯 **完整功能**: 与 Python SDK 功能对等
- 🔌 **易于集成**: 简洁的 API 设计，易于集成到 Gin、Echo、Fiber 等 Web 框架
- 🧠 **智能去重**: 自动检测和合并相似记忆
- 📊 **多种存储**: 支持 OceanBase、SQLite、PostgreSQL
- 🤖 **多代理支持**: 完善的多代理记忆管理和协作
- 📈 **Ebbinghaus 算法**: 内置遗忘曲线算法，智能管理记忆强度

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
    memory, err := client.Add(ctx, "用户喜欢Python编程",
        powermem.WithUserID("user123"),
        powermem.WithMetadata(map[string]interface{}{
            "category": "preference",
        }),
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("添加记忆: %d\n", memory.ID)
    
    // 4. 搜索记忆
    results, err := client.Search(ctx, "用户的偏好是什么",
        powermem.WithUserID("user123"),
        powermem.WithLimit(5),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    for _, mem := range results {
        fmt.Printf("- %s (相关度: %.3f)\n", mem.Content, mem.Score)
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

INTELLIGENCE_ENABLED=true
```

然后在代码中加载：

```go
// 从环境变量加载配置
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
memory, err := client.Add(ctx, "记忆内容",
    powermem.WithUserID("user123"),
    powermem.WithAgentID("agent1"),
    powermem.WithMetadata(map[string]interface{}{"key": "value"}),
)

// 搜索记忆
results, err := client.Search(ctx, "查询关键词",
    powermem.WithUserID("user123"),
    powermem.WithLimit(10),
    powermem.WithMinScore(0.7),
)

// 更新记忆
updated, err := client.Update(ctx, memoryID, "新的内容")

// 删除记忆
err := client.Delete(ctx, memoryID)

// 获取所有记忆（支持分页）
memories, err := client.GetAll(ctx,
    powermem.WithUserID("user123"),
    powermem.WithLimit(100),
    powermem.WithOffset(0),
)

// 删除所有记忆
err := client.DeleteAll(ctx, powermem.WithUserID("user123"))
```

### 2. 智能去重

启用智能去重功能，自动检测和合并相似记忆：

```go
// 配置智能记忆
config.Intelligence = &powermem.IntelligenceConfig{
    Enabled:             true,
    DecayRate:           0.1,  // 遗忘曲线衰减率
    ReinforcementFactor: 0.3,  // 强化因子
    DuplicateThreshold:  0.95, // 去重相似度阈值
}

// 添加记忆时启用智能去重
memory1, _ := client.Add(ctx, "用户喜欢Python编程",
    powermem.WithUserID("user123"),
    powermem.WithInfer(true), // 启用去重
)

// 尝试添加相似记忆，会自动合并
memory2, _ := client.Add(ctx, "用户喜欢Python开发",
    powermem.WithUserID("user123"),
    powermem.WithInfer(true),
)

// memory1.ID == memory2.ID 表示记忆被合并
```

### 3. 多代理支持

支持多个 AI 代理共享或隔离记忆：

```go
// Agent1 添加私有记忆（只有自己能看到）
_, err := client.Add(ctx, "Agent1的私有数据",
    powermem.WithAgentID("agent1"),
    powermem.WithUserID("user123"),
    powermem.WithScope(powermem.ScopePrivate), // 私有作用域
)

// Agent2 添加代理组共享记忆（组内代理都能看到）
_, err := client.Add(ctx, "共享知识",
    powermem.WithAgentID("agent2"),
    powermem.WithUserID("user123"),
    powermem.WithScope(powermem.ScopeAgentGroup), // 代理组共享
)

// Agent3 添加全局共享记忆（所有代理都能看到）
_, err := client.Add(ctx, "全局配置信息",
    powermem.WithAgentID("agent3"),
    powermem.WithUserID("user123"),
    powermem.WithScope(powermem.ScopeGlobal), // 全局共享
)

// Agent1 搜索（只能看到自己的私有记忆 + 共享记忆）
results, _ := client.Search(ctx, "数据",
    powermem.WithAgentID("agent1"),
    powermem.WithUserID("user123"),
)
```

### 4. 高级搜索

支持基于元数据的过滤搜索：

```go
// 添加带元数据的记忆
_, err := client.Add(ctx, "完成了Go语言学习",
    powermem.WithUserID("user123"),
    powermem.WithMetadata(map[string]interface{}{
        "type":       "achievement",
        "importance": "high",
        "date":       "2024-01-15",
    }),
)

// 带过滤器的搜索
results, err := client.Search(ctx, "学习",
    powermem.WithUserID("user123"),
    powermem.WithFilters(map[string]interface{}{
        "type": "achievement", // 只搜索成就类型的记忆
    }),
    powermem.WithMinScore(0.8), // 最低相关度分数
    powermem.WithLimit(10),
)

for _, mem := range results {
    fmt.Printf("- [%.3f] %s\n", mem.Score, mem.Content)
    fmt.Printf("  元数据: %v\n", mem.Metadata)
}
```

## 🏗️ 架构设计

```
powermem-go/
├── pkg/
│   ├── powermem/          # 核心客户端
│   │   ├── memory.go      # Memory 客户端主接口
│   │   ├── config.go      # 配置管理
│   │   ├── types.go       # 核心类型定义
│   │   ├── options.go     # 选项模式
│   │   └── errors.go      # 错误定义
│   ├── storage/           # 存储层
│   │   ├── base.go        # 存储接口定义
│   │   └── oceanbase/     # OceanBase 向量存储实现
│   │       ├── client.go  # 客户端实现
│   │       ├── vector.go  # 向量操作
│   │       └── utils.go   # 工具函数
│   ├── llm/               # LLM 提供商
│   │   ├── base.go        # LLM 接口
│   │   └── openai/        # OpenAI 实现
│   ├── embedder/          # Embedder 提供商
│   │   ├── base.go        # Embedder 接口
│   │   └── openai/        # OpenAI Embedding 实现
│   └── intelligence/      # 智能功能
│       ├── dedup.go       # 去重逻辑
│       └── ebbinghaus.go  # Ebbinghaus 遗忘曲线
└── examples/              # 示例代码
    ├── basic/             # 基础使用示例
    ├── advanced/          # 高级功能示例
    └── multi_agent/       # 多代理协作示例
```

## 📖 示例代码

查看 `examples/` 目录获取完整示例：

### 基础使用

```bash
cd examples/basic
go run main.go
```

### 高级功能（智能去重、过滤搜索等）

```bash
cd examples/advanced
go run main.go
```

### 多代理协作场景

```bash
cd examples/multi_agent
go run main.go
```

## 🔧 配置选项

### LLM 提供商

**当前支持：**

- OpenAI (`openai`)

**即将支持：**

- 通义千问 (`qwen`)
- Anthropic (`anthropic`)
- Google Gemini (`gemini`)
- Ollama (`ollama`)

### Embedder 提供商

**当前支持：**

- OpenAI Embeddings (`openai`)

**即将支持：**

- 通义千问 Embeddings (`qwen`)
- HuggingFace (`huggingface`)
- Ollama (`ollama`)

### 向量存储

**当前支持：**

- OceanBase (`oceanbase`) - 推荐，支持向量索引和混合搜索

**即将支持：**

- SQLite (`sqlite`) - 轻量级本地存储
- PostgreSQL + pgvector (`postgres`) - 开源向量数据库

## 🧪 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./pkg/powermem

# 运行测试并显示覆盖率
go test -cover ./...

# 运行集成测试
go test -tags=integration ./...
```

## 🚀 性能优化

PowerMem Go SDK 针对高并发场景进行了优化：

- 使用读写锁（`sync.RWMutex`）保护共享资源
- 支持 context 取消和超时控制
- 连接池管理数据库连接
- 批量操作支持
- 向量索引优化（HNSW、IVF）

## 🤝 贡献

我们欢迎所有形式的贡献！

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

详见 [CONTRIBUTING.md](../../CONTRIBUTING.md)

## 📄 许可证

本项目采用 Apache License 2.0 许可证 - 详见 [LICENSE](../../LICENSE) 文件

## 🔗 相关链接

- [PowerMem Python SDK](https://github.com/oceanbase/powermem)
- [GitHub Issue #143](https://github.com/oceanbase/powermem/issues/143)
- [在线文档](https://powermem.oceanbase.com)
- [OceanBase 数据库](https://github.com/oceanbase/oceanbase)
- [API 参考](https://pkg.go.dev/github.com/oceanbase/powermem-go)

## 💬 社区与支持

- 💬 提交 Issue: [GitHub Issues](https://github.com/oceanbase/powermem/issues)
- 💬 参与讨论: [GitHub Discussions](https://github.com/oceanbase/powermem/discussions)
- 📧 邮件列表: <powermem@oceanbase.com>
- 🐦 Twitter: [@OceanBase](https://twitter.com/OceanBase)

## 🙏 致谢

感谢所有为 PowerMem 项目做出贡献的开发者！

特别感谢：

- OceanBase 团队提供强大的向量数据库支持
- Go 社区提供优秀的开源库

---

<div align="center">
Made with ❤️ by the OceanBase Team
</div>
