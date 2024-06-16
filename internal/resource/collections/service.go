package collections

import (
	"database/sql"
	"fmt"
	"github.com/EugeneNail/actum/internal/resource/activities"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db}
}

func (service *Service) HasDuplicate(name string, userId int) (bool, error) {
	var count int

	err := service.db.
		QueryRow(`SELECT COUNT(*) FROM collections WHERE user_id = ? AND LOWER(name) = ?`, userId, name).
		Scan(&count)
	if err != nil {
		return false, fmt.Errorf("collections.HasDuplicate(): %w", err)
	}

	return count > 0, nil
}

func (service *Service) ExceedsLimit(limit int, userId int) (bool, error) {
	var count int

	err := service.db.QueryRow(
		`SELECT COUNT(*) FROM collections WHERE user_id = ?`,
		userId).
		Scan(&count)
	if err != nil {
		return false, fmt.Errorf("collections.ExceedsLimit(): %w", err)
	}

	return count >= limit, nil
}

func (service *Service) CollectCollections(userId int) ([]*Collection, error) {
	collections, err := service.fetchCollections(userId)
	if err != nil {
		return collections, fmt.Errorf("collections.CollectCollections(): %w", err)
	}

	activities, err := service.fetchActivities(userId)
	if err != nil {
		return collections, fmt.Errorf("collections.CollectCollections(): %w", err)
	}

	service.assignToCollections(collections, activities)

	return collections, nil
}

func (service *Service) fetchCollections(userId int) ([]*Collection, error) {
	items := make([]*Collection, 0)

	rows, err := service.db.Query(`SELECT id, name, color, user_id FROM collections WHERE user_id = ?`, userId)
	defer rows.Close()
	if err != nil {
		return items, fmt.Errorf("fetchCollections(): %w", err)
	}

	for rows.Next() {
		collection := Collection{}
		if err := rows.Scan(&collection.Id, &collection.Name, &collection.Color, &collection.UserId); err != nil {
			return items, fmt.Errorf("fetchCollections(): %w", err)
		}
		items = append(items, &collection)
	}

	return items, nil
}

func (service *Service) fetchActivities(userId int) ([]activities.Activity, error) {
	items := make([]activities.Activity, 0)

	rows, err := service.db.Query(`SELECT id, name, icon, collection_id, user_id FROM activities WHERE user_id = ?`, userId)
	defer rows.Close()
	if err != nil {
		return items, fmt.Errorf("fetchActivities(): %w", err)
	}

	for rows.Next() {
		var activity activities.Activity
		if err := rows.Scan(&activity.Id, &activity.Name, &activity.Icon, &activity.CollectionId, &activity.UserId); err != nil {
			return items, fmt.Errorf("fetchActivities(): %w", err)
		}
		items = append(items, activity)
	}

	return items, nil
}

func (service *Service) assignToCollections(initialCollections []*Collection, activities []activities.Activity) {
	collections := make(map[int]*Collection, len(initialCollections))

	for _, collection := range initialCollections {
		collections[collection.Id] = collection
	}

	for _, activity := range activities {
		collection := collections[activity.CollectionId]
		collection.Activities = append(collection.Activities, activity)
	}
}
