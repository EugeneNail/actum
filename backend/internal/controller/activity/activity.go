package activity

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/model/activities"
	"github.com/EugeneNail/actum/internal/model/collections"
	"strings"
)

func hasDuplicateActivity(name string, collectionId int) (bool, error) {
	collection, err := collections.Find(collectionId)
	if err != nil {
		return false, fmt.Errorf("hasDuplicateActivity(): %w", err)
	}

	activities, err := getCollectionActivities(collection.Id)
	if err != nil {
		return false, fmt.Errorf("hasDuplicateActivity(): %w", err)
	}

	for _, activity := range activities {
		if strings.ToLower(activity.Name) == strings.ToLower(name) {
			return true, nil
		}
	}

	return false, nil
}

func getCollectionActivities(collectionId int) ([]activities.Activity, error) {
	var collectionActivities []activities.Activity

	db, err := mysql.Connect()
	defer db.Close()
	if err != nil {
		return collectionActivities, fmt.Errorf("collections.Activities(): %w", err)
	}

	rows, err := db.Query(`SELECT * FROM activities WHERE collection_id = ?`, collectionId)
	defer rows.Close()
	if err != nil {
		return collectionActivities, fmt.Errorf("collections.Activities(): %w", err)
	}

	for rows.Next() {
		var activity activities.Activity
		err := rows.Scan(&activity.Id, &activity.Name, &activity.Icon, &activity.UserId, &activity.CollectionId)
		if err != nil {
			return collectionActivities, fmt.Errorf("collections.Activities(): %w", err)
		}
		collectionActivities = append(collectionActivities, activity)
	}

	return collectionActivities, nil
}
