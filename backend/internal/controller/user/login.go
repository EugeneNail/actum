package user

import (
	"encoding/json"
	"github.com/EugeneNail/actum/internal/model/users"
	"github.com/EugeneNail/actum/internal/service/controller"
	"github.com/EugeneNail/actum/internal/service/validation"
	"net/http"
)

type loginInput struct {
	Email    string `json:"email" rules:"required|email|max:100"`
	Password string `json:"password" rules:"required|min:8|max:100|mixedCase|regex:^\\S+$"`
}

func Login(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	input, err := controller.Parse[loginInput](writer, request)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	validationErrors, err := validation.Perform(input)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(validationErrors) > 0 {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		if err := encoder.Encode(validationErrors); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	user, err := users.FindBy("email", input.Email)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.Id == 0 || user.Password != hashPassword(input.Password) {
		writer.WriteHeader(http.StatusUnauthorized)
		err := encoder.Encode(map[string]string{"email": "Incorrect email address or password"})
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}
