package collections

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/database/mysql"
)

type Collection struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	UserId int    `json:"userId"`
}

func Find(id int) (Collection, error) {
	var collection Collection

	db, err := mysql.Connect()
	if err != nil {
		return collection, fmt.Errorf("collection.Find(): %w", err)
	}
	defer db.Close()

	err = db.
		QueryRow("SELECT * FROM collections WHERE id = ?", id).
		Scan(&collection.Id, &collection.Name, &collection.UserId)
	if err != nil {
		return collection, fmt.Errorf("collection.Find(): %w", err)
	}

	return collection, nil
}

func (collection *Collection) Save() error {
	db, err := mysql.Connect()
	if err != nil {
		return fmt.Errorf("groups.Save(): %w", err)
	}
	defer db.Close()

	result, err := db.Exec(`
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
