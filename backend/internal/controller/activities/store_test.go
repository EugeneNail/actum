package activities

import (
	"github.com/EugeneNail/actum/internal/database/resource/activities"
	"github.com/EugeneNail/actum/internal/service/tests"
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"testing"
)

func TestStoreValidData(t *testing.T) {
	client, database := startup.Activities(t)

	client.
		Post("/api/collections", `{
			"name": "Привычки",
			"color": 3
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Привычки",
			"color":   3,
			"user_id": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "Чистка зубов",
			"icon": 100,
			"collectionId": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Чистка зубов",
			"icon":          100,
			"collection_id": 1,
			"user_id":       1,
		})
}

func TestStoreInvalidData(t *testing.T) {
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
			"name": "Очень длинное название активности",
			"icon": 1001,
			"collectionId": -99
		}`).
		AssertStatus(http.StatusUnprocessableEntity).
		AssertHasValidationErrors([]string{"name", "icon", "collectionId"})

	database.AssertCount("activities", 0)
}

func TestStoreDuplicate(t *testing.T) {
	client, database := startup.Activities(t)

	client.
		Post("/api/collections", `{
			"name": "Сон",
			"color": 3
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Сон",
			"color":   3,
			"user_id": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "Хорошо поспал",
			"icon": 700,
			"collectionId": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Хорошо поспал",
			"icon":          700,
			"collection_id": 1,
			"user_id":       1,
		})

	client.
		Post("/api/activities", `{
			"name": "ХоРоШо поспаЛ",
			"icon": 700,
			"collectionId": 1
		}`).
		AssertStatus(http.StatusConflict).
		AssertHasValidationErrors([]string{"name"})

	database.
		AssertCount("activities", 1).
		AssertHas("activities", map[string]any{
			"name":          "Хорошо поспал",
			"icon":          700,
			"collection_id": 1,
			"user_id":       1,
		})
}

func TestStoreToSomeonesCollection(t *testing.T) {
	client, database := startup.Activities(t)

	client.
		Post("/api/collections", `{
			"name":"Хозяйство",
			"color": 3
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Хозяйство",
			"color":   3,
			"user_id": 1,
		})

	client.ChangeUser()
	client.
		Post("/api/activities", `{
			"name": "Косил траву",
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
			"name": "Тренировки",
			"color": 3
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Тренировки",
			"color":   3,
			"user_id": 1,
		})

	client.
		Post("/api/activities", `{
			"name": "Многоповторные",
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
			"name":"Здоровье",
			"color": 3
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Здоровье",
			"user_id": 1,
		})

	activities.NewFactory(1, 1).Make(20).Insert()

	database.
		AssertCount("activities", 20).
		AssertHas("activities", map[string]any{"collection_id": 1})

	client.
		Post("/api/activities", `{
			"name": "Мыл руки",
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
		{"name", "Short", "Бег"},
		{"name", "Long", "Название  активности"},
		{"name", "One word", "Душ"},
		{"name", "Multiple words", "Встал рано"},
		{"name", "Numbers", "Встал в 6 утра"},
		{"name", "Only numbers", "123534"},
		{"name", "Dash", "Работал 9-10 часов"},
		{"icon", "First group", 100},
		{"icon", "Ninth group", 903},
		{"icon", "Third group", 333},
		{"collectionId", "Existent collection", 1},
	})

	tests.AssertValidationFail[storeInput](t, []tests.ValidationTest{
		{"name", "Too short", "Не"},
		{"name", "Too long", "Очень длинное название"},
		{"name", "Has comma", "Спать, спать и спать"},
		{"name", "Period", "Лучше. Быстрее."},
		{"name", "Other symbols", "[]/\\?!"},
		{"icon", "Zero group", 99},
		{"icon", "Negative group", -100},
		{"icon", "Nonexistent group", 1001},
		{"collectionId", "Nonexistent collection", 2},
	})
}
