package startup

import (
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/tests"
	"testing"
)

func Users(t *testing.T) (tests.Client, tests.Database) {
	env.Load()

	t.Cleanup(func() {
		tests.Check(mysql.Truncate(tests.DB, "users"))
	})

	return tests.NewClientWithoutAuth(t), tests.NewDatabase(t)
}

func Collections(t *testing.T) (tests.Client, tests.Database) {
	env.Load()

	t.Cleanup(func() {
		tests.Check(mysql.TruncateMany(tests.DB, []string{"users", "collections"}))
	})

	client := tests.NewClient(t)
	database := tests.NewDatabase(t)

	database.AssertHas("users", map[string]any{
		"id": 1,
	})

	return client, database
}

func Activities(t *testing.T) (tests.Client, tests.Database) {
	env.Load()

	t.Cleanup(func() {
		tests.Check(mysql.TruncateMany(tests.DB, []string{"activities", "collections", "users"}))
	})

	client := tests.NewClient(t)
	database := tests.NewDatabase(t)

	database.AssertHas("users", map[string]any{
		"id": 1,
	})

	return client, database
}

func Records(t *testing.T) (tests.Client, tests.Database) {
	env.Load()

	t.Cleanup(func() {
		tests.Check(mysql.TruncateMany(tests.DB, []string{"records", "activities", "collections", "users"}))
	})

	client := tests.NewClient(t)
	database := tests.NewDatabase(t)

	database.AssertHas("users", map[string]any{
		"id": 1,
	})

	return client, database
}
