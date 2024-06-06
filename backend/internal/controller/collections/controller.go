package collections

import (
	"database/sql"
	"github.com/EugeneNail/actum/internal/database/resource/collections"
)

type Controller struct {
	db      *sql.DB
	dao     *collections.DAO
	service *collections.Service
}

func New(db *sql.DB, dao *collections.DAO, service *collections.Service) Controller {
	return Controller{db, dao, service}
}
