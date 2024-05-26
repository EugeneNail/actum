package activities

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/database/resource/activities"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/response"
	"github.com/EugeneNail/actum/internal/service/validation"
	"net/http"
)

type storeInput struct {
	Name         string `json:"name" rules:"required|min:3|max:20|sentence"`
	Icon         string `json:"icon" rules:"required|max:25|regex:^[0-9a-zA-Z_]+$"`
	CollectionId int    `json:"collectionId" rules:"required|exists:collections,id"`
}

func (controller *Controller) Store(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)
	validator := validation.NewValidator[storeInput]()

	errors, input, err := validator.Validate(request)
	if err != nil {
		response.Send(err, http.StatusBadRequest)
		return
	}

	if len(errors) > 0 {
		response.Send(errors, http.StatusUnprocessableEntity)
		return
	}

	collection, err := controller.collectionRepo.Find(input.CollectionId)
	if err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	if collection.Id == 0 {
		message := fmt.Sprintf("Collection %d not found", input.CollectionId)
		response.Send(message, http.StatusNotFound)
		return
	}

	user := jwt.GetUser(request)
	if collection.UserId != user.Id {
		response.Send("You are not allowed to manage other people's collections", http.StatusForbidden)
		return
	}

	exceededLimit, err := controller.exceededLimit(collection.Id, user.Id)
	if err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	if exceededLimit {
		response.Send("You can have only 20 activities per collection", http.StatusConflict)
		return
	}

	hasDuplicate, err := controller.hasDuplicate(input.Name, user.Id)
	if err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	if hasDuplicate {
		response.Send(map[string]string{"name": "Activity already exists"}, http.StatusConflict)
		return
	}

	activity := activities.New(input.Name, input.Icon, input.CollectionId, user.Id)

	err = controller.activityRepo.Save(&activity)
	if err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	response.Send(activity.Id, http.StatusCreated)
	log.Info("User", user.Id, "created activity", activity.Id, "of collection", collection.Id)
}
