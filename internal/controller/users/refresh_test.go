package users

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/resource/users"
	"github.com/EugeneNail/actum/internal/service/hash"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/refresh"
	"github.com/EugeneNail/actum/internal/service/tests"
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"github.com/EugeneNail/actum/internal/service/uuid"
	"net/http"
	"testing"
	"time"
)

func TestRefreshValidData(t *testing.T) {
	client, database := startup.Users(t)

	var tokens tests.Tokens
	client.
		Post("/api/users", `{
			"name": "Базилио",
			"email": "basil1984@gmail.com",
			"password": "Strong123",
			"passwordConfirmation": "Strong123"
		}`).
		AssertStatus(http.StatusCreated).
		ReadData(&tokens)

	database.
		AssertCount("users", 1).
		AssertHas("users", map[string]any{
			"id":    1,
			"name":  "Базилио",
			"email": "basil1984@gmail.com",
		}).
		AssertCount("user_refresh_tokens", 1).
		AssertHas("user_refresh_tokens", map[string]any{
			"uuid":    hash.New(tokens.Refresh),
			"user_id": 1,
		})

	isValid, err := refresh.NewService(tests.DB).IsValid(tokens.Refresh, 1)
	tests.Check(err)
	if !isValid {
		t.Errorf("Returned refresh token is invalid")
		return
	}

	var accessToken string
	client.
		Post("/api/users/refresh-token", fmt.Sprintf(`{
			"uuid": "%s",
			"userId": 1
			}`, tokens.Refresh),
		).
		AssertStatus(http.StatusOK).
		ReadData(&accessToken)

	if !jwt.IsValid(accessToken) {
		t.Errorf("Returned access token is invalid")
		return
	}
}

func TestRefreshInvalidData(t *testing.T) {
	client, database := startup.Users(t)

	database.AssertCount("users", 0)

	client.
		Post("/api/users/refresh-token", `{
			"uuid": "12345678-1234-1234-1234-1234567812349",
			"userId": 1
		}`).
		AssertStatus(http.StatusUnprocessableEntity).
		AssertHasValidationErrors([]string{"uuid", "userId"})
}

func TestRefreshExpiredToken(t *testing.T) {
	client, database := startup.Users(t)

	var tokens tests.Tokens
	client.
		Post("/api/users", `{
			"name": "Космос",
			"email": "lookatme@iamhere.cc",
			"password": "Strong123",
			"passwordConfirmation": "Strong123"
		}`).
		AssertStatus(http.StatusCreated).
		ReadData(&tokens)

	database.
		AssertCount("users", 1).
		AssertHas("users", map[string]any{
			"id":    1,
			"name":  "Космос",
			"email": "lookatme@iamhere.cc",
		}).
		AssertCount("user_refresh_tokens", 1).
		AssertHas("user_refresh_tokens", map[string]any{
			"uuid":    hash.New(tokens.Refresh),
			"user_id": 1,
		})

	isValid, err := refresh.NewService(tests.DB).IsValid(tokens.Refresh, 1)
	tests.Check(err)
	if !isValid {
		t.Errorf("Returned refresh token is invalid")
		return
	}

	_, err = tests.DB.Exec(
		`UPDATE user_refresh_tokens SET expired_at = ? WHERE user_id = ?`,
		time.Now().Add(time.Hour*-1), 1,
	)
	tests.Check(err)

	client.
		Post("/api/users/refresh-token", fmt.Sprintf(`{
			"uuid": "%s",
			"userId": 1
			}`, tokens.Refresh),
		).
		AssertStatus(http.StatusUnauthorized)
}

func TestRefreshValidation(t *testing.T) {
	client, database := startup.Users(t)
	client.ChangeUser()

	database.
		AssertHas("users", map[string]any{
			"id": 1,
		}).
		AssertLacks("users", map[string]any{
			"id": 2,
		})

	tests.AssertValidationSuccess[users.refreshInput](t, []tests.ValidationTest{
		{"uuid", "Uuid 1", "6d56e42c-7f75-42d9-be95-2493652937d5"},
		{"uuid", "Uuid 2", "39bef25d-8a12-46c8-92bc-12ddf369de34"},
		{"uuid", "Uuid 3", "990ade10-74f0-42ad-bcc8-a1452647dde5"},
		{"uuid", "Uuid 4", "cd185189-7e53-4d87-9c5a-de976736cec7"},
		{"uuid", "Uuid 5", "563db045-0723-4333-b9e9-6b6ec6a8e61b"},
		{"uuid", "Uuid 6", uuid.New()},
		{"userId", "Existing user", 1},
	})

	tests.AssertValidationFail[users.refreshInput](t, []tests.ValidationTest{
		{"uuid", "Short", "d1f1ef3b-e8e9-436d-a6b3-036c715b756"},
		{"uuid", "Long", "162b519a-f026-4b9e-a2e9-0143be261a8da"},
		{"uuid", "Invalid characters", "f34bf4bd-7376-4mdf-8da6-74e3t6609535"},
		{"uuid", "Invalid format", "4a90-30a8c58a-48a3-8813-e68ad139a51e"},
		{"userId", "Nonexistent user", 2},
	})
}
