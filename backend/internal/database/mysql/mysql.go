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
		return fmt.Errorf("mysql.SafeTrucate(): %w", err)
	}

	_, err = db.Exec(`DELETE FROM ` + table)
	if err != nil {
		return fmt.Errorf("mysql.SafeTrucate(): %w", err)
	}

	_, err = db.Exec(fmt.Sprintf(`ALTER TABLE %s AUTO_INCREMENT = 0`, table))
	if err != nil {
		return fmt.Errorf("mysql.SafeTrucate(): %w", err)
	}

	return nil
}
