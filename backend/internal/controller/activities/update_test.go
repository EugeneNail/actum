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
			"name": "lorem ipsum",
			"color": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "lorem ipsum",
			"color":   1,
			"user_id": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "Dolor sit amet",
			"icon": 200,
			"collectionId": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Dolor sit amet",
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
			"name":          "Dolor sit amet",
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
			"name": "Collection",
			"color": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":  "Collection",
			"color": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "Look at me",
			"icon": 401,
			"collectionId": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Look at me",
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
			"name":          "Look at me",
			"icon":          401,
			"collection_id": 1,
		})
}

func TestUpdateNotFound(t *testing.T) {
	client, database := startup.Activities(t)

	client.
		Post("/api/collections", `{
			"name": "lorem ipsum",
			"color": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":  "lorem ipsum",
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
			"name": "Playing",
			"color": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":  "Playing",
			"color": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "Creating something",
			"icon": 800,
			"collectionId": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Creating something",
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
			"name":          "Creating something",
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
