package cache

import (
	"sync"
	"time"
)

type Map struct {
	mu     sync.Mutex
	values map[string]*CachedValue
}

type CachedValue struct {
	lastUsed time.Time
	value    interface{}
}

func NewCacheMap() *Map {
	return &Map{
		values: make(map[string]*CachedValue),
	}
}

func (c *Map) GCCollector() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, value := range c.values {
		if value.lastUsed.Before(time.Now().Add(-1 * time.Hour)) {
			delete(c.values, key)
		}
	}
}

func (c *Map) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.values[key] = &CachedValue{
		lastUsed: time.Now(),
		value:    value,
	}
}

func (c *Map) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	value, ok := c.values[key]
	if !ok {
		return nil, false
	}
	value.lastUsed = time.Now()

	return value.value, true
}
