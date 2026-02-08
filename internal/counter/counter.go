package counter

import "sync"

// Counter stores a single integer value safely for concurrent access.
type Counter struct {
	mu    sync.RWMutex
	value int
}

func New() *Counter {
	return &Counter{}
}

func (c *Counter) Get() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

func (c *Counter) Increment() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
	return c.value
}

func (c *Counter) Decrement() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value--
	return c.value
}

func (c *Counter) Reset() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = 0
	return c.value
}
