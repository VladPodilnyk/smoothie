package ratelimiter

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisScript = redis.NewScript(`
if tonumber(redis.call('incr', KEYS[1])) == 1 then
	return redis.call('expireAt', KEYS[1], ARGV[1])
else
	return 0
end
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
		maybeError := effect()
		return maybeError
	}

	return errors.New(("Limit exceeded, please try again later."))
}

func (limiter *RateLimiter) Allow(ctx context.Context, key string) bool {
	counter, err := limiter.get(ctx, key)
	if err != nil {
		return false
	}

	limiter.inc(ctx, key)
	if counter > limiter.rate.NumberOfRequests+1 {
		return false
	}
	return true
}

func (limiter *RateLimiter) inc(ctx context.Context, key string) {
	ttl := time.Now().Add(limiter.rate.Duration).Unix()
	_, err := redisScript.Run(ctx, limiter.client, []string{key}, ttl).Int()
	if err != nil {
		panic(err)
	}
}

func (limiter *RateLimiter) get(ctx context.Context, key string) (uint, error) {
	result := limiter.client.Get(ctx, key)
	if result.Err() != nil {
		return 0, result.Err()
	}
	maybeInt, err := strconv.Atoi(result.Val())
	if err != nil {
		return 0, err
	}
	return uint(maybeInt), nil
}
