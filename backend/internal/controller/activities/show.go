package activities

import (
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/middleware/routing"
	"github.com/EugeneNail/actum/internal/service/response"
	"net/http"
	"strconv"
)

func (controller *Controller) Show(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	id, err := strconv.Atoi(routing.GetVariable(request, 0))
	if err != nil {
		response.Send(err, http.StatusBadRequest)
		return
	}

	activity, err := controller.activityDAO.Find(id)
	if err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	if activity.Id == 0 {
		response.Send(nil, http.StatusNotFound)
		return
	}

	user := jwt.GetUser(request)
	if activity.UserId != user.Id {
		response.Send("You are not allowed to view other people's activities", http.StatusForbidden)
		return
	}

	response.Send(activity, http.StatusOK)
	log.Info("User", user.Id, "fetched activity", activity.Id)
}