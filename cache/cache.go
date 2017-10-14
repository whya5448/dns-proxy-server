package cache

type Cache interface {

	Get(key interface{}) interface{}
	Put(key, value interface{})

}
