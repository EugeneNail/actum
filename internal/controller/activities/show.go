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

func (controller *Controller) Show(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	id, err := strconv.Atoi(routing.GetVariable(request, 0))
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Show(): %w", err), http.StatusBadRequest)
		return
	}

	activity, err := controller.activityDAO.Find(id)
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Show(): %w", err), http.StatusInternalServerError)
		return
	}

	if activity.Id == 0 {
		response.Send(nil, http.StatusNotFound)
		return
	}

	user := jwt.GetUser(request)
	if activity.UserId != user.Id {
		response.Send("Вы не можете использовать чужую активность.", http.StatusForbidden)
		return
	}

	response.Send(activity, http.StatusOK)
	log.Info("User", user.Id, "fetched activity", activity.Id)
}
