package entity

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

func NewUser(name string, email string, password string) User {
	return User{0, name, email, password}
}
