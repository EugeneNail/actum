package user

import (
	"github.com/EugeneNail/actum/internal/model/users"
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/tests"
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"testing"
)

func TestLoginValidData(t *testing.T) {
	client := startup.UsersLogin(t)

	response := client.Post("/api/users", `{
		"name": "John",
		"email": "jodame3394@agafx.com",
		"password": "Strong123",
		"passwordConfirmation": "Strong123"
	}`)
	response.AssertStatus(http.StatusCreated)

	response = client.Post("/api/users/login", `{
		"email": "jodame3394@agafx.com",
		"password": "Strong123"
	}`)

	response.AssertStatus(http.StatusOK)
	response.AssertHasToken()
}

func TestLoginInvalidData(t *testing.T) {
	client := startup.UsersLogin(t)

	response := client.Post("/api/users/login", `{
		"email": "yibewek618goulink.com",
		"password": "v9"
	}`)

	response.AssertStatus(http.StatusUnprocessableEntity)
	response.AssertHasValidationErrors([]string{"email", "password"})
	response.AssertHasNoToken()
}

func TestLoginIncorrectEmail(t *testing.T) {
	client := startup.UsersLogin(t)

	client.Post("/api/users", `{
		"email": "doleya5976@agafx.com",
		"password": "w24V,KY$f2YSIPQ"
	}`)

	user, err := users.Find(1)
	tests.Check(err)

	response := client.Post("/api/users/login", `{
		"email": "doley5976@agafx.com",
		"password": "w24V,KY$f2YSIPQ"
	}`)

	response.AssertStatus(http.StatusUnauthorized)
	response.AssertHasValidationErrors([]string{"email"})
	response.AssertHasNoToken()
	tests.AssertUserIsUntouched(user, t)
}

func TestLoginIncorrectPassword(t *testing.T) {
	client := startup.UsersLogin(t)

	client.Post("/api/users", `{
		"email": "pleonius@sentimentdate.com",
		"password": "L00k@tmEImHer3"
	}`)

	user, err := users.Find(1)
	tests.Check(err)

	response := client.Post("/api/users/login", `{
		"email": "pleonius@sentimentdate.com",
		"password": "Lo0k@tmEImHer3"
	}`)

	response.AssertStatus(http.StatusUnauthorized)
	response.AssertHasValidationErrors([]string{"email"})
	response.AssertHasNoToken()
	tests.AssertUserIsUntouched(user, t)
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
