package main

import (
	activityController "github.com/EugeneNail/actum/internal/controller/activities"
	collectionController "github.com/EugeneNail/actum/internal/controller/collections"
	recordController "github.com/EugeneNail/actum/internal/controller/records"
	userController "github.com/EugeneNail/actum/internal/controller/users"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/database/resource/activities"
	"github.com/EugeneNail/actum/internal/database/resource/collections"
	"github.com/EugeneNail/actum/internal/database/resource/records"
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

	userDAO := users.NewDAO(db)
	userController := userController.New(userDAO)
	routing.Post("/api/users", userController.Store)
	routing.Post("/api/users/login", userController.Login)

	collectionDAO := collections.NewDAO(db)
	collectionController := collectionController.New(db, collectionDAO)
	routing.Post("/api/collections", collectionController.Store)
	routing.Put("/api/collections/:id", collectionController.Update)
	routing.Delete("/api/collections/:id", collectionController.Destroy)
	routing.Get("/api/collections/:id", collectionController.Show)
	routing.Get("/api/collections", collectionController.Index)

	activityDAO := activities.NewDAO(db)
	activityController := activityController.New(db, activityDAO, collectionDAO)
	routing.Post("/api/activities", activityController.Store)
	routing.Put("/api/activities/:id", activityController.Update)
	routing.Delete("/api/activities/:id", activityController.Destroy)
	routing.Get("/api/activities/:id", activityController.Show)

	recordDAO := records.NewDAO(db)
	recordController := recordController.New(db, recordDAO, activityDAO)
	routing.Post("/api/records", recordController.Store)
	routing.Put("/api/records/:id", recordController.Update)
	routing.Get("/api/records/:id", recordController.Show)

	handler := middleware.BuildPipeline(db, []middleware.Middleware{
		middleware.SetHeaders,
		middleware.Authenticate,
		routing.Middleware,
	})

	err = http.ListenAndServe(":"+os.Getenv("APP_PORT"), handler)
	if err != nil {
		panic(err)
	}
}
