package database

import "database/sql"

var db *sql.DB

func GetDB() *sql.DB {
	return db
}

func SetDB(d *sql.DB) {
	db = d
}
