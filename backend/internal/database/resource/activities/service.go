package activities

import (
	"database/sql"
	"fmt"
	"strings"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db}
}

func (service *Service) HasDuplicate(name string, userId int) (bool, error) {
	var count int

	err := service.db.
		QueryRow(`SELECT COUNT(*) FROM activities WHERE user_id = ? AND LOWER(name) = ?`, userId, strings.ToLower(name)).
		Scan(&count)
	if err != nil {
		return false, fmt.Errorf("activities.HasDuplicate(): %w", err)
	}

	return count > 0, nil
}

func (service *Service) ExceedsLimit(limit int, collectionId int, userId int) (bool, error) {
	var count int

	err := service.db.
		QueryRow(`SELECT COUNT(*) FROM activities WHERE user_id = ? AND collection_id = ?`, userId, collectionId).
		Scan(&count)
	if err != nil {
		return false, fmt.Errorf("activities.ExceedsLimit(): %w", err)
	}

	return count >= limit, nil
}
