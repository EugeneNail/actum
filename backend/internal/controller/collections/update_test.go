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
			"name": "Exercises",
			"color": 2
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Exercises",
			"color":   2,
			"user_id": 1,
		})

	client.
		Put("/api/collections/1", `{
			"name": "Sport",
			"color": 3
		}`).
		AssertStatus(http.StatusNoContent)

	database.
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Sport",
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
			"name": "Wor",
			"color": 2
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Wor",
			"color":   2,
			"user_id": 1,
		})

	client.
		Put("/api/collections/2", `{
			"name": "Work",
			"color": 3
		}`).
		AssertStatus(http.StatusNotFound)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Wor",
			"color":   2,
			"user_id": 1,
		})
}

func TestUpdateSomeoneElsesCollection(t *testing.T) {
	client, database := startup.Collections(t)

	client.
		Post("/api/collections", `{
			"name": "Travelling",
			"color": 2
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Travelling",
			"color":   2,
			"user_id": 1,
		})

	client.ChangeUser()

	client.
		Put("/api/collections/1", `{
			"name": "LolIHackedYou",
			"color": 5
		}`).
		AssertStatus(http.StatusForbidden)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Travelling",
			"color":   2,
			"user_id": 1,
		})
}

func TestUpdateValidation(t *testing.T) {
	tests.AssertValidationSuccess[storeInput](t, []tests.ValidationTest{
		{"name", "Short", "Sit"},
		{"name", "One word", "Sleeping"},
		{"name", "Multiple words", "Wake up early"},
		{"name", "Numbers", "From 8 to 9"},
		{"name", "Only numbers", "234 4524"},
		{"name", "Dash", "Too long-short"},
		{"name", "Long", "Making my cat nice"},
		{"color", "Color 1", 1},
		{"color", "Color 2", 2},
		{"color", "Color 3", 3},
		{"color", "Color 4", 4},
		{"color", "Color 5", 5},
		{"color", "Color 6", 6},
	})

	tests.AssertValidationFail[storeInput](t, []tests.ValidationTest{
		{"name", "Too short", "Mb"},
		{"name", "Too long", "The quick brown fox jumps"},
		{"name", "Has comma", "Eating, sleeping and working"},
		{"name", "Has period", "Today is today. Tomorrow is tomorrow"},
		{"name", "Has other symbols", "@'!?;"},
		{"color", "Less than min", 0},
		{"color", "Nonexistent", 7},
	})
}
