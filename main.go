package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

type service struct {
	cache *lru.Cache
}

type resizeRequest struct {
	URLs   []string `json:"urls"`
	Width  uint     `json:"width"`
	Height uint     `json:"height"`
}

type resizeResult struct {
	Result string `json:"result"`
	URL    string `json:"url,omitempty"`
	Cached bool   `json:"cached"`
}

const (
	proto    = "http://"
	hostport = "localhost:8080"
	success  = "success"
	failure  = "failure"
)

// Periodically checks all the entries in the cache and removes those with a lifetime greater than 5 minutes.
func filterCacheJob(cache *lru.Cache) {
	for  {
		select {
		case <- time.After(2 * time.Second):

			//fmt.Printf("The size of the cache is (%v)\n", len(cache.Keys()))
			//fmt.Printf("The key 1 (%v)", cache.Keys()[0])

			keys := cache.Keys()

			for i := 0; i < len(keys); i++ {
				key := keys[i].(string)

				fmt.Printf("Processing key for purge (%v)\n", key)

				cachVal, _ := cache.Get(key)

				if cachVal != nil {
					if (time.Now()).Sub(cachVal.(cacheEntry).createdAt).Seconds() > 10 {
						fmt.Println("Removing cache entry (%v)", string(key))
						cache.Remove(string(key))
					}
				}
			}
		}
	}
}

func main() {
	cache, err := lru.New(1024)
	if err != nil {
		log.Panicf("Faild to create cache: %v", err)
	}

	go filterCacheJob(cache)

	svc := &service{
		cache: cache,
	}

	mux := http.NewServeMux()
	mux.Handle("/v1/resize", svc.resizeHandler())
	mux.Handle("/v1/image/", svc.getImageHandler())
	address := hostport

	log.Print("Listening on ", hostport)
	// When running on docker mac, can't listen only on localhost
	panic(http.ListenAndServe(address, mux))
}
