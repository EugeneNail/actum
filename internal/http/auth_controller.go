package controller

import (
	"context"
	"fmt"
	"github.com/EugeneNail/actum/internal/database/repository"
	"github.com/EugeneNail/actum/internal/infrastructure/errors"
	"github.com/EugeneNail/actum/internal/infrastructure/response"
	"github.com/EugeneNail/actum/internal/infrastructure/validation"
	"github.com/EugeneNail/actum/internal/service/auth"
	"net/http"
)

type AuthController struct {
	auth auth.Client
	repo *repository.UserRepository
}

func NewAuthController(auth auth.Client, repo *repository.UserRepository) AuthController {
	return AuthController{auth, repo}
}

type registerInput struct {
	Name                 string `json:"name" rules:"required|word|min:3|max:20"`
	Email                string `json:"email" rules:"required|email|max:100|unique:users,email"`
	Password             string `json:"password" rules:"required|min:8|max:100|mixedCase|regex:^\\S+$"`
	PasswordConfirmation string `json:"passwordConfirmation" rules:"required|min:8|max:100|mixedCase|regex:^\\S+$"`
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	resp := response.NewSender(w)

	vErrs, input, err := validation.NewValidator[registerInput]().Validate(r)
	if err != nil {
		resp.Send(errors.Wrap(err, "failed to validate"), http.StatusInternalServerError)
		return
	}

	if len(vErrs) > 0 {
		resp.Send(vErrs, http.StatusUnprocessableEntity)
		return
	}

	id, err := c.auth.Register(r.Context(), input.Name, input.Email, input.Password)
	if err != nil {
		resp.Send(errors.Wrap(err, "failed to register the user"), http.StatusInternalServerError)
		return
	}

	resp.Send(id, http.StatusCreated)
}

type loginInput struct {
	Email    string `json:"email" rules:"required|email|max:100"`
	Password string `json:"password" rules:"required|min:8|max:100|mixedCase|regex:^\\S+$"`
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	resp := response.NewSender(w)

	vErrs, input, err := validation.NewValidator[loginInput]().Validate(r)
	if err != nil {
		resp.Send(fmt.Errorf("UserController.Login(): %w", err), http.StatusBadRequest)
		return
	}

	if len(vErrs) > 0 {
		resp.Send(vErrs, http.StatusUnprocessableEntity)
		return
	}

	accessToken, refreshToken, err := c.auth.Login(context.TODO(), input.Email, input.Password)
	if err != nil {
		resp.Send(errors.Wrap(err, "failed to log in"), http.StatusInternalServerError)
		return
	}

	if accessToken == "" || refreshToken == "" {
		resp.Send(map[string]string{"email": "Неверные адрес почты или пароль."}, http.StatusUnprocessableEntity)
		return
	}

	resp.Send(map[string]string{"access": accessToken, "refresh": refreshToken}, http.StatusOK)
}
