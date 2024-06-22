package mysql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	))
	if err != nil {
		return db, fmt.Errorf("mysql.Connect(): %w", err)
	}

	db.SetMaxOpenConns(150)
	db.SetMaxIdleConns(150)
	db.SetConnMaxLifetime(time.Second * 5)
	db.SetConnMaxIdleTime(time.Second * 10)

	return db, nil
}

func MustConnect() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	))
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(150)
	db.SetMaxIdleConns(150)
	db.SetConnMaxLifetime(time.Second * 5)
	db.SetConnMaxIdleTime(time.Second * 10)

	return db
}

func Truncate(db *sql.DB, table string) error {
	_, err := db.Exec(`DELETE FROM ` + table)
	if err != nil {
		return fmt.Errorf("mysql.Truncate(): %w", err)
	}

	_, err = db.Exec(fmt.Sprintf(`ALTER TABLE %s AUTO_INCREMENT = 0`, table))
	if err != nil {
		return fmt.Errorf("mysql.Truncate(): %w", err)
	}

	return nil
}

func TruncateMany(db *sql.DB, tables []string) error {
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("mysql.TruncateMany(): %w", err)
	}

	for _, table := range tables {
		_, err := tx.Exec(`DELETE FROM ` + table)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return fmt.Errorf("mysql.TruncateMany(): %w", err)
			}
			return fmt.Errorf("mysql.TruncateMany(): %w", err)
		}

		_, err = tx.Exec(fmt.Sprintf(`ALTER TABLE %s AUTO_INCREMENT = 0`, table))
		if err != nil {

			if err := tx.Rollback(); err != nil {
				return fmt.Errorf("mysql.TruncateMany(): %w", err)
			}
			return fmt.Errorf("mysql.TruncateMany(): %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("mysql.TruncateMany(): %w", err)
	}

	return nil
}
