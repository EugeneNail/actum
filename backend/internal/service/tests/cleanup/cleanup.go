package cleanup

import (
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/service/tests"
)

func LoginUsers() {
	err := mysql.Truncate("users")
	tests.Check(err)
}

func StoreUsers() {
	err := mysql.Truncate("users")
	tests.Check(err)
}
