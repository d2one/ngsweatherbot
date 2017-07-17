package main

import (
	"database/sql"
	"log"
	"time"
)

func readCache(cacheKey string) (CachedItem, error) {
	sqlReadOne := `
	SELECT id, cache_key, cache_value, ttl, ttl_lock FROM weather_cache
	WHERE cache_key = ?`
	log.Printf("	SELECT id, cache_key, cache_value, ttl, ttl_lock FROM weather_cache WHERE cache_key = ?" + cacheKey)
	row := db.QueryRow(sqlReadOne, cacheKey)
	item := CachedItem{}
	err := row.Scan(&item.id, &item.cache_key, &item.cache_value, &item.ttl, &item.ttl_lock)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("MOY FOUND")
		return CachedItem{}, nil
	case err != nil:
		log.Printf(err.Error())
		return CachedItem{}, err
	}

	unixTime := time.Now().Unix()
	log.Printf("%d", item.ttl)
	log.Printf("%d", unixTime)

	if item.ttl > unixTime {

		return item, nil
	}

	if item.ttl < unixTime && item.ttl_lock < unixTime+60*10 && item.ttl_lock > 0 {
		log.Printf("CACHE FOUND")
		return item, nil
	}
	lockCacheKey(cacheKey)
	return CachedItem{}, nil
}

func lockCacheKey(cacheKey string) error {
	log.Printf("LOCK CACHE")
	sqlAddItem := `
	INSERT OR REPLACE INTO weather_cache(ttl_lock, cache_key) VALUES (?,?)`

	stmt, err := db.Prepare(sqlAddItem)
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

	stmt, err := db.Prepare(sqlAddItem)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if ttl == 0 {
		ttl = time.Now().Unix() + 60*10
	}

	item := CachedItem{
		cache_value: cacheValue,
		cache_key:   cacheKey,
		ttl:         ttl,
		ttl_lock:    0,
	}
	_, err = stmt.Exec(item.cache_key, item.cache_value, item.ttl, item.ttl_lock)
	return err
}
