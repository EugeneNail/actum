package activities

import (
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"testing"
)

func TestShow(t *testing.T) {
	client, database := startup.Activities(t)

	client.
		Post("/api/collections", `{
			"name": "Работа",
			"color": 3
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Работа",
			"color":   3,
			"user_id": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "Сон",
			"icon": 421,
			"collectionId": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Сон",
			"icon":          421,
			"user_id":       1,
			"collection_id": 1,
		})

	client.
		Get("/api/activities/1").
		AssertStatus(http.StatusOK)
}

func TestShowInvalidId(t *testing.T) {
	client, _ := startup.Activities(t)

	client.
		Get("/api/activities/one").
		AssertStatus(http.StatusBadRequest)
}

func TestShowSomeonesActivity(t *testing.T) {
	client, database := startup.Activities(t)

	client.
		Post("/api/collections", `{
			"name": "Готовка",
			"color": 3
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Готовка",
			"color":   3,
			"user_id": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "Испекла пирог",
			"icon": 777,
			"collectionId": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Испекла пирог",
			"icon":          777,
			"user_id":       1,
			"collection_id": 1,
		})

	client.ChangeUser()
	client.
		Get("/api/activities/1").
		AssertStatus(http.StatusForbidden)
}

func TestShowNotFound(t *testing.T) {
	client, database := startup.Activities(t)

	client.
		Post("/api/collections", `{
			"name": "Коллекционирование",
			"color": 3
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Коллекционирование",
			"color":   3,
			"user_id": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "Бег",
			"icon": 645,
			"collectionId": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Бег",
			"icon":          645,
			"user_id":       1,
			"collection_id": 1,
		})

	client.
		Get("/api/activities/2").
		AssertStatus(http.StatusNotFound)
}
