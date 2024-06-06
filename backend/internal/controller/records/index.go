package records

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/response"
	"github.com/EugeneNail/actum/internal/service/validation"
	"net/http"
	"time"
)

type indexInput struct {
	Cursor string `json:"cursor" rules:"required|date"`
}

type ShortRecord struct {
	Id          int               `json:"id"`
	Date        string            `json:"date"`
	Mood        int               `json:"mood"`
	Notes       string            `json:"notes"`
	Collections []ShortCollection `json:"collections"`
}

type ShortCollection struct {
	Id         int             `json:"id"`
	Name       string          `json:"name"`
	Color      int             `json:"color"`
	Activities []ShortActivity `json:"activities"`
}

type ShortActivity struct {
	RecordId     int    `json:"recordId"`
	CollectionId int    `json:"collectionId"`
	Icon         int    `json:"icon"`
	Name         string `json:"name"`
}

func (controller *Controller) Index(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	errors, input, err := validation.NewValidator[indexInput]().Validate(request)
	if err != nil {
		response.Send(err, http.StatusBadRequest)
		return
	}

	if len(errors) > 0 {
		response.Send(errors, http.StatusUnprocessableEntity)
		return
	}

	user := jwt.GetUser(request)
	end, err := time.Parse("2006-01-02", input.Cursor)
	if err != nil {
		response.Send(err, http.StatusBadRequest)
		return
	}
	start := end.Add(time.Hour * 24 * 14 * -1)

	records, err := controller.fetchRecords(start, end, user.Id)
	if err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	if err := controller.fetchCollections(records, user.Id); err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	if err := controller.fetchActivities(records, user.Id); err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	response.Send(records, http.StatusOK)
	log.Info("User", user.Id, "fetched", len(records), "records from", start.Format("2006-01-02"), "to", input.Cursor)
}

func (controller *Controller) fetchRecords(start time.Time, end time.Time, userId int) ([]*ShortRecord, error) {
	var records []*ShortRecord

	rows, err := controller.db.Query(
		`SELECT id, mood, date, notes FROM records WHERE user_id = ? AND date > ? AND date <= ?`,
		userId, start, end,
	)
	defer rows.Close()
	if err != nil {
		return records, fmt.Errorf("records.fetchRecords(): %w", err)
	}

	for rows.Next() {
		var record ShortRecord

		if err := rows.Scan(&record.Id, &record.Mood, &record.Date, &record.Notes); err != nil {
			return records, fmt.Errorf("records.fetchRecords(): %w", err)
		}
		record.Collections = []ShortCollection{}
		records = append(records, &record)
	}

	return records, nil
}

func (controller *Controller) fetchCollections(records []*ShortRecord, userId int) error {
	rows, err := controller.db.Query(
		`SELECT id, name, color FROM collections WHERE user_id = ?`,
		userId,
	)
	defer rows.Close()
	if err != nil {
		return fmt.Errorf("records.fetchCollections(): %w", err)
	}

	for rows.Next() {
		var collection ShortCollection

		if err := rows.Scan(&collection.Id, &collection.Name, &collection.Color); err != nil {
			return fmt.Errorf("records.fetchCollections(): %w", err)
		}

		for _, record := range records {
			record.Collections = append(record.Collections, collection)
		}
	}

	return nil
}

func (controller *Controller) fetchActivities(records []*ShortRecord, userId int) error {
	query, values := controller.prepareActivitiesQuery(records, userId)
	rows, err := controller.db.Query(query, values...)
	defer rows.Close()
	if err != nil {
		return fmt.Errorf("records.fetchActivities(): %w", err)
	}

	for rows.Next() {
		var activity ShortActivity
		if err := rows.Scan(&activity.RecordId, &activity.CollectionId, &activity.Name, &activity.Icon); err != nil {
			return fmt.Errorf("records.fetchActivities(): %w", err)
		}
		controller.assignToRecords(records, activity)
	}

	return nil
}

func (controller *Controller) prepareActivitiesQuery(records []*ShortRecord, userId int) (string, []any) {
	var placeholders string
	values := make([]any, len(records)+1)
	values[0] = userId

	for i, record := range records {
		values[i+1] = record.Id
		placeholders += "?,"
	}
	placeholders = "(" + placeholders[:len(placeholders)-1] + ")"

	query := `
		SELECT records.id AS record_id,
	       collection_id,
	       name,
	       icon
		FROM activities
		     JOIN records_activities
		          ON activities.id = records_activities.activity_id
		     JOIN records
		          ON records_activities.record_id = records.id
		WHERE activities.user_id = ?
		  AND records.id IN ` + placeholders

	return query, values
}

func (controller *Controller) assignToRecords(records []*ShortRecord, activity ShortActivity) {
	for _, record := range records {
		if activity.RecordId != record.Id {
			continue
		}

		for i, collection := range record.Collections {
			if activity.CollectionId != collection.Id {
				continue
			}
			record.Collections[i].Activities = append(record.Collections[i].Activities, activity)
		}
	}
}
