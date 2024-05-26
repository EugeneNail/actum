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

func Destroy(writer http.ResponseWriter, request *http.Request) {
	controller := controller.New[any](writer)
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

	if err := collection.Delete(); err != nil {
		controller.Response(err, http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
	log.Info("User", user.Id, "deleted collection", collection.Id)
}
