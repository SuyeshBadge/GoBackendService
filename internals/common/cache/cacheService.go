package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheService struct {
	client *redis.Client
}

// NewCacheService creates a new CacheService instance that uses the provided Redis connection details.
// The CacheService provides a set of methods for interacting with a Redis cache.
func NewCacheService(addr string, password string, db int) *CacheService {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &CacheService{client: client}
}

// Set stores the provided value in the cache with the given key and expiration duration.
// If the value cannot be marshaled to JSON, an error is returned.
// If the value cannot be stored in the cache, the error from the cache client is returned.
func (c *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, jsonValue, expiration).Err()
}

func (c *CacheService) Get(ctx context.Context, key string, target interface{}) error {
	jsonValue, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(jsonValue), target)
}

// Delete removes the given key from the cache.
// If the key does not exist, the operation is a no-op.
func (c *CacheService) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Exists checks if the given key exists in the cache.
// It returns true if the key exists, and false otherwise.
// If an error occurs while checking the existence of the key, the error is also returned.
func (c *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, key).Result()
	return result == 1, err
}

// Increment increments the value stored in the cache for the given key.
// If the key does not exist, it will be created with an initial value of 1.
// Returns an error if the operation fails.
func (c *CacheService) Increment(ctx context.Context, key string) error {
	return c.client.Incr(ctx, key).Err()
}

// Decrement decrements the value stored in the cache for the given key.
// If the key does not exist, it will be created with a value of 0 and then decremented.
// The decremented value is returned as an error.
func (c *CacheService) Decrement(ctx context.Context, key string) error {
	return c.client.Decr(ctx, key).Err()
}

// Close closes the underlying cache client connection.
func (c *CacheService) Close() error {
	return c.client.Close()
}
