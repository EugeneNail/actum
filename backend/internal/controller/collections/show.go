package collections

import (
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
		response.Send(err, http.StatusBadRequest)
		return
	}

	collection, err := controller.repository.Find(id)
	if err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	if collection.Id == 0 {
		response.Send(nil, http.StatusNotFound)
		return
	}

	user := jwt.GetUser(request)
	if collection.UserId != user.Id {
		response.Send("You are not allowed to view other people's collections", http.StatusForbidden)
		return
	}

	response.Send(collection, http.StatusOK)
	log.Info("User", user.Id, "fetched collection", collection.Id)
}
