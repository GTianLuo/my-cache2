package mycache

import (
	my_cache2 "my-cache2"
	"my-cache2/evict"
	"sync"
)

type mainCache struct {
	cache    *evict.Cache
	maxBytes uint64
	mu       sync.Mutex
}

func (c *mainCache) get(key string) (my_cache2.BytesValue, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cache == nil {
		c.cache = evict.NewCache(c.maxBytes)
		c.cache.DeleteExpired()
	}
	value, ok := c.cache.Get(key)
	if !ok {
		return my_cache2.BytesValue{}, false
	}
	return value.(my_cache2.BytesValue), true
}

func (c *mainCache) add(key string, value my_cache2.BytesValue, expire int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cache == nil {
		c.cache = evict.NewCache(c.maxBytes)
		c.cache.DeleteExpired()
	}
	return c.cache.Add(key, value, expire)
}
