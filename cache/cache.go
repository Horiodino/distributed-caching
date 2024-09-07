package cache

import (
	"fmt"
	"sync"
)

type Cache struct {
	lock sync.RWMutex
	data map[string][]byte
}

func NewCache() Cache {
	return Cache{
		data: make(map[string][]byte),
	}
}

func (c *Cache) Delete(key []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.data, string(key))

	return nil
}

func (c *Cache) Has(key []byte) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	_, ok := c.data[string(key)]
	return ok
}

func (c *Cache) Set(key []byte) ([]byte, error) {
	c.lock.RLock()
	c.lock.RUnlock()

	keystr := string(key)

	val, ok := c.data[keystr]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}

	return val, nil
}
