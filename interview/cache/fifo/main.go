package main

import (
	"fmt"
	"sync"
)

type FIFOCache struct {
	mu       sync.Mutex
	capacity int
	order    []string
	data     map[string]int
	index    map[string]int
}

func New(capacity int) *FIFOCache {
	if capacity <= 0 {
		panic("capacity must be greater than 0")
	}
	return &FIFOCache{
		capacity: capacity,
		data:     make(map[string]int),
		index:    make(map[string]int),
		order:    make([]string, 0, capacity),
	}
}

func (c *FIFOCache) Get(key string) (int, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.data[key]
	return v, ok
}

func (c *FIFOCache) Put(key string, value int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.data[key]; ok {
		c.data[key] = value
		return
	}

	if len(c.order) == c.capacity {
		oldest := c.order[0]
		delete(c.data, oldest)
		delete(c.index, oldest)
		c.order = c.order[1:]
		for i, k := range c.order {
			c.index[k] = i
		}
	}

	c.data[key] = value
	c.order = append(c.order, key)
	c.index[key] = len(c.order) - 1
}

func (c *FIFOCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.data)
}

func main() {
	c := New(3)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)
	fmt.Println(c.Get("a")) // 1, true
	c.Put("d", 4)           // evicts "a"
	_, ok := c.Get("a")
	fmt.Println("a exists?", ok) // false
	fmt.Println(c.Get("b"))      // 2, true

	// Overwrite existing doesn't change order
	c.Put("b", 20)
	fmt.Println(c.Get("b")) // 20, true
}
