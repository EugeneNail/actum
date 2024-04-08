package main

import (
	"flag"
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

	switch subcommand := os.Args[1]; subcommand {
	case "migration":
		createMigration()
	default:
		panic("unknown subcommand")
	}

}

func createMigration() {
	name := *flag.
		NewFlagSet("migration", flag.ExitOnError).
		String("name", "nameless_migration", "migration name")
	now := time.Now().Unix()
	migrations := filepath.Join(os.Getenv("APP_PATH"), "internal", "database", "migrations")
	up := fmt.Sprintf("%s/%d.%s.%s.sql", migrations, now, name, "up")
	down := fmt.Sprintf("%s/%d.%s.%s.sql", migrations, now, name, "down")
	_, err := os.Create(up)
	check(err)
	_, err = os.Create(down)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
