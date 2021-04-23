package lruCache

type chanLRUCache struct {
	operations chan func(*chanLRUCache)
	data       cacheType
	usageOrder []cacheKey
	maxSize    int
}

func NewLRUChan(size int) *chanLRUCache {
	c := &chanLRUCache{
		operations: make(chan func(*chanLRUCache)),
		data:       make(cacheType),
		usageOrder: make([]cacheKey, 0, size),
		maxSize:    size,
	}

	go c.opLoop()

	return c
}

func (c *chanLRUCache) opLoop() {
	for op := range c.operations {
		op(c)
	}
}

func (c *chanLRUCache) Set(key cacheKey, data cacheValue) {
	c.operations <- func(obj *chanLRUCache) {
		if _, ok := obj.data[key]; ok {
			moveToHead(obj.usageOrder, key)
		} else {
			if c.isFull() {
				delete(obj.data, obj.usageOrder[0])
				obj.usageOrder = obj.usageOrder[1:]
			}

			obj.usageOrder = append(obj.usageOrder, key)
		}

		obj.data[key] = data
	}
}

func (c *chanLRUCache) Item(key cacheKey) (cacheValue, bool) {
	value := make(chan cacheValue, 1)
	exist := make(chan bool, 1)
	c.operations <- func(obj *chanLRUCache) {
		val, ok := obj.data[key]
		value <- val
		exist <- ok

		if ok {
			moveToHead(obj.usageOrder, key)
		}
	}

	return <-value, <-exist
}

func (c *chanLRUCache) isFull() bool {
	return len(c.usageOrder) == c.maxSize
}
