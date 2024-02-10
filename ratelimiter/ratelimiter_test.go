package ratelimiter

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

var testRedisOptions = redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
}

func verifyRedisConnection(t *testing.T) {
	client := redis.NewClient(&testRedisOptions)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		t.Fatalf("Could not connect to Redis: %s", err)
	}
}

func fatalError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Fatal error: %s", err.Error())
	}
}

func TestAllowRequestWithinSpecifiedRate(t *testing.T) {
	verifyRedisConnection(t)

	limiter := New(&testRedisOptions, Rate{NumberOfRequests: 1, Duration: 5 * time.Second})
	firstTry, err := limiter.Allow(context.Background(), "test-key")
	fatalError(t, err)

	time.Sleep(5 * time.Second)

	secondTry, err := limiter.Allow(context.Background(), "test-key")
	fatalError(t, err)
	if !firstTry || !secondTry {
		t.Errorf("Expected both requests to be allowed, bug got rate limited.")
	}
}

func TestDoNotAllowRequestThatExceedsLimit(t *testing.T) {
	verifyRedisConnection(t)

	limiter := New(&testRedisOptions, Rate{NumberOfRequests: 1, Duration: 10 * time.Second})
	limiter.Allow(context.Background(), "test-key")
	isAllowed, err := limiter.Allow(context.Background(), "test-key")
	fatalError(t, err)

	if isAllowed {
		t.Errorf("Expected request to be rate limited, but it was allowed.")
	}
}

func TestExecFunctionWithingSpecifiedRate(t *testing.T) {
	verifyRedisConnection(t)
	limiter := New(&testRedisOptions, Rate{NumberOfRequests: 1, Duration: 5 * time.Second})

	counter := 0
	incrementCounter := func() error {
		counter += 1
		return nil
	}

	limiter.Exec(context.Background(), "test-key", incrementCounter)
	time.Sleep(5 * time.Second)
	err := limiter.Exec(context.Background(), "test-key", incrementCounter)

	if err != nil {
		t.Errorf("Expected request to be allowed, but got rate limited.")
	}

	if counter != 2 {
		t.Errorf("Expected effect to be executed twice, but got %d", counter)
	}
}

func TestRejectFunctionExecutionIfRateIsExceeded(t *testing.T) {
	verifyRedisConnection(t)
	limiter := New(&testRedisOptions, Rate{NumberOfRequests: 1, Duration: 5 * time.Second})

	counter := 0
	incrementCounter := func() error {
		counter += 1
		return nil
	}

	limiter.Exec(context.Background(), "test-key", incrementCounter)
	err := limiter.Exec(context.Background(), "test-key", incrementCounter)

	if err == nil {
		t.Errorf("Expected request to be allowed, but got rate limited.")
	}
}
