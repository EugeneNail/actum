package collections

import (
	"database/sql"
	"fmt"
	"github.com/EugeneNail/actum/internal/database/resource/collections"
	"github.com/EugeneNail/actum/internal/database/resource/users"
)

type Controller struct {
	db  *sql.DB
	dao *collections.DAO
}

func New(db *sql.DB, dao *collections.DAO) Controller {
	return Controller{db, dao}
}

func (controller *Controller) hasDuplicate(name string, user users.User) (bool, error) {
	var count int

	err := controller.db.QueryRow(
		`SELECT COUNT(*) FROM collections WHERE user_id = ? AND LOWER(name) = ?`,
		user.Id, name).
		Scan(&count)
	if err != nil {
		return false, fmt.Errorf("collections.hadDuplicate(): %w", err)
	}

	return count > 0, nil
}

func (controller *Controller) exceededLimit(userId int) (bool, error) {
	var count int

	err := controller.db.QueryRow(
		`SELECT COUNT(*) FROM collections WHERE user_id = ?`,
		userId).
		Scan(&count)
	if err != nil {
		return false, fmt.Errorf("collections.exceededLimit(): %w", err)
	}

	if count >= 15 {
		return true, nil
	}

	return false, nil
}