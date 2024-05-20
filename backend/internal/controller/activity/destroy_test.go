package activity

import (
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"testing"
)

func TestDestroy(t *testing.T) {
	client, database := startup.ActivitiesDestroy(t)

	client.
		Post("/api/collections", `{
			"name": "fabulous collection"
		}`).
		AssertStatus(http.StatusCreated)

	database.AssertHas("collections", map[string]any{
		"name":    "fabulous collection",
		"user_id": 1,
	})

	client.
		Post("/api/activities", `{
			"name": "Fabulous activity",
			"icon": "product",
			"collectionId": 1 
		}`).
		AssertStatus(http.StatusCreated)

	database.AssertHas("activities", map[string]any{
		"name":          "Fabulous activity",
		"icon":          "product",
		"collection_id": 1,
		"user_id":       1,
	})

	client.Delete("/api/activities/1").AssertStatus(http.StatusNoContent)
	database.AssertCount("activities", 0)

}

func TestDestroyInvalidId(t *testing.T) {
	client, database := startup.ActivitiesDestroy(t)

	client.
		Post("/api/collections", `{
			"name": "fabulous collection"
		}`).
		AssertStatus(http.StatusCreated)

	database.AssertHas("collections", map[string]any{
		"name":    "fabulous collection",
		"user_id": 1,
	})

	client.
		Post("/api/activities", `{
			"name": "Fabulous activity",
			"icon": "product",
			"collectionId": 1 
		}`).
		AssertStatus(http.StatusCreated)

	database.AssertHas("activities", map[string]any{
		"name":          "Fabulous activity",
		"icon":          "product",
		"collection_id": 1,
		"user_id":       1,
	})

	client.Delete("/api/activities/one").AssertStatus(http.StatusBadRequest)
}

func TestDestroyNotFound(t *testing.T) {
	client, database := startup.ActivitiesDestroy(t)

	database.AssertCount("activities", 0)
	client.Delete("/api/activities/1").AssertStatus(http.StatusNotFound)
}

func TestDestroySomeoneElsesActivity(t *testing.T) {
	client, database := startup.ActivitiesDestroy(t)

	client.
		Post("/api/collections", `{
			"name": "fabulous collection"
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "fabulous collection",
			"user_id": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "Fabulous activity",
			"icon": "product",
			"collectionId": 1 
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Fabulous activity",
			"icon":          "product",
			"collection_id": 1,
		})

	client.ChangeUser()
	client.Delete("/api/activities/1").AssertStatus(http.StatusForbidden)
}
