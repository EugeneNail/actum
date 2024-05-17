package activities

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/database/mysql"
)

type Activity struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Icon         string `json:"icon"`
	UserId       int    `json:"userId"`
	CollectionId int    `json:"collectionId"`
}

func Find(id int) (Activity, error) {
	var activity Activity

	db, err := mysql.Connect()
	defer db.Close()
	if err != nil {
		return activity, fmt.Errorf("activities.Find(): %w", err)
	}

	rows, err := db.Query(`SELECT * FROM activities WHERE id = ?`, id)
	defer rows.Close()
	if err != nil {
		return activity, fmt.Errorf("activities.Find(): %w", err)
	}

	for rows.Next() {
		err := rows.Scan(&activity.Id, &activity.Name, &activity.Icon, &activity.UserId, &activity.CollectionId)
		if err != nil {
			return activity, fmt.Errorf("activities.Find(): %w", err)
		}
	}

	return activity, nil
}

func (activity *Activity) Save() error {
	db, err := mysql.Connect()
	defer db.Close()
	if err != nil {
		return fmt.Errorf("activities.Save(): %w", err)
	}

	result, err := db.Exec(`
		INSERT INTO activities 
		    (name, icon, collection_id, user_id)
		VALUES 
		    (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			name = VALUES(name),
			icon = VALUES(icon),
			collection_id = VALUES(collection_id),
			user_id = VALUES(user_id);
	`, activity.Name, activity.Icon, activity.CollectionId, activity.UserId)

	if err != nil {
		return fmt.Errorf("activities.Save(): %w", err)
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("activities.Save(): %w", err)
	}

	if lastInsertId != 0 {
		activity.Id = int(lastInsertId)
	}

	return nil
}
