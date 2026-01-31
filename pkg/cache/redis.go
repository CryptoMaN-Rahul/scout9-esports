package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache implements caching using Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache client
func NewRedisCache(addr string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     "", // no password by default
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &RedisCache{client: client}, nil
}

// Get retrieves a value from cache
func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, fmt.Errorf("redis get: %w", err)
	}
	return val, nil
}

// Set stores a value in cache with TTL
func (c *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if err := c.client.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}
	return nil
}

// Delete removes a value from cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis delete: %w", err)
	}
	return nil
}

// DeletePattern removes all keys matching a pattern
func (c *RedisCache) DeletePattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
			return fmt.Errorf("redis delete pattern: %w", err)
		}
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("redis scan: %w", err)
	}
	return nil
}

// Close closes the Redis connection
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// RateLimitTracker tracks API rate limits using Redis
type RateLimitTracker struct {
	cache  *RedisCache
	prefix string
}

// NewRateLimitTracker creates a new rate limit tracker
func NewRateLimitTracker(cache *RedisCache, prefix string) *RateLimitTracker {
	return &RateLimitTracker{
		cache:  cache,
		prefix: prefix,
	}
}

// Increment increments the request count for a key
func (t *RateLimitTracker) Increment(ctx context.Context, key string, window time.Duration) (int64, error) {
	fullKey := fmt.Sprintf("%s:%s", t.prefix, key)

	pipe := t.cache.client.Pipeline()
	incr := pipe.Incr(ctx, fullKey)
	pipe.Expire(ctx, fullKey, window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("rate limit increment: %w", err)
	}

	return incr.Val(), nil
}

// GetCount gets the current request count for a key
func (t *RateLimitTracker) GetCount(ctx context.Context, key string) (int64, error) {
	fullKey := fmt.Sprintf("%s:%s", t.prefix, key)

	val, err := t.cache.client.Get(ctx, fullKey).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("rate limit get: %w", err)
	}

	return val, nil
}

// IsAllowed checks if a request is allowed under the rate limit
func (t *RateLimitTracker) IsAllowed(ctx context.Context, key string, limit int64, window time.Duration) (bool, error) {
	count, err := t.Increment(ctx, key, window)
	if err != nil {
		return false, err
	}
	return count <= limit, nil
}
