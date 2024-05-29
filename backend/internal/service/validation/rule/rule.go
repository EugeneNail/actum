package rule

import (
	"errors"
	"fmt"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
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
	case "exists":
		return newExistsRuleFunc(rule)
	case "mixedCase":
		return mixedCase
	case "before":
		return newBeforeRuleFunc(rule)
	case "after":
		return newAfterRuleFunc(rule)
	case "today":
		return newBeforeRuleFunc([]string{"", time.Now().Format("2006-01-02")})
	case "integer":
		return integer
	default:
		return func(string, any) (error, error) { return nil, fmt.Errorf("unknown rule %s", rule[0]) }
	}
}

func required(name string, value any) (error, error) {
	if value == 0 || value == "" {
		return fmt.Errorf("The %s field is required", name), nil
	}

	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() == reflect.Slice && reflectValue.Len() == 0 {
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
		return fmt.Errorf("The %s field be one word, only containing letters", name), nil
	}

	return nil, nil
}

func sentence(name string, value any) (error, error) {
	match, err := regexp.MatchString("^[a-zA-Z0-9 -]+$", value.(string))

	if err != nil {
		return nil, fmt.Errorf("sentence(): %w", err)
	}

	if !match {
		return fmt.Errorf("The %s field must be a sentence containing letters, number, spaces, slashes or dashes", name), nil
	}

	return nil, nil
}

func date(name string, value any) (error, error) {
	message := fmt.Sprintf("The %s field must be a valid YYYY-MM-DD date value", name)

	match, err := regexp.MatchString("^20[0-9]{2}-(0[0-9]|1[0-2])-([0-2][0-9]|3[0-1])$", value.(string))
	if err != nil {
		return errors.New(message), nil
	}

	if !match {
		return errors.New(message), nil
	}

	_, err = time.Parse("2006-01-02", value.(string))
	if err != nil {
		return fmt.Errorf(message), nil
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
		rows, err := db.Query(query, value)
		defer rows.Close()
		if err != nil {
			return nil, fmt.Errorf("newUniqueRuleFunc(): %w", err)
		}

		var count int
		for rows.Next() {
			err := rows.Scan(&count)

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

func mixedCase(name string, value any) (error, error) {
	hasLower := false
	hasUpper := false

	for _, char := range value.(string) {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsLower(char) {
			hasLower = true
		}
	}

	if !(hasLower && hasUpper) {
		return fmt.Errorf("The %s field must contain at least one uppercase and one lowercase letter", name), nil
	}

	return nil, nil
}

func newExistsRuleFunc(rule []string) RuleFunc {
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
		defer db.Close()
		if err != nil {
			return nil, fmt.Errorf("newExistsRuleFunc(): %w", err)
		}

		var count int
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?", table, column)
		rows, err := db.Query(query, value)
		defer rows.Close()
		if err != nil {
			return nil, fmt.Errorf("newExistsRuleFunc(): %w", err)
		}

		for rows.Next() {
			err := rows.Scan(&count)
			if err != nil {
				return nil, fmt.Errorf("newExistsRuleFunc(): %w", err)
			}
		}

		if count < 1 {
			return errors.New("The selected value is not found"), nil
		}

		return nil, nil
	}
}

func newBeforeRuleFunc(rule []string) RuleFunc {
	conditionDate, err := time.Parse("2006-01-02", rule[1])
	if err != nil {
		return func(string, any) (error, error) {
			return nil, err
		}
	}
	conditionDate.Add(time.Second*60*60*24 - 1)

	return func(name string, value any) (error, error) {
		validationError, err := date(name, value)
		if validationError != nil || err != nil {
			return validationError, err
		}

		inputDate, err := time.Parse("2006-01-02", value.(string))
		if err != nil {
			return nil, fmt.Errorf("newBeforeRuleFunc(): %w", err)
		}

		if inputDate.After(conditionDate) {
			return errors.New(fmt.Sprintf("The %s field must be a date before %s", name, rule[1])), nil
		}

		return nil, nil
	}
}

func newAfterRuleFunc(rule []string) RuleFunc {
	conditionDate, err := time.Parse("2006-01-02", rule[1])
	if err != nil {
		return func(string, any) (error, error) {
			return nil, err
		}
	}

	return func(name string, value any) (error, error) {
		validationError, err := date(name, value)
		if validationError != nil || err != nil {
			return validationError, err
		}

		inputDate, err := time.Parse("2006-01-02", value.(string))
		if err != nil {
			return nil, fmt.Errorf("newBeforeRuleFunc(): %w", err)
		}

		if inputDate.Before(conditionDate) {
			return errors.New(fmt.Sprintf("The %s field must be a date after %s", name, rule[1])), nil
		}

		return nil, nil
	}
}

func integer(name string, value any) (error, error) {
	switch value.(type) {
	case float32:
		if int(value.(float32)) != value {
			return fmt.Errorf("The %s field must be an integer", name), nil
		}
	case float64:
		if int(value.(float64)) != value {
			return fmt.Errorf("The %s field must be an integer", name), nil
		}
	}

	return nil, nil
}
