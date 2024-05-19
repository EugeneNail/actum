package activity

import (
	"github.com/EugeneNail/actum/internal/model/activities"
	"github.com/EugeneNail/actum/internal/service/tests"
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"strings"
	"testing"
)

func TestStoreValidData(t *testing.T) {
	client, database := startup.ActivitiesStore(t)

	client.
		Post("/api/collections", `{
			"name": "Habits"
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Habits",
			"user_id": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "Clean teeth",
			"icon": "add",
			"collectionId": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Clean teeth",
			"icon":          "add",
			"collection_id": 1,
			"user_id":       1,
		})
}

func TestStoreInvalidData(t *testing.T) {
	client, database := startup.ActivitiesStore(t)

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
			"name": "Very long name of the activity",
			"icon": "Trend up",
			"collectionId": -99
		}`).
		AssertStatus(http.StatusUnprocessableEntity).
		AssertHasValidationErrors([]string{"name", "collectionId"})

	database.AssertCount("activities", 0)
}

func TestStoreDuplicate(t *testing.T) {
	client, database := startup.ActivitiesStore(t)

	client.
		Post("/api/collections", `{
			"name": "Sleep"
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Sleep",
			"user_id": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "Sleep good",
			"icon": "add",
			"collectionId": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Sleep good",
			"icon":          "add",
			"collection_id": 1,
			"user_id":       1,
		})

	client.
		Post("/api/activities", `{
			"name": "sleEp goOd",
			"icon": "add",
			"collectionId": 1
		}`).
		AssertStatus(http.StatusConflict).
		AssertHasValidationErrors([]string{"name"})

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Sleep good",
			"icon":          "add",
			"collection_id": 1,
			"user_id":       1,
		})
}

func TestStoreToSomeonesCollection(t *testing.T) {
	client, database := startup.ActivitiesStore(t)

	client.
		Post("/api/collections", `{
			"name":"Household"
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Household",
			"user_id": 1,
		})

	client.ChangeUser()
	client.
		Post("/api/activities", `{
			"name": "Cut grass",
			"icon": "Person",
			"collectionId": 1
		}`).
		AssertStatus(403)

	database.AssertCount("activities", 0)
}

func TestStoreToNonexistentCollection(t *testing.T) {
	client, database := startup.ActivitiesStore(t)

	client.
		Post("/api/collections", `{
			"name": "Workout"
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Workout",
			"user_id": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "Lightweight",
			"icon": "run",
			"collectionId": 2
		}`).
		AssertStatus(http.StatusUnprocessableEntity).
		AssertHasValidationErrors([]string{"collectionId"})

	database.AssertCount("activities", 0)
}

func TestStoreTooMany(t *testing.T) {
	client, database := startup.ActivitiesStore(t)

	client.
		Post("/api/collections", `{
			"name":"Health"
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
			"icon": "hand",
			"collectionId": 1
		}`).
		AssertStatus(http.StatusConflict)

	database.AssertCount("activities", 20)
}

func TestStoreValidation(t *testing.T) {
	client, database := startup.ActivitiesStore(t)

	client.
		Post("/api/collections", `{
			"name": "Test Collection"
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"id": 1,
		})

	tests.AssertValidationSuccess[storeInput](t, []tests.ValidationTest{
		{"Short", "name", "Run"},
		{"Long", "name", "VeryEnoughLongName"},
		{"One word", "name", "Washing"},
		{"Multiple words", "name", "Wake up early"},
		{"Numbers", "name", "Wake p at 6 am"},
		{"Only numbers", "name", "123534"},
		{"Dash", "name", "Work for 9-10 hours"},
		{"One word", "icon", "add"},
		{"Multiple words", "icon", "local_gas_station"},
		{"Mixed case", "icon", "Trending_Up"},
		{"Existent collection", "collectionId", 1},
	})

	tests.AssertValidationFail[storeInput](t, []tests.ValidationTest{
		{"Too short", "name", "Ha"},
		{"Too long", "name", strings.Repeat("Very", 5) + "LongName"},
		{"Has comma", "name", "Sleep, sleep and sleep"},
		{"Period", "name", "Better. Faster. Stronger."},
		{"Other symbols", "name", "[]/\\?!"},
		{"Has space", "icon", "Bug report"},
		{"Has dash", "icon", "card-giftcard"},
		{"Nonexistent collection", "collectionId", 2},
	})
}
