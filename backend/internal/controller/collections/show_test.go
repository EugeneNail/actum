package collections

import (
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"testing"
)

func TestShow(t *testing.T) {
	client, database := startup.Collections(t)

	client.
		Post("/api/collections", `{
			"name": "Workout",
			"color": 6
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("collections", 1).
		AssertHas("collections", map[string]any{
			"name":    "Workout",
			"color":   6,
			"user_id": 1,
		})

	client.
		Get("/api/collections/1").
		AssertStatus(http.StatusOK)
}
