package activities

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/infrastructure/log"
	"github.com/EugeneNail/actum/internal/infrastructure/middleware/routing"
	"github.com/EugeneNail/actum/internal/infrastructure/response"
	"github.com/EugeneNail/actum/internal/infrastructure/validation"
	"github.com/EugeneNail/actum/internal/service/auth/jwt"
	"net/http"
	"strconv"
)

type updateInput struct {
	Name string `json:"name" rules:"required|min:3|max:20|sentence"`
	Icon int    `json:"icon" rules:"required|min:100|max:1000"`
}

func (controller *Controller) Update(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)
	validator := validation.NewValidator[updateInput]()

	id, err := strconv.Atoi(routing.GetVariable(request, 0))
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Update(): %w", err), http.StatusBadRequest)
		return
	}

	activity, err := controller.activityDAO.Find(id)
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Update(): %w", err), http.StatusInternalServerError)
		return
	}

	if activity.Id == 0 {
		message := fmt.Sprintf("Активность %d не найдена.", activity.Id)
		response.Send(message, http.StatusNotFound)
		return
	}

	user := jwt.GetUser(request)
	if activity.UserId != user.Id {
		response.Send("Вы не можете изменить чужую активность.", http.StatusForbidden)
		return
	}

	errors, input, err := validator.Validate(request)
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Update(): %w", err), http.StatusBadRequest)
		return
	}

	if len(errors) > 0 {
		response.Send(errors, http.StatusUnprocessableEntity)
		return
	}

	activityBefore := fmt.Sprintf("%+v", activity)
	activity.Name = input.Name
	activity.Icon = input.Icon
	if err := controller.activityDAO.Save(&activity); err != nil {
		response.Send(fmt.Errorf("ActivityController.Update(): %w", err), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
	log.Info("User", user.Id, "changed activity", activityBefore, "to", fmt.Sprintf("%+v", activity))
}
