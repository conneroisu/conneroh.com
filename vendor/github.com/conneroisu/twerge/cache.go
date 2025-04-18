package twerge

import (
	"fmt"
	"sync"
)

// Cache is a simple thread-safe cache.
type Cache struct {
	mu  sync.Mutex
	Map map[string]CacheValue
}

// CacheValue is a cache value.
type CacheValue struct {
	Generated string
	Merged    string
}

// NewCache creates a new cache.
func NewCache() *Cache {
	return &Cache{
		mu:  sync.Mutex{},
		Map: make(map[string]CacheValue),
	}
}

// NewCacheFromMaps creates a new cache from the given maps.
func NewCacheFromMaps(raw, merged map[string]string) *Cache {
	ir := make(map[string]CacheValue)
	for k, v := range raw {
		ir[k] = CacheValue{
			Generated: v,
			Merged:    merged[k],
		}
	}
	return &Cache{
		mu:  sync.Mutex{},
		Map: ir,
	}
}

// Get returns the value for the given key in the merged cache.
func (c *Cache) Get(key string) (string, string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Map[key].Generated, c.Map[key].Merged
}

// Set sets the value for the given key in the merged cache.
func (c *Cache) Set(raw, merged string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	generated := fmt.Sprintf("tw-%d", len(c.Map)+1)
	c.Map[raw] = CacheValue{
		Generated: generated,
		Merged:    merged,
	}
	return generated
}

// All returns an iterator over the string keys and CacheValue values in the cache.
func (c *Cache) All() func(yield func(string, CacheValue) bool) {
	return func(yield func(string, CacheValue) bool) {
		c.mu.Lock()
		defer c.mu.Unlock()

		for k, v := range c.Map {
			if !yield(k, v) {
				return
			}
		}
	}
}
