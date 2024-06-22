package records

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/service/auth/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/middleware/routing"
	"github.com/EugeneNail/actum/internal/service/response"
	"net/http"
	"strconv"
)

func (controller *Controller) Show(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	id, err := strconv.Atoi(routing.GetVariable(request, 0))
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Show(): %w", err), http.StatusBadRequest)
		return
	}

	record, err := controller.recordDAO.Find(id)
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Show(): %w", err), http.StatusInternalServerError)
		return
	}

	if record.Id == 0 {
		response.Send(fmt.Sprintf("Запись %d не найдена", id), http.StatusNotFound)
		return
	}

	user := jwt.GetUser(request)
	if record.UserId != user.Id {
		response.Send("Вы не можете использовать чужую запись.", http.StatusForbidden)
		return
	}

	activities, err := controller.recordService.FetchIdsOfActivities(record.Id)
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Show(): %w", err), http.StatusInternalServerError)
		return
	}

	photos, err := controller.recordService.FetchNamesOfPhotos(record.Id)
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Show(): %w", err), http.StatusInternalServerError)
		return
	}

	output := make(map[string]any, 7)
	output["date"] = record.Date.Format("2006-01-02")
	output["mood"] = record.Mood
	output["weather"] = record.Weather
	output["notes"] = record.Notes
	output["id"] = record.Id
	output["activities"] = activities
	output["photos"] = photos

	response.Send(output, http.StatusOK)
	log.Info("User", user.Id, "fetched record", record.Id)
}
