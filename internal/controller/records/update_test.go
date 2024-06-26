package records

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/database/resource/activities"
	"github.com/EugeneNail/actum/internal/database/resource/collections"
	"github.com/EugeneNail/actum/internal/service/fake"
	"github.com/EugeneNail/actum/internal/service/tests"
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"strings"
	"testing"
)

func TestUpdateValidData(t *testing.T) {
	client, database := startup.Records(t)

	collections.NewFactory(1).Make(1).Insert()
	activities.NewFactory(1, 1).Make(5).Insert()
	var name string
	client.
		Post("/api/photos", fmt.Sprintf(`{
			"image": "%s"
		}`, fake.Base64Image())).
		AssertStatus(http.StatusCreated).
		ReadData(&name)

	database.
		AssertCount("collections", 1).
		AssertCount("activities", 5).
		AssertCount("photos", 1)

	client.
		Post("/api/records", `{
			"mood": 1,
			"weather": 1,
			"date": "2024-01-01",
			"notes": "Иисус, жги!",
			"activities": [1, 2, 3, 4, 5]
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("records", 1).
		AssertHas("records", map[string]any{
			"mood":    1,
			"weather": 1,
			"date":    "2024-01-01",
			"notes":   "Иисус, жги!",
			"user_id": 1,
		}).
		AssertCount("records_activities", 5)

	for i := 1; i <= 5; i++ {
		database.AssertHas("records_activities", map[string]any{
			"record_id":   1,
			"activity_id": i,
		})
	}

	client.
		Put("/api/records/1", fmt.Sprintf(`{
			"mood": 5,
			"weather": 2,
			"notes": "Сюда смотри",
			"activities": [1, 2, 4, 5],
			"photos": ["%s"]
		}`, name)).
		AssertStatus(http.StatusNoContent)

	database.
		AssertCount("records", 1).
		AssertHas("records", map[string]any{
			"id":      1,
			"mood":    5,
			"weather": 2,
			"date":    "2024-01-01",
			"notes":   "Сюда смотри",
		}).
		AssertCount("records_activities", 4).
		AssertLacks("records_activities", map[string]any{
			"record_id":   1,
			"activity_id": 3,
		}).
		AssertHas("photos", map[string]any{
			"record_id": 1,
			"name":      name,
			"user_id":   1,
		})
}

func TestUpdateInvalidData(t *testing.T) {
	client, database := startup.Records(t)

	collections.NewFactory(1).Make(1).Insert()
	activities.NewFactory(1, 1).Make(5).Insert()

	database.AssertCount("collections", 1).AssertCount("activities", 5)

	client.
		Post("/api/records", `{
			"mood": 1,
			"weather": 3,
			"date": "2024-01-01",
			"notes": "",
			"activities": [1, 2, 3, 4, 5]
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("records", 1).
		AssertHas("records", map[string]any{
			"mood":    1,
			"weather": 3,
			"date":    "2024-01-01",
			"notes":   "",
			"user_id": 1,
		}).
		AssertCount("records_activities", 5)

	for i := 1; i <= 5; i++ {
		database.AssertHas("records_activities", map[string]any{
			"record_id":   1,
			"activity_id": i,
		})
	}

	client.
		Put("/api/records/1", `{
			"mood": 0,
			"weather": 0,
			"activities": []
		}`).
		AssertStatus(http.StatusUnprocessableEntity)

	database.
		AssertCount("records", 1).
		AssertHas("records", map[string]any{
			"mood":    1,
			"weather": 3,
			"date":    "2024-01-01",
			"notes":   "",
			"user_id": 1,
		}).
		AssertCount("records_activities", 5)
}

func TestUpdateInvalidId(t *testing.T) {
	client, database := startup.Records(t)

	client.
		Put("/api/records/one", `{
			"mood": 1,
			"weather": 1,
			"activities": [1]
		}`).
		AssertStatus(http.StatusBadRequest)

	database.
		AssertCount("records", 0).
		AssertCount("records_activities", 0)
}

func TestUpdateNonexistentActivity(t *testing.T) {
	client, database := startup.Records(t)

	database.AssertCount("records", 0)

	client.
		Put("/api/records/1", `{
			"mood": 5,
			"weather": 5,
			"notes": "Сюда смотри",
			"activities": [1, 2, 4, 5]
		}`).
		AssertStatus(http.StatusNotFound)

	database.
		AssertCount("records", 0).
		AssertCount("records_activities", 0)
}

func TestUpdateWithSomeoneElsesActivities(t *testing.T) {
	client, database := startup.Records(t)

	collections.NewFactory(1).Make(1).Insert()
	activities.NewFactory(1, 1).Make(1).Insert()
	database.AssertCount("collections", 1).AssertCount("activities", 1)

	client.ChangeUser()
	collections.NewFactory(2).Make(1).Insert()
	activities.NewFactory(2, 2).Make(5).Insert()
	database.AssertCount("collections", 2).AssertCount("activities", 6)

	client.
		Post("/api/records", `{
			"mood": 5,
			"date": "2022-12-31",
			"weather": 8,
			"notes": "I was sad",
			"activities": [2, 3, 4, 5, 6]
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("records", 1).
		AssertHas("records", map[string]any{
			"mood":    5,
			"weather": 8,
			"date":    "2022-12-31",
			"notes":   "I was sad",
			"user_id": 2,
		}).
		AssertCount("records_activities", 5)

	for i := 2; i <= 6; i++ {
		database.AssertHas("records_activities", map[string]any{
			"record_id":   1,
			"activity_id": i,
		})
	}

	client.
		Put("/api/records/1", `{
			"mood": 5,
			"weather": 9,
			"notes": "Сюда смотри",
			"activities": [1, 2, 3, 4, 5, 6]
		}`).
		AssertStatus(http.StatusNotFound)

	database.
		AssertCount("records", 1).
		AssertHas("records", map[string]any{
			"mood":    5,
			"weather": 8,
			"date":    "2022-12-31",
			"notes":   "I was sad",
			"user_id": 2,
		}).
		AssertCount("records_activities", 5).
		AssertLacks("records_activities", map[string]any{
			"record_id":   1,
			"activity_id": 1,
		})
}

func TestUpdateValidation(t *testing.T) {
	tests.AssertValidationSuccess[updateInput](t, []tests.ValidationTest{
		{"mood", "Mood 1", 1},
		{"mood", "Mood 2", 2},
		{"mood", "Mood 3", 3},
		{"mood", "Mood 4", 4},
		{"mood", "Mood 5", 5},
		{"weather", "Weather 1", 1},
		{"weather", "Weather 2", 2},
		{"weather", "Weather 3", 3},
		{"weather", "Weather 4", 4},
		{"weather", "Weather 5", 5},
		{"weather", "Weather 6", 6},
		{"weather", "Weather 7", 7},
		{"weather", "Weather 8", 8},
		{"weather", "Weather 9", 9},
		{"notes", "Empty", ""},
		{"notes", "Short", "Что-то было"},
		{"notes", "Sentence 1", "Поешь этимх мягких булок да выпей час"},
		{"notes", "Sentence 1", "Я прилетел вчера. Причем, летел слишком долго, как это бывает на боингах №707"},
		{"activities", "Activities 1", []int{1, 2, 3}},
	})

	tests.AssertValidationFail[updateInput](t, []tests.ValidationTest{
		{"mood", "Zero", 0},
		{"mood", "Non integer", 1.1},
		{"mood", "Negative", -1},
		{"mood", "Nonexistent", 6},
		{"weather", "Zero", 0},
		{"weather", "Non integer", 1.1},
		{"weather", "Negative", -1},
		{"weather", "Nonexistent", 10},
		{"notes", "Too long", strings.Repeat("Тут 15 символов", 350)},
		{"activities", "Empty", []int{}},
	})
}
