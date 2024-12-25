package gocache

import (
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
	type testCase[T any] struct {
		name     string
		value    T
		expected T
	}

	testCases := []interface{}{
		testCase[int]{
			name:     "int",
			value:    1,
			expected: 1,
		},
		testCase[string]{
			name:     "string",
			value:    "value",
			expected: "value",
		},
		testCase[struct {
			A int
			B string
		}]{
			name: "struct",
			value: struct {
				A int
				B string
			}{A: 1, B: "value"},
			expected: struct {
				A int
				B string
			}{A: 1, B: "value"},
		},
		testCase[[]int]{
			name:     "slice",
			value:    []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		testCase[map[string]int]{
			name:     "map",
			value:    map[string]int{"a": 1, "b": 2},
			expected: map[string]int{"a": 1, "b": 2},
		},
		testCase[*int]{
			name: "pointer",
			value: func() *int {
				ptr := new(int)
				*ptr = 1
				return ptr
			}(),
			expected: func() *int {
				ptr := new(int)
				*ptr = 1
				return ptr
			}(),
		},
		testCase[interface{}]{
			name:     "interface",
			value:    1,
			expected: 1,
		},
	}

	for _, tc := range testCases {
		switch tcTyped := tc.(type) {
		case testCase[int]:
			t.Run(tcTyped.name, func(t *testing.T) {
				cache := New[int](time.Second, time.Second)
				cache.Set("key", tcTyped.value, time.Second)
				v, ok := cache.Get("key")
				if !ok {
					t.Errorf("key not found")
				}
				if v != tcTyped.expected {
					t.Errorf("expected %v, got %v", tcTyped.expected, v)
				}
			})
		case testCase[string]:
			t.Run(tcTyped.name, func(t *testing.T) {
				cache := New[string](time.Second, time.Second)
				cache.Set("key", tcTyped.value, time.Second)
				v, ok := cache.Get("key")
				if !ok {
					t.Errorf("key not found")
				}
				if v != tcTyped.expected {
					t.Errorf("expected %v, got %v", tcTyped.expected, v)
				}
			})
		case testCase[struct {
			A int
			B string
		}]:
			t.Run(tcTyped.name, func(t *testing.T) {
				cache := New[struct {
					A int
					B string
				}](time.Second, time.Second)
				cache.Set("key", tcTyped.value, time.Second)
				v, ok := cache.Get("key")
				if !ok {
					t.Errorf("key not found")
				}
				if v.A != tcTyped.expected.A || v.B != tcTyped.expected.B {
					t.Errorf("expected %v, got %v", tcTyped.expected, v)
				}
			})
		case testCase[[]int]:
			t.Run(tcTyped.name, func(t *testing.T) {
				cache := New[[]int](time.Second, time.Second)
				cache.Set("key", tcTyped.value, time.Second)
				v, ok := cache.Get("key")
				if !ok {
					t.Errorf("key not found")
				}
				if !reflect.DeepEqual(v, tcTyped.expected) {
					t.Errorf("expected %v, got %v", tcTyped.expected, v)
				}
			})
		case testCase[map[string]int]:
			t.Run(tcTyped.name, func(t *testing.T) {
				cache := New[map[string]int](time.Second, time.Second)
				cache.Set("key", tcTyped.value, time.Second)
				v, ok := cache.Get("key")
				if !ok {
					t.Errorf("key not found")
				}
				if !reflect.DeepEqual(v, tcTyped.expected) {
					t.Errorf("expected %v, got %v", tcTyped.expected, v)
				}
			})
		case testCase[*int]:
			t.Run(tcTyped.name, func(t *testing.T) {
				cache := New[*int](time.Second, time.Second)
				cache.Set("key", tcTyped.value, time.Second)
				v, ok := cache.Get("key")
				if !ok {
					t.Errorf("key not found")
				}
				if v == nil || *v != *tcTyped.expected {
					t.Errorf("expected %v, got %v", *tcTyped.expected, *v)
				}
			})
		case testCase[interface{}]:
			t.Run(tcTyped.name, func(t *testing.T) {
				cache := New[interface{}](time.Second, time.Second)
				cache.Set("key", tcTyped.value, time.Second)
				v, ok := cache.Get("key")
				if !ok {
					t.Errorf("key not found")
				}
				if v != tcTyped.expected {
					t.Errorf("expected %v, got %v", tcTyped.expected, v)
				}
			})
		default:
			t.Errorf("unknown test case type")
		}
	}
}

func TestCache_Get(t *testing.T) {
	cache := New[int](time.Second, time.Second)
	cache.Set("key", 1, time.Second)
	value, ok := cache.Get("key")
	if !ok {
		t.Error("key not found")
	}
	if value != 1 {
		t.Error("value is not 1")
	}
}

func TestCache_Delete(t *testing.T) {
	cache := New[int](time.Second, time.Second)
	cache.Set("key", 1, time.Second)
	cache.Delete("key")
	_, ok := cache.Get("key")
	if ok {
		t.Error("key should be deleted")
	}
}

func TestCache_Cleanup(t *testing.T) {
	cache := New[int](time.Second, time.Second)
	cache.Set("key", 1, time.Second)
	time.Sleep(time.Second)
	_, ok := cache.Get("key")
	if ok {
		t.Error("key should be deleted")
	}
}
