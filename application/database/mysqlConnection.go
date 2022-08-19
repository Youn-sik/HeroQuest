package database

import (
	"database/sql"
)

func NewMysqlConnection() *sql.DB {
	db, err := sql.Open("mysql", "root:master@tcp(192.168.35.154:3306)/heroquest")
	if err != nil {
		panic(err)
	}
	return db
}
