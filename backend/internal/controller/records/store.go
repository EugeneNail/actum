package records

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/database/resource/records"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/response"
	"github.com/EugeneNail/actum/internal/service/validation"
	"net/http"
)

type storeInput struct {
	Mood       int    `json:"mood" rules:"required|integer|min:1|max:5"`
	Weather    int    `json:"weather" rules:"required|integer|min:1|max:9"`
	Notes      string `json:"notes" rules:"max:5000"`
	Date       string `json:"date" rules:"required|date|today|after:2020-01-01"`
	Activities []int  `json:"activities" rules:"required"`
}

func (controller *Controller) Store(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	errors, input, err := validation.NewValidator[storeInput]().Validate(request)
	if err != nil {
		response.Send(err, http.StatusBadRequest)
		return
	}

	if len(errors) > 0 {
		response.Send(errors, http.StatusUnprocessableEntity)
		return
	}

	user := jwt.GetUser(request)
	isDateTaken, err := controller.recordService.IsDateTaken(input.Date, user.Id)
	if err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	if isDateTaken {
		message := fmt.Sprintf("У вас уже есть запись на дату %s.", input.Date)
		response.Send(map[string]any{"date": message}, http.StatusConflict)
		return
	}

	allExist, missingActivities, err := controller.activityService.CheckExistence(input.Activities, user.Id)
	if err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	if !allExist {
		errors := map[string]any{"activities": fmt.Sprintf("Активности %v не найдены.", missingActivities)}
		response.Send(errors, http.StatusNotFound)
		return
	}

	record, err := records.New(input.Mood, input.Weather, input.Date, input.Notes, user.Id)
	if err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	if err := controller.recordDAO.Save(&record); err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	if err = controller.recordDAO.SyncRelations(record.Id, input.Activities); err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	response.Send(record.Id, http.StatusCreated)
	log.Info("User", user.Id, "created record", record.Id, "for date", input.Date, "with", len(input.Activities), "activities")
}
