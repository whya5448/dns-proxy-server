package cache

type Cache interface {

	GetName() string
	Get(key interface{}) interface{}
	Put(key, value interface{})

}