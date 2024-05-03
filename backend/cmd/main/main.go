package main

import (
	"github.com/EugeneNail/actum/internal/controller/user"
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/routing"
	"net/http"
	"os"
)

func main() {
	env.Load()
	routing.Post("/api/users", user.Store)
	routing.Post("/api/users/login", user.Login)
	go log.RotateFiles()
	err := http.ListenAndServe(":"+os.Getenv("APP_PORT"), routing.Serve())

	if err != nil {
		panic(err)
	}
}
