package validation

import (
	"errors"
	"github.com/EugeneNail/actum/internal/service/validation/rule"
	"reflect"
	"strings"
)

type field struct {
	name      string
	value     any
	ruleFuncs []rule.RuleFunc
}

func Perform(data any) (map[string]string, error) {
	fields := extractFields(data)
	validationErrors := validate(fields)

	if len(validationErrors) > 0 {
		return validationErrors, errors.New("unprocessable entity")
	}

	return validationErrors, nil
}

func extractFields(data any) []field {
	structFields := reflect.VisibleFields(reflect.TypeOf(data))
	v := reflect.ValueOf(data)
	var fields = make([]field, 0, len(structFields))

	for _, structField := range structFields {
		name := structField.Tag.Get("json")
		pipeRules := structField.Tag.Get("rules")
		value := v.FieldByName(structField.Name).Interface()
		fields = append(fields, newField(name, value, pipeRules))
	}

	return fields
}

func newField(name string, value any, pipeRules string) field {
	var rules = make([]rule.RuleFunc, 0)

	for _, pipeRule := range strings.Split(pipeRules, "|") {
		rules = append(rules, rule.Extract(pipeRule))
	}

	return field{name, value, rules}
}

func validate(fields []field) map[string]string {
	errors := make(map[string]string)

	for _, field := range fields {
	ruleLoop:
		for _, ruleFunc := range field.ruleFuncs {
			err := ruleFunc(field.name, field.value)

			if err != nil {
				errors[field.name] = err.Error()
				break ruleLoop
			}
		}
	}

	return errors
}
