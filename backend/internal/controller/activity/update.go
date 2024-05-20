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

type updateInput struct {
	Name string `json:"name" rules:"required|min:3|max:20|sentence"`
	Icon string `json:"icon" rules:"required|max:25|regex:^[0-9a-zA-Z_]+$"`
}

func Update(writer http.ResponseWriter, request *http.Request) {
	controller := controller.New[updateInput](writer)

	id, err := strconv.Atoi(routing.GetVariable(request, 0))
	if err != nil {
		controller.Response(err, http.StatusBadRequest)
		return
	}

	activity, err := activities.Find(id)
	if err != nil {
		controller.Response(err, http.StatusInternalServerError)
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

	if !controller.Validate(request) {
		return
	}

	activityBefore := fmt.Sprintf("%+v", activity)
	activity.Name = controller.Input.Name
	activity.Icon = controller.Input.Icon
	if err := activity.Save(); err != nil {
		controller.Response(err, http.StatusInternalServerError)
		return
	}

	activityAfter := fmt.Sprintf("%+v", activity)
	writer.WriteHeader(http.StatusNoContent)
	log.Info("User", user.Id, "changed activity", activityBefore, "to", activityAfter)
}
