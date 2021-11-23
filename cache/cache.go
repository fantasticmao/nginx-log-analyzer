package cache

type Key interface{}

type Cache interface {
	Len() int

	Get(key Key) interface{}

	Put(key Key, value interface{})
}
