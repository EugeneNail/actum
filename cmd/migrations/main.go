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
	"time"
)

type migration struct {
	name string
	path string
}

func main() {
	env.Load()
	createTable()
	createDirectory()
	args := os.Args[1:]

	if len(args) == 0 {
		printError("Expected create, apply, rollback or refresh subcommand")
		return
	}

	switch args[0] {
	case "create":
		createMigration()
	case "apply":
		applyMigrations()
	case "rollback":
		rollbackCommand := flag.NewFlagSet("rollback", flag.ExitOnError)
		steps := rollbackCommand.Int("steps", -1, "number of migrations to rollback")
		err := rollbackCommand.Parse(os.Args[2:])
		check(err)
		rollback(*steps)
	case "refresh":
		rollback(-1)
		applyMigrations()
	default:
		printError("Expected create, apply, rollback or refresh subcommand")
	}

}

func createTable() {
	db, err := mysql.Connect()
	check(err)
	query := ` 
	CREATE TABLE IF NOT EXISTS migrations (
	    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		migration VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	);`
	_, err = db.Exec(query)
	check(err)
}

func createDirectory() {
	pathToMigrations := filepath.Join(
		os.Getenv("APP_PATH"), "internal", "database", "migrations",
	)
	_, err := os.Stat(pathToMigrations)
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(pathToMigrations, 0755)
		check(err)
	}
}

func createMigration() {
	if len(os.Args) < 3 {
		printError("Expected name of migration")
		return
	}

	name := os.Args[2]
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

func applyMigrations() {
	var migrations = getMigrations("up")
	appliedMigrations := getAppliedMigrations()
	var isMigrated bool

	for _, migration := range migrations {
		if !slices.Contains(appliedMigrations, migration.name) {
			isMigrated = true
			applyMigration(migration)
		}
	}

	if !isMigrated {
		fmt.Println("Nothing to migrate")
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
		applyMigration(currentMigration)
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

func applyMigration(migration migration) {
	file, err := os.ReadFile(migration.path)
	check(err)
	queries := strings.Split(string(file), ";")
	db, err := mysql.Connect()
	check(err)
	transaction, err := db.BeginTx(context.Background(), &sql.TxOptions{})
	check(err)

	for _, query := range queries {
		_, err = transaction.Exec(query)
		checkTransaction(err, migration, transaction)
	}

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

func printError(message string) {
	fmt.Println("\033[31m" + message + "\033[0m")
}
