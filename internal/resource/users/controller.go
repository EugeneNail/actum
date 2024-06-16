package users

import (
	"database/sql"
	"fmt"
	"github.com/EugeneNail/actum/internal/service/hash"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/refresh"
	"github.com/EugeneNail/actum/internal/service/response"
	"github.com/EugeneNail/actum/internal/service/validation"
	"net/http"
	"strings"
)

type Controller struct {
	db             *sql.DB
	dao            *DAO
	refreshService *refresh.Service
}

func NewController(db *sql.DB, dao *DAO, refreshService *refresh.Service) Controller {
	return Controller{db, dao, refreshService}
}

type loginInput struct {
	Email    string `json:"email" rules:"required|email|max:100"`
	Password string `json:"password" rules:"required|min:8|max:100|mixedCase|regex:^\\S+$"`
}
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

	user := New(
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
func (controller *Controller) Logout(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	user := jwt.GetUser(request)
	if err := controller.refreshService.Unset(user.Id); err != nil {
		response.Send(fmt.Errorf("UserController.Logout(): %w", err), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
	log.Info("User", user.Id, "logged out")

}

type refreshInput struct {
	Uuid   string `json:"uuid" rules:"required|min:36|max:36|regex:^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"`
	UserId int    `json:"userId" rules:"required|exists:users,id"`
}

func (controller *Controller) RefreshToken(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	errors, input, err := validation.NewValidator[refreshInput]().Validate(request)
	if err != nil {
		response.Send(fmt.Errorf("UserController.Refresh(): %w", err), http.StatusBadRequest)
		return
	}

	if len(errors) > 0 {
		response.Send(errors, http.StatusUnprocessableEntity)
		return
	}

	isValid, err := controller.refreshService.IsValid(input.Uuid, input.UserId)
	if err != nil {
		response.Send(fmt.Errorf("UserController.Refresh(): %w", err), http.StatusInternalServerError)
		return
	}

	if !isValid {
		response.Send("Токен обновления неправильный или истек", http.StatusUnauthorized)
		return
	}

	jwt, err := jwt.Make(input.UserId)
	if err != nil {
		response.Send(fmt.Errorf("UserController.Refresh(): %w", err), http.StatusInternalServerError)
		return
	}

	response.Send(jwt, http.StatusOK)
	log.Info("User", input.UserId, "refreshed token")
	return
}
