package records

import (
	"github.com/EugeneNail/actum/internal/database/resource/activities"
	"github.com/EugeneNail/actum/internal/database/resource/collections"
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

	database.AssertCount("collections", 1).AssertCount("activities", 5)

	client.
		Post("/api/records", `{
			"mood": 1,
			"date": "2024-01-01",
			"notes": "There is not one to save the humanity",
			"activities": [1, 2, 3, 4, 5]
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("records", 1).
		AssertHas("records", map[string]any{
			"mood":    1,
			"date":    "2024-01-01",
			"notes":   "There is not one to save the humanity",
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
			"mood": 5,
			"notes": "Look over here",
			"activities": [1, 2, 4, 5]
		}`).
		AssertStatus(http.StatusNoContent)

	database.
		AssertCount("records", 1).
		AssertHas("records", map[string]any{
			"mood":  5,
			"date":  "2024-01-01",
			"notes": "Look over here",
		}).
		AssertCount("records_activities", 4).
		AssertLacks("records_activities", map[string]any{
			"record_id":   1,
			"activity_id": 3,
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
			"date": "2024-01-01",
			"notes": "",
			"activities": [1, 2, 3, 4, 5]
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("records", 1).
		AssertHas("records", map[string]any{
			"mood":    1,
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
			"activities": []
		}`).
		AssertStatus(http.StatusUnprocessableEntity)

	database.
		AssertCount("records", 1).
		AssertHas("records", map[string]any{
			"mood":    1,
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
			"notes": "Look over here",
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
			"notes": "I was sad",
			"activities": [2, 3, 4, 5, 6]
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("records", 1).
		AssertHas("records", map[string]any{
			"mood":    5,
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
			"notes": "Look over here",
			"activities": [1, 2, 3, 4, 5, 6]
		}`).
		AssertStatus(http.StatusNotFound)

	database.
		AssertCount("records", 1).
		AssertHas("records", map[string]any{
			"mood":    5,
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
		{"notes", "Empty", ""},
		{"notes", "Short", "Something good happened"},
		{"notes", "Sentence 1", "The quick brown fox jumps over the lazy dog."},
		{"notes", "Sentence 1", "Today i was arriving from Mexico on airplane №34"},
		{"activities", "Activities 1", []int{1, 2, 3}},
	})

	tests.AssertValidationFail[updateInput](t, []tests.ValidationTest{
		{"mood", "Zero", 0},
		{"mood", "Non integer", 1.1},
		{"mood", "Negative", -1},
		{"mood", "Nonexistent", 6},
		{"notes", "Too long", strings.Repeat("Here are 23 characters", 250)},
		{"activities", "Empty", []int{}},
	})
}
