package lru

import "container/list"

type LruCache interface {
	// Якщо наш кеш вже повний (ми досягли нашого capacity)
	// то має видалитись той елемент, який ми до якого ми доступались (читали) найдавніше
	Put(key, value string)
	Get(key string) (string, bool)
}

type LruCacheImpl struct {
	capacity int
	cache map[string]*list.Element
	list *list.List
}

type entry struct {
	key   string
	value string
}

func NewLruCache(capacity int) LruCache {
	return LruCacheImpl{
		capacity: capacity,
		cache: make(map[string]*list.Element),
		list: list.New(),
	}
}

func (c LruCacheImpl) Get(key string) (string, bool) {
	if elem, found := c.cache[key]; found {
		c.list.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}
	return "", false
}

func (c LruCacheImpl) Put(key, value string) {
	if elem, found := c.cache[key]; found {
		elem.Value.(*entry).value = value
		c.list.MoveToFront(elem)
		return
	}

	if c.list.Len() >= c.capacity {
		oldest := c.list.Back()
		if oldest != nil {
			c.list.Remove(oldest)
			kv := oldest.Value.(*entry)
			delete(c.cache, kv.key)
		}
	}

	newEntry := &entry{key: key, value: value}
	listElement := c.list.PushFront(newEntry)
	c.cache[key] = listElement
}
