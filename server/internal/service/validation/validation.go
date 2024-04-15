package validation

import (
	"github.com/EugeneNail/actum/internal/service/validation/rule"
	"reflect"
	"strings"
)

type field struct {
	name  string
	value any
	rules []rule.Rule
}

type Validation struct {
	fields []field
	errors map[string]string
}

func New(data any) *Validation {
	validation := &Validation{make([]field, 0), map[string]string{}}
	fields := extractFields(data)
	validation.errors = validate(fields)

	return validation
}

func validate(fields []field) map[string]string {
	errors := make(map[string]string)

	for _, field := range fields {
	ruleLoop:
		for _, rule := range field.rules {
			err := rule.Apply(field.name, field.value)

			if err != nil {
				errors[field.name] = err.Error()
				break ruleLoop
			}
		}
	}

	return errors
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
	var rules = make([]rule.Rule, 0)

	for _, pipeRule := range strings.Split(pipeRules, "|") {
		rules = append(rules, rule.Determine(pipeRule))
	}

	return field{name, value, rules}
}

func (this *Validation) Errors() map[string]string {
	return this.errors
}

func (this *Validation) IsFailed() bool {
	return len(this.errors) > 0
}
