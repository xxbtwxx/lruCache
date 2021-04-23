package lruCache

import "sync"

type mutexLRUCache struct {
	data       cacheType
	mu         *sync.Mutex
	usageOrder []cacheKey
	maxSize    int
}

func NewLRUMutex(size int) *mutexLRUCache {
	c := &mutexLRUCache{
		data:       make(cacheType),
		usageOrder: make([]cacheKey, 0, size),
		maxSize:    size,
		mu:         new(sync.Mutex),
	}

	return c
}

func (c *mutexLRUCache) Set(key cacheKey, data cacheValue) {
	c.mu.Lock()

	if _, ok := c.data[key]; ok {
		moveToHead(c.usageOrder, key)
	} else {
		if c.isFull() {
			delete(c.data, c.usageOrder[0])
			c.usageOrder = c.usageOrder[1:]
		}

		c.usageOrder = append(c.usageOrder, key)
	}

	c.data[key] = data

	c.mu.Unlock()
}

func (c *mutexLRUCache) Item(key cacheKey) (cacheValue, bool) {
	c.mu.Lock()

	val, ok := c.data[key]
	if ok {
		moveToHead(c.usageOrder, key)
	}

	c.mu.Unlock()

	return val, ok
}

func (c *mutexLRUCache) isFull() bool {
	return len(c.usageOrder) == c.maxSize
}
