package records

import (
	"database/sql"
	"fmt"
	"github.com/EugeneNail/actum/internal/database/resource/activities"
	"github.com/EugeneNail/actum/internal/database/resource/records"
)

type Controller struct {
	db          *sql.DB
	recordDAO   *records.DAO
	activityDAO *activities.DAO
}

func New(db *sql.DB, recordDAO *records.DAO, activityDAO *activities.DAO) Controller {
	return Controller{db, recordDAO, activityDAO}
}

func (controller *Controller) checkExistence(needle []int, userId int) (bool, []int, error) {
	existing, err := controller.activityDAO.ListIn(needle, userId)
	if err != nil {
		return false, []int{}, fmt.Errorf("checkExistence(): %w", err)
	}

	missing := controller.findMissing(needle, existing)
	return len(missing) == 0, missing, nil
}

func (controller *Controller) findMissing(needle []int, haystack []activities.Activity) []int {
	var missing []int

	hashtable := make(map[int]struct{}, len(haystack))
	for _, activity := range haystack {
		hashtable[activity.Id] = struct{}{}
	}

	for _, needleId := range needle {
		_, found := hashtable[needleId]
		if !found {
			missing = append(missing, needleId)
		}
	}

	return missing

}

func (controller *Controller) ownsEach(activityIds []int, userId int) (bool, error) {
	var count int
	var placeholders string
	values := make([]any, 1+len(activityIds))
	values[0] = userId

	for i, id := range activityIds {
		values[i+1] = id
		placeholders += "?,"
	}
	placeholders = "(" + placeholders[:len(placeholders)-1] + ")"

	err := controller.db.QueryRow(`SELECT COUNT(*) FROM activities WHERE user_id = ? AND id IN`+placeholders, values...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("activities.ListIn(): %w", err)
	}

	return count == len(activityIds), nil
}
