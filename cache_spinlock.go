package lruCache

import (
	"github.com/xxbtwxx/lruCache/spinlock"
)

type spinlockLRUCache struct {
	data       cacheType
	spinlock   *spinlock.Lock
	usageOrder []cacheKey
	maxSize    int
}

func NewLRUSpinlock(size int) *spinlockLRUCache {
	c := &spinlockLRUCache{
		data:       make(cacheType),
		usageOrder: make([]cacheKey, 0, size),
		maxSize:    size,
		spinlock:   spinlock.New(),
	}

	return c
}

func (c *spinlockLRUCache) Set(key cacheKey, data cacheValue) {
	c.spinlock.Lock()

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

	c.spinlock.Unlock()
}

func (c *spinlockLRUCache) Item(key cacheKey) (cacheValue, bool) {
	c.spinlock.Lock()

	val, ok := c.data[key]
	if ok {
		moveToHead(c.usageOrder, key)
	}

	c.spinlock.Unlock()

	return val, ok
}

func (c *spinlockLRUCache) isFull() bool {
	return len(c.usageOrder) == c.maxSize
}
