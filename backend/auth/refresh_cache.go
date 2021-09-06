package auth

import (
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/filebrowser/filebrowser/v3/log"
)

// InMemoryAuthRefreshCache used by authenticator to minimize repeatable token refreshes
type InMemoryAuthRefreshCache struct {
	cache *ristretto.Cache
}

func NewInMemoryAuthRefreshCache() *InMemoryAuthRefreshCache {
	memCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 10000,
		MaxCost:     1000,
		BufferItems: 64,
	})

	if err != nil {
		log.Fatalf("Failed to init cache: %s", err)
	}

	return &InMemoryAuthRefreshCache{cache: memCache}
}

// Get implements cache getter with key converted to string
func (c *InMemoryAuthRefreshCache) Get(key interface{}) (interface{}, bool) {
	return c.cache.Get(key)
}

// Set implements cache setter with key converted to string
func (c *InMemoryAuthRefreshCache) Set(key, value interface{}) {
	c.cache.SetWithTTL(key, value, 1, 5*time.Second)
}
