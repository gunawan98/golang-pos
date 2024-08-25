package app

import (
	"database/sql"
	"time"

	"github.com/gunawan98/golang-restfull-api/helper"
)

func NewDB() *sql.DB {

	db, err := sql.Open("mysql", "root:4IjEyA@tcp(localhost:3306)/pos?parseTime=True")
	helper.PanicIfError(err)

	db.SetMaxOpenConns(25)
	db.SetConnMaxIdleTime(5)

	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

// C:\Users\Hidayat\task\golang-pos\app\database.go
// migrate -database "mysql://root:4IjEyA@tcp(localhost:3306)/pos" -path db/migrations up
