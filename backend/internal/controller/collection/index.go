package collection

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/controller"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/model/activities"
	"github.com/EugeneNail/actum/internal/model/collections"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"net/http"
)

type output = []outputItem

type outputItem struct {
	collections.Collection
	Activities []activities.Activity `json:"activities"`
}

type Item struct {
	Collection collections.Collection
	Activities []activities.Activity
}

func Index(writer http.ResponseWriter, request *http.Request) {
	controller := controller.New[any](writer)
	user := jwt.GetUser(request)

	collectionsWithActivities, err := getOutput(user.Id)
	if err != nil {
		controller.Response(err, http.StatusInternalServerError)
		return
	}

	controller.Response(collectionsWithActivities, http.StatusOK)
	log.Info("User", user.Id, "indexed", len(collectionsWithActivities), "collections")
}

func getOutput(userId int) (output, error) {
	var output output

	db, err := mysql.Connect()
	defer db.Close()
	if err != nil {
		return output, fmt.Errorf("getOutput(): %w", err)
	}

	query := `
		SELECT * 
		FROM collections 
		    LEFT JOIN activities 
		        ON collections.id = activities.collection_id 
		WHERE collections.user_id = ?
		ORDER BY collections.id
	`

	rows, err := db.Query(query, userId)
	defer rows.Close()
	if err != nil {
		return output, fmt.Errorf("getOutput(): %w", err)
	}

	itemMap := map[int]*Item{}

	for rows.Next() {
		var collection collections.Collection
		var activity activities.Activity

		err := rows.Scan(
			&collection.Id,
			&collection.Name,
			&collection.UserId,
			&activity.Id,
			&activity.Name,
			&activity.Icon,
			&activity.UserId,
			&activity.CollectionId,
		)
		if err != nil {
			return output, fmt.Errorf("getOutput(): %w", err)
		}

		if item, exists := itemMap[collection.Id]; exists {
			if activity.Id != 0 {
				item.Activities = append(item.Activities, activity)
			}
		} else {
			item = &Item{collection, []activities.Activity{}}
			if activity.Id != 0 {
				item.Activities = append(item.Activities, activity)
			}
			itemMap[collection.Id] = item
		}
	}

	for _, item := range itemMap {
		output = append(output, outputItem{item.Collection, item.Activities})
	}

	return output, nil
}
