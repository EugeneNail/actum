package collections

import (
	"github.com/EugeneNail/actum/internal/service/tests"
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"testing"
)

func TestUpdateValidData(t *testing.T) {
	client, database := startup.CollectionsUpdate(t)

	client.
		Post("/api/collections", `{
			"name": "Exercises"
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Exercises",
			"user_id": 1,
		})

	client.
		Put("/api/collections/1", `{
			"name": "Sport"
		}`).
		AssertStatus(http.StatusNoContent)

	database.
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Sport",
			"user_id": 1,
		}).
		AssertCount("collections", 1)
}

func TestUpdateInvalidData(t *testing.T) {
	client, database := startup.CollectionsUpdate(t)

	client.
		Post("/api/collections", `{
			"name": "30"
		}`).
		AssertStatus(http.StatusUnprocessableEntity).
		AssertHasValidationErrors([]string{"name"})

	database.
		AssertCount("collections", 0).
		AssertLacks("collections", map[string]any{
			"name": "30",
		})
}

func TestUpdateNotFound(t *testing.T) {
	client, database := startup.CollectionsUpdate(t)

	client.
		Post("/api/collections", `{
			"name": "Wor"
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Wor",
			"user_id": 1,
		})

	client.
		Put("/api/collections/2", `{
			"name": "Work"
		}`).
		AssertStatus(http.StatusNotFound)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Wor",
			"user_id": 1,
		})
}

func TestUpdateDuplicate(t *testing.T) {
	client, database := startup.CollectionsUpdate(t)

	client.
		Post("/api/collections", `{
			"name": "sport"	
		}`).
		AssertStatus(http.StatusCreated)

	client.
		Put("/api/collections/1", `{
			"name": "SpOrt"	
		}`).
		AssertStatus(http.StatusConflict).
		AssertHasValidationErrors([]string{"name"})

	database.
		AssertCount("collections", 1).
		AssertLacks("collections", map[string]any{
			"name": "SpOrt",
		})
}

func TestUpdateSomeoneElsesCollection(t *testing.T) {
	client, database := startup.CollectionsUpdate(t)

	client.
		Post("/api/collections", `{
			"name": "Travelling"
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Travelling",
			"user_id": 1,
		})

	client.ChangeUser()

	client.
		Put("/api/collections/1", `{
			"name": "LolIHackedYou"
		}`).
		AssertStatus(http.StatusForbidden)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"id":      1,
			"name":    "Travelling",
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
	})

	tests.AssertValidationFail[storeInput](t, []tests.ValidationTest{
		{"Too short", "name", "Mb"},
		{"Too long", "name", "The quick brown fox jumps"},
		{"Has comma", "name", "Eating, sleeping and working"},
		{"Has period", "name", "Today is today. Tomorrow is tomorrow"},
		{"Has other symbols", "name", "@'!?;"},
	})
}
