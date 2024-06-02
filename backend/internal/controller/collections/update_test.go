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
		{"Short", "name", "Sit"},
		{"One word", "name", "Sleeping"},
		{"Multiple words", "name", "Wake up early"},
		{"Numbers", "name", "From 8 to 9"},
		{"Only numbers", "name", "234 4524"},
		{"Dash", "name", "Too long-short"},
		{"Long", "name", "Making my cat nice"},
		{"Color 1", "color", 1},
		{"Color 2", "color", 2},
		{"Color 3", "color", 3},
		{"Color 4", "color", 4},
		{"Color 5", "color", 5},
		{"Color 6", "color", 6},
	})

	tests.AssertValidationFail[storeInput](t, []tests.ValidationTest{
		{"Too short", "name", "Mb"},
		{"Too long", "name", "The quick brown fox jumps"},
		{"Has comma", "name", "Eating, sleeping and working"},
		{"Has period", "name", "Today is today. Tomorrow is tomorrow"},
		{"Has other symbols", "name", "@'!?;"},
		{"Less than min", "color", 0},
		{"Nonexistent", "color", 7},
	})
}
