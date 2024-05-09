package tests

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/model/users"
	"testing"
)

func AssertUserIsUntouched(user users.User, t *testing.T) {
	dbUser, err := users.Find(1)
	Check(err)
	if dbUser.Name != user.Name {
		t.Errorf(`field "name" has been corrupted`)
	}

	if dbUser.Email != user.Email {
		t.Errorf(`field "email" has been corrupted`)
	}

	if dbUser.Password != user.Password {
		t.Errorf(`field "password" has been corrupted`)
	}
}

func AssertTableIsEmpty(table string, t *testing.T) {
	db, err := mysql.Connect()
	Check(err)
	defer db.Close()

	var rows int
	err = db.
		QueryRow(`SELECT COUNT(*) FROM ` + table).
		Scan(&rows)
	Check(err)

	if rows > 0 {
		t.Errorf("The table %s is expected to be empty, got %d rows instead", table, rows)
	}
}

func AssertDatabaseHas(table string, entity map[string]any, t *testing.T) {
	db, err := mysql.Connect()
	Check(err)
	defer db.Close()

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

	var count int
	err = db.
		QueryRow(query).
		Scan(&count)
	Check(err)

	if count == 0 {
		t.Errorf("The table %s does not contain an entity %+v", table, entity)
	}
}

func AssertDatabaseCount(table string, expected int, t *testing.T) {
	db, err := mysql.Connect()
	Check(err)
	defer db.Close()

	var count int
	err = db.
		QueryRow(`SELECT COUNT(*) FROM ` + table).
		Scan(&count)
	Check(err)

	if count != expected {
		t.Errorf("The %s table must have %d rows, got %d", table, expected, count)
	}

}
