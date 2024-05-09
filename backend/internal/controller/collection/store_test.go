package collection

import (
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/tests"
	"github.com/EugeneNail/actum/internal/service/tests/cleanup"
	"net/http"
	"testing"
)

func TestValidData(t *testing.T) {
	env.Load()
	t.Cleanup(cleanup.StoreCollections)
	client := tests.NewClient(t)

	response := client.Post("/api/collections", `{
		"name": "Sport"	
	}`)

	response.AssertStatus(http.StatusCreated)
	tests.AssertDatabaseCount("collections", 1, t)
	tests.AssertDatabaseHas("collections", map[string]any{
		"name": "Sport",
	}, t)
}

func TestStoreUnauthorized(t *testing.T) {
	env.Load()
	t.Cleanup(cleanup.StoreCollections)
	client := tests.NewClientWithoutAuth(t)

	response := client.Post("/api/collections", `{
		"name": "Sport"	
	}`)

	response.AssertStatus(http.StatusUnauthorized)
	tests.AssertTableIsEmpty("collections", t)
}

func TestStoreInvalidData(t *testing.T) {
	env.Load()
	t.Cleanup(cleanup.StoreCollections)
	client := tests.NewClient(t)

	response := client.Post("/api/collections", `{
		"name": "Sp"	
	}`)

	response.AssertStatus(http.StatusUnprocessableEntity)
	response.AssertHasValidationErrors([]string{"name"})
	tests.AssertTableIsEmpty("collections", t)
}

func TestStoreDuplicate(t *testing.T) {
	env.Load()
	t.Cleanup(cleanup.StoreCollections)
	client := tests.NewClient(t)

	client.Post("/api/collections", `{
		"name": "Sport"	
	}`)

	response := client.Post("/api/collections", `{
		"name": "Sport"	
	}`)

	response.AssertStatus(http.StatusConflict)
	response.AssertHasValidationErrors([]string{"name"})
	tests.AssertDatabaseCount("collections", 1, t)
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
