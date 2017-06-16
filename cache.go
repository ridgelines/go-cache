package cache

import (
	"sort"
	"time"
)

// CachedValue is a type for cache value
type CachedValue interface{}

// CachedItems is a map for CachedValue
type CachedItems map[string]CachedValue

// CachedExpiries is a map for Timer
type CachedExpiries map[string]*time.Timer

// A Cache is a thread-safe store for fast item storage and retrieval
type Cache struct {
	itemOps   chan func(CachedItems)
	expiryOps chan func(CachedExpiries)
}

// New returns an empty cache
func New() *Cache {
	c := &Cache{
		itemOps:   make(chan func(CachedItems)),
		expiryOps: make(chan func(CachedExpiries)),
	}

	go c.loopItemOps()
	go c.loopExpiryOps()
	return c
}

func (c *Cache) loopItemOps() {
	items := CachedItems{}
	for op := range c.itemOps {
		op(items)
	}
}

func (c *Cache) loopExpiryOps() {
	expiries := CachedExpiries{}
	for op := range c.expiryOps {
		op(expiries)
	}
}

// Add inserts an entry into the cache at the specified key.
// If an entry already exists at the specified key, it will be overwritten
func (c *Cache) Add(key string, val CachedValue) {
	c.itemOps <- func(items CachedItems) {
		items[key] = val
	}
}

// Addf inserts an entry into the cache at the specified key with an expiry.
// If an entry already exists at the specified key, the value and expiry will be overwritten
func (c *Cache) Addf(key string, val CachedValue, expiry time.Duration) {
	c.Add(key, val)

	c.expiryOps <- func(expiries CachedExpiries) {
		if timer, ok := expiries[key]; ok {
			timer.Stop()
		}

		expiries[key] = time.AfterFunc(expiry, func() { c.Delete(key) })
	}
}

// Clear removes all entries from the cache
func (c *Cache) Clear() {
	c.itemOps <- func(items CachedItems) {
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
	c.itemOps <- func(items CachedItems) {
		if _, ok := items[key]; ok {
			delete(items, key)
		}
	}
}

// Get retrieves an entry at the specified key
func (c *Cache) Get(key string) CachedValue {
	result := make(chan CachedValue, 1)
	c.itemOps <- func(items CachedItems) {
		result <- items[key]
	}

	return <-result
}

// Getf retrieves an entry at the specified key.
// Returns bool specifying if the entry exists
func (c *Cache) Getf(key string) (CachedValue, bool) {
	result := make(chan CachedValue, 1)
	exists := make(chan bool, 1)
	c.itemOps <- func(items CachedItems) {
		v, ok := items[key]
		result <- v
		exists <- ok
	}

	return <-result, <-exists
}

// Items retrieves all entries in the cache
func (c *Cache) Items() CachedItems {
	result := make(chan CachedItems, 1)
	c.itemOps <- func(items CachedItems) {
		result <- items
	}

	return <-result
}

// Keys retrieves a sorted list of all keys in the cache
func (c *Cache) Keys() []string {
	result := make(chan []string, 1)
	c.itemOps <- func(items CachedItems) {
		keys := make([]string, 0, len(items))
		for k := range items {
			keys = append(keys, k)
		}

		sort.Strings(keys)
		result <- keys
	}

	return <-result
}
