package records

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/database/resource/activities"
	"github.com/EugeneNail/actum/internal/database/resource/collections"
	"github.com/EugeneNail/actum/internal/resource/records"
	"github.com/EugeneNail/actum/internal/service/fake"
	"github.com/EugeneNail/actum/internal/service/tests"
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestStoreValidData(t *testing.T) {
	client, database := startup.Records(t)

	collections.NewFactory(1).Make(1).Insert()
	activities.NewFactory(1, 1).Make(5).Insert().List()
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

	notes := fake.Text()
	client.
		Post("/api/records", fmt.Sprintf(`{
			"mood": 1,
			"notes": "%s",
			"weather": 7,
			"date": "2024-05-29",
			"activities": [1,2,3,5],
			"photos": ["%s"]
		}`, notes, name)).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("records_activities", 4).
		AssertCount("records", 1).
		AssertHas("records", map[string]any{
			"id":      1,
			"mood":    1,
			"weather": 7,
			"notes":   notes,
			"date":    "2024-05-29",
			"user_id": 1,
		}).
		AssertCount("photos", 1).
		AssertHas("photos", map[string]any{
			"name":      name,
			"record_id": 1,
			"user_id":   1,
		})

	for _, activityId := range []int{1, 2, 3, 5} {
		database.AssertHas("records_activities", map[string]any{
			"record_id":   1,
			"activity_id": activityId,
		})
	}
}

func TestStoreInvalidData(t *testing.T) {
	client, database := startup.Records(t)

	collections.NewFactory(1).Make(1).Insert()
	activities.NewFactory(1, 1).Make(5).Insert().List()

	database.AssertCount("collections", 1).AssertCount("activities", 5)

	notes := strings.Repeat("abcdefghij", 5000/10+1)
	client.
		Post("/api/records", fmt.Sprintf(`{
			"mood": 0,
			"weather": 4,
			"notes": "%s",
			"date": "2024-35-29",
			"photos": []
		}`, notes)).
		AssertStatus(http.StatusUnprocessableEntity)

	database.
		AssertCount("records_activities", 0).
		AssertCount("records", 0).
		AssertLacks("records", map[string]any{
			"mood":    0,
			"weather": 4,
			"notes":   notes,
			"date":    "2024-35-29",
			"user_id": 1,
		})
}

func TestStoreMissingActivity(t *testing.T) {
	client, database := startup.Records(t)

	collections.NewFactory(1).Make(1).Insert()
	activities.NewFactory(1, 1).Make(2).Insert().List()

	database.AssertCount("collections", 1).AssertCount("activities", 2)

	client.
		Post("/api/records", `{
			"mood": 1,
			"weather": 2,
			"date": "2020-10-29",
			"activities": [1,2,3]
		}`).
		AssertStatus(http.StatusNotFound)

	database.AssertCount("records", 0).AssertCount("records_activities", 0)
}

func TestStoreConflictDate(t *testing.T) {
	client, database := startup.Records(t)

	collections.NewFactory(1).Make(1).Insert()
	activities.NewFactory(1, 1).Make(5).Insert().List()

	database.AssertCount("collections", 1).AssertCount("activities", 5)

	notes := fake.Text()
	client.
		Post("/api/records", fmt.Sprintf(`{
			"mood": 1,
			"weather": 3,
			"notes": "%s",
			"date": "2022-01-01",
			"activities": [1,2,3,4,5],
			"photos": []
		}`, notes)).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("records_activities", 5).
		AssertCount("records", 1).
		AssertHas("records", map[string]any{
			"mood":    1,
			"weather": 3,
			"notes":   notes,
			"date":    "2022-01-01",
			"user_id": 1,
		})

	for _, activityId := range []int{1, 2, 3, 4, 5} {
		database.AssertHas("records_activities", map[string]any{
			"record_id":   1,
			"activity_id": activityId,
		})
	}

	client.
		Post("/api/records", fmt.Sprintf(`{
			"mood": 1,
			"weather": 4,
			"notes": "%s",
			"date": "2022-01-01",
			"activities": [1,2,3,4,5]
		}`, notes)).
		AssertStatus(http.StatusConflict)

	database.
		AssertCount("records_activities", 5).
		AssertCount("records", 1).
		AssertHas("records", map[string]any{
			"mood":    1,
			"weather": 3,
			"notes":   notes,
			"date":    "2022-01-01",
			"user_id": 1,
		})

	for _, activityId := range []int{1, 2, 3, 4, 5} {
		database.AssertHas("records_activities", map[string]any{
			"record_id":   1,
			"activity_id": activityId,
		})
	}
}

func TestStoreValidation(t *testing.T) {
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().Add(time.Hour * 24 * -1).Format("2006-01-02")
	yearAgo := time.Now().Add(time.Hour * 24 * 365 * -1).Format("2006-01-02")

	tests.AssertValidationSuccess[records.storeInput](t, []tests.ValidationTest{
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
		{"notes", "Short", fake.SentenceLength(1, 2)},
		{"notes", "Average", fake.Sentence()},
		{"notes", "Long", fake.Text()},
		{"date", "Today", today},
		{"date", "Yesterday", yesterday},
		{"date", "Long ago", "2020-01-01"},
		{"date", "Year ago", yearAgo},
		{"date", "Just date", "2024-05-25"},
		{"date", "New Year", "2023-12-31"},
		{"date", "After the New Year", "2024-01-01"},
		{"activities", "Activities", []int{1, 2, 3}},
	})

	tests.AssertValidationFail[records.storeInput](t, []tests.ValidationTest{
		{"mood", "Mood 1", 0},
		{"mood", "Mood 2", -1},
		{"mood", "Mood 3", 6},
		{"mood", "Mood 4", 66},
		{"mood", "Mood 5", 0.3},
		{"weather", "Zero", 0},
		{"weather", "Non integer", 1.1},
		{"weather", "Negative", -1},
		{"weather", "Nonexistent", 10},
		{"notes", "Too long", strings.Repeat("Something", 600)},
		{"date", "Too long ago", "2019-12-31"},
		{"date", "Invalid year", "024-13-25"},
		{"date", "Invalid month", "2024-13-25"},
		{"date", "Invalid date", "2023-12-33"},
		{"date", "Tomorrow", time.Now().Add(time.Hour * 24).Format("2006-01-02")},
		{"activities", "Empty activities", []int{}},
	})
}
