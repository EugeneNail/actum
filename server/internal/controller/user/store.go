package user

import (
	"encoding/json"
	"github.com/EugeneNail/actum/internal/service/controller"
	"github.com/EugeneNail/actum/internal/service/validation"
	"net/http"
)

type input struct {
	Name                 string `json:"name" rules:"required|min:3|max:20"`
	Email                string `json:"email" rules:"required|email|max:100"`
	Password             string `json:"password" rules:"required|min:8|max:100"`
	PasswordConfirmation string `json:"passwordConfirmation" rules:"required|min:8|max:100"`
}

func Store(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	input, err := controller.Parse[input](writer, request)

	if err != nil {
		return
	}

	if errors, err := validation.New(input); err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)

		if err := encoder.Encode(errors); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	}
}
