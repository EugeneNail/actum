package tests

import (
	"github.com/EugeneNail/actum/internal/service/validation/rule"
	"reflect"
	"strings"
	"testing"
)

type ValidationTest struct {
	Field string
	Name  string
	Value any
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

	var ruleFuncs []rule.Func

	validationRules := structField.Tag.Get("rules")
	if len(validationRules) == 0 {
		return 0
	}

	for _, validationRule := range strings.Split(validationRules, "|") {
		ruleFuncs = append(ruleFuncs, rule.Extract(validationRule))
	}

	errorCount := 0

	for _, ruleFunc := range ruleFuncs {
		validationError, err := ruleFunc(test.Value)
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
