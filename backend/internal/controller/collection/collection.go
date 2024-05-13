package collection

import (
	"github.com/EugeneNail/actum/internal/model/users"
	"strings"
)

func hasDuplicateCollection(name string, user users.User) (bool, error) {
	collections, err := user.Collections()
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
