package collection

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/resource/collections"
	"github.com/EugeneNail/actum/internal/resource/users"
	"strings"
)

func hasDuplicateCollection(name string, user users.User) (bool, error) {
	collections, err := getUserCollections(user.Id)
	if err != nil {
		return false, err
	}

	for _, collection := range collections {
		if strings.ToLower(collection.Name) == strings.ToLower(name) {
			return true, nil
		}
	}

	return false, nil
}

func getUserCollections(userId int) ([]collections.Collection, error) {
	var userCollections []collections.Collection

	db, err := mysql.Connect()
	defer db.Close()
	if err != nil {
		return userCollections, fmt.Errorf("user.Collections(): %w", err)
	}

	rows, err := db.Query(`SELECT * FROM collections WHERE user_id = ?`, userId)
	defer rows.Close()
	if err != nil {
		return userCollections, fmt.Errorf("user.Collections(): %w", err)
	}

	for rows.Next() {
		var collection collections.Collection

		err := rows.Scan(&collection.Id, &collection.Name, &collection.UserId)
		if err != nil {
			return userCollections, fmt.Errorf("user.Collections(): %w", err)
		}

		userCollections = append(userCollections, collection)
	}

	return userCollections, nil
}
