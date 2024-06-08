package validation

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/service/validation/rule"
	"reflect"
	"strings"
)

type field struct {
	name      string
	value     any
	ruleFuncs []rule.Func
}

func Perform(data any) (map[string]string, error) {
	fields := extractFields(data)
	validationErrors := make(map[string]string)

	for _, field := range fields {
	currentFieldLoop:
		for _, applyRule := range field.ruleFuncs {
			validationError, err := applyRule(field.value)

			if err != nil {
				return nil, fmt.Errorf("validate(): %w", err)
			}

			if len(validationError) > 0 {
				validationErrors[field.name] = validationError
				break currentFieldLoop
			}
		}
	}

	return validationErrors, nil
}

func extractFields(data any) []field {
	structFields := reflect.VisibleFields(reflect.TypeOf(data))
	v := reflect.ValueOf(data)
	var fields = make([]field, 0, len(structFields))

	for _, structField := range structFields {
		pipeRules := structField.Tag.Get("rules")

		if len(pipeRules) > 0 {
			name := structField.Tag.Get("json")
			value := v.FieldByName(structField.Name).Interface()
			fields = append(fields, newField(name, value, pipeRules))
		}
	}

	return fields
}

func newField(name string, value any, pipeRules string) field {
	var rules = make([]rule.Func, 0)

	for _, pipeRule := range strings.Split(pipeRules, "|") {
		rules = append(rules, rule.Extract(pipeRule))
	}

	return field{name, value, rules}
}
