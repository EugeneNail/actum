package main

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/service/env"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	env.Load()

	if len(os.Args) == 1 {
		fmt.Println("subcommand required")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "migration":
		createMigration()
	default:
		panic("unknown subcommand")
	}

}

func createMigration() {
	name := "nameless_migration"
	if len(os.Args) > 2 {
		name = os.Args[2]
	}
	now := time.Now().Unix()
	pathToMigrations := filepath.Join(
		os.Getenv("APP_PATH"), "internal", "database", "migrations",
	)
	_, err := os.Stat(pathToMigrations)

	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(pathToMigrations, 0755)
		check(err)
	}

	up := fmt.Sprintf("%s/%d.%s.%s.sql", pathToMigrations, now, name, "up")
	down := fmt.Sprintf("%s/%d.%s.%s.sql", pathToMigrations, now, name, "down")
	_, err = os.Create(up)
	check(err)
	_, err = os.Create(down)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
