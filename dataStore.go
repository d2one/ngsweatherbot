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
	CREATE TABLE IF NOT EXISTS user_data(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chat_id INT NOT NULL UNIQUE,
		city_alias TEXT,
		city_title TEXT,
		notifications_next_run INTEGER NULL,
		forecast_type CHAR(64) DEFAULT short,
		created_at INTEGER
	);
	CREATE UNIQUE INDEX IF NOT EXISTS chat_id ON user_data (chat_id);
	`

	_, err = ds.DB.Exec(sqlTable)
	return err
}

func (ds *DataStore) initUser(ChatID int64) error {
	sqlAddItem := `INSERT OR REPLACE INTO user_data(chat_id, created_at) values(?, strftime('%s', 'now'))`

	stmt, err := ds.DB.Prepare(sqlAddItem)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ChatID)
	return err
}

func (ds *DataStore) setForecastType(ChatID int64, forecastType string) error {
	sqlAddItem := `UPDATE user_data SET forecast_type = ? WHERE chat_id = ?`
	log.Println(forecastType)
	stmt, err := ds.DB.Prepare(sqlAddItem)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(forecastType, ChatID)
	return err
}

func (ds *DataStore) saveUserCity(ChatID int64, city *City) error {
	sqlAddItem := `UPDATE user_data SET city_alias = ?, city_title = ? WHERE chat_id = ?`

	stmt, err := ds.DB.Prepare(sqlAddItem)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(city.Alias, city.Title, ChatID)
	return err
}

func (ds *DataStore) saveUserNotification(UserNotification UserNotification) error {
	sqlAddItem := `UPDATE user_data SET notifications_next_run = ? WHERE chat_id = ?`
	stmt, err := ds.DB.Prepare(sqlAddItem)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(UserNotification.NextRun, UserNotification.ChatID)
	return err
}

func (ds *DataStore) deleteUserNotification(chatID int64) error {
	sqlAddItem := `UPDATE user_data SET notifications_next_run=NULL WHERE chat_id = ?`
	stmt, err := ds.DB.Prepare(sqlAddItem)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(chatID)
	return err
}

func (ds *DataStore) getCronUserNotification() ([]*UserNotification, error) {
	sqlRead := `
	SELECT chat_id, notifications_next_run FROM user_data
	WHERE notifications_next_run IS NOT NULL AND notifications_next_run <= strftime('%s', 'now')`

	rows, err := ds.DB.Query(sqlRead)
	if err != nil {
		return nil, err
	}
	var uns []*UserNotification
	for rows.Next() {
		un := new(UserNotification)
		if err := rows.Scan(&un.ChatID, &un.NextRun); err != nil {
			return nil, err
		}
		uns = append(uns, un)
	}
	return uns, nil
}

func (ds *DataStore) getUserNotification(chatID int) (*UserNotification, error) {
	sqlRead := `
	SELECT chat_id, notifications_next_run FROM user_data
	WHERE chat_id = ?`

	row := ds.DB.QueryRow(sqlRead, chatID)
	item := &UserNotification{}
	err := row.Scan(&item.ChatID, &item.NextRun)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return item, nil
	}
}

func (ds *DataStore) getUserCity(UserID int64) (*UserCity, error) {
	sqlReadOne := `
	SELECT chat_id, city_alias, city_title,forecast_type FROM user_data
	WHERE chat_id = ?
	ORDER BY created_at DESC`

	row := ds.DB.QueryRow(sqlReadOne, UserID)
	item := new(UserCity)
	err := row.Scan(&item.ChatID, &item.CityAlias, &item.CityTitle, &item.ForecastType)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return item, nil
	}
}
