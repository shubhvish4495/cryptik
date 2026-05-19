// Package cache provides a thread-safe in-memory caching implementation with expiration support.
package cache

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Package level variables for singleton pattern implementation
var (
	instance Cache     // Singleton instance of the cache
	once     sync.Once // Ensures thread-safe initialization
)

// Cache interface defines the methods that must be implemented by any cache implementation.
type Cache interface {
	// Get retrieves a value from the cache by its key.
	// Returns the value and a boolean indicating if the key exists and is not expired.
	Get(key string) (any, bool)

	// Set stores a value in the cache with the specified key and expiration time.
	// The expiration time should be a Unix timestamp.
	Set(key string, value any, expiration int64) error

	// Delete removes a value from the cache by its key.
	Delete(key string)

	// Exists checks if a key exists in the cache and is not expired.
	Exists(key string) bool

	// Clear removes all entries from the cache.
	Clear()

	// RemoveExpiredEntries removes all expired entries from the cache.
	RemoveExpiredEntries()
}

// CacheEntry represents a single entry in the cache with its data and expiration time.
type CacheEntry struct {
	Data   any   // The actual data stored in the cache
	Expiry int64 // Unix timestamp when this entry expires
}

// IsExpired checks if the cache entry has expired based on current time.
func (c CacheEntry) IsExpired() bool {
	return c.Expiry < time.Now().Unix()
}

// CacheInstance implements the Cache interface using an in-memory map.
type CacheInstance struct {
	data map[string]CacheEntry // Internal storage for cache entries
	mu   sync.RWMutex          // RWMutex for thread-safe operations
}

// RemoveExpiredEntries removes all expired entries from the cache.
// This method is called periodically by the cleanup goroutine.
func (c *CacheInstance) RemoveExpiredEntries() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, entry := range c.data {
		if entry.IsExpired() {
			delete(c.data, key)
		}
	}
}

// Set stores a value in the cache with the specified expiration time.
// Returns an error if the key is empty.
func (c *CacheInstance) Set(key string, value any, expiration int64) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = CacheEntry{
		Data:   value,
		Expiry: expiration,
	}
	return nil
}

// Get retrieves a value from the cache.
// Returns the value and true if the key exists and is not expired,
// otherwise returns nil and false.
func (c *CacheInstance) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, exists := c.data[key]
	if !exists || entry.IsExpired() {
		return nil, false
	}
	return entry.Data, true
}

// Delete removes an entry from the cache by its key.
func (c *CacheInstance) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// Exists checks if a key exists in the cache and is not expired.
func (c *CacheInstance) Exists(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, exists := (c.data)[key]
	if !exists {
		return false
	}
	return !entry.IsExpired()
}

// Clear removes all entries from the cache.
func (c *CacheInstance) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]CacheEntry)
}

// GetCache returns the singleton instance of the cache.
// If the instance doesn't exist, it creates one and starts the cleanup routine.
func GetCache(ctx context.Context) Cache {
	once.Do(func() {
		instance = &CacheInstance{
			data: make(map[string]CacheEntry),
			mu:   sync.RWMutex{},
		}

		ticker := time.NewTicker(1 * time.Minute)

		// Start a goroutine to periodically remove expired entries
		go func() {
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					instance.RemoveExpiredEntries()
				}
			}
		}()
	})
	return instance
}
