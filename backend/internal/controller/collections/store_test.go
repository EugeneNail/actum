package collections

import (
	"github.com/EugeneNail/actum/internal/database/resource/collections"
	"github.com/EugeneNail/actum/internal/service/tests"
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"testing"
)

func TestValidData(t *testing.T) {
	client, database := startup.CollectionsStore(t)

	client.
		Post("/api/collections", `{
			"name": "Sport"	
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name": "Sport",
		})
}

func TestStoreUnauthorized(t *testing.T) {
	client, database := startup.CollectionsStore(t)
	client.UnsetToken()

	client.
		Post("/api/collections", `{
			"name": "Sport"	
		}`).
		AssertStatus(http.StatusUnauthorized)

	database.
		AssertEmpty("collections").
		AssertLacks("collections", map[string]any{
			"name": "Sport",
		})
}

func TestStoreInvalidData(t *testing.T) {
	client, database := startup.CollectionsStore(t)

	client.
		Post("/api/collections", `{
			"name": "Sp"	
		}`).
		AssertStatus(http.StatusUnprocessableEntity).
		AssertHasValidationErrors([]string{"name"})

	database.
		AssertEmpty("collections").
		AssertLacks("collections", map[string]any{
			"name": "Sp",
		})
}

func TestStoreDuplicate(t *testing.T) {
	client, database := startup.CollectionsStore(t)

	client.
		Post("/api/collections", `{
			"name": "Sport"	
		}`).
		AssertStatus(http.StatusCreated)

	client.
		Post("/api/collections", `{
			"name": "SpOrt"	
		}`).
		AssertStatus(http.StatusConflict).
		AssertHasValidationErrors([]string{"name"})

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Sport",
			"user_id": 1,
		})
}

func TestStoreTooMany(t *testing.T) {
	client, database := startup.CollectionsStore(t)

	database.AssertCount("users", 1).AssertHas("users", map[string]any{"id": 1})

	collections.NewFactory(1).Make(15).Insert()
	database.AssertCount("collections", 15)

	client.
		Post("/api/collections", `{
			"name": "Do something"	
		}`).
		AssertStatus(http.StatusConflict)

	database.
		AssertCount("collections", 15).
		AssertLacks("collections", map[string]any{
			"name": "Do something",
		})
}

func TestStoreValidation(t *testing.T) {
	tests.AssertValidationSuccess[storeInput](t, []tests.ValidationTest{
		{"Short", "name", "Run"},
		{"One word", "name", "Sport"},
		{"Multiple words", "name", "Cut nails"},
		{"Numbers", "name", "Gaming for 8 hours"},
		{"Only numbers", "name", "1263 123 6662 123"},
		{"Dash", "name", "Sleep for 3-4 hours"},
		{"Long", "name", "Go to the store for"},
	})

	tests.AssertValidationFail[storeInput](t, []tests.ValidationTest{
		{"Too short", "name", "Be"},
		{"Too long", "name", "The quick brown fox jumps"},
		{"Has comma", "name", "Work tomorrow, today"},
		{"Has period", "name", "Run. Sleep."},
		{"Has other symbols", "name", "@'!?;"},
	})
}
