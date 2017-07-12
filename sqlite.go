package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil { panic(err) }
	if db == nil { panic("db nil") }
	return db
}

func CreateTable(db *sql.DB) {
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

	_, err := db.Exec(sql_table)
	if err != nil { panic(err) }
}

func StoreItem(db *sql.DB, items []UserCity) {
	sqlAddItem := `
	INSERT OR REPLACE INTO user_city(
		user_id,
		chat_id,
		city_alias,
		created_at
	) values( ?, ?, ?, CURRENT_TIMESTAMP)
	`

	stmt, err := db.Prepare(sqlAddItem)
	if err != nil { panic(err) }
	defer stmt.Close()

	for _, item := range items {
		_, err2 := stmt.Exec(item.user_id, item.chat_id, item.city_alias)
		if err2 != nil { panic(err2) }
	}
}

func ReadItem(db *sql.DB, user_id int64) UserCity {
	sqlReadOne := `
	SELECT id, user_id, chat_id, city_alias FROM user_city
	WHERE user_id = ?
	ORDER BY datetime(created_at) DESC`

	row := db.QueryRow(sqlReadOne, user_id)
	item := UserCity{}
	err := row.Scan(&item.id, &item.user_id, &item.chat_id, &item.city_alias)
	if err != nil { panic(err) }
	return item
}