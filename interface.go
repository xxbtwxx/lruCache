package lruCache

type cacheInterface interface {
	Set(cacheKey, cacheValue)
	Item(cacheKey) (cacheValue, bool)
}
