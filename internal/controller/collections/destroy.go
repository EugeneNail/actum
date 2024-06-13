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

func (controller *Controller) Destroy(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	id, err := strconv.Atoi(routing.GetVariable(request, 0))
	if err != nil {
		response.Send(fmt.Errorf("UserController.Destroy(): %w", err), http.StatusBadRequest)
		return
	}

	collection, err := controller.dao.Find(id)
	if err != nil {
		response.Send(fmt.Errorf("UserController.Destroy(): %w", err), http.StatusInternalServerError)
		return
	}

	if collection.Id == 0 {
		response.Send(
			fmt.Sprintf("Коллекция %d не найдена.", id),
			http.StatusNotFound,
		)
		return
	}

	user := jwt.GetUser(request)
	if user.Id != collection.UserId {
		response.Send("Вы можете не удалить чужую коллекцию.", http.StatusForbidden)
		return
	}

	if err := controller.dao.Delete(collection); err != nil {
		response.Send(fmt.Errorf("UserController.Destroy(): %w", err), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
	log.Info("User", user.Id, "deleted collection", collection.Id)
}
