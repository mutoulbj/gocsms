package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mutoulbj/gocsms/internal/config"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewRedisCache(cfg *config.RedisConfig) (*Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Cache{client: rdb}, nil
}

func (c *Cache) Client() *redis.Client {
	return c.client
}

func ProviderRedisClient(cache *Cache) *redis.Client {
	return cache.Client()
}

func (c *Cache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.client.Set(ctx, key, data, expiration).Err()
}

func (c *Cache) Get(ctx context.Context, key string, dest any) error {
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("cache not found")
		}
		return fmt.Errorf("failed to get value: %w", err)
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}
	return count > 0, nil
}

func (c *Cache) Close() error {
	if err := c.client.Close(); err != nil {
		return fmt.Errorf("failed to close Redis client: %w", err)
	}
	return nil
}
