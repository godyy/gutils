package leaky_bucket

import (
	"sync/atomic"
	"time"
)

// atomicInt64Limiter 是基于int64原子操作的速率限制器。
type atomicInt64Limiter struct {
	// lint:ignore U1000 Padding 未使用，但为了保持这个速率限制器的性能，
	// 在与其他频繁访问的内存共存时它是至关重要的。
	prepadding [64]byte // 缓存行大小 = 64；创建的目的是避免伪共享。
	state      int64    // 下一个权限问题的 Unix 纳秒时间。
	// lint:ignore U1000 类似于预填充。
	postpadding [56]byte // 缓存行大小 - 状态大小 = 64 - 8；创建的目的是避免伪共享。

	perRequest time.Duration
	maxSlack   time.Duration
	clock      Clock
}

// NewAtomicInt64Based 返回一个基于int64原子操作的速率限制器。
func NewAtomicInt64Based(rate int, opts ...Option) Limiter {
	config := buildConfig(opts)
	perRequest := config.per / time.Duration(rate)
	l := &atomicInt64Limiter{
		perRequest: perRequest,
		maxSlack:   time.Duration(config.slack) * perRequest,
		clock:      config.clock,
	}
	atomic.StoreInt64(&l.state, 0)
	return l
}

func (l *atomicInt64Limiter) Take() time.Time {
	var (
		newTimeOfNextPermissionIssue int64
		now                          int64
	)
	for {
		now = l.clock.Now().UnixNano()
		timeOfNextPermissionIssue := atomic.LoadInt64(&l.state)

		switch {
		case timeOfNextPermissionIssue == 0 || (l.maxSlack == 0 && now-timeOfNextPermissionIssue > int64(l.perRequest)):
			// 如果这是我们的第一次调用，或者 t.maxSlack == 0，我们需要将问题时间缩小到当前时间。
			newTimeOfNextPermissionIssue = now
		case l.maxSlack > 0 && now-timeOfNextPermissionIssue > int64(l.maxSlack)+int64(l.perRequest):
			// 自上次 Take 调用以来已经过去了很多纳秒，
			// 我们将限制最大累积时间为 maxSlack。
			newTimeOfNextPermissionIssue = now - int64(l.maxSlack)
		default:
			// 计算我们的权限被授予的时间。
			newTimeOfNextPermissionIssue = timeOfNextPermissionIssue + int64(l.perRequest)
		}

		if atomic.CompareAndSwapInt64(&l.state, timeOfNextPermissionIssue, newTimeOfNextPermissionIssue) {
			break
		}
	}

	sleepDuration := time.Duration(newTimeOfNextPermissionIssue - now)
	if sleepDuration > 0 {
		l.clock.Sleep(sleepDuration)
		return time.Unix(0, newTimeOfNextPermissionIssue)
	}
	// 无需休眠,则返回当前时间。
	return time.Unix(0, now)
}
