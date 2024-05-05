package controller

import (
	"encoding/json"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/validation"
	"net/http"
)

func Parse[T any](request *http.Request) (T, error) {
	input := new(T)
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(input)

	if err != nil {
		return *input, err
	}

	return *input, nil
}

func GetInput[T any](writer http.ResponseWriter, request *http.Request, encoder *json.Encoder) (input T, ok bool) {
	input, err := Parse[T](request)

	if err != nil {
		WriteError(writer, err)
		return
	}

	validationErrors, err := validation.Perform(input)
	if err != nil {
		WriteError(writer, err)
		return
	}

	if len(validationErrors) > 0 {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		if err := encoder.Encode(validationErrors); err != nil {
			WriteError(writer, err)
		}
		return
	}

	ok = true
	return
}

func WriteError(writer http.ResponseWriter, err error) {
	http.Error(writer, err.Error(), http.StatusBadRequest)
	log.Error(err)
}
