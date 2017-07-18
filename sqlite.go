package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// db *DB
type DB struct {
	DB *sql.DB
}

func (db *DB) init() error {
	db.DB, err = sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		return err
	}
	if err := db.initTables(); err != nil {
		return err
	}
	return nil
}

func (db *DB) initTables() error {
	// create table if not exists
	sqlTable := `
	CREATE TABLE IF NOT EXISTS user_city(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INT NOT NULL UNIQUE,
		chat_id INT NOT NULL UNIQUE,
		city_alias TEXT,
		created_at INTEGER
	);
	
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
	CREATE TABLE IF NOT EXISTS user_notifications(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INT NOT NULL UNIQUE,
		chat_id INT NOT NULL UNIQUE,
		next_run INTEGER NOT NULL,
		created_at INTEGER
	);
	`

	if _, err := db.DB.Exec(sqlTable); err != nil {
		return err
	}
	return nil
}

func (db *DB) saveUserCity(UserID int64, CityAlias string) error {
	sqlAddItem := `
	INSERT OR REPLACE INTO user_city(
		user_id,
		chat_id,
		city_alias,
		created_at
	) values( ?, ?, ?, strftime('%s', 'now'))`

	stmt, err := db.DB.Prepare(sqlAddItem)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(UserID, UserID, CityAlias)
	return err
}

func (db *DB) saveUserNotification(UserNotification UserNotification) error {
	sqlAddItem := `
	INSERT OR REPLACE INTO user_notifications(
		user_id,
		chat_id,
		next_run,
		created_at
	) values( ?, ?, ?, strftime('%s', 'now'))
	`
	stmt, err := db.DB.Prepare(sqlAddItem)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(UserNotification.UserID, UserNotification.ChatID, UserNotification.NextRun)
	return err
}

func (db *DB) getCronUserNotification() ([]*UserNotification, error) {
	sqlRead := `
	SELECT id, user_id, chat_id, next_run FROM user_notifications
	WHERE next_run <= strftime('%s', 'now')`

	rows, err := db.DB.Query(sqlRead)
	if err != nil {
		return nil, err
	}
	var uns []*UserNotification
	for rows.Next() {
		un := new(UserNotification)
		if err := rows.Scan(&un.ID, &un.UserID, &un.ChatID, &un.NextRun); err != nil {
			return nil, err
		}
		uns = append(uns, un)
	}
	return uns, nil
}

func (db *DB) getUserCity(UserID int64) (*UserCity, error) {
	sqlReadOne := `
	SELECT id, user_id, chat_id, city_alias FROM user_city
	WHERE user_id = ?
	ORDER BY created_at DESC`

	row := db.DB.QueryRow(sqlReadOne, UserID)
	item := new(UserCity)
	err := row.Scan(&item.ID, &item.UserID, &item.ChatID, &item.CityAlias)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return item, nil
	}
}
