package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3307)/fileserver?charset=utf8mb4")
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(1000)
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func DBConn() *sql.DB {
	return db
}
