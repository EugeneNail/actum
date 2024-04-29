package rule

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type RuleFunc func(string, any) (error, error)

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
	case "unique":
		return newUniqueRuleFunc(rule)
	default:
		return func(string, any) (error, error) { return nil, fmt.Errorf("unknown rule %s", rule[0]) }
	}
}

func required(name string, value any) (error, error) {
	if value == 0 || value == "" {
		return fmt.Errorf("The %s field is required", name), nil
	}

	return nil, nil
}

func accepted(name string, value any) (error, error) {
	if value != true {
		return fmt.Errorf("The %s field must be accepted", name), nil
	}

	return nil, nil
}

func word(name string, value any) (error, error) {
	match, err := regexp.MatchString("^[a-zA-Z]+$", value.(string))

	if err != nil || !match {
		return fmt.Errorf("The %s field be one word", name), nil
	}

	return nil, nil
}

func sentence(name string, value any) (error, error) {
	match, err := regexp.MatchString("^[a-zA-Z0-9 -/]+$", value.(string))

	if err != nil {
		return nil, fmt.Errorf("sentence(): %w", err)
	}

	if !match {
		return fmt.Errorf("The %s field must be a sentence containing letters, number, spaces, slashes or dashes", name), nil
	}

	return nil, nil
}

func date(name string, value any) (error, error) {
	_, err := time.Parse(time.DateTime, value.(string))

	if err != nil {
		return fmt.Errorf("The %s field must be a datetime value", name), nil
	}

	return nil, nil
}

func email(name string, value any) (error, error) {
	match, err := regexp.MatchString("^\\S+@\\S+\\.\\S+$", value.(string))

	if err != nil {
		return nil, fmt.Errorf("email(): %w", err)
	}

	if !match {
		return fmt.Errorf("The %s field must be a valid email address", name), nil
	}

	return nil, nil
}

func newMinRuleFunc(rule []string) RuleFunc {
	limit := 0
	if len(rule) == 2 {
		parsed, err := strconv.Atoi(rule[1])

		if err == nil {
			limit = parsed
		}
	}

	return func(name string, value any) (error, error) {
		switch value.(type) {
		case string:
			if len(value.(string)) < limit {
				return fmt.Errorf("The %s field must be at least %d characters", name, limit), nil
			}
		case int:
			if value.(int) < limit {
				return fmt.Errorf("The %s field must be greater than %d", name, limit), nil
			}
		}

		return nil, nil
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

	return func(name string, value any) (error, error) {
		switch value.(type) {
		case string:
			if len(value.(string)) > limit {
				return fmt.Errorf("The %s field must not be greater than %d characters", name, limit), nil
			}
		case int:
			if value.(int) > limit {
				return fmt.Errorf("The %s field must be less than %d", name, limit), nil
			}
		}

		return nil, nil
	}
}

func newRegexRuleFunc(rule []string) RuleFunc {
	pattern := ".*"
	if len(rule) == 2 {
		pattern = rule[1]
	}

	return func(name string, value any) (error, error) {
		match, err := regexp.MatchString(pattern, value.(string))

		if err != nil {
			return nil, fmt.Errorf("newRegexRuleFunc(): %w", err)
		}

		if !match {
			return fmt.Errorf("The %s field format is invalid", name), nil
		}

		return nil, nil
	}
}

func newUniqueRuleFunc(rule []string) RuleFunc {
	table, column := "", ""

	if len(rule) == 2 {
		args := strings.Split(rule[1], ",")
		switch len(args) {
		case 1:
			table, column = args[0], ""
		case 2:
			table, column = args[0], args[1]
		}
	}

	return func(name string, value any) (error, error) {
		db, err := mysql.Connect()
		if err != nil {
			return nil, fmt.Errorf("newUniqueRuleFunc(): %w", err)
		}

		query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?", table, column)
		result, err := db.Query(query, value)
		if err != nil {
			return nil, fmt.Errorf("newUniqueRuleFunc(): %w", err)
		}

		var count int
		for result.Next() {
			err := result.Scan(&count)

			if err != nil {
				return nil, fmt.Errorf("newUniqueRuleFunc(): %w", err)
			}
		}

		if count > 0 {
			return fmt.Errorf("The %s has already been taken", name), nil
		}

		return nil, nil
	}
}
