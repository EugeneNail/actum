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

	test.AssertHasToken(response, t)
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

	test.AssertHasNoToken(response, t)
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

	test.AssertHasNoToken(response, t)
}

func TestStoreValidation(t *testing.T) {
	env.Load()

	successes := []test.ValidationTest{
		{"Name 1", "name", "Joe"},
		{"Name 2", "name", "John"},
		{"Name 3", "name", "William"},
		{"Name 4", "name", "Bartholomew"},
		{"Name 5", "name", "Benjamin"},
		{"Email 1", "email", "user@domain.com"},
		{"Email 2", "email", "user.a.user@domain.com"},
		{"Email 3", "email", "user@106list.org"},
		{"Email 4", "email", "user_user@domain.com"},
		{"Email 5", "email", "user@domain-12.ru"},
		{"Email 6", "email", "user_user.a.user@domain.com"},
		{"Password 1", "password", "Strong123"},
		{"Password 2", "password", "VeryStrongP@ssw0rd32186"},
		{"Password 3", "password", "J7<}9*G?a\\-0"},
		{"Password 4", "password", "/Pb/>BX<82rQvW4tq!'9i1@0(e7Kzq/F?RnP<iq:ob;h#l,'%q"},
	}
	for _, tableTest := range successes {
		t.Run(tableTest.Name, func(t *testing.T) {
			test.AssertValidationSuccess[storeInput](tableTest, t)
		})
	}

	fails := []test.ValidationTest{
		{"Empty name", "name", ""},
		{"Too short name", "name", "Jo"},
		{"Too long name", "name", strings.Repeat("Very", 5) + "LongName"},
		{"Name has numbers", "name", "John1"},
		{"Name has symbols", "name", "John's"},
		{"Name has only numbers", "name", "123"},
		{"Name has only symbols", "name", "/*-+"},
		{"Empty email", "email", ""},
		{"Email has no mail", "email", "user@"},
		{"Email has no address", "email", "@domain.com"},
		{"Email has spaces", "email", "user domain.com"},
		{"Email has no domain", "email", "Veryuser@domain."},
		{"Empty password", "password", ""},
		{"Too short password", "password", "Short"},
		{"Too long password", "password", strings.Repeat("Very", 25) + "LongPassword"},
		{"Password has only numbers", "password", "12345678"},
		{"Password has only lowercase", "password", "nomixedcase"},
		{"Password has only uppercase", "password", "NOMIXEDCASE"},
		{"Password has spaces", "password", "With spaces"},
	}
	for _, tableTest := range fails {
		t.Run(tableTest.Name, func(t *testing.T) {
			test.AssertValidationFail[storeInput](tableTest, t)
		})
	}
}
