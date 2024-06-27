package hw04lrucache

type Key string

type Pair struct {
	key   Key
	value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewPair(k Key, v interface{}) *Pair {
	return &Pair{
		key:   k,
		value: v,
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lru *lruCache) Clear() {
	lru.items = make(map[Key]*ListItem, lru.capacity)
	lru.queue = NewList()
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	val, ok := lru.items[key]
	if ok {
		item, _ := val.Value.(Pair)
		item.value = value
		val.Value = item
		lru.queue.MoveToFront(val)
	} else {
		if lru.queue.Len() == lru.capacity {
			back := lru.queue.Back()
			lru.queue.Remove(back)
			delete(lru.items, back.Value.(Pair).key)
		}
		lru.items[key] = lru.queue.PushFront(*NewPair(key, value))
	}

	return ok
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	val, ok := lru.items[key]
	if ok {
		lru.queue.MoveToFront(val)
		return val.Value.(Pair).value, true
	}
	return nil, false
}
