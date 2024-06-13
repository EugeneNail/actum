package activities

import (
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/service/fake"
	"github.com/EugeneNail/actum/internal/service/tests"
)

type Factory struct {
	userId       int
	collectionId int
	activities   []Activity
}

func NewFactory(userId int, collectionId int) *Factory {
	return &Factory{userId, collectionId, []Activity{}}
}

func (factory *Factory) Make(count int) *Factory {
	factory.activities = make([]Activity, count)

	for i := 0; i < count; i++ {
		activity := Activity{
			0,
			fake.SentenceLength(3, 4),
			fake.Icon(),
			factory.collectionId,
			factory.userId,
		}
		factory.activities[i] = activity
	}

	return factory
}

func (factory *Factory) Insert() *Factory {
	const columnsCount = 4
	var placeholders string
	values := make([]any, len(factory.activities)*columnsCount)

	for i, activity := range factory.activities {
		placeholders += "(?, ?, ?, ?),"
		values[columnsCount*i+0] = activity.Name
		values[columnsCount*i+1] = activity.Icon
		values[columnsCount*i+2] = activity.UserId
		values[columnsCount*i+3] = activity.CollectionId
	}
	placeholders = placeholders[:len(placeholders)-1]

	db, err := mysql.Connect()
	defer db.Close()
	tests.Check(err)

	_, err = db.Exec(
		`INSERT INTO activities (name, icon, user_id, collection_id) VALUES`+placeholders,
		values...,
	)
	tests.Check(err)

	return factory
}

func (factory *Factory) Get(index int) Activity {
	return factory.activities[index]
}

func (factory *Factory) List() []Activity {
	return factory.activities
}
