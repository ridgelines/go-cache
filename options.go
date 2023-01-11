package cache

import "time"

// A SetOption will perform logic after a set action completes
type SetOption[T any] func(c *Cache[T], key string, val T)

// Expire is a SetOption that will cause the entry to expire after the specified duration
func Expire[T any](expiry time.Duration) SetOption[T] {
	return func(c *Cache[T], key string, val T) {
		c.expiryOps <- func(expiries map[string]*time.Timer) {
			if timer, ok := expiries[key]; ok {
				timer.Stop()
			}

			expiries[key] = time.AfterFunc(expiry, func() { c.Delete(key) })
		}
	}
}
