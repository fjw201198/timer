# timer
简单定时器for go

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
