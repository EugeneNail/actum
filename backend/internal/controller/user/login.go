package user

import (
	"github.com/EugeneNail/actum/internal/controller"
	"github.com/EugeneNail/actum/internal/model/users"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"net/http"
	"strings"
)

type loginInput struct {
	Email    string `json:"email" rules:"required|email|max:100"`
	Password string `json:"password" rules:"required|min:8|max:100|mixedCase|regex:^\\S+$"`
}

func Login(writer http.ResponseWriter, request *http.Request) {
	controller := controller.New[loginInput](writer)

	isValid := controller.Validate(request)
	if !isValid {
		return
	}

	user, err := users.FindBy("email", strings.ToLower(controller.Input.Email))
	if err != nil {
		controller.Response(err, http.StatusInternalServerError)
		return
	}

	if user.Id == 0 || user.Password != hashPassword(controller.Input.Password) {
		controller.Response(map[string]string{"email": "Incorrect email address or password"}, http.StatusUnauthorized)
		return
	}

	token, err := jwt.Make(user)
	if err != nil {
		controller.Response(err, http.StatusInternalServerError)
		return
	}

	controller.Response(token, http.StatusOK)
	log.Info("User", user.Id, "logged in")
}
