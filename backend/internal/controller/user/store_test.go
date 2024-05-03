package user

import (
	"encoding/json"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/model/users"
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/test"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestStoreValidData(t *testing.T) {
	env.Load()
	t.Cleanup(cleanup)
	response, err := http.Post(getUrl(), "application/json", strings.NewReader(`{
		"name": "John",
		"email": "blank@gmail.com",
		"password": "Strong123",
		"passwordConfirmation": "Strong123"
	}`))
	check(err)

	if response.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", response.StatusCode)
	}

	count, err := mysql.GetRowCount("users")
	check(err)
	if count != 1 {
		t.Errorf("expected 1 row, got %d", count)
		return
	}

	user, err := users.Find(1)
	check(err)
	if user.Name != "John" {
		t.Errorf("expected the name field John, got %s", user.Name)
	}

	if user.Email != "blank@gmail.com" {
		t.Errorf("expected the email field blank@gmail.com, got %s", user.Email)
	}

	hashedPassword := hashPassword("Strong123")
	if user.Password != hashedPassword {
		t.Errorf("expected the password field %s, got %s", hashedPassword, user.Name)
	}
}

func TestStoreInvalidData(t *testing.T) {
	env.Load()
	t.Cleanup(cleanup)
	response, err := http.Post(getUrl(), "application/json", strings.NewReader(`{
		"name": "Jo",
		"email": "blankgmail.com",
		"password": "String1",
		"passwordConfirmation": ""
	}`))
	check(err)

	if response.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", response.StatusCode)
	}

	var validationMessages map[string]string
	data, err := io.ReadAll(response.Body)
	check(err)
	err = json.Unmarshal(data, &validationMessages)
	check(err)
	for _, field := range []string{"name", "email", "password", "passwordConfirmation"} {
		if _, exists := validationMessages[field]; !exists {
			t.Errorf(`expected validation error for field "%s" to be present`, field)
		}
	}

	count, err := mysql.GetRowCount("users")
	check(err)
	if count != 0 {
		t.Errorf("expected no created rows, got %d", count)
		return
	}
}

func TestStoreDuplicateEmail(t *testing.T) {
	env.Load()
	t.Cleanup(cleanup)
	url := getUrl()
	input := `{
		"name": "John",
		"email": "blank@gmail.com",
		"password": "Strong123",
		"passwordConfirmation": "Strong123"
	}`
	_, err := http.Post(url, "application/json", strings.NewReader(input))
	check(err)
	response, err := http.Post(url, "application/json", strings.NewReader(input))
	check(err)

	if response.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", response.StatusCode)
		return
	}

	validationMessages := make(map[string]string)
	data, err := io.ReadAll(response.Body)
	check(err)
	err = json.Unmarshal(data, &validationMessages)
	check(err)
	if _, exists := validationMessages["email"]; !exists {
		t.Error("expected validation error for the email field to be present")
	}

	count, err := mysql.GetRowCount("users")
	check(err)
	if count != 1 {
		t.Errorf("expected 1 row, got %d", count)
	}
}

func TestStoreValidation(t *testing.T) {
	env.Load()

	successes := []test.Field{
		{"name", "Joe"},
		{"name", "John"},
		{"name", "William"},
		{"name", "Bartholomew"},
		{"name", "Benjamin"},
		{"email", "user@domain.com"},
		{"email", "user.a.user@domain.com"},
		{"email", "user@106list.org"},
		{"email", "user_user@domain.com"},
		{"email", "user@domain-12.ru"},
		{"email", "user_user.a.user@domain.com"},
		{"password", "Strong123"},
		{"password", "VeryStrongP@ssw0rd32186"},
		{"password", "J7<}9*G?a\\-0"},
		{"password", "/Pb/>BX<82rQvW4tq!'9i1@0(e7Kzq/F?RnP<iq:ob;h#l,'%q"},
	}
	for _, field := range successes {
		test.AssertValidationSuccess[storeInput](field, t)
	}

	fails := []test.Field{
		{"name", ""},
		{"name", "Jo"},
		{"name", strings.Repeat("Very", 5) + "LongName"},
		{"name", "John1"},
		{"name", "John's"},
		{"name", "123"},
		{"name", "/*-+"},
		{"email", ""},
		{"email", "user@"},
		{"email", "@domain.com"},
		{"email", "user domain.com"},
		{"email", "Veryuser@domain."},
		{"password", ""},
		{"password", "Short"},
		{"password", "12345678"},
		{"password", "nomixedcase"},
		{"password", "NOMIXEDCASE"},
		{"password", "With spaces"},
		{"password", strings.Repeat("Very", 25) + "LongPassword"},
	}
	for _, field := range fails {
		test.AssertValidationFail[storeInput](field, t)
	}
}
