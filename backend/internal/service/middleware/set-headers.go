package middleware

import (
	"database/sql"
	"net/http"
	"strings"
)

func SetHeaders(_ *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		setOrigin(writer, request)

		if request.Method == "OPTIONS" {
			writer.WriteHeader(200)
			return
		}

		next.ServeHTTP(writer, request)
	})
}

func setOrigin(writer http.ResponseWriter, request *http.Request) {
	origin := request.Header.Get("Origin")
	if strings.Contains(origin, "192.168.") || strings.Contains(origin, "localhost") {
		writer.Header().Set("Access-Control-Allow-Origin", origin)
	}
}
