package user

import (
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/model/users"
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/tests"
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
	tests.Check(err)

	tests.AssertStatus(response, http.StatusCreated, t)

	count, err := mysql.GetRowCount("users")
	tests.Check(err)
	if count != 1 {
		t.Errorf("expected 1 row, got %d", count)
		return
	}

	user, err := users.Find(1)
	tests.Check(err)
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

	tests.AssertHasToken(response, t)
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
	tests.Check(err)

	tests.AssertStatus(response, http.StatusUnprocessableEntity, t)
	tests.AssertHasValidationErrors(response, []string{"name", "email", "password", "passwordConfirmation"}, t)

	count, err := mysql.GetRowCount("users")
	tests.Check(err)
	if count != 0 {
		t.Errorf("expected no created rows, got %d", count)
		return
	}

	tests.AssertHasNoToken(response, t)
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
	tests.Check(err)
	response, err := http.Post(url, "application/json", strings.NewReader(input))
	tests.Check(err)

	tests.AssertStatus(response, http.StatusUnprocessableEntity, t)
	tests.AssertHasValidationErrors(response, []string{"email"}, t)

	count, err := mysql.GetRowCount("users")
	tests.Check(err)
	if count != 1 {
		t.Errorf("expected 1 row, got %d", count)
	}

	tests.AssertHasNoToken(response, t)
}

func TestStoreValidation(t *testing.T) {
	env.Load()

	tests.AssertValidationSuccess[storeInput](t, []tests.ValidationTest{
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
	})

	tests.AssertValidationFail[storeInput](t, []tests.ValidationTest{
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
	})
}
