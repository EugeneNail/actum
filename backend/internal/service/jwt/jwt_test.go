package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/EugeneNail/actum/internal/model/users"
	"github.com/EugeneNail/actum/internal/service/env"
	"strings"
	"testing"
	"time"
)

func TestBuildHeader(t *testing.T) {
	header, err := buildHeader()
	check(err)
	encoded := base64.URLEncoding.EncodeToString([]byte(`{"alg":"SH256","typ":"JWT"}`))
	if header != encoded {
		t.Errorf("Header is invalid")
	}
}

func TestBuildPayload(t *testing.T) {
	user := users.New("John", "blank@mail.com", "")
	payload, err := buildPayload(user)
	check(err)
	expires := time.Now().Add(time.Hour).Unix()
	json := fmt.Sprintf(`{"id":0,"name":"John","exp":%d}`, expires)

	encoded := base64.URLEncoding.EncodeToString([]byte(json))
	if payload != encoded {
		t.Errorf("Payload is invalid")
	}
}

func TestBuildSignature(t *testing.T) {
	user := users.New("John", "blank@mail.com", "")
	base64Header := base64.URLEncoding.EncodeToString([]byte(`{"alg":"SH256","typ":"JWT"}`))
	expires := time.Now().Add(time.Hour).Unix()
	jsonPayload := fmt.Sprintf(`{"id":0,"name":"John","exp":%d}`, expires)
	base64Payload := base64.URLEncoding.EncodeToString([]byte(jsonPayload))
	signature, err := buildSignature(user)
	check(err)

	signatureBytes := hmac.
		New(sha256.New, []byte(env.Get("JWT_SALT"))).
		Sum([]byte(base64Header + "." + base64Payload))
	recreatedSignature := base64.URLEncoding.EncodeToString(signatureBytes)

	if signature != recreatedSignature {
		t.Errorf("Signature expected to be %s, got %s", recreatedSignature, signature)
	}
}

func TestMake(t *testing.T) {
	header := base64.URLEncoding.EncodeToString([]byte(`{"alg":"SH256","typ":"JWT"}`))
	expires := time.Now().Add(time.Hour).Unix()
	jsonPayload := fmt.Sprintf(`{"id":0,"name":"John","exp":%d}`, expires)
	payload := base64.URLEncoding.EncodeToString([]byte(jsonPayload))

	signatureBytes := hmac.
		New(sha256.New, []byte(env.Get("JWT_SALT"))).
		Sum([]byte(header + "." + payload))
	signature := base64.URLEncoding.EncodeToString(signatureBytes)

	user := users.New("John", "blank@mail.com", "")

	token, err := Make(user)
	check(err)
	parts := strings.Split(token, ".")

	if len(parts) != 3 {
		t.Errorf("Expected token to have 3 parts, got %d", len(parts))
	}

	if parts[0] != header {
		t.Errorf("Expected header to be %s, got %s", header, parts[0])
	}

	if parts[1] != payload {
		t.Errorf("Expected payload to be %s, got %s", payload, parts[1])
	}

	if parts[2] != signature {
		t.Errorf("Expected signature to be %s, got %s", signature, parts[2])
	}
}

func TestIsValid(t *testing.T) {
	user := users.New("John", "blank@mail.com", "")

	token := "header.payload,signature"
	if IsValid(token) {
		t.Error("Token must have 3 parts")
	}

	token, err := Make(user)
	check(err)
	if !IsValid(token) {
		t.Errorf("A valid token %s is considered invalid", token)
	}

	token += "noise"
	if IsValid(token) {
		t.Errorf("An invalid token %s is considered valid", token)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
