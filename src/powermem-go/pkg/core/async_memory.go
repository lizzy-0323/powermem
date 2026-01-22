package core

import (
	"context"
	"sync"
)

// AsyncClient 异步 PowerMem 客户端
// 提供与同步 Client 相同的功能，但所有操作都在独立的 goroutine 中执行
// 适合需要并发处理多个操作的场景
type AsyncClient struct {
	*Client
	wg sync.WaitGroup
}

// NewAsyncClient 创建新的异步 PowerMem 客户端
// 参数:
//   - cfg: PowerMem 配置
//
// 返回:
//   - *AsyncClient: 异步客户端实例
//   - error: 如果配置无效或初始化失败则返回错误
func NewAsyncClient(cfg *Config) (*AsyncClient, error) {
	client, err := NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &AsyncClient{
		Client: client,
	}, nil
}

// AddAsync 异步添加记忆
// 在独立的 goroutine 中执行，通过 channel 返回结果
// 参数:
//   - ctx: 上下文，用于控制请求生命周期
//   - content: 要添加的记忆内容
//   - opts: 可选的添加选项（UserID、AgentID、Metadata 等）
//
// 返回:
//   - <-chan *MemoryResult: 接收结果的 channel，包含 Memory 和 error
func (ac *AsyncClient) AddAsync(ctx context.Context, content string, opts ...AddOption) <-chan *MemoryResult {
	resultChan := make(chan *MemoryResult, 1)
	ac.wg.Add(1)

	go func() {
		defer ac.wg.Done()
		memory, err := ac.Client.Add(ctx, content, opts...)
		resultChan <- &MemoryResult{
			Memory: memory,
			Error:  err,
		}
		close(resultChan)
	}()

	return resultChan
}

// SearchAsync 异步搜索记忆
// 在独立的 goroutine 中执行，通过 channel 返回结果
// 参数:
//   - ctx: 上下文，用于控制请求生命周期
//   - query: 搜索查询文本
//   - opts: 可选的搜索选项（UserID、AgentID、Limit、Filters 等）
//
// 返回:
//   - <-chan *AsyncSearchResult: 接收搜索结果的 channel，包含 Memories 和 error
func (ac *AsyncClient) SearchAsync(ctx context.Context, query string, opts ...SearchOption) <-chan *AsyncSearchResult {
	resultChan := make(chan *AsyncSearchResult, 1)
	ac.wg.Add(1)

	go func() {
		defer ac.wg.Done()
		memories, err := ac.Client.Search(ctx, query, opts...)
		resultChan <- &AsyncSearchResult{
			Memories: memories,
			Error:    err,
		}
		close(resultChan)
	}()

	return resultChan
}

// GetAsync 异步获取记忆
// 在独立的 goroutine 中执行，通过 channel 返回结果
// 参数:
//   - ctx: 上下文，用于控制请求生命周期
//   - id: 记忆 ID
//
// 返回:
//   - <-chan *MemoryResult: 接收结果的 channel，包含 Memory 和 error
func (ac *AsyncClient) GetAsync(ctx context.Context, id int64) <-chan *MemoryResult {
	resultChan := make(chan *MemoryResult, 1)
	ac.wg.Add(1)

	go func() {
		defer ac.wg.Done()
		memory, err := ac.Client.Get(ctx, id)
		resultChan <- &MemoryResult{
			Memory: memory,
			Error:  err,
		}
		close(resultChan)
	}()

	return resultChan
}

// UpdateAsync 异步更新记忆
// 在独立的 goroutine 中执行，通过 channel 返回结果
// 参数:
//   - ctx: 上下文，用于控制请求生命周期
//   - id: 记忆 ID
//   - content: 新的记忆内容
//
// 返回:
//   - <-chan *MemoryResult: 接收结果的 channel，包含 Memory 和 error
func (ac *AsyncClient) UpdateAsync(ctx context.Context, id int64, content string) <-chan *MemoryResult {
	resultChan := make(chan *MemoryResult, 1)
	ac.wg.Add(1)

	go func() {
		defer ac.wg.Done()
		memory, err := ac.Client.Update(ctx, id, content)
		resultChan <- &MemoryResult{
			Memory: memory,
			Error:  err,
		}
		close(resultChan)
	}()

	return resultChan
}

// DeleteAsync 异步删除记忆
// 在独立的 goroutine 中执行，通过 channel 返回结果
// 参数:
//   - ctx: 上下文，用于控制请求生命周期
//   - id: 记忆 ID
//
// 返回:
//   - <-chan error: 接收错误的 channel，如果删除成功则为 nil
func (ac *AsyncClient) DeleteAsync(ctx context.Context, id int64) <-chan error {
	errChan := make(chan error, 1)
	ac.wg.Add(1)

	go func() {
		defer ac.wg.Done()
		err := ac.Client.Delete(ctx, id)
		errChan <- err
		close(errChan)
	}()

	return errChan
}

// GetAllAsync 异步获取所有记忆
// 在独立的 goroutine 中执行，通过 channel 返回结果
// 参数:
//   - ctx: 上下文，用于控制请求生命周期
//   - opts: 可选的获取选项（UserID、AgentID、Limit、Offset 等）
//
// 返回:
//   - <-chan *AsyncGetAllResult: 接收结果的 channel，包含 Memories 和 error
func (ac *AsyncClient) GetAllAsync(ctx context.Context, opts ...GetAllOption) <-chan *AsyncGetAllResult {
	resultChan := make(chan *AsyncGetAllResult, 1)
	ac.wg.Add(1)

	go func() {
		defer ac.wg.Done()
		memories, err := ac.Client.GetAll(ctx, opts...)
		resultChan <- &AsyncGetAllResult{
			Memories: memories,
			Error:    err,
		}
		close(resultChan)
	}()

	return resultChan
}

// DeleteAllAsync 异步删除所有记忆
// 在独立的 goroutine 中执行，通过 channel 返回结果
// 参数:
//   - ctx: 上下文，用于控制请求生命周期
//   - opts: 可选的删除选项（UserID、AgentID 等）
//
// 返回:
//   - <-chan error: 接收错误的 channel，如果删除成功则为 nil
func (ac *AsyncClient) DeleteAllAsync(ctx context.Context, opts ...DeleteAllOption) <-chan error {
	errChan := make(chan error, 1)
	ac.wg.Add(1)

	go func() {
		defer ac.wg.Done()
		err := ac.Client.DeleteAll(ctx, opts...)
		errChan <- err
		close(errChan)
	}()

	return errChan
}

// Wait 等待所有异步操作完成
// 用于在程序退出前确保所有 goroutine 都已完成
func (ac *AsyncClient) Wait() {
	ac.wg.Wait()
}

// Close 关闭异步客户端
// 会先等待所有异步操作完成，然后关闭底层客户端
func (ac *AsyncClient) Close() error {
	ac.Wait()
	return ac.Client.Close()
}

// MemoryResult 内存操作结果
type MemoryResult struct {
	Memory *Memory
	Error  error
}

// AsyncSearchResult 异步搜索结果
type AsyncSearchResult struct {
	Memories []*Memory
	Error    error
}

// AsyncGetAllResult 异步获取所有结果
type AsyncGetAllResult struct {
	Memories []*Memory
	Error    error
}
