package middleware

import (
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

func BuildPipeline(handler http.Handler, middlewares []Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}

func SetResponseHeaders(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		writer.Header().Set("Access-Control-Allow-Methods", "*")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if request.Method == "OPTIONS" {
			writer.WriteHeader(200)
			return
		}

		handler.ServeHTTP(writer, request)
	})
}

type route struct {
	method string
	path   string
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		excludedRoutes := []route{
			{"POST", "/api/users/login"},
			{"POST", "/api/users"},
		}

		for _, route := range excludedRoutes {
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
			log.Error(err)
			return
		}

		if payload.Exp < time.Now().Unix() {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
