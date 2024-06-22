package tests

import (
	"encoding/json"
	"fmt"
	"github.com/EugeneNail/actum/internal/infrastructure/env"
	"math/rand/v2"
	"net/http"
	"strings"
	"testing"
)

type Tokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type Client struct {
	t     *testing.T
	token string
}

func NewClient(t *testing.T) (client Client) {
	client.t = t

	client.
		Post("/api/auth/register", `{
			"name": "John",
			"email": "blank@gmail.com",
			"password": "Strong123",
			"passwordConfirmation": "Strong123"
		}`).
		AssertStatus(http.StatusCreated)

	response := client.
		Post("/api/auth/login", `{
			"email": "blank@gmail.com",
			"password": "Strong123"
		}`).
		AssertStatus(http.StatusOK)

	client.setToken(response)

	return
}

func (client *Client) setToken(response *Response) {
	decoder := json.NewDecoder(response.Body)
	var tokens Tokens
	err := decoder.Decode(&tokens)
	Check(err)
	client.token = "Bearer " + tokens.Access
}

func NewClientWithoutAuth(t *testing.T) (client Client) {
	client.t = t
	return
}

func (client *Client) Get(path string) *Response {
	url := "http://127.0.0.1:" + env.Get("APP_PORT") + path
	request, err := http.NewRequest("GET", url, nil)
	Check(err)

	if len(client.token) > 0 {
		request.Header.Set("Authorization", client.token)
	}
	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	Check(err)

	return &Response{response, client.t}
}

func (client *Client) Post(path string, json string) *Response {
	url := "http://127.0.0.1:" + env.Get("APP_PORT") + path
	body := strings.NewReader(json)
	request, err := http.NewRequest("POST", url, body)
	Check(err)

	if len(client.token) > 0 {
		request.Header.Set("Authorization", client.token)
	}
	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	Check(err)

	return &Response{response, client.t}
}

func (client *Client) Put(path string, json string) *Response {
	url := "http://127.0.0.1:" + env.Get("APP_PORT") + path
	body := strings.NewReader(json)
	request, err := http.NewRequest("PUT", url, body)
	Check(err)
	if len(client.token) > 0 {
		request.Header.Set("Authorization", client.token)
	}
	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	Check(err)

	return &Response{response, client.t}
}

func (client *Client) Delete(path string) *Response {
	url := "http://127.0.0.1:" + env.Get("APP_PORT") + path
	body := strings.NewReader("")
	request, err := http.NewRequest("DELETE", url, body)
	Check(err)
	if len(client.token) > 0 {
		request.Header.Set("Authorization", client.token)
	}
	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	Check(err)

	return &Response{response, client.t}
}

func (client *Client) ChangeUser() {
	names := []string{"Sara", "John", "Sam", "Donald", "William"}
	emails := []string{"blank@mail.com", "hackerman106@gmail.com", "abcdefg@108list.org", "name.surname@bing.xorg", "jaja@yahoo.cc"}

	input := fmt.Sprintf(`{
		"name": "%s",
		"email": "%s",
		"password": "Strong123",
		"passwordConfirmation": "Strong123"
	}`, names[rand.IntN(4)], emails[rand.IntN(4)])

	response := client.
		Post("/api/users", input).
		AssertStatus(http.StatusCreated)

	client.setToken(response)

}

func (client *Client) UnsetToken() {
	client.token = ""
}
