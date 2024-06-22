package middleware

import (
	"github.com/EugeneNail/actum/internal/service/env"
	"net/http"
	"path/filepath"
	"strings"
)

func RedirectToFrontend(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		if !strings.HasPrefix(request.URL.Path, "/api") {
			frontendPath := filepath.Join(env.Get("APP_PATH"), "frontend", "build")

			if _, err := http.Dir(frontendPath).Open(request.URL.Path); err != nil {
				http.ServeFile(writer, request, filepath.Join(frontendPath, "index.html"))
				return
			}

			http.ServeFile(writer, request, filepath.Join(frontendPath, request.URL.Path))
			return
		}

		handler.ServeHTTP(writer, request)
	})
}
