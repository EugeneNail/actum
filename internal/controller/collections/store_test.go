package collections

import (
	"github.com/EugeneNail/actum/internal/database/resource/collections"
	"github.com/EugeneNail/actum/internal/infrastructure/tests"
	"github.com/EugeneNail/actum/internal/infrastructure/tests/startup"
	"net/http"
	"testing"
)

func TestValidData(t *testing.T) {
	client, database := startup.Collections(t)

	client.
		Post("/api/collections", `{
			"name": "Спорт",
			"color": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":  "Спорт",
			"color": 1,
		})
}

func TestStoreUnauthorized(t *testing.T) {
	client, database := startup.Collections(t)
	client.UnsetToken()

	client.
		Post("/api/collections", `{
			"name": "Спорт",	
			"color": 1
		}`).
		AssertStatus(http.StatusUnauthorized)

	database.
		AssertEmpty("collections").
		AssertLacks("collections", map[string]any{
			"name":  "Спорт",
			"color": 1,
		})
}

func TestStoreInvalidData(t *testing.T) {
	client, database := startup.Collections(t)

	client.
		Post("/api/collections", `{
			"name": "Sp",	
			"color": 0
		}`).
		AssertStatus(http.StatusUnprocessableEntity).
		AssertHasValidationErrors([]string{"name"})

	database.
		AssertEmpty("collections").
		AssertLacks("collections", map[string]any{
			"name":  "Sp",
			"color": 0,
		})
}

func TestStoreDuplicate(t *testing.T) {
	client, database := startup.Collections(t)

	client.
		Post("/api/collections", `{
			"name": "Спорт",	
			"color": 1
		}`).
		AssertStatus(http.StatusCreated)

	client.
		Post("/api/collections", `{
			"name": "СпОрт",	
			"color": 1
		}`).
		AssertStatus(http.StatusConflict).
		AssertHasValidationErrors([]string{"name"})

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Спорт",
			"color":   1,
			"user_id": 1,
		})
}

func TestStoreTooMany(t *testing.T) {
	client, database := startup.Collections(t)

	database.AssertCount("users", 1).AssertHas("users", map[string]any{"id": 1})

	collections.NewFactory(1).Make(15).Insert()
	database.AssertCount("collections", 15)

	client.
		Post("/api/collections", `{
			"name": "Делать что-то",	
			"color": 2
		}`).
		AssertStatus(http.StatusConflict)

	database.
		AssertCount("collections", 15).
		AssertLacks("collections", map[string]any{
			"name":  "Делать что-то",
			"color": 2,
		})
}

func TestStoreValidation(t *testing.T) {
	tests.AssertValidationSuccess[storeInput](t, []tests.ValidationTest{
		{"name", "Short", "Бег"},
		{"name", "One word", "Спорт"},
		{"name", "Multiple words", "Уборка дома"},
		{"name", "Numbers", "Играл 8 часов"},
		{"name", "Only numbers", "1263 123 6662 123"},
		{"name", "Dash", "Спал 3-4 часа"},
		{"name", "Long", "Ходил туда в магазин"},
		{"color", "Color 1", 1},
		{"color", "Color 2", 2},
		{"color", "Color 3", 3},
		{"color", "Color 4", 4},
		{"color", "Color 5", 5},
		{"color", "Color 6", 6},
	})

	tests.AssertValidationFail[storeInput](t, []tests.ValidationTest{
		{"name", "Too short", "Be"},
		{"name", "Too long", "Поешь этих мягких булок"},
		{"name", "Has comma", "Вчера, сегодня"},
		{"name", "Has period", "Бегать. Спать."},
		{"name", "Has other symbols", "@'!?;"},
		{"color", "Less than min", 0},
		{"color", "Nonexistent", 7},
	})
}
