package rule

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type RuleFunc func(string, any) error

func Extract(pipeRule string) RuleFunc {
	rule := strings.Split(pipeRule, ":")
	switch rule[0] {
	case "required":
		return required
	case "accepted":
		return accepted
	case "word":
		return word
	case "sentence":
		return sentence
	case "date":
		return date
	case "email":
		return email
	case "min":
		return newMinRuleFunc(rule)
	case "max":
		return newMaxRuleFunc(rule)
	case "regex":
		return newRegexRuleFunc(rule)
	default:
		return func(string, any) error { return nil }
	}
}

func required(name string, value any) error {
	if value == 0 || value == "" || value == nil {
		return fmt.Errorf("The %s field is required", name)
	}

	return nil
}

func accepted(name string, value any) error {
	if value != true {
		return fmt.Errorf("The %s field must be accepted", name)
	}

	return nil
}

func word(name string, value any) error {
	match, err := regexp.MatchString("^[a-zA-Z]+$", value.(string))

	if err != nil || !match {
		return fmt.Errorf("The %s field be one word", name)
	}

	return nil
}

func sentence(name string, value any) error {
	match, err := regexp.MatchString("^[a-zA-Z0-9 -/]+$", value.(string))
	if err != nil || !match {
		return fmt.Errorf("The %s field must be a sentence containing letters, number, spaces, slashes or dashes", name)
	}

	return nil
}

func date(name string, value any) error {
	_, err := time.Parse(time.DateTime, value.(string))

	if err != nil {
		return fmt.Errorf("The %s field must be a datetime value", name)
	}

	return nil
}

func email(name string, value any) error {
	match, err := regexp.MatchString("^\\S+@\\S+\\.\\S+$", value.(string))

	if err != nil || !match {
		return fmt.Errorf("The %s field must be a valid email address", name)
	}

	return nil
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
				return fmt.Errorf("The %s field must be at least %d characters", name, limit)
			}
		case int:
			if value.(int) < limit {
				return fmt.Errorf("The %s field must be greater than %d", name, limit)
			}
		}

		return nil
	}
}

func newMaxRuleFunc(rule []string) RuleFunc {
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
			if len(value.(string)) > limit {
				return fmt.Errorf("The %s field must not be greater than %d characters", name, limit)
			}
		case int:
			if value.(int) > limit {
				return fmt.Errorf("The %s field must be less than %d", name, limit)
			}
		}

		return nil
	}
}

func newRegexRuleFunc(rule []string) RuleFunc {
	pattern := ".*"

	if len(rule) == 2 {
		pattern = rule[1]
	}

	return func(name string, value any) error {
		match, err := regexp.MatchString(pattern, value.(string))

		if err != nil || !match {
			return fmt.Errorf("The %s field format is invalid", name)
		}

		return nil
	}
}
