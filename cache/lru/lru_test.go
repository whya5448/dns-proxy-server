package lru

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPutAndGetSuccess(t *testing.T){

	cache := NewLRUCache("test1", -1, -1);
	cache.Put("key1", "value1");

	assert.Equal(t, "value1", cache.Get("key1").(string))

}

func TestPutAndGetSuccessSizeLimited(t *testing.T){

	cache := NewLRUCache("test1", 2, -1);
	cache.Put("key1", "value1");
	cache.Put("key2", "value2");
	cache.Put("key3", "value3");

	assert.Nil(t, cache.Get("key1"))
	assert.Equal(t, "value2", cache.Get("key2").(string))
	assert.Equal(t, "value3", cache.Get("key3").(string))

}
