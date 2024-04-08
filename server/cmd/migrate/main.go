package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/service/env"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type migration struct {
	name string
	path string
}

func main() {
	env.Load()
	createTable()
	rollbackCommand := flag.NewFlagSet("rollback", flag.ExitOnError)
	steps := rollbackCommand.Int("steps", -1, "number of migrations to rollback")

	switch {
	case len(os.Args) == 1:
		migrate()
	case os.Args[1] == "rollback":
		err := rollbackCommand.Parse(os.Args[2:])
		check(err)
		rollback(*steps)
	default:
		fmt.Println("invalid command arguments")
		os.Exit(1)
	}
}

func rollback(steps int) {
	migrations := getMigrations("down")
	appliedMigrations := getAppliedMigrations()
	var isRolledBack bool

	for i := len(appliedMigrations) - 1; i >= 0; i-- {
		if steps == 0 {
			break
		}
		currentMigration := getCurrentMigration(migrations, appliedMigrations[i])
		apply(currentMigration)
		apply(currentMigration)
		steps--
		isRolledBack = true
	}

	if !isRolledBack {
		fmt.Println("Nothing to rollback")
	}
}

func getCurrentMigration(migrations []migration, appliedMigration string) migration {
	var currentMigration migration

	for _, migration := range migrations {
		if migration.name == appliedMigration {
			currentMigration = migration
			break
		}
	}

	if currentMigration.name == "" {
		log.Fatal("no migration file found for " + appliedMigration)
	}

	return currentMigration
}

func migrate() {
	var migrations = getMigrations("up")
	appliedMigrations := getAppliedMigrations()
	var isMigrated bool

	for _, migration := range migrations {
		if !slices.Contains(appliedMigrations, migration.name) {
			isMigrated = true
			apply(migration)
		}
	}

	if !isMigrated {
		fmt.Println("Nothing to migrate")
	}
}

func apply(migration migration) {
	query, err := os.ReadFile(migration.path)
	check(err)
	db, err := mysql.Connect()
	check(err)
	transaction, err := db.BeginTx(context.Background(), &sql.TxOptions{})
	check(err)
	_, err = transaction.Exec(string(query))
	checkTransaction(err, migration, transaction)

	if strings.Contains(migration.path, ".up.") {
		_, err := transaction.Exec("INSERT INTO migrations (migration) VALUES (?)", migration.name)
		checkTransaction(err, migration, transaction)
	}

	if strings.Contains(migration.path, ".down.") {
		_, err := transaction.Exec("DELETE FROM migrations WHERE migration = ?", migration.name)
		checkTransaction(err, migration, transaction)
	}
	err = transaction.Commit()
	check(err)
	fmt.Println("DONE", migration.name)
}

func getAppliedMigrations() []string {
	db, err := mysql.Connect()
	check(err)
	var appliedMigrations []string
	rows, err := db.Query("SELECT migration FROM migrations")
	check(err)

	defer rows.Close()

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		check(err)
		appliedMigrations = append(appliedMigrations, name)
	}

	return appliedMigrations
}

func getMigrations(direction string) []migration {
	migrationsDirectory := filepath.Join(
		os.Getenv("APP_PATH"), "internal", "database", "migrations",
	)
	files, err := os.ReadDir(migrationsDirectory)
	check(err)
	var migrations = make([]migration, 0, len(files)/2)
	suffix := "." + direction + ".sql"

	for _, file := range files {
		if strings.HasSuffix(file.Name(), suffix) {
			migration := migration{
				name: strings.Replace(file.Name(), suffix, "", -1),
				path: filepath.Join(migrationsDirectory, file.Name()),
			}
			migrations = append(migrations, migration)
		}
	}

	return migrations
}

func createTable() {
	db, err := mysql.Connect()
	check(err)
	query := ` 
	CREATE TABLE IF NOT EXISTS migrations (
	    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		migrations VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	);`
	_, err = db.Exec(query)
	check(err)
}

func checkTransaction(err error, migration migration, transaction *sql.Tx) {
	if err != nil {
		transaction.Rollback()
		fmt.Println("FAIL", migration.name, "\n")
		log.Fatal(err)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)

	}
}
