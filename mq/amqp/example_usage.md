# AMQP 配置使用示例

本文档展示了如何使用修正后的 `mergeConfigs` 函数和 `mq.Open()` 方法。

## 基本用法

### 1. 使用默认配置

```go
package main

import (
    "log"
    "github.com/hysios/x/mq"
)

func main() {
    // 使用空配置，全部使用默认值
    driver, err := mq.Open("amqp", mq.Config{})
    if err != nil {
        log.Fatal(err)
    }
    defer driver.Close()
    
    // DefaultConfig 的值：
    // URL: "amqp://guest:guest@localhost:5672/"
    // ExchangeName: "events"
    // PublishTimeout: 5 * time.Second
    // Durable: true
    // QueueName: "" (空字符串)
}
```

### 2. 部分配置覆盖

```go
// 只设置 URL，其他使用默认值
driver, err := mq.Open("amqp", mq.Config{
    "url": "amqp://myuser:mypass@myserver.com:5672/",
})

// 结果配置：
// URL: "amqp://myuser:mypass@myserver.com:5672/" (用户设置)
// ExchangeName: "events" (默认值)
// PublishTimeout: 5s (默认值)
// Durable: true (默认值)
```

### 3. 零值覆盖默认值

```go
// 使用零值覆盖默认的非零值
driver, err := mq.Open("amqp", mq.Config{
    "durable": false,         // 覆盖默认的 true
    "publish_timeout": "0s",  // 覆盖默认的 5s
})

// 结果配置：
// URL: "amqp://guest:guest@localhost:5672/" (默认值)
// ExchangeName: "events" (默认值)
// PublishTimeout: 0s (用户设置的零值)
// Durable: false (用户设置的零值)
```

### 4. 完整自定义配置

```go
driver, err := mq.Open("amqp", mq.Config{
    "url":             "amqp://prod:secret@prod.rabbitmq.com:5672/",
    "exchange_name":   "production_events",
    "queue_name":      "service_queue",
    "publish_timeout": "30s",
    "durable":         true,
})
```

## 真实场景示例

### 开发环境

```go
// 开发环境：使用本地 RabbitMQ，全默认配置
devDriver, err := mq.Open("amqp", mq.Config{})
```

### 测试环境

```go
// 测试环境：自定义 exchange，其他保持默认
testDriver, err := mq.Open("amqp", mq.Config{
    "exchange_name": "test_events",
})
```

### 生产环境

```go
// 生产环境：完整自定义配置
prodDriver, err := mq.Open("amqp", mq.Config{
    "url":             "amqp://prod_user:prod_pass@prod.rabbitmq.com:5672/",
    "exchange_name":   "prod_events",
    "publish_timeout": "60s",
    "durable":         true,
})
```

### 高性能场景

```go
// 高性能场景：关闭持久化以提高性能
fastDriver, err := mq.Open("amqp", mq.Config{
    "durable":         false, // 零值覆盖默认的 true
    "publish_timeout": "1s",  // 较短的超时
})
```

## 配置字段说明

| 字段名            | 类型   | 默认值                                 | 说明              |
|-------------------|--------|----------------------------------------|-----------------|
| `url`             | string | `"amqp://guest:guest@localhost:5672/"` | RabbitMQ 连接 URL |
| `exchange_name`   | string | `"events"`                             | Exchange 名称     |
| `queue_name`      | string | `""` (空字符串)                        | 队列名称          |
| `publish_timeout` | string | `"5s"`                                 | 发布超时时间      |
| `durable`         | bool   | `true`                                 | 是否持久化        |

## 重要特性

### ✅ 零值覆盖支持

修正后的 `mergeConfigs` 函数支持用户设置的零值正确覆盖默认的非零值：

```go
// 这会正确地将 durable 设置为 false，而不是保持默认的 true
mq.Open("amqp", mq.Config{
    "durable": false,
})
```

### ✅ 键名匹配

配置键名必须精确匹配（区分大小写）：

```go
// ✅ 正确
mq.Config{"url": "..."}

// ❌ 错误 - 大小写不匹配
mq.Config{"URL": "..."}
```

### ✅ 安全的 nil 处理

函数安全地处理 nil 配置：

```go
// 不会 panic，使用默认配置
mq.Open("amqp", nil)
```

## 测试验证

所有功能都有完整的测试覆盖：

- ✅ DefaultConfig URL 的正确使用
- ✅ 用户配置覆盖默认配置  
- ✅ 零值正确覆盖非零默认值
- ✅ nil 值的安全处理
- ✅ 键名大小写敏感性
- ✅ 部分配置合并
- ✅ 空配置处理
- ✅ 真实使用场景

运行测试：
```bash
go test ./mq/amqp/ -v
```
