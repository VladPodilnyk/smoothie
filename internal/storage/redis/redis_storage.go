package storage

import (
	"context"
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

type redisStorage struct {
	client *redis.Client
	ttl    time.Duration
}

type RedisConfig struct {
	addr     string
	password string
	db       int
}

func NewRedisStorage(config *redis.Options) redisStorage {
	return redisStorage{client: redis.NewClient(config)}
}

func (s redisStorage) Inc(ctx context.Context, key string) {
	redisScript.Run(ctx, s.client, []string{key}, time.Now().Add(s.ttl).Unix())
}

func (s redisStorage) Get(ctx context.Context, key string) (uint, error) {
	result := s.client.Get(ctx, key)
	if result.Err() != nil {
		return 0, result.Err()
	}
	maybeInt, err := strconv.Atoi(result.Val())
	if err != nil {
		return 0, err
	}
	return uint(maybeInt), nil
}
