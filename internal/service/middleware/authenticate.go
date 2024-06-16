package middleware

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/response"
	"net/http"
	"regexp"
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
		{"POST", "/api/users/refresh-token"},
		{"GET", "/api/photos/[0-9a-zA-Z_.-]+"},
	}
}

func Authenticate(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		response := response.NewSender(writer)

		for _, route := range unprotectedRoutes {
			pathMatched, err := regexp.MatchString("^"+route.path+"$", request.RequestURI)
			if err != nil {
				response.Send(fmt.Errorf("AuthenticateMiddleware: %w", err), http.StatusInternalServerError)
				return
			}

			if route.method == request.Method && pathMatched {
				next.ServeHTTP(writer, request)
				return
			}
		}

		parts := strings.Split(request.Header.Get("Authorization"), " ")
		if len(parts) < 2 {
			response.Send("Токен отсутствует", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		if !jwt.IsValid(token) {
			response.Send("Неправильный токен", http.StatusUnauthorized)
			return
		}

		payload, err := jwt.ExtractPayload(token)
		if err != nil {
			response.Send(err, http.StatusBadRequest)
			return
		}

		if payload.Exp < time.Now().Unix() {
			response.Send("Срок действия токена истек", http.StatusUnauthorized)
			return
		}

		userExists, err := userExists(db, payload.Id)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		if !userExists {
			writer.WriteHeader(http.StatusUnauthorized)
			message := []byte(fmt.Sprintf(`"Пользователь %d не найден"`, payload.Id))

			if _, err := writer.Write(message); err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		ctx := context.WithValue(request.Context(), jwt.CtxKey("userExists"), payload.Id)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func userExists(db *sql.DB, userId int) (bool, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE id = ? LIMIT 1`, userId).
		Scan(&count)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, fmt.Errorf("authenticate.userExists(): %w", err)
	}

	return count > 0, nil
}
