package lru

import (
	"fmt"
	"testing"
)

func TestLRUCache(t *testing.T) {
	cache := New(3)
	cache.Put(1, 1)
	cache.Put(2, 2)
	cache.Put(3, 3)
	cache.Put(4, 4)
	t.Log(cache.Print())
	fmt.Println("Get 1:", cache.Get(1))
	t.Log(cache.Print())
	fmt.Println("Get 2:", cache.Get(2))
	t.Log(cache.Print())
	cache.Put(5, 5)
	t.Log(cache.Print())
}
