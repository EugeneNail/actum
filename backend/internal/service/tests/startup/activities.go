package startup

import (
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/tests"
	"testing"
)

func ActivitiesStore(t *testing.T) (tests.Client, tests.Database) {
	env.Load()

	t.Cleanup(func() {
		err := mysql.Truncate(tests.DB, "activities")
		tests.Check(err)
		err = mysql.Truncate(tests.DB, "collections")
		tests.Check(err)
		err = mysql.Truncate(tests.DB, "users")
		tests.Check(err)
	})

	return tests.NewClient(t), tests.NewDatabase(t)
}

func ActivitiesShow(t *testing.T) (tests.Client, tests.Database) {
	env.Load()

	t.Cleanup(func() {
		err := mysql.Truncate(tests.DB, "activities")
		tests.Check(err)
		err = mysql.Truncate(tests.DB, "collections")
		tests.Check(err)
		err = mysql.Truncate(tests.DB, "users")
		tests.Check(err)
	})

	return tests.NewClient(t), tests.NewDatabase(t)
}

func ActivitiesUpdate(t *testing.T) (tests.Client, tests.Database) {
	env.Load()

	t.Cleanup(func() {
		err := mysql.Truncate(tests.DB, "activities")
		tests.Check(err)
		err = mysql.Truncate(tests.DB, "collections")
		tests.Check(err)
		err = mysql.Truncate(tests.DB, "users")
		tests.Check(err)
	})

	return tests.NewClient(t), tests.NewDatabase(t)
}

func ActivitiesDestroy(t *testing.T) (tests.Client, tests.Database) {
	env.Load()

	t.Cleanup(func() {
		err := mysql.Truncate(tests.DB, "activities")
		tests.Check(err)
		err = mysql.Truncate(tests.DB, "collections")
		tests.Check(err)
		err = mysql.Truncate(tests.DB, "users")
		tests.Check(err)
	})

	return tests.NewClient(t), tests.NewDatabase(t)
}
