package counter

import (
	"sync"
	"time"
)

// Limiter 计数器限流器
type Limiter struct {
	mu        sync.Mutex
	limit     int           // 时间窗口内的最大操作次数
	interval  time.Duration // 时间窗口长度
	counter   int           // 当前计数
	startTime time.Time     // 时间窗口起始时间
}

// New 创建一个计数器限流器
func New(limit int, interval time.Duration) *Limiter {
	return &Limiter{
		limit:     limit,
		interval:  interval,
		startTime: time.Now(),
	}
}

// Allow 检查是否允许操作。如果允许，则增加计数器并返回 true，否则返回 false
func (rl *Limiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	// 如果当前时间超过时间窗口，重置计数器
	if now.Sub(rl.startTime) > rl.interval {
		rl.counter = 0
		rl.startTime = now
	}

	// 如果计数器小于限制值，允许操作
	if rl.counter < rl.limit {
		rl.counter++
		return true
	}

	// 否则，不允许操作
	return false
}

// Remaining 返回当前时间窗口内剩余的操作次数
func (rl *Limiter) Remaining() int {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	// 如果当前时间超过时间窗口，重置计数器
	if now.Sub(rl.startTime) > rl.interval {
		rl.counter = 0
		rl.startTime = now
	}

	return rl.limit - rl.counter
}
