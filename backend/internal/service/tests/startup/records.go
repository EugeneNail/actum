package startup

import (
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/tests"
	"testing"
)

func RecordsStore(t *testing.T) (tests.Client, tests.Database) {
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
