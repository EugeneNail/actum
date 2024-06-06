package records

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/middleware/routing"
	"github.com/EugeneNail/actum/internal/service/response"
	"github.com/EugeneNail/actum/internal/service/validation"
	"net/http"
	"strconv"
)

type updateInput struct {
	Mood       int    `json:"mood" rules:"required|integer|min:1|max:5"`
	Notes      string `json:"notes" rules:"max:5000"`
	Activities []int  `json:"activities" rules:"required"`
}

func (controller *Controller) Update(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	id, err := strconv.Atoi(routing.GetVariable(request, 0))
	if err != nil {
		response.Send(err, http.StatusBadRequest)
		return
	}

	record, err := controller.recordDAO.Find(id)
	if err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	if record.Id == 0 {
		response.Send(fmt.Sprintf("Record %d not found", id), http.StatusNotFound)
		return
	}

	errors, input, err := validation.NewValidator[updateInput]().Validate(request)
	if err != nil {
		response.Send(err, http.StatusBadRequest)
		return
	}

	if len(errors) > 0 {
		response.Send(errors, http.StatusUnprocessableEntity)
		return
	}

	user := jwt.GetUser(request)
	allExist, missingActivities, err := controller.activityService.CheckExistence(input.Activities, user.Id)
	if err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	if !allExist {
		errors := map[string]any{"activities": fmt.Sprintf("Activities %v not found", missingActivities)}
		response.Send(errors, http.StatusNotFound)
		return
	}

	recordBefore := record
	record.Notes = input.Notes
	record.Mood = input.Mood
	if err := controller.recordDAO.Save(&record); err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	if err := controller.recordDAO.SyncRelations(record.Id, input.Activities); err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
	log.Info("User", user.Id, "changed record", fmt.Sprintf("%+v", recordBefore), "to", fmt.Sprintf("%+v", record))
}
