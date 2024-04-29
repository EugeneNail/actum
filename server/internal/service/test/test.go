package test

import (
	"github.com/EugeneNail/actum/internal/service/validation/rule"
	"reflect"
	"strings"
	"testing"
)

type Field struct {
	Name  string
	Value any
}

func AssertValidationSuccess[T any](field Field, t *testing.T) {
	errorCount := getValidationErrorCount[T](field)
	if errorCount > 0 {
		t.Errorf("%s: should success at value [%s]", field.Name, field.Value)
	}
}

func AssertValidationFail[T any](field Field, t *testing.T) {
	errorCount := getValidationErrorCount[T](field)
	if errorCount == 0 {
		t.Errorf("%s: should fail at value [%s] ", field.Name, field.Value)
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

func getValidationErrorCount[T any](field Field) int {
	structField := getStructFieldByName[T](field.Name)

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
		validationError, err := ruleFunc(field.Name, field.Value)
		if err != nil {
			panic(err)
		}

		if validationError != nil {
			errorCount++
		}
	}

	return errorCount
}
