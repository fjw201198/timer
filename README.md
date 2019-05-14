# timer
简单定时器for go

## 常量
```go
// TIMER_QUEUE_SIZE 最大超时事件队列长度
TIMER_QUEUE_SIZE = 1024
```

## 方法
1. 创建定时器对象
```go
func NewTimer() *Timer
```

2. 设置超时事件
```go
func (tm *Timer) SetTimeout(seconds int, callback func()) (timerId int, err error)
```

3. 取消超时事件
```go
func (tm *Timer) KillTimer(timerId int)
```
