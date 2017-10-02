package lru

import (
	. "github.com/mageddo/dns-proxy-server/log"
	"github.com/mageddo/dns-proxy-server/cache"
	"container/list"
)

type LRUCache struct {
	name string
	size int
	timeout int64
	cache map[interface{}]interface{}
	keyset *list.List
}

func (c *LRUCache) GetName() string {
	return c.name
}

func (c *LRUCache) Get(key interface{}) interface{} {
	return c.cache[key]
}

func (c *LRUCache) Put(key, value interface{}) {
	if c.size > 0 {
		if c.keyset.Len() == c.size {
			lastKey := c.keyset.Back()
			LOGGER.Debugf("status=size-limit-reached, size=%d, key=%v, keyToRemove=%v", c.size, key, lastKey.Value)
			c.keyset.Remove(lastKey)
			delete(c.cache, lastKey.Value)
		}
		c.keyset.PushFront(key)
	}
	c.cache[key] = value
}

//
// Creates a LRU cache
// size is the maximum size of the cache, -1 if it is unlimited
// time time in millis before cache expires, -1 if it is unlimited
//
func NewLRUCache(name string, size int, timeout int64) cache.Cache {
	return &LRUCache{name, size, timeout, make(map[interface{}]interface{}), list.New()}
}

