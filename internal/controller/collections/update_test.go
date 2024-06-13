package collections

import (
	"github.com/EugeneNail/actum/internal/service/tests"
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"testing"
)

func TestUpdateValidData(t *testing.T) {
	client, database := startup.Collections(t)

	client.
		Post("/api/collections", `{
			"name": "Упражнения",
			"color": 2
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Упражнения",
			"color":   2,
			"user_id": 1,
		})

	client.
		Put("/api/collections/1", `{
			"name": "Спорт",
			"color": 3
		}`).
		AssertStatus(http.StatusNoContent)

	database.
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Спорт",
			"color":   3,
			"user_id": 1,
		}).
		AssertCount("collections", 1)
}

func TestUpdateInvalidData(t *testing.T) {
	client, database := startup.Collections(t)

	client.
		Post("/api/collections", `{
			"name": "30",
			"color": 77
		}`).
		AssertStatus(http.StatusUnprocessableEntity).
		AssertHasValidationErrors([]string{"name"})

	database.
		AssertCount("collections", 0).
		AssertLacks("collections", map[string]any{
			"name":  "30",
			"color": 77,
		})
}

func TestUpdateNotFound(t *testing.T) {
	client, database := startup.Collections(t)

	client.
		Post("/api/collections", `{
			"name": "Раб",
			"color": 2
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Раб",
			"color":   2,
			"user_id": 1,
		})

	client.
		Put("/api/collections/2", `{
			"name": "Работа",
			"color": 3
		}`).
		AssertStatus(http.StatusNotFound)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Раб",
			"color":   2,
			"user_id": 1,
		})
}

func TestUpdateSomeoneElsesCollection(t *testing.T) {
	client, database := startup.Collections(t)

	client.
		Post("/api/collections", `{
			"name": "Путешествия",
			"color": 2
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Путешествия",
			"color":   2,
			"user_id": 1,
		})

	client.ChangeUser()

	client.
		Put("/api/collections/1", `{
			"name": "Взломал тебя",
			"color": 5
		}`).
		AssertStatus(http.StatusForbidden)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Путешествия",
			"color":   2,
			"user_id": 1,
		})
}

func TestUpdateValidation(t *testing.T) {
	tests.AssertValidationSuccess[storeInput](t, []tests.ValidationTest{
		{"name", "Short", "Сон"},
		{"name", "One word", "Сидение"},
		{"name", "Multiple words", "Вставание рано"},
		{"name", "Numbers", "От 8 до 9"},
		{"name", "Only numbers", "234 4524"},
		{"name", "Dash", "С дефисом - во"},
		{"name", "Long", "Принарядил   котейку"},
		{"color", "Color 1", 1},
		{"color", "Color 2", 2},
		{"color", "Color 3", 3},
		{"color", "Color 4", 4},
		{"color", "Color 5", 5},
		{"color", "Color 6", 6},
	})

	tests.AssertValidationFail[storeInput](t, []tests.ValidationTest{
		{"name", "Too short", "Мб"},
		{"name", "Too long", "Тут написано что-то очень длинное"},
		{"name", "Has comma", "Ем, сплю, работаю"},
		{"name", "Has period", "Вчера. Сегодня"},
		{"name", "Has other symbols", "@'!?;"},
		{"color", "Less than min", 0},
		{"color", "Nonexistent", 7},
	})
}
