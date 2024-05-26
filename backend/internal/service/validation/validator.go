package validation

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Validator[T any] struct {
	input T
}

func NewValidator[T any]() *Validator[T] {
	return &Validator[T]{}
}

func (validator *Validator[any]) Validate(request *http.Request) (map[string]string, any, error) {
	validationErrors := make(map[string]string)

	if err := validator.parse(request); err != nil {
		return validationErrors, validator.input, fmt.Errorf("validation.Validate(): %w", err)
	}

	validationErrors, err := Perform(validator.input)
	if err != nil {
		return validationErrors, validator.input, fmt.Errorf("validation.Validate(): %w", err)
	}

	return validationErrors, validator.input, nil
}

func (validator *Validator[any]) parse(request *http.Request) error {
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&validator.input); err != nil {
		return fmt.Errorf("parse(): %w", err)
	}

	return nil
}
