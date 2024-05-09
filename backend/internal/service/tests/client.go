package tests

import (
	"github.com/EugeneNail/actum/internal/service/env"
	"net/http"
	"strings"
	"testing"
)

type Client struct {
	t     *testing.T
	token string
}

func NewClient(t *testing.T) (client Client) {
	client.t = t

	response := client.Post("/api/users", `{
		"name": "John",
		"email": "blank@gmail.com",
		"password": "Strong123",
		"passwordConfirmation": "Strong123"
	}`)

	response.AssertStatus(http.StatusCreated)

	return
}

func NewClientWithoutAuth(t *testing.T) (client Client) {
	client.t = t
	return
}

func (client *Client) Post(path string, json string) Response {
	url := "http://127.0.0.1:" + env.Get("APP_PORT") + path
	body := strings.NewReader(json)
	request, err := http.NewRequest("POST", url, body)
	Check(err)
	request.Header.Set("Cookie", "Access-Token="+client.token)
	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	Check(err)

	for _, cookie := range response.Cookies() {
		if cookie.Name == "Access-Token" {
			client.token = cookie.Value
			break
		}
	}

	return Response{response, client.t}
}
