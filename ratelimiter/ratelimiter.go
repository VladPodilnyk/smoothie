package ratelimiter

// holds strategy and configuration values
type RateLimiter struct{}

func New() *RateLimiter {
	return &RateLimiter{}
}

// extecutes a given function according to the rate limit stratagy
func (limiter *RateLimiter) Exec(effect func(any) any) {}

func (limiter *RateLimiter) Allow() bool {
	return false
}
