package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func BuildPipeline(middlewares []Middleware) http.Handler {
	emptyHandler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
	handler := middlewares[0](emptyHandler)

	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}
