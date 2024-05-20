package collection

import (
	"database/sql"
	"fmt"
	"github.com/EugeneNail/actum/internal/controller"
	"github.com/EugeneNail/actum/internal/database"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/model/activities"
	"github.com/EugeneNail/actum/internal/model/collections"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"net/http"
)

type Group struct {
	Collection collections.Collection
	Activities []activities.Activity
}

type row struct {
	collectionId         database.NullableInt
	collectionName       database.NullableString
	collectionUserId     database.NullableInt
	activityId           database.NullableInt
	activityName         database.NullableString
	activityIcon         database.NullableString
	activityUserId       database.NullableInt
	activityCollectionId database.NullableInt
}

func Index(writer http.ResponseWriter, request *http.Request) {
	controller := controller.New[any](writer)
	user := jwt.GetUser(request)

	collections, err := getCollections(user.Id)
	if err != nil {
		controller.Response(err, http.StatusInternalServerError)
		return
	}

	controller.Response(collections, http.StatusOK)
	log.Info("User", user.Id, "indexed", len(collections), "collections")
}

func getCollections(userId int) ([]map[string]any, error) {
	var collections []map[string]any

	db, err := mysql.Connect()
	defer db.Close()
	if err != nil {
		return collections, fmt.Errorf("getCollections(): %w", err)
	}

	rows, err := fetchData(db, userId)
	defer rows.Close()
	if err != nil {
		return collections, fmt.Errorf("getCollections(): %w", err)
	}

	groups, err := mapDataToGroups(rows)
	if err != nil {
		return collections, fmt.Errorf("getCollections(): %w", err)
	}

	for _, group := range groups {
		collections = append(collections, map[string]any{
			"id":         group.Collection.Id,
			"name":       group.Collection.Name,
			"userId":     group.Collection.UserId,
			"activities": group.Activities,
		})
	}

	return collections, nil
}

func fetchData(db *sql.DB, userId int) (*sql.Rows, error) {
	query := `
		SELECT * 
		FROM collections 
		    LEFT JOIN activities 
		        ON collections.id = activities.collection_id 
		WHERE collections.user_id = ?
		ORDER BY collections.id, activities.id 
	`

	rows, err := db.Query(query, userId)
	if err != nil {
		return rows, fmt.Errorf("fetchData(): %w", err)
	}

	return rows, nil
}

func mapDataToGroups(rows *sql.Rows) (map[int]*Group, error) {
	groups := map[int]*Group{}

	for rows.Next() {
		collection, activity, err := scanRow(rows)
		if err != nil {
			return groups, fmt.Errorf("mapDataToGroups(): %w", err)
		}

		if group, exists := groups[collection.Id]; exists {
			if activity.Id != 0 {
				group.Activities = append(group.Activities, activity)
			}
		} else {
			group = &Group{collection, []activities.Activity{}}
			if activity.Id != 0 {
				group.Activities = append(group.Activities, activity)
			}
			groups[collection.Id] = group
		}
	}

	return groups, nil
}

func scanRow(rows *sql.Rows) (collections.Collection, activities.Activity, error) {
	var collection collections.Collection
	var activity activities.Activity

	row := row{}
	err := rows.Scan(
		&row.collectionId,
		&row.collectionName,
		&row.collectionUserId,
		&row.activityId,
		&row.activityName,
		&row.activityIcon,
		&row.activityUserId,
		&row.activityCollectionId,
	)
	if err != nil {
		return collection, activity, fmt.Errorf("scanRow(): %w", err)
	}

	collection = collections.Collection{
		int(row.collectionId),
		string(row.collectionName),
		int(row.collectionUserId),
	}

	activity = activities.Activity{
		int(row.activityId),
		string(row.activityName),
		string(row.activityIcon),
		int(row.activityUserId),
		int(row.activityCollectionId),
	}

	return collection, activity, nil
}
