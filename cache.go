package cache

import (
	"sort"
	"time"
)

// A Cache is a thread-safe store for fast item storage and retrieval
type Cache struct {
	itemOps   chan func(map[string]interface{})
	expiryOps chan func(map[string]*time.Timer)
}

// New returns an empty cache
func New() *Cache {
	c := &Cache{
		itemOps:   make(chan func(map[string]interface{})),
		expiryOps: make(chan func(map[string]*time.Timer)),
	}

	go c.loopItemOps()
	go c.loopExpiryOps()
	return c
}

func (c *Cache) loopItemOps() {
	items := map[string]interface{}{}
	for op := range c.itemOps {
		op(items)
	}
}

func (c *Cache) loopExpiryOps() {
	expiries := map[string]*time.Timer{}
	for op := range c.expiryOps {
		op(expiries)
	}
}

// Add inserts an entry into the cache at the specified key.
// If an entry already exists at the specified key, it will be overwritten
func (c *Cache) Add(key string, val interface{}) {
	c.itemOps <- func(items map[string]interface{}) {
		items[key] = val
	}
}

// Addf inserts an entry into the cache at the specified key with an expiry.
// If an entry already exists at the specified key, the value and expiry will be overwritten
func (c *Cache) Addf(key string, val interface{}, expiry time.Duration) {
	c.Add(key, val)

	c.expiryOps <- func(expiries map[string]*time.Timer) {
		if timer, ok := expiries[key]; ok {
			timer.Stop()
		}

		expiries[key] = time.AfterFunc(expiry, func() { c.Delete(key) })
	}
}

// Clear removes all entries from the cache
func (c *Cache) Clear() {
	c.itemOps <- func(items map[string]interface{}) {
		for key := range items {
			delete(items, key)
		}
	}
}

// ClearEvery clears the cache on a loop after the specified duration
func (c *Cache) ClearEvery(d time.Duration) *time.Ticker {
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
func (c *Cache) Delete(key string) {
	c.itemOps <- func(items map[string]interface{}) {
		if _, ok := items[key]; ok {
			delete(items, key)
		}
	}
}

// Get retrieves an entry at the specified key
func (c *Cache) Get(key string) interface{} {
	result := make(chan interface{}, 1)
	c.itemOps <- func(items map[string]interface{}) {
		result <- items[key]
	}

	return <-result
}

// Getf retrieves an entry at the specified key.
// Returns bool specifying if the entry exists
func (c *Cache) Getf(key string) (interface{}, bool) {
	result := make(chan interface{}, 1)
	exists := make(chan bool, 1)
	c.itemOps <- func(items map[string]interface{}) {
		v, ok := items[key]
		result <- v
		exists <- ok
	}

	return <-result, <-exists
}

// Items retrieves all entries in the cache
func (c *Cache) Items() map[string]interface{} {
	result := make(chan map[string]interface{}, 1)
	c.itemOps <- func(items map[string]interface{}) {
		cp := map[string]interface{}{}
		for key, val := range items {
			cp[key] = val
		}

		result <- cp
	}

	return <-result
}

// Keys retrieves a sorted list of all keys in the cache
func (c *Cache) Keys() []string {
	result := make(chan []string, 1)
	c.itemOps <- func(items map[string]interface{}) {
		keys := make([]string, 0, len(items))
		for k := range items {
			keys = append(keys, k)
		}

		sort.Strings(keys)
		result <- keys
	}

	return <-result
}
