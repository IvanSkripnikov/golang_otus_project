package database

import (
	"database/sql"
	"fmt"
	"os"

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
	env := func(key, defaultValue string) string {
		if value := os.Getenv(key); value != "" {
			return value
		}
		return defaultValue
	}

	// host := "db"
	host := "localhost"
	user := env("MYSQL_USER", "user")
	pass := env("MYSQL_PASSWORD", "pass")
	prot := env("MYSQL_PROT", "tcp")
	addr := env("MYSQL_ADDR", host+":3306")
	dbname := env("MYSQL_DATABASE", "test")
	netAddr := fmt.Sprintf("%s(%s)", prot, addr)
	dsn := fmt.Sprintf("%s:%s@%s/%s?timeout=30s", user, pass, netAddr, dbname)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.SendToFatalLog(fmt.Sprintf("DB connection has been failed. %s", err.Error()))
	}

	logger.SendToInfoLog("connected!!")

	return db
}
