package main

import (
	"database/sql"
	"log"
	"time"
)

func readCache(cacheKey string) (*CachedItem, error) {
	sqlReadOne := `
	SELECT id, cache_key, cache_value, ttl, ttl_lock FROM weather_cache
	WHERE cache_key = ?`

	row := db.DB.QueryRow(sqlReadOne, cacheKey)
	item := new(CachedItem)
	err := row.Scan(&item.ID, &item.CacheKey, &item.CacheValue, &item.TTL, &item.TTLLock)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}

	unixTime := time.Now().Unix()
	if item.TTL > unixTime {
		return item, nil
	}

	if item.TTL < unixTime && item.TTLLock < unixTime+60*10 && item.TTLLock > 0 {
		return item, nil
	}
	lockCacheKey(cacheKey)
	return nil, nil
}

func lockCacheKey(cacheKey string) error {
	log.Printf("LOCK CACHE")
	sqlAddItem := `
	INSERT OR REPLACE INTO weather_cache(ttl_lock, cache_key) VALUES (?,?)`

	stmt, err := db.DB.Prepare(sqlAddItem)
	if err != nil {
		return err
	}
	defer stmt.Close()
	ttl := time.Now().Unix() + 60*10
	_, err = stmt.Exec(ttl, cacheKey)
	return err
}

func saveCache(cacheKey string, cacheValue string, ttl int64) error {

	sqlAddItem := `
	INSERT OR REPLACE INTO weather_cache(
		cache_key, cache_value, ttl, ttl_lock
	) VALUES (?, ?, ?, ?)`

	stmt, err := db.DB.Prepare(sqlAddItem)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if ttl == 0 {
		ttl = time.Now().Unix() + 60*10
	}
	_, err = stmt.Exec(cacheKey, cacheValue, ttl, 0)
	return err
}
