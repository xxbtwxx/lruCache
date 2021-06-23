package lruCache

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChanLRU(t *testing.T) {
	cache := NewLRUActorModel(4)
	testCache(cache, t)
}

func TestMutexLRU(t *testing.T) {
	cache := NewLRUMutex(4)
	testCache(cache, t)
}

func TestSpinlockLRU(t *testing.T) {
	cache := NewLRUSpinlock(4)
	testCache(cache, t)
}

func testCache(cache cacheInterface, t *testing.T) {
	cache.Set(1, 1)
	val, ok := cache.Item(1)
	require.Equal(t, val, cacheValue(1))
	require.Equal(t, ok, true)
	cache.Set(2, 2)
	cache.Set(3, 3)
	cache.Set(4, 4)
	cache.Set(5, 5)

	val, ok = cache.Item(2)
	require.Equal(t, val, cacheValue(2))
	require.Equal(t, ok, true)

	val, ok = cache.Item(3)
	require.Equal(t, val, cacheValue(3))
	require.Equal(t, ok, true)

	val, ok = cache.Item(1)
	require.Equal(t, val, cacheValue(0))
	require.Equal(t, ok, false)

	val, ok = cache.Item(6)
	require.Equal(t, val, cacheValue(0))
	require.Equal(t, ok, false)
}

var opsPerGoroutine = 5
var key cacheKey = 1
var value cacheValue = 1

func benchmarkCache(count int, c cacheInterface, b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < count; i++ {
			go func() {
				for j := 0; j < opsPerGoroutine; j++ {
					c.Set(key, value)
					c.Item(key)
				}
			}()
		}
	}
}

var cacheSize = 4
var numCPU = runtime.NumCPU()

func BenchmarkActorModel1(b *testing.B) {
	cache := NewLRUActorModel(cacheSize)
	benchmarkCache(1, cache, b)
}
func BenchmarkMutex1(b *testing.B) {
	cache := NewLRUMutex(cacheSize)
	benchmarkCache(1, cache, b)
}
func BenchmarkSpinlock1(b *testing.B) {
	cache := NewLRUSpinlock(cacheSize)
	benchmarkCache(1, cache, b)
}

func BenchmarkActorModelNumCPU(b *testing.B) {
	cache := NewLRUActorModel(cacheSize)
	benchmarkCache(numCPU, cache, b)
}
func BenchmarkMutexNumCPU(b *testing.B) {
	cache := NewLRUMutex(cacheSize)
	benchmarkCache(numCPU, cache, b)
}
func BenchmarkSpinlockNumCPU(b *testing.B) {
	cache := NewLRUSpinlock(cacheSize)
	benchmarkCache(numCPU, cache, b)
}

func BenchmarkActorModel10(b *testing.B) {
	cache := NewLRUActorModel(cacheSize)
	benchmarkCache(10, cache, b)
}
func BenchmarkMutex10(b *testing.B) {
	cache := NewLRUMutex(cacheSize)
	benchmarkCache(10, cache, b)
}
func BenchmarkSpinlock10(b *testing.B) {
	cache := NewLRUSpinlock(cacheSize)
	benchmarkCache(10, cache, b)
}

func BenchmarkActorModel50(b *testing.B) {
	cache := NewLRUActorModel(cacheSize)
	benchmarkCache(50, cache, b)
}
func BenchmarkMutex50(b *testing.B) {
	cache := NewLRUMutex(cacheSize)
	benchmarkCache(50, cache, b)
}
func BenchmarkSpinlock50(b *testing.B) {
	cache := NewLRUSpinlock(cacheSize)
	benchmarkCache(50, cache, b)
}
