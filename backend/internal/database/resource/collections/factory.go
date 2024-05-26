package collections

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/database/resource/activities"
	"github.com/EugeneNail/actum/internal/service/fake"
	"github.com/EugeneNail/actum/internal/service/tests"
	"strings"
)

type Factory struct {
	userId      int
	collections []Collection
}

func NewFactory(userId int) *Factory {
	return &Factory{userId, []Collection{}}
}

func (factory *Factory) Make(count int) *Factory {
	for i := 0; i < count; i++ {
		collection := Collection{0, fake.SentenceLength(1, 3), factory.userId, make([]activities.Activity, 0)}
		factory.collections = append(factory.collections, collection)
	}

	return factory
}

func (factory *Factory) Insert() *Factory {
	var placeholders []string
	var values []any

	for _, collection := range factory.collections {
		placeholders = append(placeholders, "(?, ?)")
		values = append(values, collection.Name, factory.userId)
	}

	query := fmt.Sprintf("INSERT INTO collections (name, user_id) VALUES %s", strings.Join(placeholders, ", "))

	db, err := mysql.Connect()
	defer db.Close()
	tests.Check(err)

	_, err = db.Exec(query, values...)
	tests.Check(err)

	return factory
}

func (factory *Factory) Get(index int) Collection {
	return factory.collections[index]
}

func (factory *Factory) List() []Collection {
	return factory.collections
}
