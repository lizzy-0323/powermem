# PowerMem Go SDK 测试

本目录包含 PowerMem Go SDK 的单元测试。

## 测试结构

```
tests/
├── core/           # 核心功能测试
│   ├── config_test.go
│   ├── types_test.go
│   ├── errors_test.go
│   ├── options_test.go
│   └── convert_test.go
├── storage/        # 存储层测试
│   └── sqlite_test.go
└── intelligence/   # 智能功能测试
    ├── dedup_test.go
    └── ebbinghaus_test.go
```

## 运行测试

### 运行所有测试

```bash
go test ./tests/... -v
```

### 运行特定包的测试

```bash
# 运行核心功能测试
go test ./tests/core/... -v

# 运行存储层测试
go test ./tests/storage/... -v

# 运行智能功能测试
go test ./tests/intelligence/... -v
```

### 运行特定测试

```bash
# 运行特定测试函数
go test ./tests/core/... -v -run TestConfigValidate

# 运行所有配置相关测试
go test ./tests/core/... -v -run TestConfig
```

### 测试覆盖率

```bash
# 生成覆盖率报告
go test ./tests/... -coverprofile=coverage.out

# 查看覆盖率报告
go tool cover -html=coverage.out
```

## 测试说明

### 单元测试

单元测试用于测试各个模块的独立功能：

- **core**: 配置管理、类型定义、错误处理、选项模式
- **storage**: 存储层的 CRUD 操作和搜索功能
- **intelligence**: 智能去重和 Ebbinghaus 遗忘曲线算法

### 测试依赖

测试使用以下依赖：

- `github.com/stretchr/testify` - 测试断言库
- SQLite - 用于存储层测试（内存数据库）

### 注意事项

1. **SQLite 测试**: SQLite 测试会创建临时数据库文件，测试完成后会自动清理
2. **环境变量**: 某些测试可能需要设置环境变量，但大部分测试使用 mock 数据
3. **并发安全**: 测试确保所有操作都是并发安全的

## 添加新测试

添加新测试时，请遵循以下规范：

1. 测试文件命名：`*_test.go`
2. 测试函数命名：`Test*`
3. 使用 `testing.T` 进行测试
4. 使用 `testify/assert` 进行断言
5. 使用 `testify/require` 进行必须通过的断言

示例：

```go
func TestMyFeature(t *testing.T) {
    // Arrange
    input := "test"
    
    // Act
    result := MyFunction(input)
    
    // Assert
    assert.Equal(t, "expected", result)
}
```
