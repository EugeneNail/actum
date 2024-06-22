package users

import (
	"github.com/EugeneNail/actum/internal/infrastructure/tests/startup"
	"net/http"
	"testing"
)

func TestLogout(t *testing.T) {
	client, database := startup.Users(t)
	client.ChangeUser()

	database.
		AssertCount("user_refresh_tokens", 1).
		AssertHas("user_refresh_tokens", map[string]any{
			"user_id": 1,
		}).
		AssertCount("users", 1).
		AssertHas("users", map[string]any{
			"id": 1,
		})

	client.
		Post("/api/users/logout", "").
		AssertStatus(http.StatusNoContent)

	database.AssertCount("user_refresh_tokens", 0)
}
