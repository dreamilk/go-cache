package gocache

import (
	"testing"
)

func TestLRUCache(t *testing.T) {
	cache := NewLRUCache[int](3)
	cache.Set("key1", 1)
	cache.Set("key2", 2)
	cache.Set("key3", 3)

	if cache.Count() != 3 {
		t.Errorf("Expected cache to have 3 items, got %d", cache.Count())
	}

	value, ok := cache.Get("key1")
	if !ok {
		t.Errorf("Expected key1 to be in cache")
	}
	if value != 1 {
		t.Errorf("Expected key1 to have value 1, got %d", value)
	}
	// key1 key3 key2

	cache.Delete("key2")
	if cache.Count() != 2 {
		t.Errorf("Expected cache to have 2 items, got %d", cache.Count())
	}
	// key1 key3

	cache.Set("key4", 4)
	if cache.Count() != 3 {
		t.Errorf("Expected cache to have 3 items, got %d", cache.Count())
	}
	// key4 key1 key3

	cache.Set("key5", 5)
	if cache.Count() != 3 {
		t.Errorf("Expected cache to have 3 items, got %d", cache.Count())
	}
	// key5 key4 key1

	_, ok = cache.Get("key3")
	if ok {
		t.Errorf("Expected key3 to be evicted")
	}
}

func TestLRUCache_UpdateExistingKey(t *testing.T) {
	cache := NewLRUCache[string](2)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key1", "updated1") // 更新已存在的key

	if cache.Count() != 2 {
		t.Errorf("Expected 2 items, got %d", cache.Count())
	}

	value, ok := cache.Get("key1")
	if !ok || value != "updated1" {
		t.Errorf("Expected updated value 'updated1', got %v", value)
	}
}

func TestLRUCache_EvictionOrder(t *testing.T) {
	cache := NewLRUCache[int](3)

	// 按顺序添加元素
	cache.Set("key1", 1) // [key1]
	cache.Set("key2", 2) // [key2, key1]
	cache.Set("key3", 3) // [key3, key2, key1]

	// 访问 key1，使其变为最近使用
	cache.Get("key1") // [key1, key3, key2]

	// 添加新元素，应该淘汰 key2
	cache.Set("key4", 4) // [key4, key1, key3]

	if _, ok := cache.Get("key2"); ok {
		t.Error("key2 should have been evicted")
	}

	// 验证其他键是否存在
	if _, ok := cache.Get("key1"); !ok {
		t.Error("key1 should still exist")
	}
	if _, ok := cache.Get("key3"); !ok {
		t.Error("key3 should still exist")
	}
	if _, ok := cache.Get("key4"); !ok {
		t.Error("key4 should exist")
	}
}

func TestLRUCache_ClearCache(t *testing.T) {
	cache := NewLRUCache[int](3)

	cache.Set("key1", 1)
	cache.Set("key2", 2)
	cache.Set("key3", 3)

	cache.Clear()

	if cache.Count() != 0 {
		t.Errorf("Expected empty cache after clear, got count: %d", cache.Count())
	}

	// 验证所有键都已被删除
	for _, key := range []string{"key1", "key2", "key3"} {
		if _, ok := cache.Get(key); ok {
			t.Errorf("Key %s should not exist after clear", key)
		}
	}
}

func TestLRUCache_DifferentTypes(t *testing.T) {
	// 测试不同类型的缓存
	t.Run("string cache", func(t *testing.T) {
		cache := NewLRUCache[string](2)
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")

		if v, ok := cache.Get("key1"); !ok || v != "value1" {
			t.Error("Failed string cache test")
		}
	})

	t.Run("struct cache", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		cache := NewLRUCache[Person](2)
		person1 := Person{"Alice", 25}
		person2 := Person{"Bob", 30}

		cache.Set("p1", person1)
		cache.Set("p2", person2)

		if v, ok := cache.Get("p1"); !ok || v != person1 {
			t.Error("Failed struct cache test")
		}
	})
}
