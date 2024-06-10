package refresh

import (
	"database/sql"
	"fmt"
	"github.com/EugeneNail/actum/internal/service/hash"
	"github.com/EugeneNail/actum/internal/service/uuid"
	"time"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db}
}

func (service *Service) MakeToken(userId int) (string, error) {
	uuid := uuid.New()
	hashedUuid := hash.New(uuid)
	expiration := time.Now().Add(time.Hour * 24 * 7)

	_, err := service.db.Exec(
		`INSERT INTO user_refresh_tokens(uuid, user_id, expired_at) VALUES(?, ?, ?)`,
		hashedUuid, userId, expiration,
	)
	if err != nil {
		return "", fmt.Errorf("refresh.MakeToken(): %w", err)
	}

	return uuid, nil
}

func (service *Service) IsValid(token string, userId int) (bool, error) {
	var count int

	err := service.db.QueryRow(`
		SELECT COUNT(*) 
		FROM user_refresh_tokens 
		WHERE uuid = ? AND user_id = ? AND CURRENT_TIMESTAMP <= expired_at
		`,
		hash.New(token), userId,
	).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("refresh.IsValid(): %w", err)
	}

	return count > 0, nil
}

func (service *Service) Unset(userId int) error {
	_, err := service.db.Exec(`DELETE FROM user_refresh_tokens WHERE user_id = ?`, userId)
	if err != nil {
		return fmt.Errorf("refresh.Unset(): %w", err)
	}

	return nil
}
