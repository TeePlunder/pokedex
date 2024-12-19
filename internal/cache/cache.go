package cache

import (
	"sync"
	"time"
)

type Cache struct {
	mutex    sync.Mutex
	entries  map[string]cacheEntry
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	return &Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
	}
}

func (cache *Cache) Add(key string, val []byte) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}
