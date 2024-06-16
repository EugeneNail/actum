package collections

import (
	"database/sql"
	"fmt"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/middleware/routing"
	"github.com/EugeneNail/actum/internal/service/response"
	"github.com/EugeneNail/actum/internal/service/validation"
	"net/http"
	"strconv"
)

type Controller struct {
	db      *sql.DB
	dao     *DAO
	service *Service
}

func NewController(db *sql.DB, dao *DAO, service *Service) Controller {
	return Controller{db, dao, service}
}

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

	collection := New(input.Name, input.Color, user.Id)
	if err := controller.dao.Save(&collection); err != nil {
		response.Send(fmt.Errorf("CollectionController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	response.Send(collection.Id, http.StatusCreated)
	log.Info("User", user.Id, "created collection", collection.Id)
}

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
