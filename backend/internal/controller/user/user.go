package user

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/tests"
)

func hashPassword(password string) string {
	bytes := []byte(env.Get("PASSWORD_SALT") + password)
	hash := sha256.New().Sum(bytes)

	return base64.StdEncoding.EncodeToString(hash)
}

func cleanup() {
	err := mysql.Truncate("users")
	tests.Check(err)
}

func getUrl() string {
	return "http://127.0.0.1:" + env.Get("APP_PORT") + "/api/users"
}
