package activity

import (
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"testing"
)

func TestUpdateValidData(t *testing.T) {
	client, database := startup.ActivitiesUpdate(t)

	client.
		Post("/api/collections", `{
			"name": "lorem ipsum"
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name": "lorem ipsum",
		})

	client.
		Post("/api/activities", `{
			"name": "Dolor sit amet",
			"icon": "add",
			"collectionId": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Dolor sit amet",
			"icon":          "add",
			"collection_id": 1,
		})

	client.
		Put("/api/activities/1", `{
			"name": "Lorem",
			"icon": "Test"
		}`).
		AssertStatus(http.StatusNoContent)

	database.
		AssertCount("activities", 1).
		AssertLacks("activities", map[string]any{
			"name":          "Dolor sit amet",
			"icon":          "add",
			"collection_id": 1,
		}).
		AssertHas("activities", map[string]any{
			"name":          "Lorem",
			"icon":          "Test",
			"collection_id": 1,
		})
}

func TestUpdateInvalidData(t *testing.T) {
	client, database := startup.ActivitiesUpdate(t)

	client.
		Post("/api/collections", `{
			"name": "Collection"
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name": "Collection",
		})

	client.
		Post("/api/activities", `{
			"name": "Look at me",
			"icon": "hello",
			"collectionId": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Look at me",
			"icon":          "hello",
			"collection_id": 1,
		})

	client.
		Put("/api/activities/1", `{
			"name": "Lo",
			"icon": "gas-station"
		}`).
		AssertStatus(http.StatusUnprocessableEntity)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Look at me",
			"icon":          "hello",
			"collection_id": 1,
		})
}

func TestUpdateNotFound(t *testing.T) {
	client, database := startup.ActivitiesUpdate(t)

	client.
		Post("/api/collections", `{
			"name": "lorem ipsum"
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name": "lorem ipsum",
		})

	database.AssertCount("activities", 0)

	client.
		Put("/api/activities/1", `{
			"name": "Lorem",
			"icon": "Test"
		}`).
		AssertStatus(http.StatusNotFound)

	database.
		AssertCount("activities", 0).
		AssertLacks("activities", map[string]any{
			"name":          "Lorem",
			"icon":          "Test",
			"collection_id": 1,
		})
}

func TestUpdateSomeoneElses(t *testing.T) {
	client, database := startup.ActivitiesUpdate(t)

	client.
		Post("/api/collections", `{
			"name": "Playing"
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name": "Playing",
		})

	client.
		Post("/api/activities", `{
			"name": "Creating something",
			"icon": "pencil",
			"collectionId": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Creating something",
			"icon":          "pencil",
			"collection_id": 1,
		})

	client.ChangeUser()
	client.
		Put("/api/activities/1", `{
			"name": "Lorem",
			"icon": "Test"
		}`).
		AssertStatus(http.StatusForbidden)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Creating something",
			"icon":          "pencil",
			"collection_id": 1,
		})
}

func TestUpdateInvalidId(t *testing.T) {
	client, database := startup.ActivitiesUpdate(t)

	client.
		Put("/api/activities/one", `{
			"name": "Lorem",
			"icon": "Test"
		}`).
		AssertStatus(http.StatusBadRequest)

	database.AssertCount("activities", 0)
}
