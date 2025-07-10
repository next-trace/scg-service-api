// Package cache provides caching functionality.
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	appcache "github.com/next-trace/scg-service-api/application/cache"
	applogger "github.com/next-trace/scg-service-api/application/logger"
)

// cacheEntry represents an entry in the memory cache.
type cacheEntry struct {
	value      interface{}
	expiration time.Time
}

// isExpired returns true if the entry has expired.
func (e *cacheEntry) isExpired() bool {
	if e.expiration.IsZero() {
		return false
	}
	return time.Now().After(e.expiration)
}

// memoryAdapter implements the cache.Cache interface using an in-memory map.
type memoryAdapter struct {
	config    appcache.Config
	items     map[string]cacheEntry
	mu        sync.RWMutex
	log       applogger.Logger
	stopClean chan bool
}

// NewMemoryAdapter creates a new in-memory cache adapter.
func NewMemoryAdapter(config appcache.Config, log applogger.Logger) appcache.Cache {
	adapter := &memoryAdapter{
		config:    config,
		items:     make(map[string]cacheEntry),
		log:       log,
		stopClean: make(chan bool),
	}

	// Start the cleanup goroutine if cleanup interval is set
	if config.CleanupInterval > 0 {
		go adapter.startCleanup()
	}

	return adapter
}

// startCleanup starts a goroutine to periodically clean up expired entries.
func (m *memoryAdapter) startCleanup() {
	ticker := time.NewTicker(m.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.cleanup()
		case <-m.stopClean:
			return
		}
	}
}

// cleanup removes expired entries from the cache.
func (m *memoryAdapter) cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for key, entry := range m.items {
		if entry.isExpired() {
			delete(m.items, key)
		}
	}
}

// Get retrieves a value from the cache.
func (m *memoryAdapter) Get(ctx context.Context, key string) (interface{}, bool) {
	if !m.config.Enabled {
		return nil, false
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, found := m.items[key]
	if !found {
		return nil, false
	}

	if entry.isExpired() {
		// Remove expired entry
		go func() {
			m.mu.Lock()
			defer m.mu.Unlock()
			delete(m.items, key)
		}()
		return nil, false
	}

	return entry.value, true
}

// GetWithType retrieves a value from the cache and unmarshals it into the provided type.
func (m *memoryAdapter) GetWithType(ctx context.Context, key string, value interface{}) bool {
	if !m.config.Enabled {
		return false
	}

	data, found := m.Get(ctx, key)
	if !found {
		return false
	}

	// Try to assign directly
	switch v := data.(type) {
	case []byte:
		// If it's a byte slice, try to unmarshal it
		if err := json.Unmarshal(v, value); err != nil {
			m.log.Error(ctx, err, "failed to unmarshal cache value")
			return false
		}
		return true
	default:
		// Try to marshal and unmarshal to convert between types
		bytes, err := json.Marshal(v)
		if err != nil {
			m.log.Error(ctx, err, "failed to marshal cache value")
			return false
		}
		if err := json.Unmarshal(bytes, value); err != nil {
			m.log.Error(ctx, err, "failed to unmarshal cache value")
			return false
		}
		return true
	}
}

// Set stores a value in the cache with the given key and TTL.
func (m *memoryAdapter) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if !m.config.Enabled {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if we've reached the maximum number of entries
	if m.config.MaxEntries > 0 && len(m.items) >= m.config.MaxEntries {
		// Remove a random entry
		for k := range m.items {
			delete(m.items, k)
			break
		}
	}

	var expiration time.Time
	if ttl > 0 {
		expiration = time.Now().Add(ttl)
	}

	m.items[key] = cacheEntry{
		value:      value,
		expiration: expiration,
	}

	return nil
}

// Delete removes a value from the cache.
func (m *memoryAdapter) Delete(ctx context.Context, key string) error {
	if !m.config.Enabled {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.items, key)
	return nil
}

// Clear removes all values from the cache.
func (m *memoryAdapter) Clear(ctx context.Context) error {
	if !m.config.Enabled {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.items = make(map[string]cacheEntry)
	return nil
}

// Has checks if a key exists in the cache.
func (m *memoryAdapter) Has(ctx context.Context, key string) bool {
	if !m.config.Enabled {
		return false
	}

	_, found := m.Get(ctx, key)
	return found
}

// GetMulti retrieves multiple values from the cache.
func (m *memoryAdapter) GetMulti(ctx context.Context, keys []string) (map[string]interface{}, []string) {
	if !m.config.Enabled {
		return nil, keys
	}

	result := make(map[string]interface{})
	var missing []string

	for _, key := range keys {
		value, found := m.Get(ctx, key)
		if found {
			result[key] = value
		} else {
			missing = append(missing, key)
		}
	}

	return result, missing
}

// SetMulti stores multiple values in the cache with the given TTL.
func (m *memoryAdapter) SetMulti(ctx context.Context, items map[string]interface{}, ttl time.Duration) error {
	if !m.config.Enabled {
		return nil
	}

	for key, value := range items {
		if err := m.Set(ctx, key, value, ttl); err != nil {
			return err
		}
	}

	return nil
}

// DeleteMulti removes multiple values from the cache.
func (m *memoryAdapter) DeleteMulti(ctx context.Context, keys []string) error {
	if !m.config.Enabled {
		return nil
	}

	for _, key := range keys {
		if err := m.Delete(ctx, key); err != nil {
			return err
		}
	}

	return nil
}

// Increment increments a counter by the given amount.
func (m *memoryAdapter) Increment(ctx context.Context, key string, amount int64) (int64, error) {
	if !m.config.Enabled {
		return 0, errors.New("cache is disabled")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	entry, found := m.items[key]
	var value int64

	if found && !entry.isExpired() {
		// Try to convert the existing value to int64
		switch v := entry.value.(type) {
		case int:
			value = int64(v)
		case int32:
			value = int64(v)
		case int64:
			value = v
		case float32:
			value = int64(v)
		case float64:
			value = int64(v)
		default:
			return 0, errors.New("value is not a number")
		}
	}

	value += amount

	// Update the cache
	m.items[key] = cacheEntry{
		value:      value,
		expiration: entry.expiration,
	}

	return value, nil
}

// Decrement decrements a counter by the given amount.
func (m *memoryAdapter) Decrement(ctx context.Context, key string, amount int64) (int64, error) {
	return m.Increment(ctx, key, -amount)
}

// Close closes the cache connection.
func (m *memoryAdapter) Close() error {
	if m.config.CleanupInterval > 0 {
		m.stopClean <- true
	}
	return nil
}
