package hash

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/EugeneNail/actum/internal/service/env"
)

func Password(password string) string {
	bytes := []byte(env.Get("PASSWORD_SALT") + password)
	hash := sha256.New().Sum(bytes)

	return base64.StdEncoding.EncodeToString(hash)
}
