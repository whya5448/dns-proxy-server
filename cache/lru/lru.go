package lru

import "github.com/mageddo/dns-proxy-server/cache"

type LRUCache struct {
	name string
	timeout int64
	size int64
	cache map[interface{}]interface{}
}

func (c *LRUCache) GetName() string {
	return c.name
}

func (c *LRUCache) Get(key interface{}) interface{} {
	return c.cache[key]
}

func (c *LRUCache) Put(key, value interface{}) {
	c.cache[key] = value
}

//
// Creates a LRU cache
// size is the maximum size of the cache, -1 if it is unlimited
// time time in millis before cache expires, -1 if it is unlimited
//
func NewLRUCache(name string, size, timeout int64) cache.Cache {
	return &LRUCache{name, timeout, size, make(map[interface{}]interface{})}
}

