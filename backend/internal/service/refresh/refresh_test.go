package refresh

import (
	"github.com/EugeneNail/actum/internal/service/hash"
	"github.com/EugeneNail/actum/internal/service/tests"
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"testing"
	"time"
)

func TestRefresh(t *testing.T) {
	_, database := startup.Users(t)
	const userId = 1

	_, err := tests.DB.Exec(
		`INSERT INTO users(id, name, email, password) VALUES (?, ?,?,?)`,
		userId, "", "", "",
	)

	database.AssertHas("users", map[string]any{
		"id": userId,
	})
	service := NewService(tests.DB)
	refreshToken, err := service.MakeToken(userId)
	tests.Check(err)

	database.
		AssertCount("user_refresh_tokens", 1).
		AssertHas("user_refresh_tokens", map[string]any{
			"uuid":    hash.New(refreshToken),
			"user_id": 1,
		})

	expiredAtStart := time.Now().Add(time.Hour * 24 * 7).Add(time.Second * -1)
	expiredAtEnd := expiredAtStart.Add(time.Second * 2)

	var dbExpiredAt time.Time
	err = tests.DB.QueryRow(`SELECT expired_at FROM user_refresh_tokens WHERE user_id = ?`, userId).Scan(&dbExpiredAt)
	tests.Check(err)

	hasCorrectExpiredAt := expiredAtStart.Before(dbExpiredAt) && dbExpiredAt.Before(expiredAtEnd)
	if !hasCorrectExpiredAt {
		t.Errorf("Invalid expiration datetime")
	}
}
