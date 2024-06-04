package activities

import (
	"github.com/EugeneNail/actum/internal/database/resource/activities"
	"github.com/EugeneNail/actum/internal/service/tests"
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"strings"
	"testing"
)

func TestStoreValidData(t *testing.T) {
	client, database := startup.Activities(t)

	client.
		Post("/api/collections", `{
			"name": "Habits",
			"color": 3
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Habits",
			"color":   3,
			"user_id": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "Clean teeth",
			"icon": 100,
			"collectionId": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Clean teeth",
			"icon":          100,
			"collection_id": 1,
			"user_id":       1,
		})
}

func TestStoreInvalidData(t *testing.T) {
	client, database := startup.Activities(t)

	client.
		Post("/api/collections", `{
			"name": "Work",
			"color": 3
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Work",
			"color":   3,
			"user_id": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "Very long name of the activity",
			"icon": 1001,
			"collectionId": -99
		}`).
		AssertStatus(http.StatusUnprocessableEntity).
		AssertHasValidationErrors([]string{"name", "collectionId"})

	database.AssertCount("activities", 0)
}

func TestStoreDuplicate(t *testing.T) {
	client, database := startup.Activities(t)

	client.
		Post("/api/collections", `{
			"name": "Sleep",
			"color": 3
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Sleep",
			"color":   3,
			"user_id": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "Sleep good",
			"icon": 700,
			"collectionId": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Sleep good",
			"icon":          700,
			"collection_id": 1,
			"user_id":       1,
		})

	client.
		Post("/api/activities", `{
			"name": "sleEp goOd",
			"icon": 700,
			"collectionId": 1
		}`).
		AssertStatus(http.StatusConflict).
		AssertHasValidationErrors([]string{"name"})

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Sleep good",
			"icon":          700,
			"collection_id": 1,
			"user_id":       1,
		})
}

func TestStoreToSomeonesCollection(t *testing.T) {
	client, database := startup.Activities(t)

	client.
		Post("/api/collections", `{
			"name":"Household",
			"color": 3
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Household",
			"color":   3,
			"user_id": 1,
		})

	client.ChangeUser()
	client.
		Post("/api/activities", `{
			"name": "Cut grass",
			"icon": 610,
			"collectionId": 1
		}`).
		AssertStatus(403)

	database.AssertCount("activities", 0)
}

func TestStoreToNonexistentCollection(t *testing.T) {
	client, database := startup.Activities(t)

	client.
		Post("/api/collections", `{
			"name": "Workout",
			"color": 3
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Workout",
			"color":   3,
			"user_id": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "Lightweight",
			"icon": 112,
			"collectionId": 2
		}`).
		AssertStatus(http.StatusUnprocessableEntity).
		AssertHasValidationErrors([]string{"collectionId"})

	database.AssertCount("activities", 0)
}

func TestStoreTooMany(t *testing.T) {
	client, database := startup.Activities(t)

	client.
		Post("/api/collections", `{
			"name":"Health",
			"color": 3
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Health",
			"user_id": 1,
		})

	activities.NewFactory(1, 1).Make(20).Insert()

	database.
		AssertCount("activities", 20).
		AssertHas("activities", map[string]any{"collection_id": 1})

	client.
		Post("/api/activities", `{
			"name": "Wash hands",
			"icon": 502,
			"collectionId": 1
		}`).
		AssertStatus(http.StatusConflict)

	database.AssertCount("activities", 20)
}

func TestStoreValidation(t *testing.T) {
	client, database := startup.Activities(t)

	client.
		Post("/api/collections", `{
			"name": "Test Collection",
			"color": 3
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"id": 1,
		})

	tests.AssertValidationSuccess[storeInput](t, []tests.ValidationTest{
		{"name", "Short", "Run"},
		{"name", "Long", "VeryEnoughLongName"},
		{"name", "One word", "Washing"},
		{"name", "Multiple words", "Wake up early"},
		{"name", "Numbers", "Wake p at 6 am"},
		{"name", "Only numbers", "123534"},
		{"name", "Dash", "Work for 9-10 hours"},
		{"icon", "First group", 100},
		{"icon", "Ninth group", 903},
		{"icon", "Third group", 333},
		{"collectionId", "Existent collection", 1},
	})

	tests.AssertValidationFail[storeInput](t, []tests.ValidationTest{
		{"name", "Too short", "Ha"},
		{"name", "Too long", strings.Repeat("Very", 5) + "LongName"},
		{"name", "Has comma", "Sleep, sleep and sleep"},
		{"name", "Period", "Better. Faster. Stronger."},
		{"name", "Other symbols", "[]/\\?!"},
		{"icon", "Zero group", 99},
		{"icon", "Negative group", -100},
		{"icon", "Nonexistent group", 1001},
		{"collectionId", "Nonexistent collection", 2},
	})
}
