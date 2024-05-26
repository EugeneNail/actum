package users

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (repo *Repository) Find(id int) (User, error) {
	var user User

	err := repo.db.QueryRow(`SELECT * FROM users WHERE id = ? LIMIT 1`, id).
		Scan(&user.Id, &user.Name, &user.Email, &user.Password)

	if err != nil && err != sql.ErrNoRows {
		return user, fmt.Errorf("users.Find(): %w", err)
	}

	return user, nil
}

func (repo *Repository) FindBy(column string, value any) (User, error) {
	var user User

	err := repo.db.QueryRow(`SELECT * FROM users WHERE `+column+` = ?`, value).
		Scan(&user.Id, &user.Name, &user.Email, &user.Password)

	if err != nil && err != sql.ErrNoRows {
		return user, fmt.Errorf("users.FindBy(): %w", err)
	}

	return user, nil
}

func (repo *Repository) Save(user *User) error {
	result, err := repo.db.Exec(`
    INSERT INTO users 
        (id, name, email, password) 
    VALUES 
        (?, ?, ?, ?)
    ON DUPLICATE KEY UPDATE 
		name = VALUES(name), 
        email = VALUES(email), 
        password = VALUES(password);
    `, user.Id, user.Name, user.Email, user.Password)

	if err != nil {
		return fmt.Errorf("users.Save(): %w", err)
	}
	id, err := result.LastInsertId()

	if err != nil {
		return fmt.Errorf("users.Save(): %w", err)
	}

	if id != 0 {
		user.Id = int(id)
	}

	return nil
}
