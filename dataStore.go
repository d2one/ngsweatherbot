package main

import (
	"database/sql"

	"log"

	_ "github.com/mattn/go-sqlite3"
)

// db *DB
type DataStore struct {
	db *sql.DB
}

// dasdas
func NewDataStore(dbPath string) *DataStore {
	log.Println(dbPath)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return nil
	}

	dataStore := &DataStore{db: db}
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

	_, err = ds.db.Exec(sqlTable)
	return err
}

func (ds *DataStore) initUser(ChatID int64) error {
	sqlAddItem := "INSERT OR REPLACE INTO user_data(chat_id, created_at) values(?, strftime('%s', 'now'))"

	stmt, err := ds.db.Prepare(sqlAddItem)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ChatID)
	return err
}

func (ds *DataStore) getUserData(chatID int64) (*UserData, error) {
	sqlRead := "SELECT  id, chat_id, city_alias, city_title, notifications_next_run, forecast_type, created_at " +
		"FROM user_data WHERE chat_id = ?"

	row := ds.db.QueryRow(sqlRead, chatID)
	userData := &UserData{}
	err := row.Scan(
		&userData.ID,
		&userData.ChatID,
		&userData.CityAlias,
		&userData.CityTitle,
		&userData.NotificationsNextRun,
		&userData.ForecastType,
		&userData.CreatedAt,
	)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}
	return userData, nil
}

func (ds *DataStore) saveUserData(userData *UserData) error {
	sqlAddItem := "INSERT OR REPLACE INTO user_data " +
		"(city_alias, city_title, notifications_next_run, forecast_type, created_at) " +
		"VALUES (?, ?, ?, ?, strftime('%s', 'now')) " +
		"WHERE chat_id = ?"

	stmt, err := ds.db.Prepare(sqlAddItem)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		userData.CityAlias,
		userData.CityTitle,
		userData.NotificationsNextRun,
		userData.ForecastType,
		userData.ChatID,
	)
	return err
}

func (ds *DataStore) saveForecastType(ChatID int64, forecastType string) error {
	sqlAddItem := "UPDATE user_data SET forecast_type = ? WHERE chat_id = ?"
	log.Println(forecastType)
	stmt, err := ds.db.Prepare(sqlAddItem)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(forecastType, ChatID)
	return err
}

func (ds *DataStore) saveUserCity(chatID int64, city *City) error {
	sqlAddItem := "UPDATE user_data SET city_alias = ?, city_title = ? WHERE chat_id = ?"

	stmt, err := ds.db.Prepare(sqlAddItem)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(city.Alias, city.Title, chatID)
	return err
}

func (ds *DataStore) saveUserNotification(chatID int64, nextRun int64) error {
	sqlAddItem := "UPDATE user_data SET notifications_next_run = ? WHERE chat_id = ?"
	stmt, err := ds.db.Prepare(sqlAddItem)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(nextRun, chatID)
	return err
}

func (ds *DataStore) deleteUserNotification(chatID int64) error {
	sqlAddItem := "UPDATE user_data SET notifications_next_run=NULL WHERE chat_id = ?"
	stmt, err := ds.db.Prepare(sqlAddItem)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(chatID)
	return err
}

func (ds *DataStore) getCronUsersNotifications() ([]*UserData, error) {
	sqlRead := "SELECT id, chat_id, city_alias, city_title, notifications_next_run, forecast_type, created_at " +
		"FROM user_data " +
		"WHERE notifications_next_run IS NOT NULL AND notifications_next_run <= strftime('%s', 'now')"
	rows, err := ds.db.Query(sqlRead)
	if err != nil {
		return nil, err
	}
	var usersData []*UserData
	for rows.Next() {
		userData := &UserData{}
		err := rows.Scan(
			&userData.ID,
			&userData.ChatID,
			&userData.CityAlias,
			&userData.CityTitle,
			&userData.NotificationsNextRun,
			&userData.ForecastType,
			&userData.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		usersData = append(usersData, userData)
	}
	return usersData, nil
}
