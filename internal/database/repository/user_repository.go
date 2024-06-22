package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/EugeneNail/actum/internal/entity"
	"github.com/EugeneNail/actum/internal/infrastructure/errors"
	"github.com/EugeneNail/actum/internal/infrastructure/hash"
	"strings"
)

type mySQLUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) UserRepository {
	return &mySQLUserRepository{db}
}

func (repo *mySQLUserRepository) Create(ctx context.Context, name string, email string, password string) (int64, error) {
	result, err := repo.db.ExecContext(
		ctx,
		`INSERT INTO users(name, email, password) VALUES (?, ?, ?)`,
		name, strings.ToLower(email), hash.New(password),
	)

	if err != nil {
		return 0, errors.Wrap(err, "failed to insert data")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "failed to retrieve id")
	}

	return id, nil
}

func (repo *mySQLUserRepository) Find(ctx context.Context, column string, value any) (entity.User, error) {
	var user entity.User

	query := fmt.Sprintf(`SELECT id, name, email, password FROM users WHERE %s = ?`, column)
	err := repo.db.QueryRowContext(ctx, query, value).
		Scan(&user.Id, &user.Name, &user.Email, &user.Password)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return user, errors.Wrap(err, "failed to retrieve the user")
	}

	return user, nil
}
