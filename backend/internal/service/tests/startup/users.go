package startup

import (
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/tests"
	"testing"
)

func UsersStore(t *testing.T) (tests.Client, tests.Database) {
	env.Load()

	t.Cleanup(func() {
		err := mysql.Truncate("users")
		tests.Check(err)
	})

	return tests.NewClientWithoutAuth(t), tests.NewDatabase(t)
}

func UsersLogin(t *testing.T) (tests.Client, tests.Database) {
	env.Load()

	t.Cleanup(func() {
		err := mysql.Truncate("users")
		tests.Check(err)
	})

	return tests.NewClientWithoutAuth(t), tests.NewDatabase(t)
}
