package lruCache

type cacheKey int
type cacheValue int

type cacheType map[cacheKey]cacheValue

func moveToHead(slice []cacheKey, key cacheKey) {
	index := elementIndex(slice, key)
	slice = append(slice[:index], slice[index+1:]...)
	slice = append(slice, key)
}

func elementIndex(slice []cacheKey, data cacheKey) int {
	var index int

	for i, val := range slice {
		if val == data {
			index = i
			break
		}
	}

	return index
}
