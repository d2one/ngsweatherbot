package main

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// Cache *DB
type InmemoryCache struct {
	c *cache.Cache
}

// Cache Cache
func NewCache() *InmemoryCache {
	c := cache.New(10*time.Minute, 15*time.Minute)
	return &InmemoryCache{
		c: c,
	}
}

func (inmemoryCache *InmemoryCache) read(cacheKey string) (*CachedItem, error) {
	if cacheValue, found := inmemoryCache.c.Get(cacheKey); found {
		cachedItem := &CachedItem{
			CacheKey:   cacheKey,
			CacheValue: cacheValue.(string),
		}
		return cachedItem, nil
	}
	return nil, nil
}

func (inmemoryCache *InmemoryCache) write(cacheKey string, cacheValue string, ttl int64) error {
	inmemoryCache.c.Set(cacheKey, cacheValue, cache.DefaultExpiration)
	return nil
}

func (inmemoryCache *InmemoryCache) delete(cacheKey string) error {
	inmemoryCache.c.Delete(cacheKey)
	return nil
}
