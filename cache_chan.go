package lruCache

type actorModelLRUCache struct {
	operations chan func()
	data       cacheType
	usageOrder []cacheKey
	maxSize    int
}

func NewLRUActorModel(size int) *actorModelLRUCache {
	c := &actorModelLRUCache{
		operations: make(chan func()),
		data:       make(cacheType),
		usageOrder: make([]cacheKey, 0, size),
		maxSize:    size,
	}

	go c.opLoop()

	return c
}

func (c *actorModelLRUCache) opLoop() {
	for op := range c.operations {
		op()
	}
}

func (c *actorModelLRUCache) Set(key cacheKey, data cacheValue) {
	c.operations <- func() {
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
	}
}

func (c *actorModelLRUCache) Item(key cacheKey) (cacheValue, bool) {
	value := make(chan cacheValue, 1)
	exist := make(chan bool, 1)
	c.operations <- func() {
		val, ok := c.data[key]
		value <- val
		exist <- ok

		if ok {
			moveToHead(c.usageOrder, key)
		}
	}

	return <-value, <-exist
}

func (c *actorModelLRUCache) isFull() bool {
	return len(c.usageOrder) == c.maxSize
}
