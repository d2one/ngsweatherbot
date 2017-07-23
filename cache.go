package main

import (
	"database/sql"
	"log"
	"time"
)

// Cache *DB
type Cache struct {
	DB *sql.DB
}

// Cache Cache
func NewCache() *Cache {
	db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Println(err)
		return nil
	}

	cache := &Cache{DB: db}
	if err := cache.initTables(); err != nil {
		log.Println(err)
		return nil
	}
	return cache
}

func (cache *Cache) initTables() error {
	// create table if not exists
	sqlTable := `
	CREATE TABLE IF NOT EXISTS weather_cache
	(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		cache_key VARCHAR NOT NULL,
		cache_value TEXT NOT NULL,
		ttl INTEGER NOT NULL,
		ttl_lock INTEGER
	);
	
	CREATE UNIQUE INDEX IF NOT EXISTS weather_cache_id_uindex ON weather_cache (id);
	CREATE UNIQUE INDEX IF NOT EXISTS weather_cache_cache_key_uindex ON weather_cache (cache_key);
	`
	_, err = cache.DB.Exec(sqlTable)
	return err
}

func (cache *Cache) read(cacheKey string) (*CachedItem, error) {
	sqlReadOne := `
	SELECT id, cache_key, cache_value, ttl, ttl_lock FROM weather_cache
	WHERE cache_key = ?`

	row := cache.DB.QueryRow(sqlReadOne, cacheKey)
	item := &CachedItem{}
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
	cache.lockKey(cacheKey)
	return nil, nil
}

func (cache *Cache) lockKey(cacheKey string) error {
	sqlAddItem := `
	INSERT OR REPLACE INTO weather_cache(ttl_lock, cache_key) VALUES (?,?)`

	stmt, err := cache.DB.Prepare(sqlAddItem)
	if err != nil {
		return err
	}
	defer stmt.Close()
	ttl := time.Now().Unix() + 60*10
	_, err = stmt.Exec(ttl, cacheKey)
	return err
}

func (cache *Cache) write(cacheKey string, cacheValue string, ttl int64) error {

	sqlAddItem := `
	INSERT OR REPLACE INTO weather_cache(
		cache_key, cache_value, ttl, ttl_lock
	) VALUES (?, ?, ?, ?)`

	stmt, err := cache.DB.Prepare(sqlAddItem)
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

func (cache *Cache) delete(cacheKey string) error {

	sqlAddItem := `
	DELETE FROM weather_cache WHERE cache_key = ?`

	stmt, err := cache.DB.Prepare(sqlAddItem)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(cacheKey)
	return err
}
