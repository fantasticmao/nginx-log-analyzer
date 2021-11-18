package cache

import "container/list"

type LruCache struct {
	capacity int
	list     *list.List
	cache    map[Key]*list.Element
}

type lruEntry struct {
	key   Key
	value interface{}
}

func NewLruCache(capacity int) *LruCache {
	return &LruCache{
		capacity: capacity,
		list:     list.New(),
		cache:    make(map[Key]*list.Element),
	}
}

func (c *LruCache) Len() int {
	if c.list == nil || c.cache == nil {
		return 0
	}
	return c.list.Len()
}

func (c *LruCache) Get(key Key) interface{} {
	if c.list == nil || c.cache == nil {
		return nil
	}

	element, ok := c.cache[key]
	if !ok {
		return nil
	}
	c.list.MoveToFront(element)
	return element.Value.(*lruEntry).value
}

func (c *LruCache) Put(key Key, value interface{}) {
	if c.list == nil || c.cache == nil {
		return
	}

	element, ok := c.cache[key]
	if ok { // update
		c.list.MoveToFront(element)
		element.Value.(*lruEntry).value = value
		return
	} else { // insert
		element = c.list.PushFront(&lruEntry{key: key, value: value})
		c.cache[key] = element
		if c.list.Len() > c.capacity { // evict
			c.removeOldest()
		}
	}
}

func (c *LruCache) removeOldest() {
	oldestElement := c.list.Back()
	c.list.Remove(oldestElement)
	delete(c.cache, oldestElement.Value.(*lruEntry).key)
}
