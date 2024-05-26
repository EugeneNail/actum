package activities

import (
	"github.com/EugeneNail/actum/internal/database/mysql"
	"github.com/EugeneNail/actum/internal/service/fake"
	"github.com/EugeneNail/actum/internal/service/tests"
	"strings"
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
	for i := 0; i < count; i++ {
		activity := Activity{
			0,
			fake.SentenceLength(3, 4),
			fake.Icon(),
			factory.collectionId,
			factory.userId,
		}
		factory.activities = append(factory.activities, activity)
	}

	return factory
}

func (factory *Factory) Insert() *Factory {

	var placeholders []string
	var values []any

	for _, activity := range factory.activities {
		placeholders = append(placeholders, "(?, ?, ?, ?)")
		values = append(values, activity.Name, activity.Icon, activity.UserId, activity.CollectionId)
	}
	query := `INSERT INTO activities (name, icon, user_id, collection_id) VALUES` + strings.Join(placeholders, ", ")

	db, err := mysql.Connect()
	defer db.Close()
	tests.Check(err)

	_, err = db.Exec(query, values...)
	tests.Check(err)

	return factory
}

func (factory *Factory) Get(index int) Activity {
	return factory.activities[index]
}

func (factory *Factory) List() []Activity {
	return factory.activities
}
