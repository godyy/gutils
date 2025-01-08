package leaky_bucket

import "time"

// Limiter 用于对某些进程进行速率限制，可能跨协程 (goroutines)。
// 该进程需要在每次迭代前调用 Take()，这可能会阻塞以限制协程的执行速度。
type Limiter interface {
	// Take 应该阻塞以确保满足每秒请求数 (RPS) 的要求。
	Take() time.Time
}

// Clock 是用于通过时钟或模拟时钟实例化速率限制器的最小必要接口，
type Clock interface {
	Now() time.Time
	Sleep(duration time.Duration)
}

// clock 是 Clock 的默认实现。
type clock struct{}

func (c *clock) Now() time.Time {
	return time.Now()
}

func (c *clock) Sleep(duration time.Duration) {
	time.Sleep(duration)
}

// defaultClock 声明默认 Clock.
var defaultClock = &clock{}

// config 用于配置一个 limiter。
type config struct {
	clock Clock
	slack int
	per   time.Duration
}

// Option 用于配置一个 Limiter。
type Option interface {
	apply(*config)
}

// buildConfig 将默认配置与选项合并。
func buildConfig(opts []Option) config {
	c := config{
		clock: defaultClock,
		slack: 10,
		per:   time.Second,
	}

	for _, opt := range opts {
		opt.apply(&c)
	}
	return c
}

type clockOption struct {
	clock Clock
}

func (o clockOption) apply(c *config) {
	c.clock = o.clock
}

// WithClock 返回一个用于创建速率限制器的选项，提供一个替代的
// Clock 实现，通常是用于测试的模拟时钟。
func WithClock(clock Clock) Option {
	return clockOption{clock: clock}
}

type slackOption int

func (o slackOption) apply(c *config) {
	c.slack = int(o)
}

// WithoutSlack 配置 limiter 以严格的方式工作，不会将之前“未使用”的请求
// 累积到未来的流量激增中。
var WithoutSlack Option = slackOption(0)

// WithSlack 配置自定义的 slack。
// Slack 允许 limiter 累积“未使用”的请求，以应对未来的流量激增。
func WithSlack(slack int) Option {
	return slackOption(slack)
}

type perOption time.Duration

func (p perOption) apply(c *config) {
	c.per = time.Duration(p)
}

// WithPer 允许配置不同时间窗口的限制。
// 默认窗口为1秒，因此当 rate=100 时会产生一个每秒 100 次请求（100 Hz）的速率限制器。
// 当 rate=2, 且提供 WithPer(60*time.Second) 选项时，将 创建一个每分钟 2 次请求的速率限制器。
func WithPer(per time.Duration) Option {
	return perOption(per)
}
