package users

import "github.com/EugeneNail/actum/internal/database/resource/users"

type Controller struct {
	dao *users.DAO
}

func New(dao *users.DAO) Controller {
	return Controller{dao}
}
