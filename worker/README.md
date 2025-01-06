# Worker 包

Worker包提供了一个灵活的工作任务处理框架，支持异步任务处理和分布式作业调度。

## 功能特性

- 支持异步任务处理
- 可配置的任务重试机制
- 内置Redis迭代器支持
- 错误处理机制
- 可扩展的worker接口

## 安装

```bash
go get github.com/hysios/x/worker
```

## 快速开始

### 基础用法
```go
package main
import (
    "context"
    "github.com/hysios/x/worker"
)

func main() {
    // 创建worker实例
    w := worker.New()
    // 注册任务处理函数
    w.Register("task-name", func(ctx context.Context, payload []byte) error {
        // 处理任务的逻辑
        return nil
    })
    // 启动worker
    if err := w.Start(context.Background()); err != nil {
        panic(err)
    }
}
```


### 使用Redis迭代器

```go
// 创建Redis迭代器
iter := worker.NewRedisIterator(redisClient, "queue-key")
// 使用迭代器处理任务
for iter.Next() {
task := iter.Value()
// 处理任务...
}

## 配置选项

Worker包支持多种配置选项：

- 重试次数
- 超时时间
- 并发数量
- 错误处理策略

## 错误处理

包内置了多种错误类型：

- `ErrWorkerNotFound`: worker未找到
- `ErrInvalidPayload`: 无效的任务数据
- `ErrTaskTimeout`: 任务执行超时

## 接口定义

Worker包定义了以下主要接口：

```go
type Worker interface {
    Start(context.Context) error
    Stop() error
    Register(name string, handler HandlerFunc)
    }
    type Iterator interface {
    Next() bool
    Value() interface{}
    Error() error
    Close() error
}
```


## 贡献指南

欢迎提交Issue和Pull Request来帮助改进这个包。

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件