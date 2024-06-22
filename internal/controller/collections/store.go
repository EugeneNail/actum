package collections

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/database/resource/collections"
	"github.com/EugeneNail/actum/internal/service/auth/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/response"
	"github.com/EugeneNail/actum/internal/service/validation"
	"net/http"
)

type storeInput struct {
	Name  string `json:"name" rules:"required|min:3|max:20|sentence"`
	Color int    `json:"color" rules:"required|min:1|max:6"`
}

func (controller *Controller) Store(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)
	validator := validation.NewValidator[storeInput]()

	errors, input, err := validator.Validate(request)
	if err != nil {
		response.Send(fmt.Errorf("CollectionController.Store(): %w", err), http.StatusBadRequest)
		return
	}

	if len(errors) > 0 {
		response.Send(errors, http.StatusUnprocessableEntity)
		return
	}

	const limit = 15
	user := jwt.GetUser(request)
	exceededLimit, err := controller.service.ExceedsLimit(limit, user.Id)
	if err != nil {
		response.Send(fmt.Errorf("CollectionController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	if exceededLimit {
		message := fmt.Sprintf("Похоже, вы превысили лимит (%d) коллекций. Попробуйте удалить старые коллекции или изменить уже имеющиеся.", limit)
		response.Send(message, http.StatusConflict)
		return
	}

	hasDuplicate, err := controller.service.HasDuplicate(input.Name, user.Id)
	if err != nil {
		response.Send(fmt.Errorf("CollectionController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	if hasDuplicate {
		response.Send(map[string]string{"name": "Collection already exists"}, http.StatusConflict)
		return
	}

	collection := collections.New(input.Name, input.Color, user.Id)
	if err := controller.dao.Save(&collection); err != nil {
		response.Send(fmt.Errorf("CollectionController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	response.Send(collection.Id, http.StatusCreated)
	log.Info("User", user.Id, "created collection", collection.Id)
}
