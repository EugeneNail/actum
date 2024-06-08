package rule

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type Func func(any) (string, error)

func Extract(pipeRule string) Func {
	rule := strings.Split(pipeRule, ":")
	switch rule[0] {
	case "required":
		return required
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
		return newBeforeRuleFunc([]string{"", time.Now().Add(time.Hour * 24).Format("2006-01-02")})
	case "integer":
		return integer
	default:
		return func(any) (string, error) { return "", fmt.Errorf("unknown rule %s", rule[0]) }
	}
}

func required(value any) (string, error) {
	if value == 0 || value == "" {
		return "Заполните это поле.", nil
	}

	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() == reflect.Slice && reflectValue.Len() == 0 {
		return "Заполните это поле.", nil
	}

	return "", nil
}

func word(value any) (string, error) {
	match, err := regexp.MatchString("^[a-zA-Zа-яА-Я]+$", value.(string))

	if err != nil || !match {
		return "Введите одно слово, состоящее только из букв.", nil
	}

	return "", nil
}

func sentence(value any) (string, error) {
	match, err := regexp.MatchString("^[a-zA-Zа-яА-Я0-9 -]+$", value.(string))

	if err != nil {
		return "", fmt.Errorf("sentence(): %w", err)
	}

	if !match {
		return "Введите предложение. Оно может содержать сколько угодно слов, цифры и тире.", nil
	}

	return "", nil
}

func date(value any) (string, error) {
	message := "Введите дату в формате YYYY-MM-DD"

	match, err := regexp.MatchString("^20[0-9]{2}-(0[0-9]|1[0-2])-([0-2][0-9]|3[0-1])$", value.(string))
	if err != nil {
		return message, nil
	}

	if !match {
		return message, nil
	}

	_, err = time.Parse("2006-01-02", value.(string))
	if err != nil {
		return message, nil
	}

	return "", nil
}

func email(value any) (string, error) {
	match, err := regexp.MatchString("^\\S+@\\S+\\.\\S+$", value.(string))

	if err != nil {
		return "", fmt.Errorf("email(): %w", err)
	}

	if !match {
		return "Введите действующий адрес электронной почты.", nil
	}

	return "", nil
}

func newMinRuleFunc(rule []string) Func {
	limit := 0
	if len(rule) == 2 {
		parsed, err := strconv.Atoi(rule[1])

		if err == nil {
			limit = parsed
		}
	}

	return func(value any) (string, error) {
		switch value.(type) {
		case string:
			if utf8.RuneCountInString(value.(string)) < limit {
				return fmt.Sprintf("Длина должная быть не менее %d символов.", limit), nil
			}
		case int:
			if value.(int) < limit {
				return fmt.Sprintf("Значение должно быть больше или равно %d.", limit), nil
			}
		}

		return "", nil
	}
}

func newMaxRuleFunc(rule []string) Func {
	limit := 0
	if len(rule) == 2 {
		parsed, err := strconv.Atoi(rule[1])

		if err == nil {
			limit = parsed
		}
	}

	return func(value any) (string, error) {
		switch value.(type) {
		case string:
			if utf8.RuneCountInString(value.(string)) > limit {
				return fmt.Sprintf("Длина не должна превышать %d символов.", limit), nil
			}
		case int:
			if value.(int) > limit {
				return fmt.Sprintf("Значение должно быть больше или равно %d.", limit), nil
			}
		}

		return "", nil
	}
}

func newRegexRuleFunc(rule []string) Func {
	pattern := ".*"
	if len(rule) == 2 {
		pattern = rule[1]
	}

	return func(value any) (string, error) {
		match, err := regexp.MatchString(pattern, value.(string))

		if err != nil {
			return "", fmt.Errorf("newRegexRuleFunc(): %w", err)
		}

		if !match {
			return "Введите данные в правильном формате.", nil
		}

		return "", nil
	}
}

func newUniqueRuleFunc(rule []string) Func {
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

	return func(value any) (string, error) {
		db, err := mysql.Connect()
		if err != nil {
			return "", fmt.Errorf("newUniqueRuleFunc(): %w", err)
		}

		query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?", table, column)
		rows, err := db.Query(query, value)
		defer rows.Close()
		if err != nil {
			return "", fmt.Errorf("newUniqueRuleFunc(): %w", err)
		}

		var count int
		for rows.Next() {
			err := rows.Scan(&count)

			if err != nil {
				return "", fmt.Errorf("newUniqueRuleFunc(): %w", err)
			}
		}

		if count > 0 {
			return "Запись с таким значением уже существует.", nil
		}

		return "", nil
	}
}

func mixedCase(value any) (string, error) {
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
		return "Введите хотя бы по одной строчной и заглавной букве.", nil
	}

	return "", nil
}

func newExistsRuleFunc(rule []string) Func {
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

	return func(value any) (string, error) {
		db, err := mysql.Connect()
		defer db.Close()
		if err != nil {
			return "", fmt.Errorf("newExistsRuleFunc(): %w", err)
		}

		var count int
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?", table, column)
		rows, err := db.Query(query, value)
		defer rows.Close()
		if err != nil {
			return "", fmt.Errorf("newExistsRuleFunc(): %w", err)
		}

		for rows.Next() {
			err := rows.Scan(&count)
			if err != nil {
				return "", fmt.Errorf("newExistsRuleFunc(): %w", err)
			}
		}

		if count < 1 {
			return "Запись с таким значением не существует.", nil
		}

		return "", nil
	}
}

func newBeforeRuleFunc(rule []string) Func {
	conditionDate, err := time.Parse("2006-01-02", rule[1])
	if err != nil {
		return func(any) (string, error) {
			return "", fmt.Errorf("newBeforeRuleFunc(): invalid rule condition date: %w", err)
		}
	}

	return func(value any) (string, error) {
		validationError, err := date(value)
		if len(validationError) != 0 || err != nil {
			return validationError, err
		}

		inputDate, err := time.Parse("2006-01-02", value.(string))
		if err != nil {
			return "Введите дату в формате YYYY-MM-DD", nil
		}

		if !inputDate.Before(conditionDate) {
			return fmt.Sprintf("Введите дату до %s", rule[1]), nil
		}
		return "", nil
	}
}

func newAfterRuleFunc(rule []string) Func {
	conditionDate, err := time.Parse("2006-01-02", rule[1])
	if err != nil {
		return func(any) (string, error) {
			return "", fmt.Errorf("newAfterRuleFunc(): invalid rule condition date: %w", err)
		}
	}

	return func(value any) (string, error) {
		validationError, err := date(value)
		if len(validationError) > 0 || err != nil {
			return validationError, err
		}

		inputDate, err := time.Parse("2006-01-02", value.(string))
		if err != nil {
			return "Введите дату в формате YYYY-MM-DD", nil
		}

		if inputDate.Before(conditionDate) {
			return fmt.Sprintf("Введите дату не ранее %s.", rule[1]), nil
		}

		return "", nil
	}
}

func integer(value any) (string, error) {
	switch value.(type) {
	case float32:
		if int(value.(float32)) != value {
			return "Введите целое число.", nil
		}
	case float64:
		if int(value.(float64)) != value {
			return "Введите целое число.", nil
		}
	}

	return "", nil
}
