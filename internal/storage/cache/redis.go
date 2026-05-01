package cache

import (
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
)

type CacheService struct {
	client *redis.Client
	tracer trace.Tracer
}

func NewCacheService(client *redis.Client, tracer trace.Tracer) *CacheService {
	return &CacheService{
		client: client,
		tracer: tracer,
	}
}
