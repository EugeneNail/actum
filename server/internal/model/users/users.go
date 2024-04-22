package users

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/database/mysql"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func find(id int) (User, error) {
	panic("not implemented")
}

func list() ([]User, error) {
	panic("not implemented")
}

func (this *User) Save() error {
	db, err := mysql.Connect()
	defer db.Close()

	if err != nil {
		return fmt.Errorf("users.Save(): %w", err)
	}
	result, err := db.Exec(`
    INSERT INTO users 
        (id, name, email, password) 
    VALUES 
        (?, ?, ?, ?)
    ON DUPLICATE KEY UPDATE 
		name = VALUES(name), 
        email = VALUES(email), 
        password = VALUES(password);
    `, this.Id, this.Name, this.Email, this.Password)

	if err != nil {
		return fmt.Errorf("users.Save(): %w", err)
	}
	id, err := result.LastInsertId()

	if err != nil {
		return fmt.Errorf("users.Save(): %w", err)
	}
	this.Id = int(id)

	return nil
}

func (this *User) delete() error {
	panic("not implemented")
}
