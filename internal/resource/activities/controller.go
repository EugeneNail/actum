package activities

import (
	"database/sql"
	"fmt"
	"github.com/EugeneNail/actum/internal/resource/collections"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/middleware/routing"
	"github.com/EugeneNail/actum/internal/service/response"
	"github.com/EugeneNail/actum/internal/service/validation"
	"net/http"
	"strconv"
)

type Controller struct {
	db              *sql.DB
	activityDAO     *DAO
	collectionDAO   *collections.DAO
	activityService *Service
}

func NewController(db *sql.DB, activityDAO *DAO, collectionDAO *collections.DAO, activityService *Service) (controller Controller) {
	controller.db = db
	controller.activityDAO = activityDAO
	controller.collectionDAO = collectionDAO
	controller.activityService = activityService

	return
}

type storeInput struct {
	Name         string `json:"name" rules:"required|min:3|max:20|sentence"`
	Icon         int    `json:"icon" rules:"required|min:100|max:1000"`
	CollectionId int    `json:"collectionId" rules:"required|exists:collections,id"`
}

func (controller *Controller) Store(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)
	validator := validation.NewValidator[storeInput]()

	errors, input, err := validator.Validate(request)
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Store(): %w", err), http.StatusBadRequest)
		return
	}

	if len(errors) > 0 {
		response.Send(errors, http.StatusUnprocessableEntity)
		return
	}

	collection, err := controller.collectionDAO.Find(input.CollectionId)
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	if collection.Id == 0 {
		message := fmt.Sprintf("Коллекция %d не найдена.", input.CollectionId)
		response.Send(message, http.StatusNotFound)
		return
	}

	user := jwt.GetUser(request)
	if collection.UserId != user.Id {
		response.Send("Вы не можете добавить активность в чужую коллекцию.", http.StatusForbidden)
		return
	}

	const limit = 20
	exceedsLimit, err := controller.activityService.ExceedsLimit(limit, collection.Id, user.Id)
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	if exceedsLimit {
		message := fmt.Sprintf("Вы превысили лимит (%d) активностей для этой коллекции. Удалите старые активности или измените уже имеющуюся.", limit)
		response.Send(message, http.StatusConflict)
		return
	}

	hasDuplicate, err := controller.activityService.HasDuplicate(input.Name, user.Id)
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	if hasDuplicate {
		response.Send(map[string]string{"name": "У вас уже есть активность с таким именем."}, http.StatusConflict)
		return
	}

	activity := New(input.Name, input.Icon, input.CollectionId, user.Id)
	err = controller.activityDAO.Save(&activity)
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	response.Send(activity.Id, http.StatusCreated)
	log.Info("User", user.Id, "created activity", activity.Id, "of collection", collection.Id)
}

type updateInput struct {
	Name string `json:"name" rules:"required|min:3|max:20|sentence"`
	Icon int    `json:"icon" rules:"required|min:100|max:1000"`
}

func (controller *Controller) Update(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)
	validator := validation.NewValidator[updateInput]()

	id, err := strconv.Atoi(routing.GetVariable(request, 0))
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Update(): %w", err), http.StatusBadRequest)
		return
	}

	activity, err := controller.activityDAO.Find(id)
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Update(): %w", err), http.StatusInternalServerError)
		return
	}

	if activity.Id == 0 {
		message := fmt.Sprintf("Активность %d не найдена.", activity.Id)
		response.Send(message, http.StatusNotFound)
		return
	}

	user := jwt.GetUser(request)
	if activity.UserId != user.Id {
		response.Send("Вы не можете изменить чужую активность.", http.StatusForbidden)
		return
	}

	errors, input, err := validator.Validate(request)
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Update(): %w", err), http.StatusBadRequest)
		return
	}

	if len(errors) > 0 {
		response.Send(errors, http.StatusUnprocessableEntity)
		return
	}

	activityBefore := fmt.Sprintf("%+v", activity)
	activity.Name = input.Name
	activity.Icon = input.Icon
	if err := controller.activityDAO.Save(&activity); err != nil {
		response.Send(fmt.Errorf("ActivityController.Update(): %w", err), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
	log.Info("User", user.Id, "changed activity", activityBefore, "to", fmt.Sprintf("%+v", activity))
}

func (controller *Controller) Show(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	id, err := strconv.Atoi(routing.GetVariable(request, 0))
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Show(): %w", err), http.StatusBadRequest)
		return
	}

	activity, err := controller.activityDAO.Find(id)
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Show(): %w", err), http.StatusInternalServerError)
		return
	}

	if activity.Id == 0 {
		response.Send(nil, http.StatusNotFound)
		return
	}

	user := jwt.GetUser(request)
	if activity.UserId != user.Id {
		response.Send("Вы не можете использовать чужую активность.", http.StatusForbidden)
		return
	}

	response.Send(activity, http.StatusOK)
	log.Info("User", user.Id, "fetched activity", activity.Id)
}

func (controller *Controller) Destroy(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	id, err := strconv.Atoi(routing.GetVariable(request, 0))
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Destroy(): %w", err), http.StatusBadRequest)
		return
	}

	activity, err := controller.activityDAO.Find(id)
	if err != nil {
		response.Send(fmt.Errorf("ActivityController.Destroy(): %w", err), http.StatusInternalServerError)
		return
	}

	if activity.Id == 0 {
		message := fmt.Sprintf("Активность %d не найдена", activity.Id)
		response.Send(message, http.StatusNotFound)
		return
	}

	user := jwt.GetUser(request)
	if activity.UserId != user.Id {
		response.Send("Вы пытаетесь использовать чужую активность.", http.StatusForbidden)
		return
	}

	if err := controller.activityDAO.Delete(activity); err != nil {
		response.Send(fmt.Errorf("ActivityController.Destroy(): %w", err), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
	log.Info("User", user.Id, "deleted activity", activity.Id)

}
