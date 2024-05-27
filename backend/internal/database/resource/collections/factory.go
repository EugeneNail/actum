package collections

import (
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/service/fake"
	"github.com/EugeneNail/actum/internal/service/tests"
)

type Factory struct {
	userId      int
	collections []Collection
}

func NewFactory(userId int) *Factory {
	return &Factory{userId, []Collection{}}
}

func (factory *Factory) Make(count int) *Factory {
	factory.collections = make([]Collection, count)

	for i := 0; i < count; i++ {
		factory.collections[i] = New(fake.SentenceLength(1, 3), factory.userId)
	}

	return factory
}

func (factory *Factory) Insert() *Factory {
	const columnsCount = 2
	var placeholders string
	values := make([]any, len(factory.collections)*columnsCount)

	for i, collection := range factory.collections {
		placeholders += "(?, ?),"
		values[columnsCount*i+0] = collection.Name
		values[columnsCount*i+1] = factory.userId
	}
	placeholders = placeholders[:len(placeholders)-1]

	db, err := mysql.Connect()
	defer db.Close()
	tests.Check(err)

	_, err = db.Exec(
		`INSERT INTO collections (name, user_id) VALUES`+placeholders,
		values...,
	)
	tests.Check(err)

	return factory
}

func (factory *Factory) Get(index int) Collection {
	return factory.collections[index]
}

func (factory *Factory) List() []Collection {
	return factory.collections
}
