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
			"name": "Антон",
			"email": "JODAME3394@agafx.com",
			"password": "Strong123",
			"passwordConfirmation": "Strong123"
		}`).
		AssertStatus(http.StatusCreated)

	database.AssertHas("users", map[string]any{
		"name":     "Антон",
		"email":    "jodame3394@agafx.com",
		"password": hash.New("Strong123"),
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
			"name":     "Антон",
			"email":    "jodame3394@agafx.com",
			"password": hash.New("Strong123"),
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
			"name": "Владислав",
			"email": "doleya5976@agafx.com",
			"password": "w24V,KY$f2YSIPQ",
			"passwordConfirmation": "w24V,KY$f2YSIPQ"
		}`).
		AssertStatus(http.StatusCreated)

	database.AssertHas("users", map[string]any{
		"name":     "Владислав",
		"email":    "doleya5976@agafx.com",
		"password": hash.New("w24V,KY$f2YSIPQ"),
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
			"name": "Петр",
			"email": "pleonius@sentimentdate.com",
			"password": "L00k@tmEImHer3",
			"passwordConfirmation": "L00k@tmEImHer3"
		}`).
		AssertStatus(http.StatusCreated)

	database.AssertHas("users", map[string]any{
		"name":     "Петр",
		"email":    "pleonius@sentimentdate.com",
		"password": hash.New("L00k@tmEImHer3"),
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
		{"email", "Email 1", "Twila_Braun-Bogisich@gmail.com"},
		{"email", "Email 2", "Noemie16@gmail.com"},
		{"email", "Email 3", "Catalina41@gmail.com"},
		{"email", "Email 4", "Effie43@gmail.com"},
		{"email", "Email 5", "Roxanne.Satterfield5@zoho.com"},
		{"email", "Email 6", "Julio_Hackett@icloud.com"},
		{"password", "Password 1", "65kYThD5"},
		{"password", "Password 2", "B^H}i,o5:iJvco"},
		{"password", "Password 3", "2J*LA5NxKnZ>1}g0Beu^:HR^Bn!6-H3izGF#o2!>"},
		{"password", "Password 4", "_d=)21YWPX@%HHbV2et:D_,MH+Y0tV,+@:^]5Ne)!vgHH%@1Ls)M.BYb7bs3t~Py^5"},
	})

	tests.AssertValidationFail[loginInput](t, []tests.ValidationTest{
		{"email", "Empty email", ""},
		{"email", "Email has no mail", "Enoch_Corwin@"},
		{"email", "Email has no separator", "Cleta_Schimmelicloud.com"},
		{"email", "Email has no domain", "Triston77@outlook."},
		{"email", "Email has spaces", "Guillermo_Haag-Goyette @aol.com"},
		{"password", "Empty password", ""},
		{"password", "Too short password", "4Ot69f,"},
		{"password", "Too long password", "xMAUg>~WAu^Ep],e5m8R,~j?Pn__Cb@)#j_F~z*806QvUDERKKi8)T0-cH.Yh!3q+6uwK10yR!+r=4+kMRX5F9BcuvfzxT6>sdLEa"},
		{"password", "Password has only numbers", "27021375891235"},
		{"password", "Password has only lowercase", "nrau}h9j1d1hux@h@_wd"},
		{"password", "Password has only uppercase", "F9.F9#)30XJXM+WHW*VYJ"},
		{"password", "Password has spaces", "n6K-N%acxa)om]oT= 8muHQ?Zs=s"},
	})
}
