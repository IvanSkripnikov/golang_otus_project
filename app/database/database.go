package database

import (
	"database/sql"
	"fmt"
	"os"

	"app/logger"

	_ "github.com/go-sql-driver/mysql"
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

	host := "db"
	host = "localhost"
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

func GetBannerEvents(bannerId, groupId, slotId int, eventType string) int {
	query := "SELECT COUNT(*) as cnt from events WHERE banner_id = ? AND group_id = ? AND slot_id = ? AND type = ?"
	stmt, err := DB.Query(query, bannerId, groupId, slotId, eventType)
	if err != nil {
		return 0
	}

	defer stmt.Close()

	count := 0

	for stmt.Next() {
		if err := stmt.Scan(&count); err != nil {
			return 0
		}
	}

	return count
}

func GetBannersForSlot(slotId int) ([]int, error) {
	query := "SELECT banner_id from relations_banner_slot WHERE slot_id = ?"
	rows, err := DB.Query(query, slotId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	banners := make([]int, 0)
	banner := 0
	for rows.Next() {
		if err = rows.Scan(&banner); err != nil {
			logger.SendToErrorLog(err.Error())
			continue
		}
		banners = append(banners, banner)
	}

	return banners, nil
}

func GetAllEvents(eventType string) int {
	query := "SELECT COUNT(*) from events WHERE type = ?"
	stmt, err := DB.Query(query, eventType)
	if err != nil {
		return 0
	}

	defer stmt.Close()

	count := 0

	for stmt.Next() {
		if err := stmt.Scan(&count); err != nil {
			return 0
		}
	}

	return count
}
