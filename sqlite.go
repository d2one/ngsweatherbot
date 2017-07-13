package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"errors"
)

func initAppDb() (*sql.DB, error)  {
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
		created_at DATETIME
	);
	`

	if _, err := db.Exec(sql_table); err != nil {
		return err
	}
	return nil
}

func saveUserCity(db *sql.DB, item UserCity) error {
	sqlAddItem := `
	INSERT OR REPLACE INTO user_city(
		user_id,
		chat_id,
		city_alias,
		created_at
	) values( ?, ?, ?, CURRENT_TIMESTAMP)
	`

	stmt, err := db.Prepare(sqlAddItem)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err2 := stmt.Exec(item.user_id, item.chat_id, item.city_alias)
	return err2
}

func getUserCity(db *sql.DB, user_id int64) (UserCity, error) {
	sqlReadOne := `
	SELECT id, user_id, chat_id, city_alias FROM user_city
	WHERE user_id = ?
	ORDER BY datetime(created_at) DESC`

	row := db.QueryRow(sqlReadOne, user_id)
	item := UserCity{}
	err := row.Scan(&item.id, &item.user_id, &item.chat_id, &item.city_alias);
	switch {
		case err == sql.ErrNoRows:
			return UserCity{}, nil
		case err != nil:
			return UserCity{}, err
		default:
			return item, nil
	}

}