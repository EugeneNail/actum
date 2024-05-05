package test

import (
	"github.com/EugeneNail/actum/internal/model/users"
	"github.com/EugeneNail/actum/internal/service/validation/rule"
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

func AssertValidationSuccess[T any](test ValidationTest, t *testing.T) {
	errorCount := getValidationErrorCount[T](test)
	if errorCount > 0 {
		t.Errorf("Must success at value  ---  %s", test.Value)
	}
}

func AssertValidationFail[T any](test ValidationTest, t *testing.T) {
	errorCount := getValidationErrorCount[T](test)
	if errorCount == 0 {
		t.Errorf("Must fail at value  ---  %s", test.Value)
	}
}

func getStructFieldByName[T any](name string) reflect.StructField {
	for _, structField := range reflect.VisibleFields(reflect.TypeOf(*new(T))) {
		if structField.Tag.Get("json") == name {
			return structField
		}
	}
	panic("struct field not found: " + name)
}

func getValidationErrorCount[T any](test ValidationTest) int {
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

func AssertHasToken(response *http.Response, t *testing.T) {
	if !hasToken(response) {
		t.Errorf("The response must have an Access-Token cookie")
	}
}

func AssertHasNoToken(response *http.Response, t *testing.T) {
	if hasToken(response) {
		t.Errorf("The response must not have an Access-Token cookie")
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

func Check(err error) {
	if err != nil {
		panic(err)
	}
}
