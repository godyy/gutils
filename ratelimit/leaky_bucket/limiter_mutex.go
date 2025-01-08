package leaky_bucket

import (
	"sync"
	"time"
)

// mutexLimiter 是一个基于互斥锁的漏桶算法实现。
type mutexLimiter struct {
	sync.Mutex
	last       time.Time
	sleepFor   time.Duration
	perRequest time.Duration
	maxSlack   time.Duration
	clock      Clock
}

// NewMutexBased 返回一个基于互斥锁的速率限制器。
func NewMutexBased(rate int, opts ...Option) Limiter {
	config := buildConfig(opts)
	perRequest := config.per / time.Duration(rate)
	l := &mutexLimiter{
		perRequest: perRequest,
		maxSlack:   -1 * time.Duration(config.slack) * perRequest,
		clock:      config.clock,
	}
	return l
}

// Take 会阻塞，以确保多次调用 Take 之间的时间平均为 per/rate。
func (l *mutexLimiter) Take() time.Time {
	l.Lock()
	defer l.Unlock()

	now := l.clock.Now()

	// 如果是第一次请求，直接通过
	if l.last.IsZero() {
		l.last = now
		return l.last
	}

	// sleepFor 计算我们应该休眠的时间，基于每次请求的预算和上一请求所花费的时间。
	// 由于请求的耗时可能超过预算，这个值可能会变成负数，并在多个请求中累加。
	l.sleepFor += l.perRequest - now.Sub(l.last)

	// 我们不应该让 sleepFor 的值变得过于负数，因为这意味着
	// 一个服务在短时间内显著变慢后，接下来的每秒请求数 (RPS) 会大幅增加。
	if l.sleepFor < l.maxSlack {
		l.sleepFor = l.maxSlack
	}

	// 如果 sleepFor 为正数，那么我们应该立即休眠。
	if l.sleepFor > 0 {
		l.clock.Sleep(l.sleepFor)
		l.last = now.Add(l.sleepFor)
		l.sleepFor = 0
	} else {
		l.last = now
	}

	return l.last
}
