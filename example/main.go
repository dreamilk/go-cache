package main

import (
	"log"
	"time"

	gocache "github.com/dreamilk/go-cache"
)

func main() {
	cache := gocache.New[int](gocache.NoExpiration, time.Second)
	// set key with value and expiration
	cache.Set("key", 1, time.Second)

	// get key
	value, ok := cache.Get("key")
	if !ok {
		log.Println("key not found")
	}
	log.Println("value:", value)
}
