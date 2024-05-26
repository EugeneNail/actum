package activities

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/middleware/routing"
	"github.com/EugeneNail/actum/internal/service/response"
	"net/http"
	"strconv"
)

func (controller *Controller) Destroy(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	id, err := strconv.Atoi(routing.GetVariable(request, 0))
	if err != nil {
		response.Send(err, http.StatusBadRequest)
		return
	}

	activity, err := controller.activityRepo.Find(id)
	if err != nil {
		response.Send(err, http.StatusBadRequest)
		return
	}

	if activity.Id == 0 {
		message := fmt.Sprintf("Activity %d not found", activity.Id)
		response.Send(message, http.StatusNotFound)
		return
	}

	user := jwt.GetUser(request)
	if activity.UserId != user.Id {
		response.Send("You are not allowed to manage other people's activities", http.StatusForbidden)
		return
	}

	if err := controller.activityRepo.Delete(activity); err != nil {
		response.Send(err, http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
	log.Info("User", user.Id, "deleted activity", activity.Id)

}
