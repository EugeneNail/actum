package users

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/database/resource/users"
	"github.com/EugeneNail/actum/internal/service/auth/jwt"
	"github.com/EugeneNail/actum/internal/service/hash"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/response"
	"github.com/EugeneNail/actum/internal/service/validation"
	"net/http"
	"strings"
)

type storeInput struct {
	Name                 string `json:"name" rules:"required|word|min:3|max:20"`
	Email                string `json:"email" rules:"required|email|max:100|unique:users,email"`
	Password             string `json:"password" rules:"required|min:8|max:100|mixedCase|regex:^\\S+$"`
	PasswordConfirmation string `json:"passwordConfirmation" rules:"required|min:8|max:100|mixedCase|regex:^\\S+$"`
}

func (controller *Controller) Store(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)
	validator := validation.NewValidator[storeInput]()

	errors, input, err := validator.Validate(request)
	if err != nil {
		response.Send(fmt.Errorf("UserController.Store(): %w", err), http.StatusBadRequest)
		return
	}

	if len(errors) > 0 {
		response.Send(errors, http.StatusUnprocessableEntity)
		return
	}

	if input.Password != input.PasswordConfirmation {
		response.Send(map[string]string{"passwordConfirmation": "Пароли не совпадают."}, http.StatusUnprocessableEntity)
		return
	}

	user := users.New(
		input.Name,
		strings.ToLower(input.Email),
		hash.New(input.Password),
	)
	if err := controller.dao.Save(&user); err != nil {
		response.Send(fmt.Errorf("UserController.Store(): %w", err), http.StatusUnprocessableEntity)
		return
	}

	accessToken, err := jwt.Make(user.Id)
	if err != nil {
		response.Send(fmt.Errorf("UserController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	refreshToken, err := controller.refreshService.MakeToken(user.Id)
	if err != nil {
		response.Send(fmt.Errorf("UserController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	response.Send(map[string]string{
		"access":  accessToken,
		"refresh": refreshToken,
	}, http.StatusCreated)
	log.Info("Created user", user.Id)
}
