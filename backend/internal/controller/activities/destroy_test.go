package activities

import (
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"testing"
)

func TestDestroy(t *testing.T) {
	client, database := startup.Activities(t)

	client.
		Post("/api/collections", `{
			"name": "fabulous collection",
			"color": 2
		}`).
		AssertStatus(http.StatusCreated)

	database.AssertHas("collections", map[string]any{
		"name":    "fabulous collection",
		"color":   2,
		"user_id": 1,
	})

	client.
		Post("/api/activities", `{
			"name": "Fabulous activity",
			"icon": 332,
			"collectionId": 1 
		}`).
		AssertStatus(http.StatusCreated)

	database.AssertHas("activities", map[string]any{
		"name":          "Fabulous activity",
		"icon":          332,
		"collection_id": 1,
		"user_id":       1,
	})

	client.Delete("/api/activities/1").AssertStatus(http.StatusNoContent)
	database.AssertCount("activities", 0)

}

func TestDestroyInvalidId(t *testing.T) {
	client, database := startup.Activities(t)

	client.
		Post("/api/collections", `{
			"name": "fabulous collection",
			"color": 2
		}`).
		AssertStatus(http.StatusCreated)

	database.AssertHas("collections", map[string]any{
		"name":    "fabulous collection",
		"color":   2,
		"user_id": 1,
	})

	client.
		Post("/api/activities", `{
			"name": "Fabulous activity",
			"icon": 323,
			"collectionId": 1 
		}`).
		AssertStatus(http.StatusCreated)

	database.AssertHas("activities", map[string]any{
		"name":          "Fabulous activity",
		"icon":          323,
		"collection_id": 1,
		"user_id":       1,
	})

	client.
		Delete("/api/activities/one").
		AssertStatus(http.StatusBadRequest)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Fabulous activity",
			"icon":          323,
			"collection_id": 1,
			"user_id":       1,
		})
}

func TestDestroyNotFound(t *testing.T) {
	client, database := startup.Activities(t)

	database.AssertCount("activities", 0)
	client.
		Delete("/api/activities/1").
		AssertStatus(http.StatusNotFound)
}

func TestDestroySomeoneElsesActivity(t *testing.T) {
	client, database := startup.Activities(t)

	client.
		Post("/api/collections", `{
			"name": "fabulous collection",
			"color": 2
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "fabulous collection",
			"color":   2,
			"user_id": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "Fabulous activity",
			"icon": 876,
			"collectionId": 1 
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Fabulous activity",
			"icon":          876,
			"collection_id": 1,
		})

	client.ChangeUser()
	client.
		Delete("/api/activities/1").
		AssertStatus(http.StatusForbidden)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Fabulous activity",
			"icon":          876,
			"collection_id": 1,
		})
}
