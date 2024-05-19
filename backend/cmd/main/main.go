package main

import (
	"github.com/EugeneNail/actum/internal/controller/activity"
	"github.com/EugeneNail/actum/internal/controller/collection"
	"github.com/EugeneNail/actum/internal/controller/user"
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/middleware"
	"github.com/EugeneNail/actum/internal/service/middleware/routing"
	"net/http"
	"os"
)

func main() {
	env.Load()
	go log.RotateFiles()

	routing.Post("/api/users", user.Store)
	routing.Post("/api/users/login", user.Login)
	routing.Post("/api/collections", collection.Store)
	routing.Put("/api/collections/:id", collection.Update)
	routing.Get("/api/collections", collection.Index)
	routing.Delete("/api/collections/:id", collection.Destroy)
	routing.Get("/api/collections/:id", collection.Show)
	routing.Post("/api/activities", activity.Store)

	handler := middleware.BuildPipeline([]middleware.Middleware{
		middleware.SetHeaders,
		middleware.Authenticate,
		routing.Middleware,
	})

	err := http.ListenAndServe(":"+os.Getenv("APP_PORT"), handler)

	if err != nil {
		panic(err)
	}
}
