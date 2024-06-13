package photos

import (
	"database/sql"
	"fmt"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db}
}

func (service *Service) CheckExistence(needlePhotos []string, userId int) (bool, []string, error) {
	if len(needlePhotos) == 0 {
		return true, []string{}, nil
	}

	photoNames, err := service.collectNamesOfPhotos(needlePhotos, userId)
	if err != nil {
		return false, photoNames, fmt.Errorf("photos.CheckExistence: failed to get names of the photos: %w", err)
	}

	var missingPhotos []string
	photoMap := make(map[string]struct{}, len(needlePhotos))
	for _, name := range photoNames {
		photoMap[name] = struct{}{}
	}

	for _, needlePhoto := range needlePhotos {
		if _, exists := photoMap[needlePhoto]; !exists {
			missingPhotos = append(missingPhotos, needlePhoto)
		}
	}

	return len(missingPhotos) == 0, missingPhotos, nil
}

func (service *Service) collectNamesOfPhotos(needlePhotos []string, userId int) ([]string, error) {
	var photoNames []string

	var placeholders string
	values := make([]any, len(needlePhotos)+1)
	values[0] = userId

	for i, photo := range needlePhotos {
		placeholders += "?,"
		values[i+1] = photo
	}

	placeholders = "(" + placeholders[:len(placeholders)-1] + ")"

	rows, err := service.db.Query(`SELECT name FROM photos WHERE user_id = ? AND name IN`+placeholders, values...)
	defer rows.Close()
	if err != nil {
		return photoNames, fmt.Errorf("collectNamesOfPhotos: failed to fetch photos: %w", err)
	}

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return photoNames, fmt.Errorf("collectNamesOfPhotos: failed to scan the name of the photo: %w", err)
		}
		photoNames = append(photoNames, name)
	}

	return photoNames, nil
}

func (service *Service) SyncRelations(recordId int, photoNames []string) error {
	if len(photoNames) == 0 {
		return nil
	}

	var placeholders string
	values := make([]any, len(photoNames)+1)
	values[0] = recordId

	for i, name := range photoNames {
		placeholders += "?,"
		values[i+1] = name
	}
	placeholders = "(" + placeholders[:len(placeholders)-1] + ")"

	if _, err := service.db.Exec(`UPDATE photos SET record_id = ? WHERE name IN`+placeholders, values...); err != nil {
		return fmt.Errorf("records.SyncRelations: failed to update relations: %w", err)
	}

	return nil
}
