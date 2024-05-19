package collections

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/model/activities"
)

type Collection struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	UserId int    `json:"userId"`
}

func Find(id int) (Collection, error) {
	var collection Collection

	db, err := mysql.Connect()
	defer db.Close()
	if err != nil {
		return collection, fmt.Errorf("collection.Find(): %w", err)
	}

	rows, err := db.Query(`SELECT * FROM collections WHERE id = ?`, id)
	defer rows.Close()
	if err != nil {
		return collection, fmt.Errorf("collection.Find(): %w", err)
	}

	for rows.Next() {
		err := rows.Scan(&collection.Id, &collection.Name, &collection.UserId)
		if err != nil {
			return collection, fmt.Errorf("collection.Find(): %w", err)
		}
	}

	return collection, nil
}

func (collection *Collection) Save() error {
	db, err := mysql.Connect()
	defer db.Close()
	if err != nil {
		return fmt.Errorf("groups.Save(): %w", err)
	}

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

func (collection *Collection) Delete() error {
	db, err := mysql.Connect()
	defer db.Close()
	if err != nil {
		return fmt.Errorf("collections.Delete(): %w", err)
	}

	_, err = db.Exec(`DELETE FROM collections WHERE id = ?`, collection.Id)

	if err != nil {
		return fmt.Errorf("collections.Delete(): %w", err)
	}

	return nil
}

func (collection *Collection) Activities() ([]activities.Activity, error) {
	var result []activities.Activity

	db, err := mysql.Connect()
	defer db.Close()
	if err != nil {
		return result, fmt.Errorf("collections.Activities(): %w", err)
	}

	rows, err := db.Query(`SELECT * FROM activities WHERE collection_id = ?`, collection.Id)
	defer rows.Close()
	if err != nil {
		return result, fmt.Errorf("collections.Activities(): %w", err)
	}

	for rows.Next() {
		var activity activities.Activity
		err := rows.Scan(&activity.Id, &activity.Name, &activity.Icon, &activity.UserId, &activity.CollectionId)
		if err != nil {
			return result, fmt.Errorf("collections.Activities(): %w", err)
		}
		result = append(result, activity)
	}

	return result, nil
}
