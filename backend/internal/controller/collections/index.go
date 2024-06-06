package collections

import (
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/response"
	"net/http"
)

func (controller *Controller) Index(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	user := jwt.GetUser(request)

	collections, err := controller.service.CollectCollections(user.Id)
	if err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	response.Send(collections, http.StatusOK)
	log.Info("User", user.Id, "indexed", len(collections), "collections")
}
