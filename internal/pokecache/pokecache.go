package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheEntries map[string]cacheEntry
	mutex        *sync.Mutex
	interval     time.Duration
	done         chan struct{}
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	newCache := Cache{
		cacheEntries: make(map[string]cacheEntry),
		mutex:        &sync.Mutex{},
		interval:     interval,
		done:         make(chan struct{}),
	}

	go newCache.reapLoop()

	return &newCache
}

func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cacheEntries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, ok := c.cacheEntries[key]
	if ok {
		return entry.val, true
	}
	return nil, false
}

func (c *Cache) reapLoop() {
	t := time.NewTicker(c.interval)
	for {
		select {
		case <-c.done:
			t.Stop()
			return
		case <-t.C:
			c.mutex.Lock()
			for k, entry := range c.cacheEntries {
				if time.Since(entry.createdAt) > c.interval {
					delete(c.cacheEntries, k)
				}
			}
			c.mutex.Unlock()
		}
	}
}
