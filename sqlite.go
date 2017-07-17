package main

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

func initAppDb() (*sql.DB, error) {
	db, err := initDb("db.sqlite3")
	if err != nil {
		return nil, err
	}
	if err := initTables(db); err != nil {
		return nil, err
	}
	return db, nil
}

func initDb(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}
	if db == nil {
		return db, errors.New("db is nil")
	}
	return db, nil
}

func initTables(db *sql.DB) error {
	// create table if not exists
	sql_table := `
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

	if _, err := db.Exec(sql_table); err != nil {
		return err
	}
	return nil
}

func saveUserCity(item UserCity) error {
	sqlAddItem := `
	INSERT OR REPLACE INTO user_city(
		user_id,
		chat_id,
		city_alias,
		created_at
	) values( ?, ?, ?, strftime('%s', 'now'))
	`

	stmt, err := db.Prepare(sqlAddItem)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(item.user_id, item.chat_id, item.city_alias)
	return err
}

func saveUserNotification(item UserNotification) error {
	sqlAddItem := `
	INSERT OR REPLACE INTO user_notifications(
		user_id,
		chat_id,
		next_run,
		created_at
	) values( ?, ?, ?, strftime('%s', 'now'))
	`
	stmt, err := db.Prepare(sqlAddItem)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(item.user_id, item.chat_id, item.next_run)
	return err
}

func getCronUserNotification() ([]UserNotification, error) {
	sqlRead := `
	SELECT id, user_id, chat_id, next_run FROM user_notifications
	WHERE next_run <= strftime('%s', 'now')`

	rows, err := db.Query(sqlRead)
	if err != nil {
		return []UserNotification{}, err
	}
	uns := make([]UserNotification, 0)
	for rows.Next() {
		un := UserNotification{}
		err = rows.Scan(&un.id, &un.user_id, &un.chat_id, &un.next_run)
		if err != nil {
			return []UserNotification{}, err
		}
		uns = append(uns, un)
	}
	return uns, nil
}

func getUserCity(user_id int64) (UserCity, error) {
	sqlReadOne := `
	SELECT id, user_id, chat_id, city_alias FROM user_city
	WHERE user_id = ?
	ORDER BY created_at DESC`

	row := db.QueryRow(sqlReadOne, user_id)
	item := UserCity{}
	err := row.Scan(&item.id, &item.user_id, &item.chat_id, &item.city_alias)
	switch {
	case err == sql.ErrNoRows:
		return UserCity{}, nil
	case err != nil:
		return UserCity{}, err
	default:
		return item, nil
	}
}
