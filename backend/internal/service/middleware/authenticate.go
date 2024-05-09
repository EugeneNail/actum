package middleware

import (
	"context"
	"fmt"
	"github.com/EugeneNail/actum/internal/model/users"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"net/http"
	"strings"
	"time"
)

type route struct {
	method string
	path   string
}

var unprotectedRoutes []route

func init() {
	unprotectedRoutes = []route{
		{"POST", "/api/users/login"},
		{"POST", "/api/users"},
	}
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		for _, route := range unprotectedRoutes {
			if route.method == request.Method && route.path == request.RequestURI {
				next.ServeHTTP(writer, request)
				return
			}
		}

		token := strings.Split(request.Header.Get("Authorization"), " ")[1]
		if len(token) == 0 {
			writer.WriteHeader(http.StatusUnauthorized)
			if _, err := writer.Write([]byte(`"Token not present"`)); err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if !jwt.IsValid(token) {
			writer.WriteHeader(http.StatusUnauthorized)
			if _, err := writer.Write([]byte(`"Token is invalid"`)); err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		payload, err := jwt.ExtractPayload(token)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		if payload.Exp < time.Now().Unix() {
			writer.WriteHeader(http.StatusUnauthorized)
			if _, err := writer.Write([]byte(`"Token has expired"`)); err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		user, err := users.Find(payload.Id)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		if user.Id == 0 {
			writer.WriteHeader(http.StatusUnauthorized)
			message := []byte(fmt.Sprintf(`"User %d not found"`, user.Id))

			if _, err := writer.Write(message); err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		ctx := context.WithValue(request.Context(), jwt.CtxKey("user"), user)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
