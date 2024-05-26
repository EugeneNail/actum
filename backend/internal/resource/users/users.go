package users

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/database/resource/collections"
	"github.com/EugeneNail/actum/internal/database/resource/users"
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

func Find(id int) (users.User, error) {
	var user users.User
	db, err := mysql.Connect()
	defer db.Close()
	if err != nil {
		return user, fmt.Errorf("users.Find(): %w", err)
	}

	rows, err := db.Query(`SELECT * FROM users WHERE id = ? LIMIT 1`, id)
	defer rows.Close()
	if err != nil {
		return user, fmt.Errorf("users.Find(): %w", err)
	}

	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password); err != nil {
			return user, fmt.Errorf("users.Find(): %w", err)
		}
	}

	return user, nil
}

func FindBy(column string, value any) (users.User, error) {
	var user users.User
	db, err := mysql.Connect()
	defer db.Close()

	if err != nil {
		return user, fmt.Errorf("users.FindBy(): %w", err)
	}

	rows, err := db.Query(`SELECT * FROM users WHERE `+column+` = ?`, value)
	defer rows.Close()
	if err != nil {
		return user, fmt.Errorf("users.FindBy(): %w", err)
	}

	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return user, fmt.Errorf("users.FindBy(): %w", err)
		}
	}

	return user, nil
}

func list() ([]users.User, error) {
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

	if id != 0 {
		this.Id = int(id)
	}

	return nil
}

func (this *User) delete() error {
	panic("not implemented")
}
