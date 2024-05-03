package user

import (
	"encoding/json"
	"github.com/EugeneNail/actum/internal/model/users"
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/test"
	"io"
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
	check(err)

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", response.StatusCode)
	}
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
	check(err)

	if response.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", response.StatusCode)
	}

	var validationErrors map[string]string
	data, err := io.ReadAll(response.Body)
	check(err)
	err = json.Unmarshal(data, &validationErrors)
	check(err)
	for _, field := range []string{"email", "password"} {
		if _, exists := validationErrors[field]; !exists {
			t.Errorf(`expected validation error for field "%s" to be present`, field)
		}
	}

	assertUserIsUntouched(user, t)
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
	check(err)

	if response.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", response.StatusCode)
	}

	var validationErrors map[string]string
	data, err := io.ReadAll(response.Body)
	check(err)
	err = json.Unmarshal(data, &validationErrors)
	check(err)
	if _, exists := validationErrors["email"]; !exists {
		t.Errorf(`expected validation error for field "email" to be present`)
	}

	assertUserIsUntouched(user, t)
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
	check(err)

	if response.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", response.StatusCode)
	}

	var validationErrors map[string]string
	data, err := io.ReadAll(response.Body)
	check(err)
	err = json.Unmarshal(data, &validationErrors)
	check(err)
	if _, exists := validationErrors["email"]; !exists {
		t.Errorf(`expected validation error for field "email" to be present`)
	}

	assertUserIsUntouched(user, t)
}

func createUser(email string, password string) users.User {
	user := users.User{0, "John", email, hashPassword(password)}
	err := user.Save()
	check(err)

	return user
}

func assertUserIsUntouched(user users.User, t *testing.T) {
	dbUser, err := users.Find(1)
	check(err)
	if dbUser.Name != user.Name {
		t.Errorf(`field "name" has been corrupted`)
	}

	if dbUser.Email != user.Email {
		t.Errorf(`field "email" has been corrupted`)
	}

	if dbUser.Password != user.Password {
		t.Errorf(`field "password" has been corrupted`)
	}
}

func TestLoginValidation(t *testing.T) {
	env.Load()

	successes := []test.Field{
		{"email", "Twila_Braun-Bogisich@gmail.com"},
		{"email", "Noemie16@gmail.com"},
		{"email", "Catalina41@gmail.com"},
		{"email", "Effie43@gmail.com"},
		{"email", "Roxanne.Satterfield5@zoho.com"},
		{"email", "Julio_Hackett@icloud.com"},
		{"password", "65kYThD5"},
		{"password", "B^H}i,o5:iJvco"},
		{"password", "2J*LA5NxKnZ>1}g0Beu^:HR^Bn!6-H3izGF#o2!>"},
		{"password", "_d=)21YWPX@%HHbV2et:D_,MH+Y0tV,+@:^]5Ne)!vgHH%@1Ls)M.BYb7bs3t~Py^5"},
	}
	for _, field := range successes {
		test.AssertValidationSuccess[storeInput](field, t)
	}

	fails := []test.Field{
		{"email", ""},
		{"email", "Enoch_Corwin@"},
		{"email", "Guillermo_Haag-Goyette @aol.com"},
		{"email", "Cleta_Schimmelicloud.com"},
		{"email", "Triston77@outlook."},
		{"password", ""},
		{"password", "4Ot69f,"},
		{"password", "27021375891235"},
		{"password", "nrau}h9j1d1hux@h@_wd"},
		{"password", "F9.F9#)30XJXM+WHW*VYJ"},
		{"password", "n6K-N%acxa)om]oT= 8muHQ?Zs=s"},
		{"password", "xMAUg>~WAu^Ep],e5m8R,~j?Pn__Cb@)#j_F~z*806QvUDERKKi8)T0-cH.Yh!3q+6uwK10yR!+r=4+kMRX5F9BcuvfzxT6>sdLEa"},
	}
	for _, field := range fails {
		test.AssertValidationFail[storeInput](field, t)
	}
}