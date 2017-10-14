package store

import (
	"github.com/mageddo/dns-proxy-server/cache/lru"
	"github.com/mageddo/dns-proxy-server/cache"
)

var c cache.Cache;
func init(){
	c = lru.New(256);
}

//
// Singleton cache
//
func GetInstance() cache.Cache {
	return c
}
