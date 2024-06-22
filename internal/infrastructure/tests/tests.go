package tests

import (
	"database/sql"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/infrastructure/env"
)

var DB *sql.DB

func init() {
	env.Load()
	db, err := mysql.Connect()
	Check(err)
	DB = db
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}
