package routing

import (
	"context"
	"database/sql"
	"net/http"
	"regexp"
	"strings"
)

var routes []route

type route struct {
	method string
	regex  *regexp.Regexp
	handle func(http.ResponseWriter, *http.Request)
}

type CtxKey string

func Middleware(_ *sql.DB, __ http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var allowedMethods []string
		for _, route := range routes {
			matches := route.regex.FindStringSubmatch(request.URL.Path)

			if len(matches) > 0 {
				if route.method != request.Method {
					allowedMethods = append(allowedMethods, route.method)
					continue
				}

				ctx := context.WithValue(request.Context(), CtxKey("variables"), matches[1:])
				route.handle(writer, request.WithContext(ctx))
				return
			}
		}

		if len(allowedMethods) > 0 {
			message := "Method not allowed. Allowed methods: " + strings.Join(allowedMethods, ", ")
			writer.Header().Set("Allow", strings.Join(allowedMethods, ", "))
			http.Error(writer, message, http.StatusMethodNotAllowed)
			return
		}
		http.NotFound(writer, request)
	})
}

func Get(pattern string, handlerFunction func(http.ResponseWriter, *http.Request)) {
	RegisterRoute("GET", pattern, handlerFunction)
}

func Post(pattern string, handlerFunction func(http.ResponseWriter, *http.Request)) {
	RegisterRoute("POST", pattern, handlerFunction)
}

func Put(pattern string, handlerFunction func(http.ResponseWriter, *http.Request)) {
	RegisterRoute("PUT", pattern, handlerFunction)
}

func Delete(pattern string, handlerFunction func(http.ResponseWriter, *http.Request)) {
	RegisterRoute("DELETE", pattern, handlerFunction)
}

func RegisterRoute(method string, pattern string, handlerFunction func(http.ResponseWriter, *http.Request)) {
	validateMethod(method)
	pattern = regexp.MustCompile(":[a-zA-Z]*").ReplaceAllString(pattern, "([0-9a-zA-Z_-]+)")
	routes = append(routes, route{
		method,
		regexp.MustCompile("^" + pattern + "$"),
		handlerFunction,
	})
}

func validateMethod(method string) {
	validMethods := []string{"POST", "GET", "HEAD", "OPTIONS", "PUT", "PATCH", "DELETE"}

	for _, validMethod := range validMethods {
		if method == validMethod {
			return
		}
	}
	panic("invalid route method: " + method)
}

func GetVariable(request *http.Request, index int) string {
	variables := request.Context().Value(CtxKey("variables")).([]string)
	return variables[index]
}
