package main

import (
	pb "github.com/EugeneNail/actum/grpc/gen/auth"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/database/repository"
	controller "github.com/EugeneNail/actum/internal/http"
	"github.com/EugeneNail/actum/internal/infrastructure/env"
	"github.com/EugeneNail/actum/internal/infrastructure/middleware"
	"github.com/EugeneNail/actum/internal/infrastructure/middleware/routing"
	"github.com/EugeneNail/actum/internal/service/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

func main() {
	env.Load()
	conn, err := grpc.NewClient("localhost:50100", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	authClient := auth.NewClient(pb.NewAuthClient(conn))

	db := mysql.MustConnect()
	userRepo := repository.NewMySQLUserRepository(db)
	userController := controller.NewAuthController(authClient, &userRepo)

	routing.Post("/api/auth/register", userController.Register)
	routing.Post("/api/auth/login", userController.Login)

	pipe := middleware.BuildPipeline([]middleware.Middleware{
		middleware.RedirectToFrontend,
		middleware.SetHeaders,
		routing.Middleware,
	})

	MustListenAndServe(pipe)
}

func MustListenAndServe(handler http.Handler) {
	if err := http.ListenAndServe("localhost:8080", handler); err != nil {
		panic(err)
	}
}
