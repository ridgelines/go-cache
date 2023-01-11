package cache

import (
	"sort"
	"time"
)

// A Cache is a thread-safe store for fast item storage and retrieval
type Cache[T any] struct {
	itemOps   chan func(map[string]T)
	expiryOps chan func(map[string]*time.Timer)
}

// New returns an empty cache
func New[T any]() *Cache[T] {
	c := &Cache[T]{
		itemOps:   make(chan func(map[string]T)),
		expiryOps: make(chan func(map[string]*time.Timer)),
	}

	go c.loopItemOps()
	go c.loopExpiryOps()
	return c
}

func (c *Cache[T]) loopItemOps() {
	items := map[string]T{}
	for op := range c.itemOps {
		op(items)
	}
}

func (c *Cache[T]) loopExpiryOps() {
	expiries := map[string]*time.Timer{}
	for op := range c.expiryOps {
		op(expiries)
	}
}

// Set will set the val into the cache at the specified key.
// If an entry already exists at the specified key, it will be overwritten.
// The options param can be used to perform logic after the entry has be inserted.
func (c *Cache[T]) Set(key string, val T, options ...SetOption[T]) {
	c.expiryOps <- func(expiries map[string]*time.Timer) {
		if timer, ok := expiries[key]; ok {
			timer.Stop()
			delete(expiries, key)
		}
	}

	c.itemOps <- func(items map[string]T) {
		items[key] = val
	}

	for _, option := range options {
		option(c, key, val)
	}
}

// Clear removes all entries from the cache
func (c *Cache[T]) Clear() {
	c.itemOps <- func(items map[string]T) {
		for key := range items {
			delete(items, key)
		}
	}
}

// ClearEvery clears the cache on a loop at the specified interval
func (c *Cache[T]) ClearEvery(d time.Duration) *time.Ticker {
	ticker := time.NewTicker(d)
	go func() {
		for range ticker.C {
			c.Clear()
		}
	}()

	return ticker
}

// Delete removes an entry from the cache at the specified key.
// If no entry exists at the specified key, no action is taken
func (c *Cache[T]) Delete(key string) {
	c.itemOps <- func(items map[string]T) {
		delete(items, key)
	}
}

// Get retrieves an entry at the specified key
func (c *Cache[T]) Get(key string) T {
	result := make(chan T, 1)
	c.itemOps <- func(items map[string]T) {
		result <- items[key]
	}

	return <-result
}

// GetOK retrieves an entry at the specified key.
// Returns bool specifying if the entry exists
func (c *Cache[T]) GetOK(key string) (T, bool) {
	result := make(chan T, 1)
	exists := make(chan bool, 1)
	c.itemOps <- func(items map[string]T) {
		v, ok := items[key]
		result <- v
		exists <- ok
	}

	return <-result, <-exists
}

// Items retrieves all entries in the cache
func (c *Cache[T]) Items() map[string]T {
	result := make(chan map[string]T, 1)
	c.itemOps <- func(items map[string]T) {
		cp := map[string]T{}
		for key, val := range items {
			cp[key] = val
		}

		result <- cp
	}

	return <-result
}

// Keys retrieves a sorted list of all keys in the cache
func (c *Cache[T]) Keys() []string {
	result := make(chan []string, 1)
	c.itemOps <- func(items map[string]T) {
		keys := make([]string, 0, len(items))
		for k := range items {
			keys = append(keys, k)
		}

		sort.Strings(keys)
		result <- keys
	}

	return <-result
}
