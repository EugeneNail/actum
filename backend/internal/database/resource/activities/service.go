package activities

import (
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
