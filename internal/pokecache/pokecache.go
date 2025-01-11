package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries map[string]cacheEntry
	mu      sync.RWMutex
}

func NewCache(interval time.Duration) *Cache {
	newCache := Cache{}
	newCache.entries = make(map[string]cacheEntry)
	go newCache.reapLoop(interval)
	return &newCache
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.mu.Lock()
		if c.entries == nil {
			continue
		}
		for key, entry := range c.entries {
			// need to check that UNIT
			if time.Since(entry.createdAt) > interval {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}

func (c *Cache) Add(key string, val []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	newEntry := cacheEntry{}
	newEntry.createdAt = time.Now()
	newEntry.val = val
	c.entries[key] = newEntry
	return nil
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val := c.entries[key].val
	if val == nil {
		return nil, false
	} else {
		return val, true
	}
}
