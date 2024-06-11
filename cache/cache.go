package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	Data      interface{}
	Timestamp time.Time
}

type Cache struct {
	items map[string]*CacheItem
	mu    sync.Mutex
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]*CacheItem),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, found := c.items[key]
	if found {
		item.Timestamp = time.Now()
	}
	return item, found
}

func (c *Cache) Set(key string, data interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &CacheItem{
		Data:      data,
		Timestamp: time.Now(),
	}
}

func (c *Cache) ClearOldItems(maxAge time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for k, v := range c.items {
		if now.Sub(v.Timestamp) > maxAge {
			delete(c.items, k)
		}
	}
}
