package activity

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/model/collections"
	"strings"
)

func hasDuplicateActivity(name string, collectionId int) (bool, error) {
	collection, err := collections.Find(collectionId)
	if err != nil {
		return false, fmt.Errorf("hasDuplicateActivity(): %w", err)
	}

	activities, err := collection.Activities()
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
