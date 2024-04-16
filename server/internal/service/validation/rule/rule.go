package rule

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Determine(pipeRule string) RuleFunc {
	rule := strings.Split(pipeRule, ":")
	switch rule[0] {
	case "required":
		return required
	case "min":
		return newMinRuleFunc(rule)
	default:
		return func(string, any) error { return nil }
	}
}

func newMinRuleFunc(rule []string) RuleFunc {
	limit := 0
	if len(rule) == 2 {
		parsed, err := strconv.Atoi(rule[1])

		if err == nil {
			limit = parsed
		}
	}

	return func(name string, value any) error {
		switch value.(type) {
		case string:
			if len(value.(string)) < limit {
				return fmt.Errorf("The %s field must be at least %d characters long", name, limit)
			}
		case int:
			if value.(int) < limit {
				return fmt.Errorf("The %s field must be greater than %d", name, limit)
			}
		}

		return nil
	}
}

func required(name string, value any) error {
	if value == 0 || value == "" || value == nil {
		return errors.New(fmt.Sprintf("The %s field is required", name))
	}

	return nil
}

type RuleFunc func(string, any) error
