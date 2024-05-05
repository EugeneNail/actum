package user

import (
	"github.com/EugeneNail/actum/internal/model/users"
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/tests"
	"net/http"
	"strings"
	"testing"
)

func TestLoginValidData(t *testing.T) {
	env.Load()
	t.Cleanup(cleanup)
	url := getUrl() + "/login"
	createUser("jodame3394@agafx.com", "Strong123")

	response, err := http.Post(url, "application/json", strings.NewReader(`{
		"email": "jodame3394@agafx.com",
		"password": "Strong123"
	}`))
	tests.Check(err)

	tests.AssertStatus(response, http.StatusOK, t)
	tests.AssertHasToken(response, t)
}

func TestLoginInvalidData(t *testing.T) {
	env.Load()
	t.Cleanup(cleanup)
	url := getUrl() + "/login"
	user := createUser("yibewek618@goulink.com", "v9&;43mV,>2BE^t")

	response, err := http.Post(url, "application/json", strings.NewReader(`{
		"email": "yibewek618goulink.com",
		"password": "v9"
	}`))
	tests.Check(err)

	tests.AssertStatus(response, http.StatusUnprocessableEntity, t)
	tests.AssertHasValidationErrors(response, []string{"email", "password"}, t)
	tests.AssertHasNoToken(response, t)
	tests.AssertUserIsUntouched(user, t)
}

func TestLoginIncorrectEmail(t *testing.T) {
	env.Load()
	t.Cleanup(cleanup)
	url := getUrl() + "/login"
	user := createUser("doleya5976@agafx.com", "w24V,KY$f2YSIPQ")

	response, err := http.Post(url, "application/json", strings.NewReader(`{
		"email": "doley5976@agafx.com",
		"password": "w24V,KY$f2YSIPQ"
	}`))
	tests.Check(err)

	tests.AssertStatus(response, http.StatusUnauthorized, t)
	tests.AssertHasValidationErrors(response, []string{"email"}, t)
	tests.AssertHasNoToken(response, t)
	tests.AssertUserIsUntouched(user, t)
}

func TestLoginIncorrectPassword(t *testing.T) {
	env.Load()
	t.Cleanup(cleanup)
	url := getUrl() + "/login"
	user := createUser("pleonius@sentimentdate.com", "L00k@tmEImHer3")

	response, err := http.Post(url, "application/json", strings.NewReader(`{
		"email": "pleonius@sentimentdate.com",
		"password": "Lo0k@tmEImHer3"
	}`))
	tests.Check(err)

	tests.AssertStatus(response, http.StatusUnauthorized, t)
	tests.AssertHasValidationErrors(response, []string{"email"}, t)
	tests.AssertHasNoToken(response, t)
	tests.AssertUserIsUntouched(user, t)
}

func createUser(email string, password string) users.User {
	user := users.User{0, "John", email, hashPassword(password)}
	err := user.Save()
	tests.Check(err)

	return user
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
