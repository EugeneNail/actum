package tests

import (
	"fmt"
	"testing"
)

type Database struct {
	t *testing.T
}

func NewDatabase(t *testing.T) Database {
	return Database{t}
}

func (database *Database) AssertEmpty(table string) *Database {
	var rows int
	err := DB.
		QueryRow(`SELECT COUNT(*) FROM ` + table).
		Scan(&rows)
	Check(err)

	if rows > 0 {
		database.t.Errorf(`The table "%s" is expected to be empty, got %d rows instead`, table, rows)
		database.t.SkipNow()
	}

	return database
}

func (database *Database) AssertHas(table string, entity map[string]any) *Database {
	if getMappedCount(table, entity) == 0 {
		database.t.Errorf(`The table "%s" does not contain entity %+v`, table, entity)
		database.t.SkipNow()
	}

	return database
}

func (database *Database) AssertLacks(table string, entity map[string]any) *Database {
	if getMappedCount(table, entity) != 0 {
		database.t.Errorf(`The table "%s" contains entity %+v`, table, entity)
		database.t.SkipNow()
	}

	return database
}

func getMappedCount(table string, entity map[string]any) (count int) {
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE`, table)
	isFirstElement := true
	for column, value := range entity {
		if isFirstElement {
			query += fmt.Sprintf(` binary %s = '%v'`, column, value)
			isFirstElement = false
			continue
		}
		query += fmt.Sprintf(` AND binary %s = '%v'`, column, value)
	}

	err := DB.
		QueryRow(query).
		Scan(&count)
	Check(err)

	return
}

func (database *Database) AssertCount(table string, expected int) *Database {
	var count int
	err := DB.
		QueryRow(`SELECT COUNT(*) FROM ` + table).
		Scan(&count)
	Check(err)

	if count != expected {
		database.t.Errorf(`The "%s" table must have %d rows, got %d`, table, expected, count)
		database.t.SkipNow()
	}

	return database
}
