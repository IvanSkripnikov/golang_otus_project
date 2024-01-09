package initiate

import (
	"fmt"
	"os"
)

func GetDatabaseConnectionString() string {
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

	return fmt.Sprintf("%s:%s@%s/%s?timeout=30s", user, pass, netAddr, dbname)
}
