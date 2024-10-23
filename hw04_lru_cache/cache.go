package hw04lrucache

import (
	"sync"
)

type Key string

type listValue struct {
	key   Key
	value interface{}
}

var mutex sync.Mutex

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

func (lc *lruCache) Set(key Key, value interface{}) bool {
	mutex.Lock()
	defer mutex.Unlock()

	if li, ok := lc.items[key]; ok {
		li.Value = listValue{key: key, value: value}
		lc.queue.MoveToFront(li)
		return true
	}

	li := lc.queue.PushFront(listValue{key: key, value: value})
	lc.items[key] = li

	if lc.queue.Len() > lc.capacity {
		tail := lc.queue.Back()
		delete(lc.items, tail.Value.(listValue).key)
		lc.queue.Remove(tail)
	}

	return false
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	if li, ok := lc.items[key]; ok {
		lc.queue.MoveToFront(li)
		return li.Value.(listValue).value, true
	}

	return nil, false
}

func (lc *lruCache) Clear() {
	mutex.Lock()
	defer mutex.Unlock()

	lc.queue = NewList()
	lc.items = make(map[Key]*ListItem, lc.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
