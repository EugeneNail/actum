package records

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/infrastructure/log"
	"github.com/EugeneNail/actum/internal/infrastructure/middleware/routing"
	"github.com/EugeneNail/actum/internal/infrastructure/response"
	"github.com/EugeneNail/actum/internal/infrastructure/validation"
	"github.com/EugeneNail/actum/internal/service/auth/jwt"
	"net/http"
	"strconv"
)

type updateInput struct {
	Mood       int      `json:"mood" rules:"required|integer|min:1|max:5"`
	Weather    int      `json:"weather" rules:"required|integer|min:1|max:9"`
	Notes      string   `json:"notes" rules:"max:5000"`
	Activities []int    `json:"activities" rules:"required"`
	Photos     []string `json:"photos"`
}

func (controller *Controller) Update(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	id, err := strconv.Atoi(routing.GetVariable(request, 0))
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Update(): %w", err), http.StatusBadRequest)
		return
	}

	record, err := controller.recordDAO.Find(id)
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Update(): %w", err), http.StatusInternalServerError)
		return
	}

	if record.Id == 0 {
		response.Send(fmt.Sprintf("Запись %d не найдена", id), http.StatusNotFound)
		return
	}

	errors, input, err := validation.NewValidator[updateInput]().Validate(request)
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Update(): %w", err), http.StatusBadRequest)
		return
	}

	if len(errors) > 0 {
		response.Send(errors, http.StatusUnprocessableEntity)
		return
	}

	user := jwt.GetUser(request)
	allExist, missingActivities, err := controller.activityService.CheckExistence(input.Activities, user.Id)
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Update(): %w", err), http.StatusInternalServerError)
		return
	}

	if !allExist {
		errors := map[string]any{"activities": fmt.Sprintf("Активности %v не найдены", missingActivities)}
		response.Send(errors, http.StatusNotFound)
		return
	}

	allExist, missingPhotos, err := controller.photoService.CheckExistence(input.Photos, user.Id)
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	if !allExist {
		errors := map[string]any{"photos": fmt.Sprintf("Фотографии %v не найдены.", missingPhotos)}
		response.Send(errors, http.StatusNotFound)
		return
	}

	recordBefore := record
	record.Notes = input.Notes
	record.Mood = input.Mood
	record.Weather = input.Weather
	if err := controller.recordDAO.Save(&record); err != nil {
		response.Send(fmt.Errorf("RecordController.Update(): %w", err), http.StatusInternalServerError)
		return
	}

	if err = controller.activityService.SyncRelations(record.Id, input.Activities); err != nil {
		response.Send(fmt.Errorf("RecordController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	if err = controller.photoService.SyncRelations(record.Id, input.Photos); err != nil {
		response.Send(fmt.Errorf("RecordController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
	log.Info("User", user.Id, "changed record", fmt.Sprintf("%+v", recordBefore), "to", fmt.Sprintf("%+v", record))
}
