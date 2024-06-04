package collections

import (
	"github.com/EugeneNail/actum/internal/database/resource/collections"
	"github.com/EugeneNail/actum/internal/service/tests"
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"testing"
)

func TestValidData(t *testing.T) {
	client, database := startup.Collections(t)

	client.
		Post("/api/collections", `{
			"name": "Sport",
			"color": 1
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":  "Sport",
			"color": 1,
		})
}

func TestStoreUnauthorized(t *testing.T) {
	client, database := startup.Collections(t)
	client.UnsetToken()

	client.
		Post("/api/collections", `{
			"name": "Sport",	
			"color": 1
		}`).
		AssertStatus(http.StatusUnauthorized)

	database.
		AssertEmpty("collections").
		AssertLacks("collections", map[string]any{
			"name":  "Sport",
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
			"name": "Sport",	
			"color": 1
		}`).
		AssertStatus(http.StatusCreated)

	client.
		Post("/api/collections", `{
			"name": "SpOrt",	
			"color": 1
		}`).
		AssertStatus(http.StatusConflict).
		AssertHasValidationErrors([]string{"name"})

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Sport",
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
			"name": "Do something",	
			"color": 2
		}`).
		AssertStatus(http.StatusConflict)

	database.
		AssertCount("collections", 15).
		AssertLacks("collections", map[string]any{
			"name":  "Do something",
			"color": 2,
		})
}

func TestStoreValidation(t *testing.T) {
	tests.AssertValidationSuccess[storeInput](t, []tests.ValidationTest{
		{"name", "Short", "Run"},
		{"name", "One word", "Sport"},
		{"name", "Multiple words", "Cut nails"},
		{"name", "Numbers", "Gaming for 8 hours"},
		{"name", "Only numbers", "1263 123 6662 123"},
		{"name", "Dash", "Sleep for 3-4 hours"},
		{"name", "Long", "Go to the store for"},
		{"color", "Color 1", 1},
		{"color", "Color 2", 2},
		{"color", "Color 3", 3},
		{"color", "Color 4", 4},
		{"color", "Color 5", 5},
		{"color", "Color 6", 6},
	})

	tests.AssertValidationFail[storeInput](t, []tests.ValidationTest{
		{"name", "Too short", "Be"},
		{"name", "Too long", "The quick brown fox jumps"},
		{"name", "Has comma", "Work tomorrow, today"},
		{"name", "Has period", "Run. Sleep."},
		{"name", "Has other symbols", "@'!?;"},
		{"color", "Less than min", 0},
		{"color", "Nonexistent", 7},
	})
}
