package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	rdb *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	return &RedisCache{
		rdb: redis.NewClient(&redis.Options{Addr: addr}),
	}
}

func (c *RedisCache) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

func (c *RedisCache) GetJSON(ctx context.Context, key string, dest any) (bool, error) {
	result, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if err := json.Unmarshal([]byte(result), dest); err != nil {
		return false, err
	}
	return true, nil
}

func (c *RedisCache) SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, data, ttl).Err()
}
