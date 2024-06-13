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
	Icon         int    `json:"icon" rules:"required|min:100|max:1000"`
	CollectionId int    `json:"collectionId" rules:"required|exists:collections,id"`
}

func (controller *Controller) Store(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)
	validator := validation.NewValidator[storeInput]()

	errors, input, err := validator.Validate(request)
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Store(): %w", err), http.StatusBadRequest)
		return
	}

	if len(errors) > 0 {
		response.Send(errors, http.StatusUnprocessableEntity)
		return
	}

	collection, err := controller.collectionDAO.Find(input.CollectionId)
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	if collection.Id == 0 {
		message := fmt.Sprintf("Коллекция %d не найдена.", input.CollectionId)
		response.Send(message, http.StatusNotFound)
		return
	}

	user := jwt.GetUser(request)
	if collection.UserId != user.Id {
		response.Send("Вы не можете добавить активность в чужую коллекцию.", http.StatusForbidden)
		return
	}

	const limit = 20
	exceedsLimit, err := controller.activityService.ExceedsLimit(limit, collection.Id, user.Id)
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	if exceedsLimit {
		message := fmt.Sprintf("Вы превысили лимит (%d) активностей для этой коллекции. Удалите старые активности или измените уже имеющуюся.", limit)
		response.Send(message, http.StatusConflict)
		return
	}

	hasDuplicate, err := controller.activityService.HasDuplicate(input.Name, user.Id)
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	if hasDuplicate {
		response.Send(map[string]string{"name": "У вас уже есть активность с таким именем."}, http.StatusConflict)
		return
	}

	activity := activities.New(input.Name, input.Icon, input.CollectionId, user.Id)
	err = controller.activityDAO.Save(&activity)
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	response.Send(activity.Id, http.StatusCreated)
	log.Info("User", user.Id, "created activity", activity.Id, "of collection", collection.Id)
}
