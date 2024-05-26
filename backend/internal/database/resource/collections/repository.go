package collections

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

func (repo *Repository) Find(id int) (Collection, error) {
	var collection Collection

	err := repo.db.QueryRow(`SELECT * FROM collections WHERE id = ?`, id).
		Scan(&collection.Id, &collection.Name, &collection.UserId)

	if err != nil && err != sql.ErrNoRows {
		return collection, fmt.Errorf("collection.Find(): %w", err)
	}

	return collection, nil
}

func (repo *Repository) Save(collection *Collection) error {
	result, err := repo.db.Exec(`
		INSERT INTO collections 
		    (id, name, user_id)
		VALUES 
		    (?, ?, ?)
		ON DUPLICATE KEY UPDATE
			id = VALUES(id),
			name = VALUES(name),
			user_id = VALUES(user_id);
	`, collection.Id, collection.Name, collection.UserId)

	if err != nil {
		return fmt.Errorf("collection.Save(): %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("collection.Save(): %w", err)
	}

	if id != 0 {
		collection.Id = int(id)
	}

	return nil
}

func (repo *Repository) Delete(collection Collection) error {
	_, err := repo.db.Exec(`DELETE FROM collections WHERE id = ?`, collection.Id)

	if err != nil {
		return fmt.Errorf("collections.Delete(): %w", err)
	}

	return nil
}
