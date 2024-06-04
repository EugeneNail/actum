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
