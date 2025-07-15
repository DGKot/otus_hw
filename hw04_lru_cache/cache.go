package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type CacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.RWMutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	item, ok := lru.items[key]
	if ok {
		item.Value.(*CacheItem).value = value
		lru.queue.MoveToFront(item)
		return true
	}
	cacheItem := &CacheItem{
		key:   key,
		value: value,
	}
	item = lru.queue.PushFront(cacheItem)
	lru.items[key] = item
	if lru.capacity < lru.queue.Len() {
		deleted := lru.queue.Back()
		lru.queue.Remove(deleted)
		delete(lru.items, deleted.Value.(*CacheItem).key)
	}
	return false
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	lru.mu.RLock()
	defer lru.mu.RUnlock()
	item, ok := lru.items[key]
	if !ok {
		return nil, false
	}
	lru.queue.MoveToFront(item)
	return item.Value.(*CacheItem).value, true
}

func (lru *lruCache) Clear() {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	lru.items = make(map[Key]*ListItem, lru.capacity)
	lru.queue = NewList()
}
