package database

import (
	"database/sql"
	"fmt"

	"github.com/IvanSkripnikov/golang_otus_project/initiate"
	"github.com/IvanSkripnikov/golang_otus_project/logger"
	_ "github.com/go-sql-driver/mysql" // nolint:nolintlint
)

var DB *sql.DB

func init() {
	DB = InitDataBase()
}

func InitDataBase() *sql.DB {
	logger.SendToInfoLog("connecting ...")

	// get environment variables

	dsn := initiate.GetDatabaseConnectionString()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.SendToFatalLog(fmt.Sprintf("DB connection has been failed. %s", err.Error()))
	}

	logger.SendToInfoLog("connected!!")

	return db
}
