package ratelimiter

import (
	"context"
	"errors"
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

// Rate represents the rate limter configuration
type Rate struct {
	NumberOfRequests uint          // Number of requests allowed
	Duration         time.Duration // Time window
}

// RateLimiter holds the redis client and the rate configuration
type RateLimiter struct {
	client *redis.Client
	rate   Rate
}

// Creates a new rate limiter with the given redis config and rate
func New(config *redis.Options, rate Rate) *RateLimiter {
	// Typically it's better to pass client as a parameter, but
	// since Smoothie stricly depends on redis, it doesn't make
	// any sense to allow users to pass their own client.
	client := redis.NewClient(config)
	return &RateLimiter{client, rate}
}

// Exec executes the given function (effect) if the rate limit is not exceeded,
// otherwise it returns an error. A user should pass a key that uniquely identifies
// the request.
func (limiter *RateLimiter) Exec(ctx context.Context, key string, effect func() error) error {
	if limiter.Allow(ctx, key) {
		maybeError := effect()
		return maybeError
	}
	return errors.New(("Limit exceeded, please try again later."))
}

// Allow returns true if the rate limit is not exceeded, otherwise it returns false.
// A user should pass a key that uniquely identifies the request.
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
