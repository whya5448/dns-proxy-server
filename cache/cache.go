package cache

type Cache interface {

	Get(key interface{}) interface{}
	ContainsKey(key interface{}) bool
	Put(key, value interface{})
	PutIfAbsent(key, value interface{}) interface{}
	Clear()

}
