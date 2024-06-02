package collections

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

func (dao *DAO) Find(id int) (Collection, error) {
	var collection Collection

	err := dao.db.QueryRow(`SELECT id, name, color, user_id FROM collections WHERE id = ?`, id).
		Scan(&collection.Id, &collection.Name, &collection.Color, &collection.UserId)

	if err != nil && err != sql.ErrNoRows {
		return collection, fmt.Errorf("collection.Find(): %w", err)
	}

	return collection, nil
}

func (dao *DAO) Save(collection *Collection) error {
	result, err := dao.db.Exec(`
		INSERT INTO collections 
		    (id, name, color, user_id)
		VALUES 
		    (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			id = VALUES(id),
			name = VALUES(name),
			color = VALUES(color),
			user_id = VALUES(user_id);
	`, collection.Id, collection.Name, collection.Color, collection.UserId)

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

func (dao *DAO) Delete(collection Collection) error {
	_, err := dao.db.Exec(`DELETE FROM collections WHERE id = ?`, collection.Id)

	if err != nil {
		return fmt.Errorf("collections.Delete(): %w", err)
	}

	return nil
}
