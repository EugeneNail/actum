package collection

import (
	"github.com/EugeneNail/actum/internal/controller"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"net/http"
)

func Index(writer http.ResponseWriter, request *http.Request) {
	controller := controller.New[any](writer)
	user := jwt.GetUser(request)

	collections, err := getUserCollections(user.Id)
	if err != nil {
		controller.Response(err, http.StatusInternalServerError)
		return
	}

	controller.Response(collections, http.StatusOK)
	log.Info("User", user.Id, "indexed", len(collections), "collections")
}
