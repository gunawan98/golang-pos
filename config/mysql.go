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

	dsn := fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?parseTime=True`, dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.DBName)
	// fmt.Println(dsn)

	db, err := sql.Open("mysql", dsn)
	helper.PanicIfError(err)

	db.SetConnMaxIdleTime(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}
