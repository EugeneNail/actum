package collection

import (
	"github.com/EugeneNail/actum/internal/controller"
	"github.com/EugeneNail/actum/internal/model/collections"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"net/http"
)

type storeInput struct {
	Name string `json:"name" rules:"required|min:3|max:20|sentence"`
}

func Store(writer http.ResponseWriter, request *http.Request) {
	controller := controller.New[storeInput](writer)

	isValid := controller.Validate(request)
	if !isValid {
		return
	}

	user := jwt.GetUser(request)
	hasDuplicateCollection, err := hasDuplicateCollection(controller.Input.Name, user)
	if err != nil {
		controller.Response(err, http.StatusInternalServerError)
		return
	}

	currentCollections, err := user.Collections()
	if err != nil {
		controller.Response(err, http.StatusInternalServerError)
		return
	}

	if len(currentCollections) >= 15 {
		controller.Response("You can have only 15 collections", http.StatusConflict)
		return
	}

	if hasDuplicateCollection {
		controller.Response(map[string]string{"name": "Collection already exists"}, http.StatusConflict)
		return
	}

	collection := collections.Collection{0, controller.Input.Name, user.Id}
	if err := collection.Save(); err != nil {
		controller.Response(err, http.StatusInternalServerError)
		return
	}

	controller.Response(collection.Id, http.StatusCreated)
	log.Info("User", user.Id, "created collection", collection.Id)
}
