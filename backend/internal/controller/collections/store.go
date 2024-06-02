package collections

import (
	"github.com/EugeneNail/actum/internal/database/resource/collections"
	"github.com/EugeneNail/actum/internal/service/jwt"
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
		response.Send(err, http.StatusBadRequest)
		return
	}

	if len(errors) > 0 {
		response.Send(errors, http.StatusUnprocessableEntity)
		return
	}

	user := jwt.GetUser(request)
	exceededLimit, err := controller.exceededLimit(user.Id)
	if err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	if exceededLimit {
		response.Send("You can have only 15 collections", http.StatusConflict)
		return
	}

	hasDuplicate, err := controller.hasDuplicate(input.Name, user)
	if err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	if hasDuplicate {
		response.Send(map[string]string{"name": "Collection already exists"}, http.StatusConflict)
		return
	}

	collection := collections.New(input.Name, input.Color, user.Id)
	if err := controller.dao.Save(&collection); err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	response.Send(collection.Id, http.StatusCreated)
	log.Info("User", user.Id, "created collection", collection.Id)
}
