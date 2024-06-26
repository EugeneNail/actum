package collections

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/middleware/routing"
	"github.com/EugeneNail/actum/internal/service/response"
	"net/http"
	"strconv"
)

func (controller *Controller) Show(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	id, err := strconv.Atoi(routing.GetVariable(request, 0))
	if err != nil {
		response.Send(fmt.Errorf("CollectionController.Show(): %w", err), http.StatusBadRequest)
		return
	}

	collection, err := controller.dao.Find(id)
	if err != nil {
		response.Send(fmt.Errorf("CollectionController.Show(): %w", err), http.StatusInternalServerError)
		return
	}

	if collection.Id == 0 {
		response.Send(nil, http.StatusNotFound)
		return
	}

	user := jwt.GetUser(request)
	if collection.UserId != user.Id {
		response.Send("Вы не можете использовать чужую коллекцию.", http.StatusForbidden)
		return
	}

	response.Send(collection, http.StatusOK)
	log.Info("User", user.Id, "fetched collection", collection.Id)
}
