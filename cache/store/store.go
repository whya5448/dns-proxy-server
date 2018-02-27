package store

import (
	"github.com/mageddo/dns-proxy-server/cache/lru"
	"github.com/mageddo/dns-proxy-server/cache"
)

var c cache.Cache
func init(){
	c = lru.New(43690); // about 1 MB considering HostnameVo struct
}

//
// Singleton cache
//
func GetInstance() cache.Cache {
	return c
}
