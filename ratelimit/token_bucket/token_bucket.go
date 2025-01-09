package token_bucket

import (
	"math"
	"strconv"
	"sync"
	"time"
)

// Clock 表示时间的流逝，可以在测试中通过模拟来实现。
type Clock interface {
	Now() time.Time
	Sleep(duration time.Duration)
}

// realClock 使用标准时间函数实现 Clock 接口。
type realClock struct{}

// Now 通过调用 time.Now 实现 Clock.Now.
func (realClock) Now() time.Time {
	return time.Now()
}

// Sleep 通过调用 time.Sleep 实现 Clock.Sleep.
func (realClock) Sleep(d time.Duration) {
	time.Sleep(d)
}

// Bucket 表示一个按预定速率填充的令牌桶。
// Bucket 上的方法可以并发调用。
type Bucket struct {
	clock Clock

	// startTime 保存了令牌桶首次创建并开始计时的时间点。
	startTime time.Time

	// capacity 保存了桶的总容量。
	capacity int64

	// quantum 表示每个时钟滴答添加的令牌数量。
	quantum int64

	// fillInterval 表示每个时钟滴答之间的时间间隔。
	fillInterval time.Duration

	// mu 保护其下方的字段。
	mu sync.Mutex

	// availableTokens 保存了与最新时钟滴答 (latestTick) 相关的可用令牌数量。
	// 当有消费者在等待令牌时，该值可能为负数。
	availableTokens int64

	// latestTick 保存了我们已知桶中令牌数量的最新时钟滴答。
	latestTick int64
}

// New 返回一个新的令牌桶，该桶以每个 fillInterval 填充一个令牌，
// 直到达到指定的最大容量。两个参数必须是正数。
// 该桶最初是满的。
func New(fillInterval time.Duration, capacity int64) *Bucket {
	return NewWithClock(fillInterval, capacity, nil)
}

// NewWithClock 与 New 相同，但注入了一个可测试的时钟接口。
func NewWithClock(fillInterval time.Duration, capacity int64, clock Clock) *Bucket {
	return NewWithQuantumAndClock(fillInterval, capacity, 1, clock)
}

// rateMargin 指定实际速率与指定速率之间允许的偏差。
// 1% 被认为是合理的。
const rateMargin = 0.01

// NewWithRate 返回一个令牌桶，该桶以每秒 rate 个令牌的速率填充，
// 直至达到指定的最大容量。
// 由于时钟分辨率的限制，在高速率下，实际速率可能与指定速率相差最多 1%。
func NewWithRate(rate float64, capacity int64) *Bucket {
	return NewWithRateAndClock(rate, capacity, nil)
}

// NewWithRateAndClock 与 NewWithRate 相同，
// 但注入了一个可测试的时钟接口。
func NewWithRateAndClock(rate float64, capacity int64, clock Clock) *Bucket {
	// 在循环中每次使用相同的桶，以节省分配开销。
	b := NewWithQuantumAndClock(1, capacity, 1, clock)
	for quantum := int64(1); quantum < 1<<50; quantum = nextQuantum(quantum) {
		fillInterval := time.Duration(1e9 * float64(quantum) / rate)
		if fillInterval <= 0 {
			continue
		}
		b.fillInterval = fillInterval
		b.quantum = quantum
		if diff := math.Abs(b.Rate() - rate); diff/rate <= rateMargin {
			return b
		}
	}
	panic("cannot find suitable quantum for " + strconv.FormatFloat(rate, 'g', -1, 64))
}

// nextQuantum 返回 q 之后要尝试的下一个量子值。
// 我们以指数形式增长量子值，但增速较慢，
// 以便更好地适配较小的数值。
func nextQuantum(q int64) int64 {
	q1 := q * 11 / 10
	if q1 == q {
		q1++
	}
	return q1
}

// NewWithQuantum 类似于 New，但允许指定 quantum,
// 每个 fillInterval 将添加 quantum 个令牌。
func NewWithQuantum(fillInterval time.Duration, capacity, quantum int64) *Bucket {
	return NewWithQuantumAndClock(fillInterval, capacity, quantum, nil)
}

// NewWithQuantumAndClock 类似于 NewWithQuantum，
// 但还具有一个 clock 参数，允许客户端模拟时间的流逝。
// 如果 clock 为 nil，将使用系统时钟。
func NewWithQuantumAndClock(fillInterval time.Duration, capacity, quantum int64, clock Clock) *Bucket {
	if clock == nil {
		clock = realClock{}
	}
	if fillInterval <= 0 {
		panic("token bucket fill interval is not > 0")
	}
	if capacity <= 0 {
		panic("token bucket capacity is not > 0")
	}
	if quantum <= 0 {
		panic("token bucket quantum is not > 0")
	}
	return &Bucket{
		clock:           clock,
		startTime:       clock.Now(),
		latestTick:      0,
		fillInterval:    fillInterval,
		capacity:        capacity,
		quantum:         quantum,
		availableTokens: capacity,
	}
}

