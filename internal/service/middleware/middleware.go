package middleware

import (
	"database/sql"
	"net/http"
)

type Middleware func(*sql.DB, http.Handler) http.Handler

func BuildPipeline(db *sql.DB, middlewares []Middleware) http.Handler {
	emptyHandler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
	handler := middlewares[0](db, emptyHandler)

	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](db, handler)
	}

	return handler
}
