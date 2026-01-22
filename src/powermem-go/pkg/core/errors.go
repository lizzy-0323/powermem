package core

import (
	"errors"
	"fmt"
)

// 预定义错误
var (
	ErrNotFound         = errors.New("memory not found")
	ErrInvalidConfig    = errors.New("invalid configuration")
	ErrConnectionFailed = errors.New("connection failed")
	ErrEmbeddingFailed  = errors.New("embedding generation failed")
	ErrDuplicateMemory  = errors.New("duplicate memory detected")
	ErrInvalidInput     = errors.New("invalid input")
	ErrStorageOperation = errors.New("storage operation failed")
	ErrLLMOperation     = errors.New("llm operation failed")
)

// MemoryError 包装的错误类型
type MemoryError struct {
	Op  string // 操作名称
	Err error  // 原始错误
}

func (e *MemoryError) Error() string {
	return fmt.Sprintf("powermem: %s: %v", e.Op, e.Err)
}

func (e *MemoryError) Unwrap() error {
	return e.Err
}

// NewMemoryError 创建新的 MemoryError
func NewMemoryError(op string, err error) error {
	if err == nil {
		return nil
	}
	return &MemoryError{
		Op:  op,
		Err: err,
	}
}
