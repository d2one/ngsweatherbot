package main

import (
	"database/sql"

	"log"

	_ "github.com/mattn/go-sqlite3"
)

// db *DB
type DataStore struct {
	DB *sql.DB
}

// dasdas
func NewDataStore() *DataStore {
	db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Println(err)
		return nil
	}

	dataStore := &DataStore{DB: db}
	if err := dataStore.initTables(); err != nil {
		log.Println(err)
		return nil
	}
	return dataStore
}

func (ds *DataStore) initTables() error {
	// create table if not exists
	sqlTable := `
	CREATE TABLE IF NOT EXISTS user_city(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INT NOT NULL UNIQUE,
		chat_id INT NOT NULL UNIQUE,
		city_alias TEXT,
		city_title TEXT,
		created_at INTEGER
	);
	
	CREATE TABLE IF NOT EXISTS user_notifications(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INT NOT NULL UNIQUE,
		chat_id INT NOT NULL UNIQUE,
		next_run INTEGER NOT NULL,
		created_at INTEGER
	);
	`

	_, err = ds.DB.Exec(sqlTable)
	return err
}

func (ds *DataStore) saveUserCity(UserID int64, city *City) error {
	sqlAddItem := `
	INSERT OR REPLACE INTO user_city(
		user_id,
		chat_id,
		city_alias,
		city_title,
		created_at
	) values( ?, ?, ?, ?, strftime('%s', 'now'))`

	stmt, err := ds.DB.Prepare(sqlAddItem)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(UserID, UserID, city.Alias, city.Title)
	return err
}

func (ds *DataStore) saveUserNotification(UserNotification UserNotification) error {
	sqlAddItem := `
	INSERT OR REPLACE INTO user_notifications(
		user_id,
		chat_id,
		next_run,
		created_at
	) values( ?, ?, ?, strftime('%s', 'now'))
	`
	stmt, err := ds.DB.Prepare(sqlAddItem)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(UserNotification.UserID, UserNotification.ChatID, UserNotification.NextRun)
	return err
}

func (ds *DataStore) getCronUserNotification() ([]*UserNotification, error) {
	sqlRead := `
	SELECT id, user_id, chat_id, next_run FROM user_notifications
	WHERE next_run <= strftime('%s', 'now')`

	rows, err := ds.DB.Query(sqlRead)
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

func (ds *DataStore) getUserCity(UserID int64) (*UserCity, error) {
	sqlReadOne := `
	SELECT id, user_id, chat_id, city_alias, city_title FROM user_city
	WHERE user_id = ?
	ORDER BY created_at DESC`

	row := ds.DB.QueryRow(sqlReadOne, UserID)
	item := new(UserCity)
	err := row.Scan(&item.ID, &item.UserID, &item.ChatID, &item.CityAlias, &item.CityTitle)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return item, nil
	}
}
