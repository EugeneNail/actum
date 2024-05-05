package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func Connect() (*sql.DB, error) {
	return sql.Open("mysql", GetDsn())
}

func GetDsn() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
}

func Truncate(table string) error {
	db, err := Connect()
	defer db.Close()

	if err != nil {
		return fmt.Errorf("can't connect: %w", err)
	}

	if _, err := db.Exec("TRUNCATE TABLE " + table); err != nil {
		return fmt.Errorf("can't truncate table: %w", err)
	}

	return nil
}
