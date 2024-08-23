package app

import (
	"database/sql"
	"time"

	"github.com/gunawan98/golang-restfull-api/helper"
)

func NewDB() *sql.DB {
	db, err := sql.Open("mysql", "root:gunawan98@tcp(localhost:3306)/retail_pos?parseTime=true")
	helper.PanicIfError(err)

	db.SetConnMaxIdleTime(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

// migrate -database "mysql://userDB:passwordDB@tcp(localhost:3306)/retail_pos" -path db/migrations up
