package ratelimiter

import (
	"context"
	"fmt"
	"math/rand"
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

// Super simple random key generator (fancy stuff is redundant here)
func randKey(size int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	key := make([]byte, size)
	for i := range key {
		key[i] = charset[random.Intn(len(charset))]
	}
	return string(key)
}

func TestAllowRequestWithinSpecifiedRate(t *testing.T) {
	verifyRedisConnection(t)

	testKey := randKey(10)
	limiter := New(&testRedisOptions, Rate{NumberOfRequests: 1, Duration: 5 * time.Second})
	firstTry := limiter.Allow(context.Background(), testKey)

	time.Sleep(5 * time.Second)

	secondTry := limiter.Allow(context.Background(), testKey)
	if !firstTry || !secondTry {
		t.Errorf("Expected both requests to be allowed, bug got rate limited.")
	}
}

func TestDoNotAllowRequestThatExceedsLimit(t *testing.T) {
	verifyRedisConnection(t)

	testKey := randKey(10)
	limiter := New(&testRedisOptions, Rate{NumberOfRequests: 1, Duration: 10 * time.Second})
	limiter.Allow(context.Background(), testKey)
	isAllowed := limiter.Allow(context.Background(), testKey)

	if isAllowed {
		t.Errorf("Expected request to be rate limited, but it was allowed.")
	}
}

func TestExecFunctionWithingSpecifiedRate(t *testing.T) {
	verifyRedisConnection(t)
	limiter := New(&testRedisOptions, Rate{NumberOfRequests: 1, Duration: 5 * time.Second})

	counter := 0
	incrementCounter := func() error {
		fmt.Println("Incrementing counter")
		counter += 1
		return nil
	}

	testKey := randKey(10)
	limiter.Exec(context.Background(), testKey, incrementCounter)
	time.Sleep(5 * time.Second)
	err := limiter.Exec(context.Background(), testKey, incrementCounter)
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

	testKey := randKey(10)
	limiter.Exec(context.Background(), testKey, incrementCounter)
	err := limiter.Exec(context.Background(), testKey, incrementCounter)

	if err == nil {
		t.Errorf("Expected request to be rate limited, but got allowed")
	}
}
