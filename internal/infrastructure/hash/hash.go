package hash

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/EugeneNail/actum/internal/infrastructure/env"
)

func New(data string) string {
	bytes := []byte(env.Get("PASSWORD_SALT") + data)
	hash := sha256.New().Sum(bytes)

	return base64.StdEncoding.EncodeToString(hash)
}
