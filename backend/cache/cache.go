package cache

import "time"

type Cache interface {
	// Get returns the value (if any) and a boolean representing whether the
	// value was found or not. The value can be nil and the boolean can be true at
	// the same time.
	Get(key interface{}) (interface{}, bool)
	// Set attempts to add the key-value item to the cache. If it returns false,
	// then the Set was dropped and the key-value item isn't added to the cache. If
	// it returns true, there's still a chance it could be dropped by the policy if
	// its determined that the key-value item isn't worth keeping, but otherwise the
	// item will be added and other items will be evicted in order to make room.
	Set(key, value interface{}, cost int64) bool
	// SetWithTTL works like Set but adds a key-value pair to the cache that will expire
	// after the specified TTL (time to live) has passed. A zero value means the value never
	// expires, which is identical to calling Set. A negative value is a no-op and the value
	// is discarded.
	SetWithTTL(key, value interface{}, cost int64, ttl time.Duration) bool
	// Del deletes the key-value item from the cache if it exists.
	Del(key interface{})
	// Clear empties the cache
	Clear()
}
