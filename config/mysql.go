package config

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gunawan98/golang-restfull-api/helper"
)

// Database set database vars
type DatabaseMysql struct {
	Host              string
	User              string
	Password          string
	DBName            string
	Port              string
	ReconnectRetry    int
	ReconnectInterval int64
	DebugMode         bool
}

func LoadMySQLConfig() DatabaseMysql {
	dbDebug, _ := strconv.ParseBool(os.Getenv("DB_DEBUG"))
	conf := DatabaseMysql{
		Host:      os.Getenv("DB_HOST"),
		User:      os.Getenv("DB_USER"),
		Password:  os.Getenv("DB_PASS"),
		DBName:    os.Getenv("DB_NAME"),
		Port:      os.Getenv("DB_PORT"),
		DebugMode: dbDebug,
	}

	return conf
}

func MySQLConnect() *sql.DB {

	dbConf := LoadMySQLConfig()
	// // db, err := sql.Open("mysql", "root:4IjEyA@tcp(localhost:3306)/pos?parseTime=true")

	dsn := fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?parseTime=True`, dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.DBName)
	fmt.Print(dsn)

	// db, err := sql.Open("mysql", "root:4IjEyA@tcp(localhost:3306)/pos?parseTime=True")
	db, err := sql.Open("mysql", "root:4IjEyA@tcp(localhost:3306)/pos?parseTime=True")
	helper.PanicIfError(err)

	db.SetMaxOpenConns(25)
	db.SetConnMaxIdleTime(5)

	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}
