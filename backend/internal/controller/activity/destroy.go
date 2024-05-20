package activity

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/controller"
	"github.com/EugeneNail/actum/internal/model/activities"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/middleware/routing"
	"net/http"
	"strconv"
)

func Destroy(writer http.ResponseWriter, request *http.Request) {
	controller := controller.New[any](writer)
	id, err := strconv.Atoi(routing.GetVariable(request, 0))
	if err != nil {
		controller.Response(err, http.StatusBadRequest)
		return
	}

	activity, err := activities.Find(id)
	if err != nil {
		controller.Response(err, http.StatusBadRequest)
		return
	}

	if activity.Id == 0 {
		message := fmt.Sprintf("Activity %d not found", activity.Id)
		controller.Response(message, http.StatusNotFound)
		return
	}

	user := jwt.GetUser(request)
	if activity.UserId != user.Id {
		controller.Response("You are not allowed to manage other people's activities", http.StatusForbidden)
		return
	}

	if err := activity.Delete(); err != nil {
		controller.Response(err, http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
	log.Info("User", user.Id, "deleted activity", activity.Id)

}
