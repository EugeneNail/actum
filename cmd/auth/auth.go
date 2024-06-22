package main

import (
	"fmt"
	pb "github.com/EugeneNail/actum/grpc/gen/auth"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/database/repository"
	"github.com/EugeneNail/actum/internal/infrastructure/env"
	"github.com/EugeneNail/actum/internal/service/auth"
	"github.com/EugeneNail/actum/internal/service/auth/refresh"
	"google.golang.org/grpc"
	"net"
)

func main() {
	env.Load()
	db := mysql.MustConnect()
	userRepo := repository.NewMySQLUserRepository(db)
	refresh := refresh.NewService(db)
	service := auth.NewAuthenticationService(userRepo, refresh)
	server := auth.NewServer(service)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAuthServer(grpcServer, server)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "localhost", 50100))
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}

}
