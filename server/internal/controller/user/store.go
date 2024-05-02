package user

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/EugeneNail/actum/internal/model/users"
	"github.com/EugeneNail/actum/internal/service/controller"
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/validation"
	"net/http"
)

type storeInput struct {
	Id                   int    `json:"id"`
	Name                 string `json:"name" rules:"required|word|min:3|max:20"`
	Email                string `json:"email" rules:"required|email|max:100|unique:users,email"`
	Password             string `json:"password" rules:"required|min:8|max:100|mixedCase|regex:^\\S+$"`
	PasswordConfirmation string `json:"passwordConfirmation" rules:"required|min:8|max:100|mixedCase|regex:^\\S+$"`
}

func Store(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	input, err := controller.Parse[storeInput](writer, request)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	validationErrors, err := validation.Perform(input)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}

	if len(validationErrors) > 0 {
		writer.WriteHeader(http.StatusUnprocessableEntity)

		if err := encoder.Encode(validationErrors); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			log.Error(err)
		}
		return
	}

	if input.Password != input.PasswordConfirmation {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		err := encoder.Encode(map[string]string{"passwordConfirmation": "Passwords do not match"})

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			log.Error(err)
		}
		return
	}

	user := users.User{input.Id, input.Name, input.Email, hashPassword(input.Password)}

	if err = user.Save(); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	writer.WriteHeader(http.StatusCreated)

	if err = encoder.Encode(user.Id); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	log.Info("Created user", user.Id)
}

func hashPassword(password string) string {
	bytes := []byte(env.Get("PASSWORD_SALT") + password)
	hash := sha256.New().Sum(bytes)

	return base64.StdEncoding.EncodeToString(hash)
}
