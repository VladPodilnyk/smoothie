package ratelimiter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisScript = redis.NewScript(`
local current = redis.call('incr', KEYS[1])

if current == 1 then
	return redis.call('expireAt', KEYS[1], ARGV[1])
end

return current
`)

type Rate struct {
	NumberOfRequests uint
	Duration         time.Duration
}

type RateLimiter struct {
	client *redis.Client
	rate   Rate
}

func New(config *redis.Options, rate Rate) *RateLimiter {
	// Typically it's better to pass client as a parameter, but
	// since Smoothie stricly depends on redis, it doesn't make
	// any sense to allow users to pass their own client.
	client := redis.NewClient(config)
	return &RateLimiter{client, rate}
}

func (limiter *RateLimiter) Exec(ctx context.Context, key string, effect func() error) error {
	if limiter.Allow(ctx, key) {
		fmt.Println("Allowing request")
		maybeError := effect()
		return maybeError
	}
	return errors.New(("Limit exceeded, please try again later."))
}

func (limiter *RateLimiter) Allow(ctx context.Context, key string) bool {
	currentValue := limiter.incrementAndGet(ctx, key)
	if currentValue > limiter.rate.NumberOfRequests {
		return false
	}
	return true
}

func (limiter *RateLimiter) incrementAndGet(ctx context.Context, key string) uint {
	ttl := time.Now().Add(limiter.rate.Duration).Unix()
	value, err := redisScript.Run(ctx, limiter.client, []string{key}, ttl).Int()
	if err != nil {
		panic(err)
	}
	return uint(value)
}
