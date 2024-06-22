package users

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/service/auth/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/response"
	"net/http"
)

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
