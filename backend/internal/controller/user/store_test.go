package user

import (
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/tests"
	"net/http"
	"strings"
	"testing"
)

func TestStoreValidData(t *testing.T) {
	env.Load()
	t.Cleanup(cleanup)

	response := tests.Post("/api/users", t, `{
		"name": "John",
		"email": "blank@gmail.com",
		"password": "Strong123",
		"passwordConfirmation": "Strong123"
	}`)

	response.AssertStatus(http.StatusCreated)

	tests.AssertDatabaseHas("users", map[string]any{
		"name":     "John",
		"email":    "blank@gmail.com",
		"password": hashPassword("Strong123"),
	}, t)

	response.AssertHasToken()
}

func TestStoreInvalidData(t *testing.T) {
	env.Load()
	t.Cleanup(cleanup)

	response := tests.Post("/api/users", t, `{
		"name": "Jo",
		"email": "blankgmail.com",
		"password": "String1",
		"passwordConfirmation": ""
	}`)

	response.AssertStatus(http.StatusUnprocessableEntity)
	response.AssertHasValidationErrors([]string{"name", "email", "password", "passwordConfirmation"})
	response.AssertHasNoToken()
	tests.AssertTableIsEmpty("users", t)
}

func TestStoreDuplicateEmail(t *testing.T) {
	env.Load()
	t.Cleanup(cleanup)

	input := `{
		"name": "John",
		"email": "blank@gmail.com",
		"password": "Strong123",
		"passwordConfirmation": "Strong123"
	}`
	tests.Post("/api/users", t, input)
	response := tests.Post("/api/users", t, input)

	response.AssertStatus(http.StatusUnprocessableEntity)
	response.AssertHasValidationErrors([]string{"email"})
	tests.AssertDatabaseCount("users", 1, t)
	response.AssertHasNoToken()
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
