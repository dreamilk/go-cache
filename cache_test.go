package gocache

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

func intPtr(i int) *int {
	return &i
}

func TestCache_Set(t *testing.T) {
	runCacheSetTests(t, "int", New[int], 1, 1)
	runCacheSetTests(t, "string", New[string], "value", "value")
	runCacheSetTests(t, "slice", New[[]int], []int{1, 2, 3}, []int{1, 2, 3})
	runCacheSetTests(t, "map", New[map[string]int], map[string]int{"a": 1, "b": 2}, map[string]int{"a": 1, "b": 2})

	type TestStruct struct {
		A int
		B string
	}
	runCacheSetTests(t, "struct", New[TestStruct], TestStruct{A: 1, B: "value"}, TestStruct{A: 1, B: "value"})

	p := intPtr(1)
	runCacheSetTests(t, "pointer", New[*int], p, p)
	runCacheSetTests(t, "interface", New[interface{}], 1, 1)
	runCacheSetTests(t, "struct_pointer", New[*TestStruct], &TestStruct{A: 1, B: "value"}, &TestStruct{A: 1, B: "value"})
}

func runCacheSetTests[T any](t *testing.T, name string, newCache func(time.Duration, time.Duration) *Cache[T], expected T, value T) {
	t.Run(name, func(t *testing.T) {
		cache := newCache(time.Second, time.Second)
		cache.Set("key", value, time.Second)
		v, ok := cache.Get("key")
		if !ok {
			t.Fatalf("key not found in cache")
		}
		if !reflect.DeepEqual(v, expected) {
			t.Fatalf("expected %v, got %v", expected, v)
		}
	})
}

func TestCache_Delete(t *testing.T) {
	runCacheDeleteTests(t, "int", New[int], 1)
	runCacheDeleteTests(t, "string", New[string], "value")
	runCacheDeleteTests(t, "slice", New[[]int], []int{1, 2, 3})
	runCacheDeleteTests(t, "map", New[map[string]int], map[string]int{"a": 1, "b": 2})

	type TestStruct struct {
		A int
		B string
	}
	runCacheDeleteTests(t, "struct", New[TestStruct], TestStruct{A: 1, B: "value"})
	runCacheDeleteTests(t, "pointer", New[*int], intPtr(1))
	runCacheDeleteTests(t, "interface", New[interface{}], 1)
	runCacheDeleteTests(t, "struct_pointer", New[*TestStruct], &TestStruct{A: 1, B: "value"})
}

func runCacheDeleteTests[T any](t *testing.T, name string, newCache func(time.Duration, time.Duration) *Cache[T], value T) {
	t.Run(name, func(t *testing.T) {
		cache := newCache(time.Second, time.Second)
		cache.Set("key", value, time.Second)
		cache.Delete("key")
		if _, ok := cache.Get("key"); ok {
			t.Fatal("expected key to be deleted from cache")
		}
	})
}

func TestCache_Cleanup(t *testing.T) {
	runCacheCleanupTests(t, "int", New[int], 1)
	runCacheCleanupTests(t, "string", New[string], "value")
	runCacheCleanupTests(t, "slice", New[[]int], []int{1, 2, 3})
	runCacheCleanupTests(t, "map", New[map[string]int], map[string]int{"a": 1, "b": 2})

	type TestStruct struct {
		A int
		B string
	}
	runCacheCleanupTests(t, "struct", New[TestStruct], TestStruct{A: 1, B: "value"})
	runCacheCleanupTests(t, "pointer", New[*int], intPtr(1))
	runCacheCleanupTests(t, "interface", New[interface{}], 1)
	runCacheCleanupTests(t, "struct_pointer", New[*TestStruct], &TestStruct{A: 1, B: "value"})
}

func runCacheCleanupTests[T any](t *testing.T, name string, newCache func(time.Duration, time.Duration) *Cache[T], value T) {
	t.Run(name, func(t *testing.T) {
		cache := newCache(time.Second, time.Second)
		cache.Set("key", value, time.Second)
		time.Sleep(2 * time.Second)
		if _, ok := cache.Get("key"); ok {
			t.Fatal("expected key to be removed by cleanup")
		}
	})
}

func TestCache_ConcurrentAccessDifferentTypes(t *testing.T) {
	cacheInt := New[int](time.Minute, time.Minute)
	cacheString := New[string](time.Minute, time.Minute)
	defer cacheInt.Close()
	defer cacheString.Close()

	var wg sync.WaitGroup
	numGoroutines := 100

	// Concurrently set and get int values
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("int_key_%d", i)
			cacheInt.Set(key, i, time.Minute)
			val, ok := cacheInt.Get(key)
			if !ok || val != i {
				t.Errorf("expected %d, got %v", i, val)
			}
		}(i)
	}

	// Concurrently set and get string values
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("string_key_%d", i)
			value := fmt.Sprintf("value_%d", i)
			cacheString.Set(key, value, time.Minute)
			val, ok := cacheString.Get(key)
			if !ok || val != value {
				t.Errorf("expected %s, got %v", value, val)
			}
		}(i)
	}

	wg.Wait()

	// Verify counts
	if cacheInt.Count() != numGoroutines {
		t.Fatalf("expected %d int items, got %d", numGoroutines, cacheInt.Count())
	}
	if cacheString.Count() != numGoroutines {
		t.Fatalf("expected %d string items, got %d", numGoroutines, cacheString.Count())
	}
}

func BenchmarkCache_HighConcurrency(b *testing.B) {
	cache := New[int](time.Minute, time.Minute)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := "concurrent_key"
			cache.Set(key, 1, time.Minute)
			if val, ok := cache.Get(key); !ok || val != 1 {
				b.Errorf("unexpected value: got %v, want %v", val, 1)
			}
		}
	})
}

func BenchmarkCache_MixedReadWrite(b *testing.B) {
	cache := New[int](time.Minute, time.Minute)
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key_%d", i)
			cache.Set(key, i, time.Minute)
			cache.Get(key)
			i++
		}
	})
}

func BenchmarkCache_VariousSizes(b *testing.B) {
	sizes := []int{1000, 10000, 100000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			cache := New[int](time.Minute, time.Minute)
			for i := 0; i < size; i++ {
				cache.Set(fmt.Sprintf("key_%d", i), i, time.Minute)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cache.Get(fmt.Sprintf("key_%d", i%size))
			}
		})
	}
}

func BenchmarkCache_LargeData(b *testing.B) {
	type LargeData struct {
		Data [1024]byte
	}

	cache := New[LargeData](time.Minute, time.Minute)
	largeData := LargeData{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key_%d", i)
		cache.Set(key, largeData, time.Minute)
		_, _ = cache.Get(key)
	}
}
