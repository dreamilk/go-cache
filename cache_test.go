package gocache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := New[int](time.Second, time.Second)
	cache.Set("key", 1, time.Second)
	value, ok := cache.Get("key")
	if !ok {
		t.Error("key not found")
	}
	if value != 1 {
		t.Error("value is not 1")
	}

	time.Sleep(time.Second)
	_, ok = cache.Get("key")
	if ok {
		t.Error("key should be expired")
	}
}
