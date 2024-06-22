package users

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/service/auth/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/response"
	"github.com/EugeneNail/actum/internal/service/validation"
	"net/http"
)

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
