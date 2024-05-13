package user

import (
	"github.com/EugeneNail/actum/internal/controller"
	"github.com/EugeneNail/actum/internal/model/users"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"net/http"
	"strings"
)

type storeInput struct {
	Name                 string `json:"name" rules:"required|word|min:3|max:20"`
	Email                string `json:"email" rules:"required|email|max:100|unique:users,email"`
	Password             string `json:"password" rules:"required|min:8|max:100|mixedCase|regex:^\\S+$"`
	PasswordConfirmation string `json:"passwordConfirmation" rules:"required|min:8|max:100|mixedCase|regex:^\\S+$"`
}

func Store(writer http.ResponseWriter, request *http.Request) {
	controller := controller.New[storeInput](writer)

	isValid := controller.Validate(request)
	if !isValid {
		return
	}

	if controller.Input.Password != controller.Input.PasswordConfirmation {
		controller.Response(map[string]string{"passwordConfirmation": "Passwords do not match"}, http.StatusUnprocessableEntity)
		return
	}

	user := users.New(
		controller.Input.Name,
		strings.ToLower(controller.Input.Email),
		hashPassword(controller.Input.Password),
	)
	if err := user.Save(); err != nil {
		controller.Response(err, http.StatusUnprocessableEntity)
		return
	}

	token, err := jwt.Make(user)
	if err != nil {
		controller.Response(err, http.StatusInternalServerError)
		return
	}

	controller.Response(token, http.StatusCreated)
	log.Info("Created user", user.Id)
}
