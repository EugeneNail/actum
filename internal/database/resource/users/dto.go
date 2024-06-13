package users

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func New(name string, email string, password string) User {
	return User{0, name, email, password}
}
