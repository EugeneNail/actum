package users

import (
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/hash"
	"github.com/EugeneNail/actum/internal/service/tests"
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"strings"
	"testing"
)

func TestStoreValidData(t *testing.T) {
	client, database := startup.Users(t)

	client.
		Post("/api/users", `{
			"name": "John",
			"email": "blank@gmail.com",
			"password": "Strong123",
			"passwordConfirmation": "Strong123"
		}`).
		AssertStatus(http.StatusCreated)

	database.AssertHas("users", map[string]any{
		"name":     "John",
		"email":    "blank@gmail.com",
		"password": hash.Password("Strong123"),
	})

}

func TestStoreInvalidData(t *testing.T) {
	client, database := startup.Users(t)

	client.
		Post("/api/users", `{
			"name": "Jo",
			"email": "blankgmail.com",
			"password": "String1",
			"passwordConfirmation": ""
		}`).
		AssertStatus(http.StatusUnprocessableEntity).
		AssertHasValidationErrors([]string{"name", "email", "password", "passwordConfirmation"})

	database.AssertEmpty("users")
}

func TestStoreDuplicateEmail(t *testing.T) {
	client, database := startup.Users(t)

	input := `{
		"name": "John",
		"email": "blank@gmail.com",
		"password": "Strong123",
		"passwordConfirmation": "Strong123"
	}`

	client.
		Post("/api/users", input).
		AssertStatus(http.StatusCreated)

	client.
		Post("/api/users", input).
		AssertStatus(http.StatusUnprocessableEntity).
		AssertHasValidationErrors([]string{"email"})

	database.AssertCount("users", 1)

}

func TestStoreValidation(t *testing.T) {
	env.Load()

	tests.AssertValidationSuccess[storeInput](t, []tests.ValidationTest{
		{"name", "Name 1", "Joe"},
		{"name", "Name 2", "John"},
		{"name", "Name 3", "William"},
		{"name", "Name 4", "Bartholomew"},
		{"name", "Name 5", "Benjamin"},
		{"email", "Email 1", "user@domain.com"},
		{"email", "Email 2", "user.a.user@domain.com"},
		{"email", "Email 3", "user@106list.org"},
		{"email", "Email 4", "user_user@domain.com"},
		{"email", "Email 5", "user@domain-12.ru"},
		{"email", "Email 6", "user_user.a.user@domain.com"},
		{"password", "Password 1", "Strong123"},
		{"password", "Password 2", "VeryStrongP@ssw0rd32186"},
		{"password", "Password 3", "J7<}9*G?a\\-0"},
		{"password", "Password 4", "/Pb/>BX<82rQvW4tq!'9i1@0(e7Kzq/F?RnP<iq:ob;h#l,'%q"},
	})

	tests.AssertValidationFail[storeInput](t, []tests.ValidationTest{
		{"name", "Empty name", ""},
		{"name", "Too short name", "Jo"},
		{"name", "Too long name", strings.Repeat("Very", 5) + "LongName"},
		{"name", "Name has numbers", "John1"},
		{"name", "Name has symbols", "John's"},
		{"name", "Name has only numbers", "123"},
		{"name", "Name has only symbols", "/*-+"},
		{"email", "Empty email", ""},
		{"email", "Email has no mail", "user@"},
		{"email", "Email has no address", "@domain.com"},
		{"email", "Email has spaces", "user domain.com"},
		{"email", "Email has no domain", "Veryuser@domain."},
		{"password", "Empty password", ""},
		{"password", "Too short password", "Short"},
		{"password", "Too long password", strings.Repeat("Very", 25) + "LongPassword"},
		{"password", "Password has only numbers", "12345678"},
		{"password", "Password has only lowercase", "nomixedcase"},
		{"password", "Password has only uppercase", "NOMIXEDCASE"},
		{"password", "Password has spaces", "With spaces"},
	})
}
