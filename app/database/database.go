package database

import (
	"database/sql"
	"fmt"

	"github.com/IvanSkripnikov/golang_otus_project/config"
	"github.com/IvanSkripnikov/golang_otus_project/logger"
	_ "github.com/go-sql-driver/mysql" // nolint:nolintlint
)

var DB *sql.DB

func init() {
	DB = InitDataBase("db")
}

func InitDataBase(host string) *sql.DB {
	logger.SendToInfoLog("connecting ...")

	// get environment variables
	dsn := config.GetDatabaseConnectionString(host)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.SendToFatalLog(fmt.Sprintf("DB connection has been failed. %s", err.Error()))
	}

	logger.SendToInfoLog("connected!!")

	DB = db

	return db
}
