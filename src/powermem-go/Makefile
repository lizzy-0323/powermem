.PHONY: help build test clean install lint fmt examples

# 默认目标
.DEFAULT_GOAL := help

# 变量定义
BINARY_NAME=powermem-go
GO=go
GOFLAGS=-v
GOTEST=$(GO) test
GOLINT=golangci-lint

help: ## 显示帮助信息
	@echo "PowerMem Go SDK - Makefile 帮助"
	@echo ""
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## 构建项目
	@echo "构建项目..."
	$(GO) build $(GOFLAGS) ./...

test: ## 运行测试
	@echo "运行测试..."
	$(GOTEST) -v -race -cover ./...

test-unit: ## 运行单元测试
	@echo "运行单元测试..."
	$(GOTEST) -v ./tests/...

test-core: ## 运行核心功能测试
	@echo "运行核心功能测试..."
	$(GOTEST) -v ./tests/core/...

test-storage: ## 运行存储层测试
	@echo "运行存储层测试..."
	$(GOTEST) -v ./tests/storage/...

test-intelligence: ## 运行智能功能测试
	@echo "运行智能功能测试..."
	$(GOTEST) -v ./tests/intelligence/...

test-coverage: ## 运行测试并生成覆盖率报告
	@echo "生成测试覆盖率报告..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

bench: ## 运行基准测试
	@echo "运行基准测试..."
	$(GOTEST) -bench=. -benchmem ./...

clean: ## 清理构建产物
	@echo "清理构建产物..."
	$(GO) clean
	rm -f coverage.out coverage.html
	rm -rf bin/

install: ## 安装依赖
	@echo "安装依赖..."
	$(GO) mod download
	$(GO) mod tidy

lint: ## 运行代码检查
	@echo "运行代码检查..."
	@if command -v $(GOLINT) >/dev/null 2>&1; then \
		$(GOLINT) run ./...; \
	else \
		echo "golangci-lint 未安装，跳过..."; \
		echo "安装: brew install golangci-lint (macOS) 或访问 https://golangci-lint.run/usage/install/"; \
	fi

fmt: ## 格式化代码
	@echo "格式化代码..."
	$(GO) fmt ./...
	goimports -w .

vet: ## 运行 go vet
	@echo "运行 go vet..."
	$(GO) vet ./...

examples: ## 构建示例程序
	@echo "构建示例程序..."
	@mkdir -p bin/examples
	$(GO) build -o bin/examples/basic ./examples/basic
	$(GO) build -o bin/examples/advanced ./examples/advanced
	$(GO) build -o bin/examples/multi_agent ./examples/multi_agent
	@echo "示例程序已构建到 bin/examples/"

run-example-basic: ## 运行基础示例
	@echo "运行基础示例..."
	$(GO) run ./examples/basic/main.go

run-example-advanced: ## 运行高级示例
	@echo "运行高级示例..."
	$(GO) run ./examples/advanced/main.go

run-example-multi: ## 运行多代理示例
	@echo "运行多代理示例..."
	$(GO) run ./examples/multi_agent/main.go

mod-update: ## 更新依赖
	@echo "更新依赖..."
	$(GO) get -u ./...
	$(GO) mod tidy

check: fmt vet lint test ## 运行所有检查（格式化、vet、lint、测试）

all: clean install fmt vet build test ## 完整构建流程
