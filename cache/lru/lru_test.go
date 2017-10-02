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
