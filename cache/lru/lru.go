package lru

import (
	"github.com/mageddo/dns-proxy-server/cache"
	"github.com/hashicorp/golang-lru"
	"github.com/mageddo/dns-proxy-server/log"
)

type LRUCache struct {
	cache *lru.Cache
}

func (c *LRUCache) ContainsKey(key interface{}) bool {
	return c.cache.Contains(key)
}

func (c *LRUCache) Get(key interface{}) interface{} {
	v, _ := c.cache.Get(key)
	return v
}

//
// Put value in cache, it doesn't have guarantee of concurrency treat
//
func (c *LRUCache) Put(key, value interface{}) {
	c.cache.Add(key, value)
}

//
// Check if value is already associated, if yes just return it, if not put the passed value and return nil
// This method must be thread safe (atomic)
//
func (c *LRUCache) PutIfAbsent(key, value interface{}) interface{} {
	if ok, _ := c.cache.ContainsOrAdd(key, value); ok {
		return c.Get(key)
	}
	return nil;
}

//
// Creates a LRU cache
// size is the maximum size of the cache, -1 if it is unlimited
//
func New(size int) cache.Cache {
	c, err := lru.New(size)
	if err != nil {
		log.LOGGER.Errorf("status=cannot-create-cache, msg=%v", err)
		return nil;
	}
	return &LRUCache{c}
}


