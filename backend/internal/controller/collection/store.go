package collection

import (
	"encoding/json"
	"github.com/EugeneNail/actum/internal/controller"
	"github.com/EugeneNail/actum/internal/model/collections"
	"github.com/EugeneNail/actum/internal/model/users"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"net/http"
)

type storeInput struct {
	Name string `json:"name" rules:"required|min:3|max:50|sentence"`
}

func Store(writer http.ResponseWriter, request *http.Request) {
	encoder := json.NewEncoder(writer)

	input, ok := controller.GetInput[storeInput](writer, request, encoder)
	if !ok {
		return
	}

	user := jwt.GetUser(request)
	hasDuplicateCollection, err := hasDuplicateCollection(input.Name, user)
	if err != nil {
		controller.WriteError(writer, err)
		return
	}

	if hasDuplicateCollection {
		writer.WriteHeader(http.StatusConflict)

		err := encoder.Encode(map[string]string{"name": "Collection already exists"})
		if err != nil {
			controller.WriteError(writer, err)
		}
		return
	}

	collection := collections.Collection{0, input.Name, user.Id}
	if err := collection.Save(); err != nil {
		controller.WriteError(writer, err)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	if err := encoder.Encode(collection.Id); err != nil {
		controller.WriteError(writer, err)
	}
}

func hasDuplicateCollection(name string, user users.User) (bool, error) {
	collections, err := user.Collections()
	if err != nil {
		return false, err
	}

	for _, collection := range collections {
		if collection.Name == name {
			return true, nil
		}
	}

	return false, nil
}
