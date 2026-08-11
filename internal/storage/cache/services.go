package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func (c *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	ctx, span := c.tracer.Start(ctx, "cache.Set")
	defer span.End()

	data, err := json.Marshal(value)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("cache.Set: marshal: %w", err)
	}

	if err := c.client.Set(ctx, key, data, expiration).Err(); err != nil {
		span.RecordError(err)
		return fmt.Errorf("cache.Set: redis set: %w", err)
	}

	return nil
}

func (c *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	ctx, span := c.tracer.Start(ctx, "cache.Get")
	defer span.End()

	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		span.RecordError(err)
		if err == redis.Nil {
			return fmt.Errorf("cache.Get: miss: %w", err)
		}
		return fmt.Errorf("cache.Get: redis get: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		span.RecordError(err)
		return fmt.Errorf("cache.Get: unmarshal: %w", err)
	}

	return nil
}

func (c *CacheService) Delete(ctx context.Context, key string) error {
	ctx, span := c.tracer.Start(ctx, "cache.Delete")
	defer span.End()

	if err := c.client.Del(ctx, key).Err(); err != nil {
		span.RecordError(err)
		return fmt.Errorf("cache.Delete: redis del: %w", err)
	}

	return nil
}
