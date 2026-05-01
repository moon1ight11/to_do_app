package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// метод для установки кэша
func (c *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	ctx, span := c.tracer.Start(ctx, "cache.Set")
	defer span.End()

	// маршалим данные, которые будем кэшировать
	data, err := json.Marshal(value)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// кэшируем
	err = c.client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// метод для получения кэша
func (c *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	ctx, span := c.tracer.Start(ctx, "cache.Get")
	defer span.End()

	// получаем данные и кэша
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		span.RecordError(err)
		if err == redis.Nil {
			return fmt.Errorf("cache miss: %w", err)
		}
		return fmt.Errorf("failed to get from cache: %w", err)
	}

	err = json.Unmarshal(data, dest)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return nil
}

// метод для удаления данных из кэша
func (c *CacheService) Delete(ctx context.Context, key string) error {
	ctx, span := c.tracer.Start(ctx, "cache.Delete")
	defer span.End()

	err := c.client.Del(ctx, key).Err()
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to delete from cache: %w", err)
	}

	return nil
}
