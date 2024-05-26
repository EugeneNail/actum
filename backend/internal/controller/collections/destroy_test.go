package collections

import (
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"testing"
)

func TestDestroy(t *testing.T) {
	client, database := startup.CollectionsDestroy(t)

	client.
		Post("/api/collections", `{
			"name": "Do something"
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Do something",
			"user_id": 1,
		})

	client.
		Delete("/api/collections/1").
		AssertStatus(http.StatusNoContent)

	database.
		AssertCount("collections", 0).
		AssertLacks("collections", map[string]any{
			"id":      1,
			"name":    "Do something",
			"user_id": 1,
		})
}

func TestDestroyNotFound(t *testing.T) {
	client, database := startup.CollectionsDestroy(t)

	client.
		Post("/api/collections", `{
			"name": "Hello"
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Hello",
			"user_id": 1,
		})

	client.
		Delete("/api/collections/2").
		AssertStatus(http.StatusNotFound)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Hello",
			"user_id": 1,
		})
}

func TestDestroySomeoneElsesCollection(t *testing.T) {
	client, database := startup.CollectionsDestroy(t)

	client.
		Post("/api/collections", `{
			"name": "Looking in a mirror"
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Looking in a mirror",
			"user_id": 1,
		})

	client.ChangeUser()
	client.
		Delete("/api/collections/1").
		AssertStatus(http.StatusForbidden)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Looking in a mirror",
			"user_id": 1,
		})
}

func TestDestroyInvalidId(t *testing.T) {
	client, _ := startup.CollectionsDestroy(t)

	client.
		Delete("/api/collections/one").
		AssertStatus(http.StatusBadRequest)
}
