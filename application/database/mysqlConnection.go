package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func NewMysqlConnection() *sql.DB {
	db, err := sql.Open("mysql", "cho:master@tcp(192.168.35.154:3306)/heroquest")
	if err != nil {
		panic(err)
	}
	return db
}
