package gocache

import (
	"sync"
	"time"
)

const (
	NoExpiration      time.Duration = -1
	DefaultExpiration time.Duration = 0
)

type Item[T any] struct {
	Object     T
	Expiration int64
}

type Cache[T any] struct {
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]Item[T]
	mu                sync.RWMutex
	stopCleanup       chan struct{}

	onEvicted func(key string, value T)
}

func New[T any](defaultExpiration, cleanupInterval time.Duration) *Cache[T] {
	cache := &Cache[T]{
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
		mu:                sync.RWMutex{},
		items:             make(map[string]Item[T]),
		stopCleanup:       make(chan struct{}),
		onEvicted:         nil,
	}
	cache.startCleanup()
	return cache
}

func NewWithEviction[T any](defaultExpiration, cleanupInterval time.Duration, onEvicted func(key string, value T)) *Cache[T] {
	cache := New[T](defaultExpiration, cleanupInterval)
	cache.onEvicted = onEvicted
	return cache
}

func (c *Cache[T]) Set(key string, value T, duration time.Duration) {
	if duration == DefaultExpiration {
		duration = c.defaultExpiration
	}

	var expiration int64
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = Item[T]{
		Object:     value,
		Expiration: expiration,
	}
}

func (c *Cache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found || (item.Expiration > 0 && time.Now().UnixNano() > item.Expiration) {
		var zeroValue T
		return zeroValue, false
	}

	return item.Object, true
}

func (c *Cache[T]) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *Cache[T]) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

func (c *Cache[T]) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.items))
	for key, item := range c.items {
		if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
			continue
		}
		keys = append(keys, key)
	}
	return keys
}

func (c *Cache[T]) Map() map[string]T {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items := make(map[string]T, len(c.items))
	for key, item := range c.items {
		if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
			continue
		}
		items[key] = item.Object
	}
	return items
}

func (c *Cache[T]) Close() {
	close(c.stopCleanup)
}

func (c *Cache[T]) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]Item[T])
}

func (c *Cache[T]) startCleanup() {
	go c.cleanup()
}

func (c *Cache[T]) cleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopCleanup:
			return
		case <-ticker.C:
			c.deleteExpired()
		}
	}
}

func (c *Cache[T]) deleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, item := range c.items {
		if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
			if c.onEvicted != nil {
				c.onEvicted(key, item.Object)
			}
			delete(c.items, key)
		}
	}
}
