package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/EugeneNail/actum/internal/model/users"
	"github.com/EugeneNail/actum/internal/service/env"
	"strings"
	"time"
)

type header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type payload struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Exp  int64  `json:"exp"`
}

func Make(user users.User) (string, error) {
	header, err := buildHeader()
	if err != nil {
		return "", fmt.Errorf("jwt.Make(): %w", err)
	}

	payload, err := buildPayload(user)
	if err != nil {
		return "", fmt.Errorf("jwt.Make(): %w", err)
	}

	signature, err := buildSignature(user)
	if err != nil {
		return "", fmt.Errorf("jwt.Make(): %w", err)
	}

	return header + "." + payload + "." + signature, nil
}

func buildHeader() (string, error) {
	header := header{"SH256", "JWT"}
	jsonHeader, err := json.Marshal(header)

	if err != nil {
		return "", fmt.Errorf("buildHeader(): %w", err)
	}

	return base64.URLEncoding.EncodeToString(jsonHeader), nil
}

func buildPayload(user users.User) (string, error) {
	expires := time.Now().Add(time.Hour).Unix()
	jsonPayload, err := json.Marshal(payload{user.Id, user.Name, expires})

	if err != nil {
		return "", fmt.Errorf("buildPayload(): %w", err)
	}

	return base64.URLEncoding.EncodeToString(jsonPayload), nil
}

func buildSignature(user users.User) (string, error) {
	header, err := buildHeader()
	if err != nil {
		return "", fmt.Errorf("buildSignature(): %w", err)
	}

	payload, err := buildPayload(user)
	if err != nil {
		return "", fmt.Errorf("buildSignature(): %w", err)
	}

	signature := hmac.
		New(sha256.New, []byte(env.Get("JWT_SALT"))).
		Sum([]byte(header + "." + payload))

	return base64.URLEncoding.EncodeToString(signature), nil
}

func IsValid(token string) bool {
	parts := strings.Split(token, ".")

	if len(parts) != 3 {
		return false
	}

	signatureBytes := hmac.
		New(sha256.New, []byte(env.Get("JWT_SALT"))).
		Sum([]byte(parts[0] + "." + parts[1]))
	recreatedSignature := base64.URLEncoding.EncodeToString(signatureBytes)

	if recreatedSignature != parts[2] {
		return false
	}

	return true
}
