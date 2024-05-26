package main

import (
	"github.com/EugeneNail/actum/internal/controller/activity"
	collectionController "github.com/EugeneNail/actum/internal/controller/collections"
	userController "github.com/EugeneNail/actum/internal/controller/users"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/database/resource/collections"
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

	collectionRepository := collections.NewRepository(db)
	collectionController := collectionController.New(db, collectionRepository)
	routing.Post("/api/collections", collectionController.Store)
	routing.Put("/api/collections/:id", collectionController.Update)
	routing.Delete("/api/collections/:id", collectionController.Destroy)
	routing.Get("/api/collections/:id", collectionController.Show)
	routing.Get("/api/collections", collectionController.Index)

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
