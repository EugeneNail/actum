package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/EugeneNail/actum/internal/database/resource/users"
	"github.com/EugeneNail/actum/internal/service/env"
	"net/http"
	"strings"
	"time"
)

type header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type Payload struct {
	Id  int   `json:"id"`
	Exp int64 `json:"exp"`
}

type CtxKey string

func GetUser(request *http.Request) users.User {
	return request.Context().Value(CtxKey("user")).(users.User)
}

func Make(userId int) (string, error) {
	header, err := buildHeader()
	if err != nil {
		return "", fmt.Errorf("jwt.Make(): %w", err)
	}

	payload, err := buildPayload(userId)
	if err != nil {
		return "", fmt.Errorf("jwt.Make(): %w", err)
	}

	signature, err := buildSignature(userId)
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

func buildPayload(userId int) (string, error) {
	expires := time.Now().Add(time.Hour).Unix()
	jsonPayload, err := json.Marshal(Payload{userId, expires})

	if err != nil {
		return "", fmt.Errorf("buildPayload(): %w", err)
	}

	return base64.URLEncoding.EncodeToString(jsonPayload), nil
}

func buildSignature(userId int) (string, error) {
	header, err := buildHeader()
	if err != nil {
		return "", fmt.Errorf("buildSignature(): %w", err)
	}

	payload, err := buildPayload(userId)
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

func ExtractPayload(token string) (Payload, error) {
	var payload Payload
	base64Payload := strings.Split(token, ".")[1]

	decoded, err := base64.URLEncoding.DecodeString(base64Payload)
	if err != nil {
		return payload, fmt.Errorf("jwt.ExtractPayload(): %w", err)
	}

	if err = json.Unmarshal(decoded, &payload); err != nil {
		return payload, fmt.Errorf("jwt.ExtractPayload(): %w", err)
	}

	return payload, nil
}
