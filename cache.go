package main

import (
	"time"

	mCache "github.com/patrickmn/go-cache"
)

// Cache *DB
type Cache struct {
	c *mCache.Cache
}

// Cache Cache
func NewCache() *Cache {
	c := mCache.New(10*time.Minute, 15*time.Minute)
	return &Cache{
		c: c,
	}
}

func (cache *Cache) read(cacheKey string) (*CachedItem, error) {
	if cacheValue, found := cache.c.Get(cacheKey); found {
		cachedItem := &CachedItem{
			CacheKey:   cacheKey,
			CacheValue: cacheValue.(string),
		}
		return cachedItem, nil
	}
	return nil, nil
}

func (cache *Cache) write(cacheKey string, cacheValue string, ttl int64) error {
	cache.c.Set(cacheKey, cacheValue, mCache.DefaultExpiration)
	return nil
}

func (cache *Cache) delete(cacheKey string) error {
	cache.c.Delete(cacheKey)
	return nil
}
