package activities

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func (repo *Repository) Find(id int) (Activity, error) {
	var activity Activity

	err := repo.db.QueryRow(`SELECT * FROM activities WHERE id = ?`, id).
		Scan(&activity.Id, &activity.Name, &activity.Icon, &activity.UserId, &activity.CollectionId)

	if err != nil && err != sql.ErrNoRows {
		return activity, fmt.Errorf("activities.Find(): %w", err)
	}

	return activity, nil
}

func (repo *Repository) Save(activity *Activity) error {
	result, err := repo.db.Exec(`
		INSERT INTO activities
		    (id, name, icon, collection_id, user_id)
		VALUES
		    (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			id = VALUES(id),
			name = VALUES(name),
			icon = VALUES(icon),
			collection_id = VALUES(collection_id),
			user_id = VALUES(user_id);
	`, activity.Id, activity.Name, activity.Icon, activity.CollectionId, activity.UserId)

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

func (repo *Repository) Delete(activity *Activity) error {
	_, err := repo.db.Exec(`DELETE FROM activities WHERE id = ?`, activity.Id)
	if err != nil {
		return fmt.Errorf("activities.Delete(): %w", err)
	}

	return nil
}
