package collection

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/controller"
	"github.com/EugeneNail/actum/internal/resource/collections"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/middleware/routing"
	"net/http"
	"strconv"
)

type updateInput struct {
	Name string `json:"name" rules:"required|min:3|max:20|sentence"`
}

func Update(writer http.ResponseWriter, request *http.Request) {
	controller := controller.New[updateInput](writer)

	id, err := strconv.Atoi(routing.GetVariable(request, 0))
	if err != nil {
		controller.Response(err, http.StatusBadRequest)
		return
	}

	collection, err := collections.Find(id)
	if err != nil {
		controller.Response(err, http.StatusInternalServerError)
		return
	}

	if collection.Id == 0 {
		controller.Response(fmt.Sprintf("Collection %d not found", id), http.StatusNotFound)
		return
	}

	user := jwt.GetUser(request)
	if user.Id != collection.UserId {
		controller.Response("You are not allowed to manage other people's collections", http.StatusForbidden)
		return
	}

	isValid := controller.Validate(request)
	if !isValid {
		return
	}

	hasDuplicateCollection, err := hasDuplicateCollection(controller.Input.Name, user)
	if err != nil {
		controller.Response(err, http.StatusInternalServerError)
		return
	}

	if hasDuplicateCollection {
		controller.Response(map[string]string{"name": "Collection already exists"}, http.StatusConflict)
		return
	}

	collection.Name = controller.Input.Name
	if err := collection.Save(); err != nil {
		controller.Response(err, http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
	log.Info("User", user.Id, "changed the name of collection", collection.Id, "to", fmt.Sprintf(`"%s"`, controller.Input.Name))
}
