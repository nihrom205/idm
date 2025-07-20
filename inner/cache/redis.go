package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	return &RedisCache{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "",
			DB:       0,
		}),
	}
}

// Set сохраняет значение с TTL (в секундах)
func (rc *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rc.client.Set(ctx, key, jsonData, expiration).Err()
}

// Get возвращает значение по ключу
func (rc *RedisCache) Get(ctx context.Context, key string, dest interface{}) bool {
	val, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	err = json.Unmarshal([]byte(val), dest)
	return err == nil
}

// Delete удаляет ключ
func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	return rc.client.Del(ctx, key).Err()
}

func (rc *RedisCache) InvalidateCache(ctx context.Context, key string) error {
	return rc.Delete(ctx, key)
}

func (rc *RedisCache) GetCacheKey(prefix string, id int64) string {
	return fmt.Sprintf("%s:%d", prefix, id)
}
