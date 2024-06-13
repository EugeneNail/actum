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
	err := controller.parse(request)

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

func (controller *Controller[anyType]) parse(request *http.Request) error {
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&controller.Input)

	if err != nil {
		return err
	}

	return nil
}
