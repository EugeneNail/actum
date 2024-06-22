package middleware

import (
	"github.com/EugeneNail/actum/internal/service/env"
	"net/http"
)

func SetHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		writer.Header().Set("Access-Control-Allow-Origin", env.Get("APP_ORIGIN"))

		if request.Method == "OPTIONS" {
			writer.WriteHeader(200)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
