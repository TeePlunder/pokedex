package cache

import (
	"sync"
	"time"
)

type Cache struct {
	mutex    sync.Mutex
	entries  map[string]cacheEntry
	interval time.Duration
	// stopCh chan struct{}
	// You don't strictly need stopCh, but it’s useful for a clean shutdown of your background goroutine.
	// Without it, the reapLoop() will continue running until the program exits.
	// If you want the ability to stop the loop at some point (for example, during application shutdown or testing),
	// having a stopCh channel allows you to signal the goroutine to exit gracefully.
	// If your application never needs to stop the reap loop, you could omit it. However, it’s good practice to include a stop mechanism for long-running goroutines.
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
	}

	if interval > 0 {
		go cache.reapLoop()
	}

	return cache
}

func (cache *Cache) Add(key string, val []byte) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (cache *Cache) Get(key string) (val []byte, ok bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	entry, ok := cache.entries[key]

	val = entry.val
	return
}

func (cache *Cache) Delete(key string) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	delete(cache.entries, key)
}

// clears all entries on interval ticks
func (cache *Cache) reapLoop() {
	ticker := time.NewTicker(cache.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cache.reap()
		}
	}
}

// removes all entries older then the interval of the cache
func (cache *Cache) reap() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	now := time.Now()
	for key, entry := range cache.entries {
		if now.Sub(entry.createdAt) > cache.interval {
			delete(cache.entries, key)
		}
	}
}
