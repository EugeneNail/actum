package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/EugeneNail/actum/internal/database/resource/users"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/response"
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
		{"POST", "/api/users/refresh-token"},
	}
}

func Authenticate(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		response := response.NewSender(writer)

		for _, route := range unprotectedRoutes {
			if route.method == request.Method && route.path == request.RequestURI {
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

		user, err := getUser(db, payload.Id)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		if user.Id == 0 {
			writer.WriteHeader(http.StatusUnauthorized)
			message := []byte(fmt.Sprintf(`"Пользователь %d не найден"`, user.Id))

			if _, err := writer.Write(message); err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		ctx := context.WithValue(request.Context(), jwt.CtxKey("user"), user)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func getUser(db *sql.DB, userId int) (users.User, error) {
	var user users.User

	err := db.QueryRow(`SELECT * FROM users WHERE id = ? LIMIT 1`, userId).
		Scan(&user.Id, &user.Name, &user.Email, &user.Password)

	if err != nil && err != sql.ErrNoRows {
		return user, fmt.Errorf("authenticate.getUser(): %w", err)
	}

	return user, nil
}
