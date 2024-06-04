package records

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/service/jwt"
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

	user := jwt.GetUser(request)
	if record.UserId != user.Id {
		response.Send("You are not allowed to view other users' records", http.StatusForbidden)
		return
	}

	activities, err := controller.fetchIdsOfActivities(record.Id)
	if err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	output := make(map[string]any, 5)
	output["date"] = record.Date.Format("2006-01-02")
	output["mood"] = record.Mood
	output["notes"] = record.Notes
	output["id"] = record.Id
	output["activities"] = activities

	response.Send(output, http.StatusOK)
	log.Info("User", user.Id, "fetched record", record.Id)
}

func (controller *Controller) fetchIdsOfActivities(recordId int) ([]int, error) {
	var ids []int

	rows, err := controller.db.Query(`
		SELECT id 
		FROM activities 
		    JOIN records_activities 
		        ON activities.id = records_activities.activity_id 
		WHERE record_id = ?`, recordId)
	defer rows.Close()
	if err != nil {
		return ids, fmt.Errorf("fetchIdsOfActivities(): %w", err)
	}

	for rows.Next() {
		var id int

		if err := rows.Scan(&id); err != nil {
			return ids, fmt.Errorf("fetchIdsOfActivities(): %w", err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}
