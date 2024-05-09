package user

import (
	"encoding/json"
	"github.com/EugeneNail/actum/internal/controller"
	"github.com/EugeneNail/actum/internal/model/users"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"net/http"
)

type storeInput struct {
	Name                 string `json:"name" rules:"required|word|min:3|max:20"`
	Email                string `json:"email" rules:"required|email|max:100|unique:users,email"`
	Password             string `json:"password" rules:"required|min:8|max:100|mixedCase|regex:^\\S+$"`
	PasswordConfirmation string `json:"passwordConfirmation" rules:"required|min:8|max:100|mixedCase|regex:^\\S+$"`
}

func Store(writer http.ResponseWriter, request *http.Request) {
	encoder := json.NewEncoder(writer)

	input, ok := controller.GetInput[storeInput](writer, request, encoder)
	if !ok {
		return
	}

	if input.Password != input.PasswordConfirmation {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		err := encoder.Encode(map[string]string{"passwordConfirmation": "Passwords do not match"})

		if err != nil {
			controller.WriteError(writer, err)
		}
		return
	}

	user := users.New(input.Name, input.Email, hashPassword(input.Password))
	if err := user.Save(); err != nil {
		controller.WriteError(writer, err)
		return
	}

	token, err := jwt.Make(user)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	if err := encoder.Encode(token); err != nil {
		controller.WriteError(writer, err)
		return
	}
	log.Info("Created user", user.Id)
}
