package middleware

import (
	"github.com/EugeneNail/actum/internal/service/jwt"
	"net/http"
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

		token, err := request.Cookie("Access-Token")
		if err != nil {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !jwt.IsValid(token.Value) {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		payload, err := jwt.ExtractPayload(token.Value)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		if payload.Exp < time.Now().Unix() {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
