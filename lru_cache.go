package gocache

import "sync"

type Node[T any] struct {
	key   string
	value T
	prev  *Node[T]
	next  *Node[T]
}

type LRUCache[T any] struct {
	mu        sync.Mutex
	capacity  int
	cache     map[string]*Node[T]
	head      *Node[T]
	tail      *Node[T]
	onEvicted func(key string, value T)
}

func NewLRUCache[T any](capacity int) *LRUCache[T] {
	if capacity <= 0 {
		panic("capacity must be greater than 0")
	}

	head := &Node[T]{}
	tail := &Node[T]{}

	head.next = tail
	tail.prev = head

	return &LRUCache[T]{
		capacity:  capacity,
		cache:     make(map[string]*Node[T], capacity),
		head:      head,
		tail:      tail,
		mu:        sync.Mutex{},
		onEvicted: nil,
	}
}

func NewLRUCacheWithEviction[T any](capacity int, onEvicted func(key string, value T)) *LRUCache[T] {
	cache := NewLRUCache[T](capacity)
	cache.onEvicted = onEvicted
	return cache
}

func (c *LRUCache[T]) moveToFront(node *Node[T]) {
	if node == c.head {
		return
	}

	node.prev.next = node.next
	node.next.prev = node.prev

	node.next = c.head.next
	node.prev = c.head

	c.head.next.prev = node
	c.head.next = node
}

func (c *LRUCache[T]) Get(key string) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, ok := c.cache[key]
	if !ok {
		var zero T
		return zero, false
	}

	c.moveToFront(node)
	return node.value, true
}

func (c *LRUCache[T]) Set(key string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, ok := c.cache[key]; ok {
		node.value = value
		c.moveToFront(node)
		return
	}

	if len(c.cache) >= c.capacity {
		if c.onEvicted != nil {
			c.onEvicted(c.tail.prev.key, c.tail.prev.value)
		}

		delete(c.cache, c.tail.prev.key)
		c.tail.prev = c.tail.prev.prev
		c.tail.prev.next = c.tail
	}

	node := &Node[T]{
		key:   key,
		value: value,
		prev:  c.head,
		next:  c.head.next,
	}

	c.cache[key] = node

	c.head.next.prev = node
	c.head.next = node
}

func (c *LRUCache[T]) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, ok := c.cache[key]
	if !ok {
		return
	}

	if c.onEvicted != nil {
		c.onEvicted(key, node.value)
	}

	delete(c.cache, key)

	node.prev.next = node.next
	node.next.prev = node.prev
}

func (c *LRUCache[T]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*Node[T], c.capacity)
	c.head.next = c.tail
	c.tail.prev = c.head
}

func (c *LRUCache[T]) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.cache)
}

func (c *LRUCache[T]) Keys() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	keys := make([]string, 0, len(c.cache))
	for key := range c.cache {
		keys = append(keys, key)
	}
	return keys
}
