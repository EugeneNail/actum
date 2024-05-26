package users

import "github.com/EugeneNail/actum/internal/database/resource/users"

type Controller struct {
	repository *users.Repository
}

func New(repository *users.Repository) Controller {
	return Controller{repository}
}
