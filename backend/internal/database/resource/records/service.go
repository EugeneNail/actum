package records

import (
	"database/sql"
	"fmt"
	"time"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db}
}

func (service *Service) CollectRecordsForCursor(cursor time.Time, userId int) ([]*IndexRecord, error) {
	start := cursor.Add(time.Hour * 24 * 14 * -1)

	records, err := service.fetchRecords(start, cursor, userId)
	if err != nil {
		return records, fmt.Errorf("recordService.CollectRecordsForCursor(): %w", err)
	}

	if len(records) == 0 {
		return records, nil
	}

	if err := service.fetchCollections(records, userId); err != nil {
		return records, fmt.Errorf("recordService.CollectRecordsForCursor(): %w", err)
	}

	if err := service.fetchActivities(records, userId); err != nil {
		return records, fmt.Errorf("recordService.CollectRecordsForCursor(): %w", err)
	}

	return records, nil
}

func (service *Service) fetchRecords(start time.Time, end time.Time, userId int) ([]*IndexRecord, error) {
	var records []*IndexRecord

	rows, err := service.db.Query(
		`SELECT id, mood, weather, date, notes FROM records WHERE user_id = ? AND date > ? AND date <= ?`,
		userId, start, end,
	)
	defer rows.Close()
	if err != nil {
		return records, fmt.Errorf("records.fetchRecords(): %w", err)
	}

	for rows.Next() {
		var record IndexRecord

		if err := rows.Scan(&record.Id, &record.Mood, &record.Weather, &record.Date, &record.Notes); err != nil {
			return records, fmt.Errorf("records.fetchRecords(): %w", err)
		}
		record.Collections = []IndexCollection{}
		records = append(records, &record)
	}

	return records, nil
}

func (service *Service) fetchCollections(records []*IndexRecord, userId int) error {
	rows, err := service.db.Query(
		`SELECT id, name, color FROM collections WHERE user_id = ?`,
		userId,
	)
	defer rows.Close()
	if err != nil {
		return fmt.Errorf("records.fetchCollections(): %w", err)
	}

	for rows.Next() {
		var collection IndexCollection

		if err := rows.Scan(&collection.Id, &collection.Name, &collection.Color); err != nil {
			return fmt.Errorf("records.fetchCollections(): %w", err)
		}

		for _, record := range records {
			record.Collections = append(record.Collections, collection)
		}
	}

	return nil
}

func (service *Service) fetchActivities(records []*IndexRecord, userId int) error {
	query, values := service.prepareActivitiesQuery(records, userId)
	rows, err := service.db.Query(query, values...)
	defer rows.Close()
	if err != nil {
		return fmt.Errorf("records.fetchActivities(): %w", err)
	}

	for rows.Next() {
		var activity IndexActivity
		if err := rows.Scan(&activity.RecordId, &activity.CollectionId, &activity.Name, &activity.Icon); err != nil {
			return fmt.Errorf("records.fetchActivities(): %w", err)
		}
		service.assignToRecords(records, activity)
	}

	return nil
}

func (service *Service) prepareActivitiesQuery(records []*IndexRecord, userId int) (string, []any) {
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

func (service *Service) assignToRecords(records []*IndexRecord, activity IndexActivity) {
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

func (service *Service) IsDateTaken(date string, userId int) (bool, error) {
	var count int

	err := service.db.
		QueryRow(`SELECT COUNT(*) FROM records WHERE user_id = ? AND date = ?`, userId, date).
		Scan(&count)
	if err != nil {
		return false, fmt.Errorf("isDateTaken(): %w", err)
	}

	return count > 0, nil
}
