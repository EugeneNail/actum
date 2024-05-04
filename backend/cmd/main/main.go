package main

import (
	"github.com/EugeneNail/actum/internal/controller/user"
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/middleware"
	"github.com/EugeneNail/actum/internal/service/routing"
	"net/http"
	"os"
)

func main() {
	env.Load()
	go log.RotateFiles()

	routing.Post("/api/users", user.Store)
	routing.Post("/api/users/login", user.Login)

	handler := middleware.BuildPipeline(routing.Serve(), []middleware.Middleware{
		middleware.SetHeaders,
		middleware.Authenticate,
	})

	err := http.ListenAndServe(":"+os.Getenv("APP_PORT"), handler)

	if err != nil {
		panic(err)
	}
}
