package activities

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Service struct {
	db  *sql.DB
	dao *DAO
}

func NewService(db *sql.DB, dao *DAO) *Service {
	return &Service{db, dao}
}

func (service *Service) HasDuplicate(name string, userId int) (bool, error) {
	var count int

	err := service.db.
		QueryRow(`SELECT COUNT(*) FROM activities WHERE user_id = ? AND LOWER(name) = ?`, userId, strings.ToLower(name)).
		Scan(&count)
	if err != nil {
		return false, fmt.Errorf("activities.HasDuplicate(): %w", err)
	}

	return count > 0, nil
}

func (service *Service) ExceedsLimit(limit int, collectionId int, userId int) (bool, error) {
	var count int

	err := service.db.
		QueryRow(`SELECT COUNT(*) FROM activities WHERE user_id = ? AND collection_id = ?`, userId, collectionId).
		Scan(&count)
	if err != nil {
		return false, fmt.Errorf("activities.ExceedsLimit(): %w", err)
	}

	return count >= limit, nil
}

func (service *Service) CheckExistence(needleActivities []int, userId int) (bool, []int, error) {
	existing, err := service.dao.ListIn(needleActivities, userId)
	if err != nil {
		return false, []int{}, fmt.Errorf("checkExistence(): %w", err)
	}

	var missing []int
	activitiesById := make(map[int]struct{}, len(existing))

	for _, activity := range existing {
		activitiesById[activity.Id] = struct{}{}
	}

	for _, needleId := range needleActivities {
		if _, found := activitiesById[needleId]; !found {
			missing = append(missing, needleId)
		}
	}

	return len(missing) == 0, missing, nil
}

func (service *Service) SyncRelations(recordId int, activities []int) error {
	tx, err := service.db.BeginTx(context.Background(), &sql.TxOptions{})
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("records.SyncRelations(): %w", err)
	}

	if err = deleteUnusedRelations(tx, recordId, activities); err != nil {
		return fmt.Errorf("records.SyncRelations(): %w", err)
	}

	if err = upsertRelations(tx, recordId, activities); err != nil {
		return fmt.Errorf("records.SyncRelations(): %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("records.SyncRelations(): %w", err)
	}

	return nil
}

func deleteUnusedRelations(tx *sql.Tx, recordId int, activities []int) error {
	var placeholders string
	values := make([]any, len(activities)+1)
	values[0] = recordId

	for i, activityId := range activities {
		values[i+1] = activityId
		placeholders += "?,"
	}
	placeholders = "(" + placeholders[:len(placeholders)-1] + ")"

	_, err := tx.Exec(
		`DELETE FROM records_activities WHERE record_id = ? AND activity_id NOT IN `+placeholders,
		values...,
	)

	if err != nil {
		return fmt.Errorf("deleteUnusedRelations(): %w", err)
	}

	return nil
}

func upsertRelations(tx *sql.Tx, recordId int, activities []int) error {
	const columnsCount = 2
	var placeholders string
	values := make([]any, columnsCount*len(activities))

	for i, activityId := range activities {
		values[columnsCount*i+0] = recordId
		values[columnsCount*i+1] = activityId
		placeholders += "(?, ?),"
	}
	placeholders = placeholders[:len(placeholders)-1]

	_, err := tx.Exec(`
		INSERT INTO records_activities (record_id, activity_id) 
		VALUES `+placeholders+` 
		ON DUPLICATE KEY UPDATE record_id = VALUES(record_id)
	`, values...,
	)

	if err != nil {
		return fmt.Errorf("upsertRelations(): %w", err)
	}

	return nil
}
