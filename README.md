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
Personally, I'd recommend copying this package and replacing `interface{}` with whatever type you need to cache. 
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
	cache := cache.New()
	
	// empty the cache every hour
	cache.ClearEvery(time.Hour)
	
	// add some items
	cache.Add("one", 1)
	cache.Add("two", 2)
	
	// add some items that will expire after 5 minutes
	cache.Addf("three", 3, time.Minute*5)
	cache.Addf("four", 4, time.Minute*5)

	fmt.Println(cache.Get("1"))
	fmt.Println(cache.Get("2"))
	
	for _, key := range cache.Keys() {
		fmt.Println(key)
	}
	
	for key, val := range cache.Items() {
		fmt.Printf("%s: %v", key, val)
	}
	
	cache.Delete("one")
	
	if val, ok := cache.Getf("two"); ok {
		fmt.Println(val)
	}
	
	cache.Clear()
}
```

## License
This work is published under the MIT license.
Please see the `LICENSE` file for details.
