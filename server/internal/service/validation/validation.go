package validation

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/service/validation/rule"
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

func (this *Validation) IsFailed() bool {
	return len(this.errors) > 0
}

func (this *Validation) Field(name string, value any, pipeRules string) *Validation {
	var rules = make([]rule.Rule, 0)

	for _, pipeRule := range strings.Split(pipeRules, "|") {
		rules = append(rules, rule.Determine(pipeRule))
	}

	field := field{name, value, rules}
	this.fields = append(this.fields, field)

	return this
}

func New() *Validation {
	return &Validation{make([]field, 0), make(map[string]string)}
}

func (this *Validation) Perform() (*Validation, error) {
	for _, field := range this.fields {
		err := validate(field)

		if err != nil {
			this.errors[field.name] = err.Error()
		}
	}

	fmt.Println(this.errors)
	return this, nil
}

func validate(field field) error {
	for _, rule := range field.rules {
		err := rule.Apply(field.name, field.value)

		if err != nil {
			return err
		}
	}

	return nil
}

//
//func (this *field) validate() error {
//	for _, constraint := range this.rule {
//		err := constraint.apply(this.name, this.value)
//
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func (this *Validation) extractRules() error {
//	for _, rule := range this.fields {
//		for _, pipeRule := range strings.Split(rule.pipeRules, "|") {
//			rule, err := extractRule(pipeRule)
//
//			if err != nil {
//				return err
//			}
//			field.rule = append(field.rule, rule)
//		}
//		fmt.Printf("Field rule: %v\n", rule.rule)
//	}
//
//	return nil
//}
