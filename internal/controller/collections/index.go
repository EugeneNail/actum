package collections

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/infrastructure/log"
	"github.com/EugeneNail/actum/internal/infrastructure/response"
	"github.com/EugeneNail/actum/internal/service/auth/jwt"
	"net/http"
)

func (controller *Controller) Index(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	user := jwt.GetUser(request)

	collections, err := controller.service.CollectCollections(user.Id)
	if err != nil {
		response.Send(fmt.Errorf("CollectionController.index(): %w", err), http.StatusInternalServerError)
		return
	}

	response.Send(collections, http.StatusOK)
	log.Info("User", user.Id, "indexed", len(collections), "collections")
}
