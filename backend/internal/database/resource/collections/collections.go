package collections

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/database/mysql"
)

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
