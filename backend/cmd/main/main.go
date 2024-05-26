package main

import (
	"github.com/EugeneNail/actum/internal/controller/activity"
	"github.com/EugeneNail/actum/internal/controller/collection"
	userController "github.com/EugeneNail/actum/internal/controller/users"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/database/resource/users"
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

	db, err := mysql.Connect()
	if err != nil {
		panic(err)
	}

	userRepository := users.NewRepository(db)
	userController := userController.New(userRepository)

	routing.Post("/api/users", userController.Store)
	routing.Post("/api/users/login", userController.Login)
	routing.Post("/api/collections", collection.Store)
	routing.Put("/api/collections/:id", collection.Update)
	routing.Get("/api/collections", collection.Index)
	routing.Delete("/api/collections/:id", collection.Destroy)
	routing.Get("/api/collections/:id", collection.Show)
	routing.Post("/api/activities", activity.Store)
	routing.Get("/api/activities/:id", activity.Show)
	routing.Put("/api/activities/:id", activity.Update)
	routing.Delete("/api/activities/:id", activity.Destroy)

	handler := middleware.BuildPipeline([]middleware.Middleware{
		middleware.SetHeaders,
		middleware.Authenticate,
		routing.Middleware,
	})

	err = http.ListenAndServe(":"+os.Getenv("APP_PORT"), handler)
	if err != nil {
		panic(err)
	}
}
