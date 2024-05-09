package controller

import (
	"encoding/json"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/validation"
	"net/http"
)

type Controller[T any] struct {
	Input   T
	encoder *json.Encoder
	writer  http.ResponseWriter
}

func New[T any](writer http.ResponseWriter) (controller Controller[T]) {
	controller.encoder = json.NewEncoder(writer)
	controller.writer = writer

	return
}

func (controller *Controller[anyType]) Response(data any, status int) {
	switch data.(type) {
	case error:
		err := data.(error)
		http.Error(controller.writer, err.Error(), status)
		log.Error(err)
	default:
		controller.writer.WriteHeader(status)
		if err := controller.encoder.Encode(data); err != nil {
			http.Error(controller.writer, err.Error(), status)
			log.Error(err)
		}
	}
}

func (controller *Controller[anyType]) Validate(request *http.Request) (ok bool) {
	err := controller.Parse(request)

	if err != nil {
		http.Error(controller.writer, err.Error(), http.StatusBadRequest)
		return
	}

	validationErrors, err := validation.Perform(controller.Input)
	if err != nil {
		http.Error(controller.writer, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}

	if len(validationErrors) > 0 {
		controller.Response(validationErrors, http.StatusUnprocessableEntity)
		return
	}

	return true
}

func WriteError(writer http.ResponseWriter, err error) {
	http.Error(writer, err.Error(), http.StatusBadRequest)
	log.Error(err)
}

func (controller *Controller[anyType]) Parse(request *http.Request) error {
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&controller.Input)

	if err != nil {
		return err
	}

	return nil
}

func Pars[T any](request *http.Request) (T, error) {
	input := new(T)
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(input)

	if err != nil {
		return *input, err
	}

	return *input, nil
}

func GetInpu[T any](writer http.ResponseWriter, request *http.Request, encoder *json.Encoder) (input T, ok bool) {
	input, err := Pars[T](request)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	validationErrors, err := validation.Perform(input)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}

	if len(validationErrors) > 0 {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		if err := encoder.Encode(validationErrors); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			log.Error(err)
		}
		return
	}

	ok = true
	return
}
