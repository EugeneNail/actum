package activities

import (
	"database/sql"
	"fmt"
	"github.com/EugeneNail/actum/internal/database/resource/activities"
	"github.com/EugeneNail/actum/internal/database/resource/collections"
	"strings"
)

type Controller struct {
	db            *sql.DB
	activityDAO   *activities.DAO
	collectionDAO *collections.DAO
}

func New(db *sql.DB, activityDAO *activities.DAO, collectionDAO *collections.DAO) (controller Controller) {
	controller.db = db
	controller.activityDAO = activityDAO
	controller.collectionDAO = collectionDAO
	return
}

func (controller *Controller) exceededLimit(collectionId int, userId int) (bool, error) {
	var count int

	err := controller.db.QueryRow(`SELECT COUNT(*) FROM activities WHERE user_id = ? AND collection_id = ?`, userId, collectionId).
		Scan(&count)

	if err != nil {
		return false, fmt.Errorf("collections.exceededLimit(): %w", err)
	}

	return count >= 20, nil
}

func (controller *Controller) hasDuplicate(name string, userId int) (bool, error) {
	var count int

	err := controller.db.QueryRow(
		`SELECT COUNT(*) FROM activities WHERE user_id = ? AND LOWER(name) = ?`,
		userId, strings.ToLower(name),
	).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("activities.hasDuplicate(): %w", err)
	}

	return count > 0, nil
}
