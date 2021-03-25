package cache

import (
	"fmt"
	"hash/fnv"
	"sync"
)

var hasher = fnv.New64a()

var (
	keyNotFoundError = "Item with key: '%s' not found"
)

// any message above, and corresponding arguments
func errorf(msg string, args ...interface{}) error {
	return fmt.Errorf(msg, args...)
}

// Cache is intended to be a general purpose K-V store.
type Cache struct {
	sync.RWMutex
	cache map[uint64][]byte
}

// Returns a New Cache instance
func New() *Cache {
	c := new(Cache)
	c.cache = make(map[uint64][]byte)
	return c
}

// Returns the data from a given key
func (c *Cache) Get(key string) (data []byte, err error) {
	c.RLock()
	defer c.RUnlock()

	k := getIndexKey(key)
	item, ok := c.cache[k]

	if !ok {
		return nil, errorf(keyNotFoundError, key)
	}

	return item, nil
}

// Sets the data with a given key
func (c *Cache) Set(key string, data []byte) {
	c.Lock()
	defer c.Unlock()

	k := getIndexKey(key)

	c.cache[k] = data
}

// Deletes a data given a key
func (c *Cache) Delete(key string) (err error) {
	c.Lock()
	defer c.Unlock()

	k := getIndexKey(key)

	if _, ok := c.cache[k]; !ok {
		return errorf(keyNotFoundError, key)
	}

	delete(c.cache, k)

	return nil
}

// Clears the cache's and filter's (if available) items
func (c *Cache) Clear() {
	c.Lock()
	c.cache = make(map[uint64][]byte)
	c.Unlock()
}

func getIndexKey(key string) uint64 {
	hasher.Reset()
	hasher.Write([]byte(key))
	return hasher.Sum64()
}
