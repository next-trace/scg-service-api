// Package cache defines the abstract interface (PORT) for caching.
package cache

import (
	"context"
	"time"
)

// Cache defines the interface for caching.
type Cache interface {
	// Get retrieves a value from the cache.
	// It returns the value and a boolean indicating whether the value was found.
	Get(ctx context.Context, key string) (interface{}, bool)

	// GetWithType retrieves a value from the cache and unmarshals it into the provided type.
	// It returns a boolean indicating whether the value was found.
	GetWithType(ctx context.Context, key string, value interface{}) bool

	// Set stores a value in the cache with the given key and TTL.
	// If ttl is 0, the value will not expire.
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Delete removes a value from the cache.
	Delete(ctx context.Context, key string) error

	// Clear removes all values from the cache.
	Clear(ctx context.Context) error

	// Has checks if a key exists in the cache.
	Has(ctx context.Context, key string) bool

	// GetMulti retrieves multiple values from the cache.
	// It returns a map of values and a slice of keys that were not found.
	GetMulti(ctx context.Context, keys []string) (map[string]interface{}, []string)

	// SetMulti stores multiple values in the cache with the given TTL.
	// If ttl is 0, the values will not expire.
	SetMulti(ctx context.Context, items map[string]interface{}, ttl time.Duration) error

	// DeleteMulti removes multiple values from the cache.
	DeleteMulti(ctx context.Context, keys []string) error

	// Increment increments a counter by the given amount.
	// It returns the new value.
	Increment(ctx context.Context, key string, amount int64) (int64, error)

	// Decrement decrements a counter by the given amount.
	// It returns the new value.
	Decrement(ctx context.Context, key string, amount int64) (int64, error)

	// Close closes the cache connection.
	Close() error
}

// StoreType defines the type of cache store.
type StoreType string

const (
	// StoreTypeMemory is an in-memory cache store.
	StoreTypeMemory StoreType = "memory"

	// StoreTypeRedis is a Redis cache store.
	StoreTypeRedis StoreType = "redis"
)

// Config holds configuration for caching.
type Config struct {
	// Enabled determines if caching is enabled.
	Enabled bool

	// StoreType is the type of cache store to use.
	StoreType StoreType

	// DefaultTTL is the default time-to-live for cache entries.
	DefaultTTL time.Duration

	// CleanupInterval is the interval at which to clean up expired entries.
	// Only applicable for memory store.
	CleanupInterval time.Duration

	// MaxEntries is the maximum number of entries in the cache.
	// Only applicable for memory store.
	MaxEntries int

	// Redis configuration
	Redis struct {
		// Address is the Redis server address.
		Address string

		// Password is the Redis server password.
		Password string

		// DB is the Redis database number.
		DB int

		// MaxRetries is the maximum number of retries for Redis operations.
		MaxRetries int

		// PoolSize is the maximum number of connections in the Redis connection pool.
		PoolSize int
	}
}

// DefaultConfig returns the default configuration for caching.
func DefaultConfig() Config {
	return Config{
		Enabled:         true,
		StoreType:       StoreTypeMemory,
		DefaultTTL:      time.Minute * 5,
		CleanupInterval: time.Minute,
		MaxEntries:      10000,
		Redis: struct {
			Address    string
			Password   string
			DB         int
			MaxRetries int
			PoolSize   int
		}{
			Address:    "localhost:6379",
			Password:   "",
			DB:         0,
			MaxRetries: 3,
			PoolSize:   10,
		},
	}
}
