# Go Cache

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/zpatrick/go-cache/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/zpatrick/go-cache)](https://goreportcard.com/report/github.com/zpatrick/go-cache)
[![Go Doc](https://godoc.org/github.com/zpatrick/go-cache?status.svg)](https://godoc.org/github.com/zpatrick/go-cache)


## Overview
Go Cache is a simple package to provide thread-safe in-memory caching in Go. 
It's my attempt to practice some of the patterns/philosophies found in these articles:

* [Do not fear first class functions](https://dave.cheney.net/2016/11/13/do-not-fear-first-class-functions)
* [Share Memory By Communicating](https://blog.golang.org/share-memory-by-communicating)

The code is [tested](https://github.com/zpatrick/go-cache/blob/master/cache_test.go), although standard caveats of using `interface{}` apply.  
Personally, I'd recommend copying this package and replacing `var T interface{}` with whatever type you need to cache. 
I may add code generation in the future to make that process easier. 

## Example
```
package main

import (
        "fmt"
	"time"
        "github.com/zpatrick/go-cache"
)

func main() {
	c := cache.New()
	
	// empty the cache every hour
	c.ClearEvery(time.Hour)
	
	// add some items
	c.Set("key1", 1)
	c.Set("key2", 2)
	
	// add some items that will expire after 5 minutes
	c.Set("key3", 3, cache.Expire(time.Minute*5))
	c.Set("key4", 4,  cache.Expire(time.Minute*5))

	fmt.Println(c.Get("key1"))
	fmt.Println(c.Get("key2"))
	
	for _, key := range c.Keys() {
		fmt.Println(key)
	}
	
	for key, val := range c.Items() {
		fmt.Printf("%s: %v", key, val)
	}
	
	c.Delete("key1")
	
	if val, ok := c.GetOK("key2"); ok {
		fmt.Println(val)
	}
	
	c.Clear()
}
```

## License
This work is published under the MIT license.
Please see the `LICENSE` file for details.
