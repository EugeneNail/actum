package users

import (
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/hash"
	"github.com/EugeneNail/actum/internal/service/tests"
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"testing"
)

func TestLoginValidData(t *testing.T) {
	client, database := startup.Users(t)

	client.
		Post("/api/users", `{
			"name": "John",
			"email": "JODAME3394@agafx.com",
			"password": "Strong123",
			"passwordConfirmation": "Strong123"
		}`).
		AssertStatus(http.StatusCreated)

	database.AssertHas("users", map[string]any{
		"name":     "John",
		"email":    "jodame3394@agafx.com",
		"password": hash.Password("Strong123"),
	})

	client.
		Post("/api/users/login", `{
			"email": "jodAME3394@agafx.com",
			"password": "Strong123"
		}`).
		AssertStatus(http.StatusOK)

	database.
		AssertCount("users", 1).
		AssertHas("users", map[string]any{
			"name":     "John",
			"email":    "jodame3394@agafx.com",
			"password": hash.Password("Strong123"),
		})
}

func TestLoginInvalidData(t *testing.T) {
	client, database := startup.Users(t)

	client.
		Post("/api/users/login", `{
			"email": "yibewek618goulink.com",
			"password": "v9"
		}`).
		AssertStatus(http.StatusUnprocessableEntity).
		AssertHasValidationErrors([]string{"email", "password"})

	database.AssertCount("users", 0)
}

func TestLoginIncorrectEmail(t *testing.T) {
	client, database := startup.Users(t)

	client.
		Post("/api/users", `{
			"name": "William",
			"email": "doleya5976@agafx.com",
			"password": "w24V,KY$f2YSIPQ",
			"passwordConfirmation": "w24V,KY$f2YSIPQ"
		}`).
		AssertStatus(http.StatusCreated)

	database.AssertHas("users", map[string]any{
		"name":     "William",
		"email":    "doleya5976@agafx.com",
		"password": hash.Password("w24V,KY$f2YSIPQ"),
	})

	client.
		Post("/api/users/login", `{
			"email": "doley5976@agafx.com",
			"password": "w24V,KY$f2YSIPQ"
		}`).
		AssertStatus(http.StatusUnauthorized).
		AssertHasValidationErrors([]string{"email"})

	database.AssertCount("users", 1)
}

func TestLoginIncorrectPassword(t *testing.T) {
	client, database := startup.Users(t)

	client.
		Post("/api/users", `{
			"name": "Antony",
			"email": "pleonius@sentimentdate.com",
			"password": "L00k@tmEImHer3",
			"passwordConfirmation": "L00k@tmEImHer3"
		}`).
		AssertStatus(http.StatusCreated)

	database.AssertHas("users", map[string]any{
		"name":     "Antony",
		"email":    "pleonius@sentimentdate.com",
		"password": hash.Password("L00k@tmEImHer3"),
	})

	client.
		Post("/api/users/login", `{
			"email": "pleonius@sentimentdate.com",
			"password": "Lo0k@tmEImHer3"
		}`).
		AssertStatus(http.StatusUnauthorized).
		AssertHasValidationErrors([]string{"email"})

	database.AssertCount("users", 1)
}

func TestLoginValidation(t *testing.T) {
	env.Load()

	tests.AssertValidationSuccess[loginInput](t, []tests.ValidationTest{
		{"Email 1", "email", "Twila_Braun-Bogisich@gmail.com"},
		{"Email 2", "email", "Noemie16@gmail.com"},
		{"Email 3", "email", "Catalina41@gmail.com"},
		{"Email 4", "email", "Effie43@gmail.com"},
		{"Email 5", "email", "Roxanne.Satterfield5@zoho.com"},
		{"Email 6", "email", "Julio_Hackett@icloud.com"},
		{"Password 1", "password", "65kYThD5"},
		{"Password 2", "password", "B^H}i,o5:iJvco"},
		{"Password 3", "password", "2J*LA5NxKnZ>1}g0Beu^:HR^Bn!6-H3izGF#o2!>"},
		{"Password 4", "password", "_d=)21YWPX@%HHbV2et:D_,MH+Y0tV,+@:^]5Ne)!vgHH%@1Ls)M.BYb7bs3t~Py^5"},
	})

	tests.AssertValidationFail[loginInput](t, []tests.ValidationTest{
		{"Empty email", "email", ""},
		{"Email has no mail", "email", "Enoch_Corwin@"},
		{"Email has no separator", "email", "Cleta_Schimmelicloud.com"},
		{"Email has no domain", "email", "Triston77@outlook."},
		{"Email has spaces", "email", "Guillermo_Haag-Goyette @aol.com"},
		{"Empty password", "password", ""},
		{"Too short password", "password", "4Ot69f,"},
		{"Too long password", "password", "xMAUg>~WAu^Ep],e5m8R,~j?Pn__Cb@)#j_F~z*806QvUDERKKi8)T0-cH.Yh!3q+6uwK10yR!+r=4+kMRX5F9BcuvfzxT6>sdLEa"},
		{"Password has only numbers", "password", "27021375891235"},
		{"Password has only lowercase", "password", "nrau}h9j1d1hux@h@_wd"},
		{"Password has only uppercase", "password", "F9.F9#)30XJXM+WHW*VYJ"},
		{"Password has spaces", "password", "n6K-N%acxa)om]oT= 8muHQ?Zs=s"},
	})
}
