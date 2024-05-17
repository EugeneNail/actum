package tests

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"testing"
)

type Database struct {
	t *testing.T
}

func NewDatabase(t *testing.T) Database {
	return Database{t}
}

func (database *Database) AssertEmpty(table string) *Database {
	db, err := mysql.Connect()
	defer db.Close()
	Check(err)

	var rows int
	err = db.
		QueryRow(`SELECT COUNT(*) FROM ` + table).
		Scan(&rows)
	Check(err)

	if rows > 0 {
		database.t.Errorf(`The table "%s" is expected to be empty, got %d rows instead`, table, rows)
	}

	return database
}

func (database *Database) AssertHas(table string, entity map[string]any) *Database {
	if getMappedCount(table, entity) == 0 {
		database.t.Errorf(`The table "%s" does not contain entity %+v`, table, entity)
	}

	return database
}

func (database *Database) AssertLacks(table string, entity map[string]any) *Database {
	if getMappedCount(table, entity) != 0 {
		database.t.Errorf(`The table "%s" contains entity %+v`, table, entity)
	}

	return database
}

func getMappedCount(table string, entity map[string]any) (count int) {
	db, err := mysql.Connect()
	defer db.Close()
	Check(err)

	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE`, table)
	isFirstElement := true
	for column, value := range entity {
		if isFirstElement {
			query += fmt.Sprintf(` %s = '%v'`, column, value)
			isFirstElement = false
			continue
		}
		query += fmt.Sprintf(` AND %s = '%v'`, column, value)
	}

	err = db.
		QueryRow(query).
		Scan(&count)
	Check(err)

	return
}

func (database *Database) AssertCount(table string, expected int) *Database {
	db, err := mysql.Connect()
	defer db.Close()
	Check(err)

	var count int
	err = db.
		QueryRow(`SELECT COUNT(*) FROM ` + table).
		Scan(&count)
	Check(err)

	if count != expected {
		database.t.Errorf(`The "%s" table must have %d rows, got %d`, table, expected, count)
	}

	return database
}
