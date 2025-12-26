package storage

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	client *redis.Client
}

func NewCache(addr, password string, db int) *Cache {
	return &Cache{client: redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: db})}
}

func (c *Cache) SetMetric(ctx context.Context, key string, payload string, ttl time.Duration) error {
	return c.client.Set(ctx, key, payload, ttl).Err()
}

func (c *Cache) Close() error {
	return c.client.Close()
}
