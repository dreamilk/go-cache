package gocache

import (
	"fmt"
	"reflect"
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

func TestCache_Set(t *testing.T) {
	// Test basic types
	t.Run("int", func(t *testing.T) {
		cache := New[int](time.Second, time.Second)
		cache.Set("key", 1, time.Second)
		v, ok := cache.Get("key")
		if !ok {
			t.Error("key not found")
		}
		if v != 1 {
			t.Errorf("expected 1, got %v", v)
		}
	})

	t.Run("string", func(t *testing.T) {
		cache := New[string](time.Second, time.Second)
		cache.Set("key", "value", time.Second)
		v, ok := cache.Get("key")
		if !ok {
			t.Error("key not found")
		}
		if v != "value" {
			t.Errorf("expected 'value', got %v", v)
		}
	})

	// Test complex types
	t.Run("slice", func(t *testing.T) {
		cache := New[[]int](time.Second, time.Second)
		expected := []int{1, 2, 3}
		cache.Set("key", expected, time.Second)
		v, ok := cache.Get("key")
		if !ok {
			t.Error("key not found")
		}
		if !reflect.DeepEqual(v, expected) {
			t.Errorf("expected %v, got %v", expected, v)
		}
	})

	t.Run("map", func(t *testing.T) {
		cache := New[map[string]int](time.Second, time.Second)
		expected := map[string]int{"a": 1, "b": 2}
		cache.Set("key", expected, time.Second)
		v, ok := cache.Get("key")
		if !ok {
			t.Error("key not found")
		}
		if !reflect.DeepEqual(v, expected) {
			t.Errorf("expected %v, got %v", expected, v)
		}
	})

	t.Run("struct", func(t *testing.T) {
		cache := New[struct {
			A int
			B string
		}](time.Second, time.Second)
		expected := struct {
			A int
			B string
		}{A: 1, B: "value"}
		cache.Set("key", expected, time.Second)
		v, ok := cache.Get("key")
		if !ok {
			t.Error("key not found")
		}
		if v != expected {
			t.Errorf("expected %v, got %v", expected, v)
		}
	})

	t.Run("pointer", func(t *testing.T) {
		cache := New[*int](time.Second, time.Second)
		expected := 1
		cache.Set("key", &expected, time.Second)
		v, ok := cache.Get("key")
		if !ok {
			t.Error("key not found")
		}
		if *v != expected {
			t.Errorf("expected %v, got %v", expected, *v)
		}
	})

	t.Run("interface", func(t *testing.T) {
		cache := New[interface{}](time.Second, time.Second)
		expected := 1
		cache.Set("key", expected, time.Second)
		v, ok := cache.Get("key")
		if !ok {
			t.Error("key not found")
		}
		if v != expected {
			t.Errorf("expected %v, got %v", expected, v)
		}
	})

	t.Run("struct_pointer", func(t *testing.T) {
		cache := New[*struct {
			A int
			B string
		}](time.Second, time.Second)
		expected := &struct {
			A int
			B string
		}{A: 1, B: "value"}
		cache.Set("key", expected, time.Second)
		v, ok := cache.Get("key")
		if !ok {
			t.Error("key not found")
		}
		if v != expected {
			t.Errorf("expected %v, got %v", expected, v)
		}
	})
}

func TestCache_Delete(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		cache := New[int](time.Second, time.Second)
		cache.Set("key", 1, time.Second)
		cache.Delete("key")
		_, ok := cache.Get("key")
		if ok {
			t.Error("key should be deleted")
		}
	})

	t.Run("string", func(t *testing.T) {
		cache := New[string](time.Second, time.Second)
		cache.Set("key", "value", time.Second)
		cache.Delete("key")
		_, ok := cache.Get("key")
		if ok {
			t.Error("key should be deleted")
		}
	})

	t.Run("slice", func(t *testing.T) {
		cache := New[[]int](time.Second, time.Second)
		expected := []int{1, 2, 3}
		cache.Set("key", expected, time.Second)
		cache.Delete("key")
		_, ok := cache.Get("key")
		if ok {
			t.Error("key should be deleted")
		}
	})

	t.Run("map", func(t *testing.T) {
		cache := New[map[string]int](time.Second, time.Second)
		expected := map[string]int{"a": 1, "b": 2}
		cache.Set("key", expected, time.Second)
		cache.Delete("key")
		_, ok := cache.Get("key")
		if ok {
			t.Error("key should be deleted")
		}
	})

	t.Run("struct", func(t *testing.T) {
		cache := New[struct {
			A int
			B string
		}](time.Second, time.Second)
		expected := struct {
			A int
			B string
		}{A: 1, B: "value"}
		cache.Set("key", expected, time.Second)
		cache.Delete("key")
		_, ok := cache.Get("key")
		if ok {
			t.Error("key should be deleted")
		}
	})

	t.Run("pointer", func(t *testing.T) {
		cache := New[*int](time.Second, time.Second)
		expected := 1
		cache.Set("key", &expected, time.Second)
		cache.Delete("key")
		_, ok := cache.Get("key")
		if ok {
			t.Error("key should be deleted")
		}
	})

	t.Run("interface", func(t *testing.T) {
		cache := New[interface{}](time.Second, time.Second)
		expected := 1
		cache.Set("key", expected, time.Second)
		cache.Delete("key")
		_, ok := cache.Get("key")
		if ok {
			t.Error("key should be deleted")
		}
	})

	t.Run("struct_pointer", func(t *testing.T) {
		cache := New[*struct {
			A int
			B string
		}](time.Second, time.Second)
		expected := &struct {
			A int
			B string
		}{A: 1, B: "value"}
		cache.Set("key", expected, time.Second)
		cache.Delete("key")
		_, ok := cache.Get("key")
		if ok {
			t.Error("key should be deleted")
		}
	})
}

