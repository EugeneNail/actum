package activities

import (
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"testing"
)

func TestUpdateValidData(t *testing.T) {
	client, database := startup.Activities(t)

	client.
		Post("/api/collections", `{
			"name": "Что-то тут написано",
			"color": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Что-то тут написано",
			"color":   1,
			"user_id": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "И тут тоже",
			"icon": 200,
			"collectionId": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "И тут тоже",
			"icon":          200,
			"collection_id": 1,
		})

	client.
		Put("/api/activities/1", `{
			"name": "Lorem",
			"icon": 201
		}`).
		AssertStatus(http.StatusNoContent)

	database.
		AssertCount("activities", 1).
		AssertLacks("activities", map[string]any{
			"name":          "И тут тоже",
			"icon":          200,
			"collection_id": 1,
		}).
		AssertHas("activities", map[string]any{
			"name":          "Lorem",
			"icon":          201,
			"collection_id": 1,
		})
}

func TestUpdateInvalidData(t *testing.T) {
	client, database := startup.Activities(t)

	client.
		Post("/api/collections", `{
			"name": "Коллекция",
			"color": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":  "Коллекция",
			"color": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "На меня смотри",
			"icon": 401,
			"collectionId": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "На меня смотри",
			"icon":          401,
			"collection_id": 1,
		})

	client.
		Put("/api/activities/1", `{
			"name": "Lo",
			"icon": 10000
		}`).
		AssertStatus(http.StatusUnprocessableEntity)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "На меня смотри",
			"icon":          401,
			"collection_id": 1,
		})
}

func TestUpdateNotFound(t *testing.T) {
	client, database := startup.Activities(t)

	client.
		Post("/api/collections", `{
			"name": "Что-то тут написано",
			"color": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":  "Что-то тут написано",
			"color": 1,
		})

	database.AssertCount("activities", 0)

	client.
		Put("/api/activities/1", `{
			"name": "Lorem",
			"icon": 421
		}`).
		AssertStatus(http.StatusNotFound)

	database.AssertCount("activities", 0)
}

func TestUpdateSomeoneElses(t *testing.T) {
	client, database := startup.Activities(t)

	client.
		Post("/api/collections", `{
			"name": "Видеоигры",
			"color": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":  "Видеоигры",
			"color": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "Майнкрафт",
			"icon": 800,
			"collectionId": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Майнкрафт",
			"icon":          800,
			"collection_id": 1,
		})

	client.ChangeUser()
	client.
		Put("/api/activities/1", `{
			"name": "Lorem",
			"icon": 900
		}`).
		AssertStatus(http.StatusForbidden)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Майнкрафт",
			"icon":          800,
			"collection_id": 1,
		})
}

func TestUpdateInvalidId(t *testing.T) {
	client, database := startup.Activities(t)

	client.
		Put("/api/activities/one", `{
			"name": "Lorem",
			"icon": 100
		}`).
		AssertStatus(http.StatusBadRequest)

	database.AssertCount("activities", 0)
}
