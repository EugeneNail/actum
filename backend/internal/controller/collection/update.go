package collection

import (
	"encoding/json"
	"fmt"
	"github.com/EugeneNail/actum/internal/controller"
	"github.com/EugeneNail/actum/internal/model/collections"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/middleware/routing"
	"net/http"
	"strconv"
)

type updateInput struct {
	Name string `json:"name" rules:"required|min:3|max:50|sentence"`
}

func Update(writer http.ResponseWriter, request *http.Request) {
	encoder := json.NewEncoder(writer)

	id, err := strconv.Atoi(routing.GetVariable(request, 0))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}

	collection, err := collections.Find(id)
	if err != nil {
		controller.WriteError(writer, err)
		return
	}

	if collection.Id == 0 {
		writer.WriteHeader(http.StatusNotFound)
		message := fmt.Sprintf("Collection %d not found", id)
		if err := encoder.Encode(message); err != nil {
			controller.WriteError(writer, err)
		}
		return
	}

	user := jwt.GetUser(request)
	if user.Id != collection.UserId {
		writer.WriteHeader(http.StatusForbidden)
		if err := encoder.Encode("You are not allowed to manage other people's collections"); err != nil {
			controller.WriteError(writer, err)
		}
		return
	}

	input, ok := controller.GetInput[updateInput](writer, request, encoder)
	if !ok {
		return
	}

	collection.Name = input.Name
	if err := collection.Save(); err != nil {
		controller.WriteError(writer, err)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
	log.Info("User", user.Id, "changed the name of collection", collection.Id, "to", fmt.Sprintf(`"%s"`, input.Name))
}
