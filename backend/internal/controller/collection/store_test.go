package collection

import (
	"github.com/EugeneNail/actum/internal/service/tests"
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"testing"
)

func TestValidData(t *testing.T) {
	client, database := startup.CollectionsStore(t)

	response := client.Post("/api/collections", `{
		"name": "Sport"	
	}`)

	response.AssertStatus(http.StatusCreated)
	database.AssertCount("collections", 1)
	database.AssertHas("collections", map[string]any{
		"name": "Sport",
	})
}

func TestStoreUnauthorized(t *testing.T) {
	client, database := startup.CollectionsStore(t)
	client.UnsetToken()

	response := client.Post("/api/collections", `{
		"name": "Sport"	
	}`)

	response.AssertStatus(http.StatusUnauthorized)
	database.AssertEmpty("collections")
}

func TestStoreInvalidData(t *testing.T) {
	client, database := startup.CollectionsStore(t)

	response := client.Post("/api/collections", `{
		"name": "Sp"	
	}`)

	response.AssertStatus(http.StatusUnprocessableEntity)
	response.AssertHasValidationErrors([]string{"name"})
	database.AssertEmpty("collections")
}

func TestStoreDuplicate(t *testing.T) {
	client, database := startup.CollectionsStore(t)

	client.Post("/api/collections", `{
		"name": "Sport"	
	}`)

	response := client.Post("/api/collections", `{
		"name": "Sport"	
	}`)

	response.AssertStatus(http.StatusConflict)
	response.AssertHasValidationErrors([]string{"name"})
	database.AssertCount("collections", 1)
}

func TestStoreValidation(t *testing.T) {
	tests.AssertValidationSuccess[storeInput](t, []tests.ValidationTest{
		{"Short", "name", "Run"},
		{"One word", "name", "Sport"},
		{"Multiple words", "name", "Cut nails"},
		{"Numbers", "name", "Gaming for 8 hours"},
		{"Only numbers", "name", "1263 123 6662 123"},
		{"Dash", "name", "Sleep for 3-4 hours"},
		{"Long", "name", "Go to the store for clothes between 6 am and noon"},
	})

	tests.AssertValidationFail[storeInput](t, []tests.ValidationTest{
		{"Too short", "name", "Be"},
		{"Too long", "name", "The quick brown fox jumps over the lazy dog under the old oak"},
		{"Has comma", "name", "Work yesterday, today and tomorrow"},
		{"Has period", "name", "Run. Sleep."},
		{"Has other symbols", "name", "@'!?;"},
	})
}
