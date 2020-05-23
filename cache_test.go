package cache

import (
	"testing"
)

func benchmark(b *testing.B) {
	c := CreateNewCache("cache_data", 1000, false)
	c.setTTL(5)
	for i := 0; i < b.N; i++ {
		k := RandStringRunes(5)
		v := RandStringRunes(8)
		c.Set(k, v, -1)
		_, _, _ = c.Get(k)
	}
	c.close()
}

func BenchmarkSimpleCache(b *testing.B) {
	benchmark(b)
}