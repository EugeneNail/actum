package activities

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/infrastructure/log"
	"github.com/EugeneNail/actum/internal/infrastructure/middleware/routing"
	"github.com/EugeneNail/actum/internal/infrastructure/response"
	"github.com/EugeneNail/actum/internal/service/auth/jwt"
	"net/http"
	"strconv"
)

func (controller *Controller) Destroy(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	id, err := strconv.Atoi(routing.GetVariable(request, 0))
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Destroy(): %w", err), http.StatusBadRequest)
		return
	}

	activity, err := controller.activityDAO.Find(id)
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Destroy(): %w", err), http.StatusInternalServerError)
		return
	}

	if activity.Id == 0 {
		message := fmt.Sprintf("Активность %d не найдена", activity.Id)
		response.Send(message, http.StatusNotFound)
		return
	}

	user := jwt.GetUser(request)
	if activity.UserId != user.Id {
		response.Send("Вы пытаетесь использовать чужую активность.", http.StatusForbidden)
		return
	}

	if err := controller.activityDAO.Delete(activity); err != nil {
		response.Send(fmt.Errorf("ActivityController.Destroy(): %w", err), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
	log.Info("User", user.Id, "deleted activity", activity.Id)

}
