package ratelimiter

import (
	"context"
	"errors"

	"github.com/vladpodilnyk/smoothie/internal/storage"
)

// holds strategy and configuration values
type RateLimiter struct {
	storage storage.Storage
	rate    Rate
}

func New() *RateLimiter {
	return &RateLimiter{}
}

func (limiter *RateLimiter) Exec(ctx context.Context, key string, effect func()) error {
	if limiter.Allow(ctx, key) {
		effect()
		return nil
	}
	return errors.New(("Limit exceeded, please try again later."))
}

func (limiter *RateLimiter) Allow(ctx context.Context, key string) bool {
	counter, err := limiter.storage.Get(ctx, key)
	if err != nil {
		return false
	}

	limiter.storage.Inc(ctx, key)
	if counter > limiter.rate.NumberOfRequests+1 {
		return false
	}
	return true
}
