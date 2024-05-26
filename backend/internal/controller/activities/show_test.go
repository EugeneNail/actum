package activities

import (
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"testing"
)

func TestShow(t *testing.T) {
	client, database := startup.ActivitiesShow(t)

	client.
		Post("/api/collections", `{
			"name": "Work"
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Work",
			"user_id": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "Sleep",
			"icon": "add",
			"collectionId": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Sleep",
			"icon":          "add",
			"user_id":       1,
			"collection_id": 1,
		})

	client.
		Get("/api/activities/1").
		AssertStatus(http.StatusOK)
}

func TestShowInvalidId(t *testing.T) {
	client, _ := startup.ActivitiesShow(t)

	client.
		Get("/api/activities/one").
		AssertStatus(http.StatusBadRequest)
}

func TestShowSomeonesActivity(t *testing.T) {
	client, database := startup.ActivitiesShow(t)

	client.
		Post("/api/collections", `{
			"name": "Cooking"
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Cooking",
			"user_id": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "Make cake",
			"icon": "cake",
			"collectionId": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Make cake",
			"icon":          "cake",
			"user_id":       1,
			"collection_id": 1,
		})

	client.ChangeUser()
	client.
		Get("/api/activities/1").
		AssertStatus(http.StatusForbidden)
}

func TestShowNotFound(t *testing.T) {
	client, database := startup.ActivitiesShow(t)

	client.
		Post("/api/collections", `{
			"name": "Kill bugs"
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Kill bugs",
			"user_id": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "Run",
			"icon": "minus",
			"collectionId": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Run",
			"icon":          "minus",
			"user_id":       1,
			"collection_id": 1,
		})

	client.
		Get("/api/activities/2").
		AssertStatus(http.StatusNotFound)
}
