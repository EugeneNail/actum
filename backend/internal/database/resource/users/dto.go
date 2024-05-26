package users

import (
	"github.com/EugeneNail/actum/internal/resource/collections"
)

type User struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	collections []collections.Collection
}

func New(name string, email string, password string) User {
	return User{0, name, email, password, []collections.Collection{}}
}
