package cache

import (
	"fmt"
	"log"
	"sync"
	"time"
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

func (c *Cache) Get(key []byte) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	val, ok := c.data[string(key)]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}

	return val, nil
}

func (c *Cache) Set(key []byte, value []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[string(key)] = value
	log.Print("setting key: ", string(key), " with value: ", string(value))
	ticker := time.NewTicker(ttl)
	go func() {
		<-ticker.C
		delete(c.data, string(key))
	}()
	return nil
}
