package collections

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/database/resource/activities"
	"github.com/EugeneNail/actum/internal/database/resource/collections"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/response"
	"net/http"
)

func (controller *Controller) Index(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	user := jwt.GetUser(request)

	collections, err := controller.fetchCollections(user.Id)
	if err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	if len(collections) == 0 {
		response.Send(collections, http.StatusOK)
		log.Info("User", user.Id, "indexed no collections")
		return
	}

	activities, err := controller.fetchActivities(user.Id)
	if err != nil {
		response.Send(err, http.StatusInternalServerError)
		return
	}

	controller.allocateActivities(collections, activities)

	response.Send(collections, http.StatusOK)
	log.Info("User", user.Id, "indexed", len(collections), "collections")
}

func (controller *Controller) fetchCollections(userId int) ([]*collections.Collection, error) {
	items := make([]*collections.Collection, 0)

	rows, err := controller.db.Query(`SELECT id, name, color, user_id FROM collections WHERE user_id = ?`, userId)
	defer rows.Close()
	if err != nil {
		return items, fmt.Errorf("fetchCollections(): %w", err)
	}

	for rows.Next() {
		collection := collections.Collection{}

		err := rows.Scan(&collection.Id, &collection.Name, &collection.Color, &collection.UserId)
		if err != nil {
			return items, fmt.Errorf("fetchCollections(): %w", err)
		}

		items = append(items, &collection)
	}

	return items, nil
}

func (controller *Controller) fetchActivities(userId int) ([]activities.Activity, error) {
	items := make([]activities.Activity, 0)

	rows, err := controller.db.Query(`SELECT id, name, icon, collection_id, user_id FROM activities WHERE user_id = ?`, userId)
	defer rows.Close()
	if err != nil {
		return items, fmt.Errorf("fetchCollections(): %w", err)
	}

	for rows.Next() {
		var activity activities.Activity

		err := rows.Scan(&activity.Id, &activity.Name, &activity.Icon, &activity.CollectionId, &activity.UserId)
		if err != nil {
			return items, fmt.Errorf("fetchCollections(): %w", err)
		}

		items = append(items, activity)
	}

	return items, nil
}

func (controller *Controller) allocateActivities(initialCollections []*collections.Collection, activities []activities.Activity) {
	collections := make(map[int]*collections.Collection, len(initialCollections))

	for _, collection := range initialCollections {
		collections[collection.Id] = collection
	}

	for _, activity := range activities {
		collection := collections[activity.CollectionId]
		collection.Activities = append(collection.Activities, activity)
	}
}
