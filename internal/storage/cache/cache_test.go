package cache

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
)

// создаем тестовый кэш с миниредис
func setupTestCache(t *testing.T) (*CacheService, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	tracer := otel.Tracer("test")

	cache := &CacheService{
		client: client,
		tracer: tracer,
	}

	return cache, mr
}

// установка значения в кэш
func TestCacheService_Set_Struct(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()

	ctx := context.Background()
	key := "test-key"
	value := map[string]string{
		"userName":  "test1",
		"userEmail": "test1",
	}

	err := cache.Set(ctx, key, value, time.Minute)

	assert.NoError(t, err)
	exists := mr.Exists(key)
	assert.True(t, exists)
}

// получение данных из кэша
func TestCacheService_Get_Success(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()

	ctx := context.Background()
	key := "test-key-2"
	expectedData := map[string]string{
		"userName":  "test2",
		"userEmail": "test2",
	}

	jsonData, err := json.Marshal(expectedData)
	require.NoError(t, err)
	mr.Set(key, string(jsonData))

	var result map[string]interface{}
	err = cache.Get(ctx, key, &result)

	assert.NoError(t, err)
	assert.Equal(t, expectedData["userName"], result["userName"])
	assert.Equal(t, expectedData["userEmail"], result["userEmail"])
}

// получение просроченного кэша
func TestCacheService_Get_Expired(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()

	ctx := context.Background()
	key := "test-key-3"
	expectedData := map[string]string{
		"userName":  "test3",
		"userEmail": "test3",
	}

	jsonData, err := json.Marshal(expectedData)
	require.NoError(t, err)

	mr.Set(key, string(jsonData))
	mr.SetTTL(key, 1*time.Second)

	var result map[string]string
	err = cache.Get(ctx, key, &result)
	assert.NoError(t, err)
	assert.Equal(t, expectedData["userName"], result["userName"])

	mr.FastForward(2 * time.Second)

	err = cache.Get(ctx, key, &result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cache miss")
}

// получение данных по несуществующему ключу
func TestCacheService_Get_CacheMiss(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()

	ctx := context.Background()
	var result string

	err := cache.Get(ctx, "KEYNOTFOUND", &result)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cache miss")
}

// удаление записи
func TestCacheService_Delete(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()

	ctx := context.Background()
	key := "test-key-4"
	value := map[string]string{
		"userName":  "test4",
		"userEmail": "test4",
	}

	jsonData, err := json.Marshal(value)
	require.NoError(t, err)
	mr.Set(key, string(jsonData))

	err = cache.Delete(ctx, key)

	assert.NoError(t, err)

	exists := mr.Exists(key)
	assert.False(t, exists)
}

// удаление записи с неправильным ключом
func TestCacheService_Delete_NonExistentKey(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()

	ctx := context.Background()

	err := cache.Delete(ctx, "KEYNOTFOUND")

	assert.NoError(t, err)
}

// полный жизненный цикл кэша
func TestCacheService_FullLifecycle(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()

	ctx := context.Background()
	key := "test-key-5"
	value := map[string]string{
		"userName":  "test5",
		"userEmail": "test5",
	}

	err := cache.Set(ctx, key, value, 10*time.Minute)
	require.NoError(t, err)

	var result map[string]interface{}
	err = cache.Get(ctx, key, &result)
	require.NoError(t, err)
	assert.Equal(t, value["userName"], result["userName"])
	assert.Equal(t, value["userEmail"], result["userEmail"])

	err = cache.Delete(ctx, key)
	require.NoError(t, err)

	err = cache.Get(ctx, key, &result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cache miss")
}
