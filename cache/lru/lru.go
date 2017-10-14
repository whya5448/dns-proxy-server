package lru

import (
	"github.com/mageddo/dns-proxy-server/cache"
	"github.com/hashicorp/golang-lru"
	"github.com/mageddo/dns-proxy-server/log"
)

type LRUCache struct {
	cache *lru.Cache
}

func (c *LRUCache) Get(key interface{}) interface{} {
	v, _ := c.cache.Get(key)
	return v
}

func (c *LRUCache) Put(key, value interface{}) {
	c.cache.Add(key, value)
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


