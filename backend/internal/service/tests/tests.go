package tests

import (
	"encoding/json"
	"fmt"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/model/users"
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/validation/rule"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type ValidationTest struct {
	Name  string
	Field string
	Value any
}

type Response struct {
	*http.Response
	t *testing.T
}

func (response *Response) AssertStatus(status int) {
	if response.StatusCode != status {
		response.t.Errorf("expected status %d, got %d", status, response.StatusCode)
	}
}

func Post(path string, t *testing.T, json string) Response {
	url := "http://127.0.0.1:" + env.Get("APP_PORT") + path
	body := strings.NewReader(json)
	response, err := http.Post(url, "application/json", body)
	Check(err)

	return Response{response, t}
}

func AssertValidation[T any](mustSuccess bool, t *testing.T, validationTests []ValidationTest) {
	for _, validationTest := range validationTests {
		t.Run(validationTest.Name, func(t *testing.T) {
			errorsCount := getValidationErrorsCount[T](validationTest)

			if mustSuccess && errorsCount > 0 {
				t.Errorf("Must success at value  ---  %s", validationTest.Value)
			}

			if !mustSuccess && errorsCount == 0 {
				t.Errorf("Must fail at value  ---  %s", validationTest.Value)
			}
		})
	}
}

func AssertValidationSuccess[T any](t *testing.T, validationTests []ValidationTest) {
	AssertValidation[T](true, t, validationTests)
}

func AssertValidationFail[T any](t *testing.T, validationTests []ValidationTest) {
	AssertValidation[T](false, t, validationTests)
}

func getValidationErrorsCount[T any](test ValidationTest) int {
	structField := getStructFieldByName[T](test.Field)

	var ruleFuncs []rule.RuleFunc

	validationRules := structField.Tag.Get("rules")
	if len(validationRules) == 0 {
		return 0
	}

	for _, validationRule := range strings.Split(validationRules, "|") {
		ruleFuncs = append(ruleFuncs, rule.Extract(validationRule))
	}

	errorCount := 0

	for _, ruleFunc := range ruleFuncs {
		validationError, err := ruleFunc(test.Field, test.Value)
		Check(err)

		if validationError != nil {
			errorCount++
		}
	}

	return errorCount
}

func getStructFieldByName[T any](name string) reflect.StructField {
	for _, structField := range reflect.VisibleFields(reflect.TypeOf(*new(T))) {
		if structField.Tag.Get("json") == name {
			return structField
		}
	}
	panic("struct field not found: " + name)
}

func (response *Response) AssertHasValidationErrors(fields []string) {
	var errors map[string]string
	data, err := io.ReadAll(response.Body)
	Check(err)
	err = json.Unmarshal(data, &errors)
	Check(err)

	for _, field := range fields {
		if _, exists := errors[field]; !exists {
			response.t.Errorf(`expected validation error for field "%s" to be present`, field)
		}
	}
}

func (response *Response) AssertHasToken() {
	if !hasToken(response.Response) {
		response.t.Errorf("The response must have an Access-Token cookie")
	}
}

func (response *Response) AssertHasNoToken() {
	if hasToken(response.Response) {
		response.t.Errorf("The response must not have an Access-Token cookie")
	}
}

func hasToken(response *http.Response) bool {
	for _, cookie := range response.Cookies() {
		if cookie.Name == "Access-Token" && len(cookie.Value) > 0 {
			return true
		}
	}

	return false
}

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

func Check(err error) {
	if err != nil {
		panic(err)
	}
}
