package photos

import (
	"database/sql"
	"fmt"
)

type DAO struct {
	db *sql.DB
}

func NewDAO(db *sql.DB) *DAO {
	return &DAO{db}
}

func (dao *DAO) Find(name string) (Photo, error) {
	var photo Photo
	err := dao.db.QueryRow(`
		SELECT id, name, record_id, user_id FROM photos WHERE name = ?`, name,
	).
		Scan(&photo.Id, &photo.Name, &photo.RecordId, &photo.UserId)
	if err != nil {
		return photo, fmt.Errorf("photos.Find(): %w", err)
	}

	return photo, nil
}

func (dao *DAO) FindBy(column string, name string) (Photo, error) {
	var photo Photo

	query := fmt.Sprintf(`SELECT id, name, record_id, user_id FROM photos WHERE %s = ?`, column)
	err := dao.db.QueryRow(query, name).
		Scan(&photo.Id, &photo.Name, &photo.RecordId, &photo.UserId)

	if err != nil && err != sql.ErrNoRows {
		return photo, fmt.Errorf("photos.Find(): %w", err)
	}

	return photo, nil
}

func (dao *DAO) Save(photo *Photo) error {
	result, err := dao.db.Exec(`
		INSERT INTO photos
			(name, record_id, user_id) 
		VALUES 
		    (?, ?, ?)
		ON DUPLICATE KEY UPDATE 
			name = VALUES(name),
			record_id = VALUES(record_id),
			user_id = VALUES(user_id)
	`, photo.Name, photo.RecordId, photo.UserId)

	if err != nil {
		return fmt.Errorf("photos.Save(): %w", err)
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("photos.Save(): %w", err)
	}

	if lastInsertId != 0 {
		photo.Id = int(lastInsertId)
	}

	return nil
}

func (dao *DAO) Delete(id int) error {
	_, err := dao.db.Exec(`DELETE FROM photos WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("photos.Delete: failed to delete photo %d:%w", id, err)
	}

	return nil
}
