package main

import (
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/resource/activities"
	"github.com/EugeneNail/actum/internal/resource/collections"
	"github.com/EugeneNail/actum/internal/resource/photos"
	"github.com/EugeneNail/actum/internal/resource/records"
	"github.com/EugeneNail/actum/internal/resource/users"
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/middleware"
	"github.com/EugeneNail/actum/internal/service/middleware/routing"
	"github.com/EugeneNail/actum/internal/service/refresh"
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
	refreshService := refresh.NewService(db)
	userController := users.NewController(db, userDAO, refreshService)
	routing.Post("/api/users", userController.Store)
	routing.Post("/api/users/login", userController.Login)
	routing.Post("/api/users/refresh-token", userController.RefreshToken)
	routing.Post("/api/users/logout", userController.Logout)

	collectionDAO := collections.NewDAO(db)
	collectionService := collections.NewService(db)
	collectionController := collections.NewController(db, collectionDAO, collectionService)
	routing.Post("/api/collections", collectionController.Store)
	routing.Put("/api/collections/:id", collectionController.Update)
	routing.Delete("/api/collections/:id", collectionController.Destroy)
	routing.Get("/api/collections/:id", collectionController.Show)
	routing.Get("/api/collections", collectionController.Index)

	activityDAO := activities.NewDAO(db)
	activityService := activities.NewService(db, activityDAO)
	activityController := activities.NewController(db, activityDAO, collectionDAO, activityService)
	routing.Post("/api/activities", activityController.Store)
	routing.Put("/api/activities/:id", activityController.Update)
	routing.Delete("/api/activities/:id", activityController.Destroy)
	routing.Get("/api/activities/:id", activityController.Show)

	photoDAO := photos.NewDAO(db)
	photoService := photos.NewService(db)
	photoController := photos.NewController(photoDAO)
	routing.Post("/api/photos", photoController.Store)
	routing.Delete("/api/photos/:name", photoController.Destroy)
	routing.Get("/api/photos/:name", photoController.Show)

	recordDAO := records.NewDAO(db)
	recordService := records.NewService(db)
	recordController := records.NewController(db, recordDAO, activityDAO, activityService, recordService, photoService)
	routing.Post("/api/records", recordController.Store)
	routing.Put("/api/records/:id", recordController.Update)
	routing.Get("/api/records/:id", recordController.Show)
	routing.Post("/api/records-list", recordController.Index)

	handler := middleware.BuildPipeline(db, []middleware.Middleware{
		middleware.RedirectToFrontend,
		middleware.SetHeaders,
		middleware.Authenticate,
		routing.Middleware,
	})

	err = http.ListenAndServe(":"+os.Getenv("APP_PORT"), handler)
	if err != nil {
		panic(err)
	}
}
