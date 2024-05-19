package collection

import (
	"github.com/EugeneNail/actum/internal/controller"
	"github.com/EugeneNail/actum/internal/model/collections"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/middleware/routing"
	"net/http"
	"strconv"
)

func Show(writer http.ResponseWriter, request *http.Request) {
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
		controller.Response(nil, http.StatusNotFound)
		return
	}

	user := jwt.GetUser(request)
	if collection.UserId != user.Id {
		controller.Response("You are not allowed to view other people's collections", http.StatusForbidden)
		return
	}

	controller.Response(collection, http.StatusOK)
	log.Info("User", user.Id, "fetched collection", collection.Id)
}