// Wait 从桶中取出 count 个令牌，并等待直到这些令牌可用。
func (b *Bucket) Wait(count int64) {
	if d := b.Take(count); d > 0 {
		b.clock.Sleep(d)
	}
}

// WaitMaxDuration 类似于 Wait，但只有在需要等待的时间不超过 maxWait 时，
// 才会从桶中取出令牌。它会报告是否有令牌被移出桶中。
// 如果没有移出任何令牌，它会立即返回。
func (b *Bucket) WaitMaxDuration(count int64, maxWait time.Duration) bool {
	d, ok := b.TakeMaxDuration(count, maxWait)
	if d > 0 {
		b.clock.Sleep(d)
	}
	return ok
}

const infinityDuration time.Duration = 0x7fffffffffffffff

// Take 从桶中取出 count 个令牌，不会阻塞。
// 它返回调用者应等待的时间，直到令牌实际可用。
// 请注意，如果请求是不可撤销的,一旦此方法承诺取出令牌，就无法将令牌返回到桶中。
func (b *Bucket) Take(count int64) time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()
	d, _ := b.take(b.clock.Now(), count, infinityDuration)
	return d
}

// TakeMaxDuration 类似于 Take，
// 不同之处在于它只有在令牌等待时间不超过 maxWait 的情况下，
// 才会从桶中取出令牌。
// 如果等待令牌变得可用的时间超过了 maxWait，
// 它什么也不做并返回 false，
// 否则，它会返回调用者应等待的时间，直到令牌实际可用，并返回 true。
func (b *Bucket) TakeMaxDuration(count int64, maxWait time.Duration) (time.Duration, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.take(b.clock.Now(), count, maxWait)
}

// TakeAvailable 从桶中立即取走最多 count 个可用令牌。
// 它返回移除的令牌数量，如果没有可用令牌，则返回零。
// 它不会阻塞。
func (b *Bucket) TakeAvailable(count int64) int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.takeAvailable(b.clock.Now(), count)
}

// takeAvailable 是 TakeAvailable 的内部版本——它将当前时间作为参数，以便于测试。
func (b *Bucket) takeAvailable(now time.Time, count int64) int64 {
	if count <= 0 {
		return 0
	}
	b.adjustAvailableTokens(b.currentTick(now))
	if b.availableTokens <= 0 {
		return 0
	}
	if count > b.availableTokens {
		count = b.availableTokens
	}
	b.availableTokens -= count
	return count
}

// Available 返回可用令牌的数量。当有消费者在等待令牌时，返回值可能为负数。
// 请注意，如果返回值大于零，并不意味着从缓冲区中取令牌的调用一定会成功，
// 因为在这段时间内可用令牌的数量可能发生了变化。此方法主要用于指标报告和调试。
func (b *Bucket) Available() int64 {
	return b.available(b.clock.Now())
}

// available 是 available 的内部版本——它将当前时间作为参数，以便于测试。
func (b *Bucket) available(now time.Time) int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.adjustAvailableTokens(b.currentTick(now))
	return b.availableTokens
}

// Capacity 返回桶的容量.
func (b *Bucket) Capacity() int64 {
	return b.capacity
}

// Rate 返回桶的填充速率，即每秒填充的令牌数.
func (b *Bucket) Rate() float64 {
	return 1e9 * float64(b.quantum) / float64(b.fillInterval)
}

// take 是 Take 的内部版本——它将当前时间作为参数，以便于测试。
func (b *Bucket) take(now time.Time, count int64, maxWait time.Duration) (time.Duration, bool) {
	if count <= 0 {
		return 0, true
	}

	tick := b.currentTick(now)
	b.adjustAvailableTokens(tick)
	available := b.availableTokens - count
	if available >= 0 {
		b.availableTokens = available
		return 0, true
	}

	// 将缺失的令牌向上舍入到最接近的量子倍数——这些令牌
	// 直到该时钟滴答时才会变得可用。

	// endTick 表示所有请求的令牌均可用的时钟滴答
	endTick := tick + (-available+b.quantum-1)/b.quantum
	endTime := b.startTime.Add(time.Duration(endTick) * b.fillInterval)
	waitTime := endTime.Sub(now)
	if waitTime > maxWait {
		return 0, false
	}
	b.availableTokens = available
	return waitTime, true
}

// currentTick 返回当前的时间滴答，从 l.startTime 开始计时。
func (b *Bucket) currentTick(now time.Time) int64 {
	return int64(now.Sub(b.startTime) / b.fillInterval)
}

// adjustAvailableTokens 调整桶中在给定时间点（相对于 l.latestTick 必须是未来时间）
// 可用的令牌数量。
func (b *Bucket) adjustAvailableTokens(tick int64) {
	lastTick := b.latestTick
	b.latestTick = tick
	if b.availableTokens >= b.capacity {
		return
	}
	b.availableTokens += (tick - lastTick) * b.quantum
	if b.availableTokens > b.capacity {
		b.availableTokens = b.capacity
	}
	return
}
