package collections

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/service/auth/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/middleware/routing"
	"github.com/EugeneNail/actum/internal/service/response"
	"github.com/EugeneNail/actum/internal/service/validation"
	"net/http"
	"strconv"
)

type updateInput struct {
	Name  string `json:"name" rules:"required|min:3|max:20|sentence"`
	Color int    `json:"color" rules:"required|min:1|max:6"`
}

func (controller *Controller) Update(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)
	validator := validation.NewValidator[updateInput]()

	id, err := strconv.Atoi(routing.GetVariable(request, 0))
	if err != nil {
		response.Send(fmt.Errorf("CollectionController.Update(): %w", err), http.StatusBadRequest)
		return
	}

	collection, err := controller.dao.Find(id)
	if err != nil {
		response.Send(fmt.Errorf("CollectionController.Update(): %w", err), http.StatusInternalServerError)
		return
	}

	if collection.Id == 0 {
		response.Send(fmt.Sprintf("Коллекция %d не найдена.", id), http.StatusNotFound)
		return
	}

	user := jwt.GetUser(request)
	if user.Id != collection.UserId {
		response.Send("Мы не можете изменить чужую коллекцию.", http.StatusForbidden)
		return
	}

	errors, input, err := validator.Validate(request)
	if err != nil {
		response.Send(fmt.Errorf("CollectionController.Update(): %w", err), http.StatusBadRequest)
		return
	}

	if len(errors) > 0 {
		response.Send(errors, http.StatusUnprocessableEntity)
		return
	}

	collection.Name = input.Name
	collection.Color = input.Color
	if err := controller.dao.Save(&collection); err != nil {
		response.Send(fmt.Errorf("CollectionController.Update(): %w", err), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
	log.Info("User", user.Id, "renamed collection", collection.Id, "to", fmt.Sprintf(`"%s"`, input.Name))
}
