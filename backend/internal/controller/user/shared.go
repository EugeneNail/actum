package user

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/service/env"
)

func hashPassword(password string) string {
	bytes := []byte(env.Get("PASSWORD_SALT") + password)
	hash := sha256.New().Sum(bytes)

	return base64.StdEncoding.EncodeToString(hash)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func cleanup() {
	err := mysql.Truncate("users")
	check(err)
}

func getUrl() string {
	return "http://127.0.0.1:" + env.Get("APP_PORT") + "/api/users"
}
