package rule

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Determine(pipeRule string) Rule {
	rule := strings.Split(pipeRule, ":")
	switch rule[0] {
	case "required":
		return required{}
	case "min":
		return newMinRule(rule)
	default:
		return nil
	}
}

func newMinRule(rule []string) Rule {
	if len(rule) == 1 {
		return min{0}
	}
	limit, err := strconv.Atoi(rule[1])

	if err != nil {
		return min{0}
	}

	return min{limit}
}

type required struct{}

func (this required) Apply(name string, value any) error {
	if value == 0 || value == "" || value == nil {
		return errors.New(fmt.Sprintf("The %s field is required", name))
	}

	return nil
}

type min struct {
	limit int
}

func (this min) Apply(name string, value any) error {
	switch value.(type) {
	case string:
		if len(value.(string)) < this.limit {
			return fmt.Errorf("The %s field must be at least %d characters long", name, this.limit)
		}
	case int:
		if value.(int) < this.limit {
			return fmt.Errorf("The %s field must be greater than %d", name, this.limit)
		}
	case []any:
		fmt.Println("It is array")
		if (len(value.([]any))) < this.limit {
			return fmt.Errorf("The %s field must have at least %d items", name, this.limit)
		}
	}

	return nil
}

type Rule interface {
	Apply(string, any) error
}
