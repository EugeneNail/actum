package photos

import (
	"github.com/EugeneNail/actum/internal/database/resource/photos"
)

type Controller struct {
	dao *photos.DAO
}

func New(dao *photos.DAO) *Controller {
	return &Controller{dao}
}
