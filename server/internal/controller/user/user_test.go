package user

import (
	"bytes"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/service/env"
	"net/http"
	"testing"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func cleanup() {
	err := mysql.Truncate("users")
	check(err)
}

func getUrl() string {
	return "http://127.0.0.1:" + env.Get("APP_PORT") + "/api/users"
}

func TestStoreSuccess(t *testing.T) {
	env.Load()
	t.Cleanup(cleanup)
	url := getUrl()
	body := bytes.NewReader([]byte(`{
		"name": "John",
		"email": "blank@gmail.com",
		"password": "Strong123",
		"passwordConfirmation": "Strong123"
	}`))
	response, err := http.Post(url, "application/json", body)
	check(err)

	if response.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", response.StatusCode)
	}

	count, err := mysql.GetRowCount("users")
	check(err)
	if count != 1 {
		t.Errorf("expected 1 row, got %d", count)
	}
}
