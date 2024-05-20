package activity

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/controller"
	"github.com/EugeneNail/actum/internal/model/activities"
	"github.com/EugeneNail/actum/internal/model/collections"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"net/http"
)

type storeInput struct {
	Name         string `json:"name" rules:"required|min:3|max:20|sentence"`
	Icon         string `json:"icon" rules:"required|max:25|regex:^[0-9a-zA-Z_]+$"`
	CollectionId int    `json:"collectionId" rules:"required|exists:collections,id"`
}

func Store(writer http.ResponseWriter, request *http.Request) {
	controller := controller.New[storeInput](writer)
	isValid := controller.Validate(request)
	if !isValid {
		return
	}

	collection, err := collections.Find(controller.Input.CollectionId)
	if err != nil {
		controller.Response(err, http.StatusInternalServerError)
		return
	}

	if collection.Id == 0 {
		message := fmt.Sprintf("Collection %d not found", controller.Input.CollectionId)
		controller.Response(message, http.StatusNotFound)
		return
	}

	user := jwt.GetUser(request)
	if collection.UserId != user.Id {
		controller.Response("You are not allowed to manage other people's collections", http.StatusForbidden)
		return
	}

	currentActivities, err := getCollectionActivities(collection.Id)
	if err != nil {
		controller.Response(err, http.StatusInternalServerError)
		return
	}

	if len(currentActivities) >= 20 {
		controller.Response("You can have only 20 activities per collection", http.StatusConflict)
		return
	}

	hasDuplicate, err := hasDuplicateActivity(controller.Input.Name, controller.Input.CollectionId)
	if err != nil {
		controller.Response(err, http.StatusInternalServerError)
		return
	}

	if hasDuplicate {
		controller.Response(map[string]string{"name": "Activity already exists"}, http.StatusConflict)
		return
	}

	activity := activities.Activity{
		0,
		controller.Input.Name,
		controller.Input.Icon,
		controller.Input.CollectionId,
		user.Id,
	}

	err = activity.Save()
	if err != nil {
		controller.Response(err, http.StatusInternalServerError)
		return
	}

	controller.Response(activity.Id, http.StatusCreated)
	log.Info("User", user.Id, "created activity", activity.Id, "of collection", collection.Id)
}
