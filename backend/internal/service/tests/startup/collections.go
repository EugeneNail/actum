package startup

import (
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/tests"
	"testing"
)

func CollectionsStore(t *testing.T) (tests.Client, tests.Database) {
	env.Load()

	t.Cleanup(func() {
		err := mysql.Truncate("users")
		tests.Check(err)
		err = mysql.Truncate("collections")
		tests.Check(err)
	})

	return tests.NewClient(t), tests.NewDatabase(t)
}

func CollectionsUpdate(t *testing.T) (tests.Client, tests.Database) {
	env.Load()

	t.Cleanup(func() {
		err := mysql.Truncate("users")
		tests.Check(err)
		err = mysql.Truncate("collections")
		tests.Check(err)
	})

	return tests.NewClient(t), tests.NewDatabase(t)
}
