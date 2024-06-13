package users

import (
	"database/sql"
	"github.com/EugeneNail/actum/internal/database/resource/users"
	"github.com/EugeneNail/actum/internal/service/refresh"
)

type Controller struct {
	db             *sql.DB
	dao            *users.DAO
	refreshService *refresh.Service
}

func New(db *sql.DB, dao *users.DAO, refreshService *refresh.Service) Controller {
	return Controller{db, dao, refreshService}
}
