package main

import (
	"errors"
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
		log.Fatal(errors.New("object is missing"))
	}
	switch object := os.Args[1]; object {
	case "migration":
		createMigration()
	}

}

func createMigration() {
	if len(os.Args) == 2 {
		log.Fatal(errors.New("name of migration is missing"))
	}

	name := os.Args[2]
	now := time.Now().Unix()
	migrations := getMigrationsDirectory()
	up := fmt.Sprintf("%s/%d.%s.%s.sql", migrations, now, name, "up")
	down := fmt.Sprintf("%s/%d.%s.%s.sql", migrations, now, name, "down")

	_, err := os.Create(up)
	if err != nil {
		log.Fatal(err)
	}
	_, err = os.Create(down)
	if err != nil {
		log.Fatal(err)
	}
}

func getMigrationsDirectory() string {
	return filepath.Join(os.Getenv("APP_PATH"), "internal", "database", "migrations")
}
