package users

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/service/auth/jwt"
	"github.com/EugeneNail/actum/internal/service/hash"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/response"
	"github.com/EugeneNail/actum/internal/service/validation"
	"net/http"
	"strings"
)

type loginInput struct {
	Email    string `json:"email" rules:"required|email|max:100"`
	Password string `json:"password" rules:"required|min:8|max:100|mixedCase|regex:^\\S+$"`
}

func (controller *Controller) Login(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)
	validator := validation.NewValidator[loginInput]()

	errors, input, err := validator.Validate(request)
	if err != nil {
		response.Send(fmt.Errorf("UserController.Login(): %w", err), http.StatusBadRequest)
		return
	}

	if len(errors) > 0 {
		response.Send(errors, http.StatusUnprocessableEntity)
		return
	}

	user, err := controller.dao.FindBy("email", strings.ToLower(input.Email))
	if err != nil {
		response.Send(fmt.Errorf("UserController.Login(): %w", err), http.StatusInternalServerError)
		return
	}

	if user.Id == 0 || user.Password != hash.New(input.Password) {
		response.Send(map[string]string{"email": "Неверные адрес почты или пароль."}, http.StatusUnauthorized)
		return
	}

	accessToken, err := jwt.Make(user.Id)
	if err != nil {
		response.Send(fmt.Errorf("UserController.Login(): %w", err), http.StatusInternalServerError)
		return
	}

	refreshToken, err := controller.refreshService.MakeToken(user.Id)
	if err != nil {
		response.Send(fmt.Errorf("UserController.Login(): %w", err), http.StatusInternalServerError)
		return
	}

	response.Send(map[string]string{
		"access":  accessToken,
		"refresh": refreshToken,
	}, http.StatusOK)
	log.Info("User", user.Id, "logged in")
}
