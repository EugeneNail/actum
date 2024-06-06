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
	"unicode/utf8"
)

type Func func(any) (error, error)

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
		return newBeforeRuleFunc([]string{"", time.Now().Format("2006-01-02")})
	case "integer":
		return integer
	default:
		return func(any) (error, error) { return nil, fmt.Errorf("unknown rule %s", rule[0]) }
	}
}

func required(value any) (error, error) {
	if value == 0 || value == "" {
		return fmt.Errorf("Заполните это поле."), nil
	}

	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() == reflect.Slice && reflectValue.Len() == 0 {
		return fmt.Errorf("Заполните это поле."), nil
	}

	return nil, nil
}

func word(value any) (error, error) {
	match, err := regexp.MatchString("^[a-zA-Zа-яА-Я]+$", value.(string))

	if err != nil || !match {
		return fmt.Errorf("Введите одно слово, состоящее только из букв."), nil
	}

	return nil, nil
}

func sentence(value any) (error, error) {
	match, err := regexp.MatchString("^[a-zA-Zа-яА-Я0-9 -]+$", value.(string))

	if err != nil {
		return nil, fmt.Errorf("sentence(): %w", err)
	}

	if !match {
		return fmt.Errorf("Введите предложение. Оно может содержать сколько угодно слов, цифры и тире."), nil
	}

	return nil, nil
}

func date(value any) (error, error) {
	message := "Введите дату в формате YYYY-MM-DD"

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

func email(value any) (error, error) {
	match, err := regexp.MatchString("^\\S+@\\S+\\.\\S+$", value.(string))

	if err != nil {
		return nil, fmt.Errorf("email(): %w", err)
	}

	if !match {
		return fmt.Errorf("Введите действующий адрес электронной почты."), nil
	}

	return nil, nil
}

func newMinRuleFunc(rule []string) Func {
	limit := 0
	if len(rule) == 2 {
		parsed, err := strconv.Atoi(rule[1])

		if err == nil {
			limit = parsed
		}
	}

	return func(value any) (error, error) {
		switch value.(type) {
		case string:
			if utf8.RuneCountInString(value.(string)) < limit {
				return fmt.Errorf("Длина должная быть не менее %d символов.", limit), nil
			}
		case int:
			if value.(int) < limit {
				return fmt.Errorf("Значение должно быть больше или равно %d.", limit), nil
			}
		}

		return nil, nil
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

	return func(value any) (error, error) {
		switch value.(type) {
		case string:
			if utf8.RuneCountInString(value.(string)) > limit {
				return fmt.Errorf("Длина не должна превышать %d символов.", limit), nil
			}
		case int:
			if value.(int) > limit {
				return fmt.Errorf("Значение должно быть больше или равно %d.", limit), nil
			}
		}

		return nil, nil
	}
}

func newRegexRuleFunc(rule []string) Func {
	pattern := ".*"
	if len(rule) == 2 {
		pattern = rule[1]
	}

	return func(value any) (error, error) {
		match, err := regexp.MatchString(pattern, value.(string))

		if err != nil {
			return nil, fmt.Errorf("newRegexRuleFunc(): %w", err)
		}

		if !match {
			return fmt.Errorf("Введите данные в правильном формате."), nil
		}

		return nil, nil
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

	return func(value any) (error, error) {
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
			return fmt.Errorf("Запись с таким значением уже существует."), nil
		}

		return nil, nil
	}
}

func mixedCase(value any) (error, error) {
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
		return fmt.Errorf("Введите хотя бы по одной строчной и заглавной букве."), nil
	}

	return nil, nil
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

	return func(value any) (error, error) {
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
			return errors.New("Запись с таким значением не существует."), nil
		}

		return nil, nil
	}
}

func newBeforeRuleFunc(rule []string) Func {
	conditionDate, err := time.Parse("2006-01-02", rule[1])
	if err != nil {
		return func(any) (error, error) {
			return nil, err
		}
	}
	conditionDate.Add(time.Second*60*60*24 - 1)

	return func(value any) (error, error) {
		validationError, err := date(value)
		if validationError != nil || err != nil {
			return validationError, err
		}

		inputDate, err := time.Parse("2006-01-02", value.(string))
		if err != nil {
			return nil, fmt.Errorf("newBeforeRuleFunc(): %w", err)
		}

		if inputDate.After(conditionDate) {
			return errors.New(fmt.Sprintf("Введите дату до %s", rule[1])), nil
		}

		return nil, nil
	}
}

func newAfterRuleFunc(rule []string) Func {
	conditionDate, err := time.Parse("2006-01-02", rule[1])
	if err != nil {
		return func(any) (error, error) {
			return nil, err
		}
	}

	return func(value any) (error, error) {
		validationError, err := date(value)
		if validationError != nil || err != nil {
			return validationError, err
		}

		inputDate, err := time.Parse("2006-01-02", value.(string))
		if err != nil {
			return nil, fmt.Errorf("newBeforeRuleFunc(): %w", err)
		}

		if inputDate.Before(conditionDate) {
			return errors.New(fmt.Sprintf("Введите дату не ранее %s.", rule[1])), nil
		}

		return nil, nil
	}
}

func integer(value any) (error, error) {
	switch value.(type) {
	case float32:
		if int(value.(float32)) != value {
			return fmt.Errorf("Введите целое число."), nil
		}
	case float64:
		if int(value.(float64)) != value {
			return fmt.Errorf("Введите целое число."), nil
		}
	}

	return nil, nil
}
