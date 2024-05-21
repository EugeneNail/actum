package collection

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/controller"
	"github.com/EugeneNail/actum/internal/database/mysql"
	act "github.com/EugeneNail/actum/internal/model/activities"
	col "github.com/EugeneNail/actum/internal/model/collections"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"net/http"
)

func Index(writer http.ResponseWriter, request *http.Request) {
	controller := controller.New[any](writer)
	user := jwt.GetUser(request)

	collections, collectionsById, err := fetchCollections(user.Id)
	if err != nil {
		controller.Response(err, http.StatusInternalServerError)
		return
	}

	if len(collections) == 0 {
		controller.Response(collections, http.StatusOK)
		log.Info("User", user.Id, "indexed no collections")
		return
	}

	if err := fetchActivities(user.Id, collectionsById); err != nil {
		controller.Response(err, http.StatusInternalServerError)
		return
	}

	controller.Response(collections, http.StatusOK)
	log.Info("User", user.Id, "indexed", len(collections), "collections")
}

func fetchCollections(userId int) ([]*col.Collection, map[int]*col.Collection, error) {
	collections := make([]*col.Collection, 0)
	collectionsById := make(map[int]*col.Collection)

	db, err := mysql.Connect()
	defer db.Close()
	if err != nil {
		return collections, collectionsById, fmt.Errorf("fetchCollections(): %w", err)
	}

	rows, err := db.Query(`SELECT * FROM collections WHERE user_id = ?`, userId)
	defer rows.Close()
	if err != nil {
		return collections, collectionsById, fmt.Errorf("fetchCollections(): %w", err)
	}

	for rows.Next() {
		collection := col.Collection{}

		err := rows.Scan(&collection.Id, &collection.Name, &collection.UserId)
		if err != nil {
			return collections, collectionsById, fmt.Errorf("fetchCollections(): %w", err)
		}

		collections = append(collections, &collection)
		collectionsById[collection.Id] = &collection
	}

	return collections, collectionsById, nil
}

func fetchActivities(userId int, collectionsById map[int]*col.Collection) error {
	db, err := mysql.Connect()
	defer db.Close()
	if err != nil {
		return fmt.Errorf("fetchCollections(): %w", err)
	}

	rows, err := db.Query(`SELECT id, name, icon, collection_id, user_id FROM activities WHERE user_id = ?`, userId)
	defer rows.Close()
	if err != nil {
		return fmt.Errorf("fetchCollections(): %w", err)
	}

	for rows.Next() {
		var activity act.Activity

		err := rows.Scan(&activity.Id, &activity.Name, &activity.Icon, &activity.CollectionId, &activity.UserId)
		if err != nil {
			return fmt.Errorf("fetchCollections(): %w", err)
		}

		collection := collectionsById[activity.CollectionId]
		collection.Activities = append(collection.Activities, activity)
	}

	return nil
}