func TestCache_Cleanup(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		cache := New[int](time.Second, time.Second)
		cache.Set("key", 1, time.Second)

		time.Sleep(time.Second * 2)

		_, ok := cache.Get("key")
		if ok {
			t.Error("key should be deleted")
		}
	})

	t.Run("string", func(t *testing.T) {
		cache := New[string](time.Second, time.Second)
		cache.Set("key", "value", time.Second)

		time.Sleep(time.Second * 2)

		_, ok := cache.Get("key")
		if ok {
			t.Error("key should be deleted")
		}
	})

	t.Run("slice", func(t *testing.T) {
		cache := New[[]int](time.Second, time.Second)
		expected := []int{1, 2, 3}
		cache.Set("key", expected, time.Second)

		time.Sleep(time.Second * 2)

		_, ok := cache.Get("key")
		if ok {
			t.Error("key should be deleted")
		}
	})

	t.Run("map", func(t *testing.T) {
		cache := New[map[string]int](time.Second, time.Second)
		expected := map[string]int{"a": 1, "b": 2}
		cache.Set("key", expected, time.Second)

		time.Sleep(time.Second * 2)

		_, ok := cache.Get("key")
		if ok {
			t.Error("key should be deleted")
		}
	})

	t.Run("struct", func(t *testing.T) {
		cache := New[struct {
			A int
			B string
		}](time.Second, time.Second)
		expected := struct {
			A int
			B string
		}{A: 1, B: "value"}
		cache.Set("key", expected, time.Second)

		time.Sleep(time.Second * 2)

		_, ok := cache.Get("key")
		if ok {
			t.Error("key should be deleted")
		}
	})

	t.Run("pointer", func(t *testing.T) {
		cache := New[*int](time.Second, time.Second)
		expected := 1
		cache.Set("key", &expected, time.Second)

		time.Sleep(time.Second * 2)

		_, ok := cache.Get("key")
		if ok {
			t.Error("key should be deleted")
		}
	})

	t.Run("interface", func(t *testing.T) {
		cache := New[interface{}](time.Second, time.Second)
		expected := 1
		cache.Set("key", expected, time.Second)

		time.Sleep(time.Second * 2)

		_, ok := cache.Get("key")
		if ok {
			t.Error("key should be deleted")
		}
	})

	t.Run("struct_pointer", func(t *testing.T) {
		cache := New[*struct {
			A int
			B string
		}](time.Second, time.Second)
		expected := &struct {
			A int
			B string
		}{A: 1, B: "value"}
		cache.Set("key", expected, time.Second)

		time.Sleep(time.Second * 2)

		_, ok := cache.Get("key")
		if ok {
			t.Error("key should be deleted")
		}
	})
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

			// 预填充缓存
			for i := 0; i < size; i++ {
				key := fmt.Sprintf("key_%d", i)
				cache.Set(key, i, time.Minute)
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				key := fmt.Sprintf("key_%d", i%size)
				cache.Get(key)
			}
		})
	}
}

func BenchmarkCache_LargeData(b *testing.B) {
	type LargeData struct {
		Data [1024]byte // 1KB 数据
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
