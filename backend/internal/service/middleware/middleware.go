package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func BuildPipeline(handler http.Handler, middlewares []Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}

func SetContentType(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(writer, request)
	})
}
