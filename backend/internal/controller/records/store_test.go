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
	"time"
)

func TestStoreValidData(t *testing.T) {
	client, database := startup.Records(t)

	collections.NewFactory(1).Make(1).Insert()
	activities.NewFactory(1, 1).Make(5).Insert().List()

	database.AssertCount("collections", 1).AssertCount("activities", 5)

	notes := fake.Text()
	client.
		Post("/api/records", fmt.Sprintf(`{
			"mood": 1,
			"notes": "%s",
			"date": "2024-05-29",
			"activities": [1,2,3,5]
		}`, notes)).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("records_activities", 4).
		AssertCount("records", 1).
		AssertHas("records", map[string]any{
			"mood":    1,
			"notes":   notes,
			"date":    "2024-05-29",
			"user_id": 1,
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
			"notes": "%s",
			"date": "2024-35-29"
		}`, notes)).
		AssertStatus(http.StatusUnprocessableEntity)

	database.
		AssertCount("records_activities", 0).
		AssertCount("records", 0).
		AssertLacks("records", map[string]any{
			"mood":    0,
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
			"notes": "%s",
			"date": "2022-01-01",
			"activities": [1,2,3,4,5]
		}`, notes)).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("records_activities", 5).
		AssertCount("records", 1).
		AssertHas("records", map[string]any{
			"mood":    1,
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
	tests.AssertValidationSuccess[storeInput](t, []tests.ValidationTest{
		{"Mood 1", "mood", 1},
		{"Mood 2", "mood", 2},
		{"Mood 3", "mood", 3},
		{"Mood 4", "mood", 4},
		{"Mood 5", "mood", 5},
		{"Short", "notes", fake.SentenceLength(1, 2)},
		{"Average", "notes", fake.Sentence()},
		{"Long", "notes", fake.Text()},
		{"Date 1", "date", "2020-01-02"},
		{"Date 2", "date", "2024-05-25"},
		{"New Year", "date", "2023-12-31"},
		{"After the New Year", "date", "2024-01-01"},
		{"Activities", "activities", []int{1, 2, 3}},
	})

	tests.AssertValidationFail[storeInput](t, []tests.ValidationTest{
		{"Mood 1", "mood", 0},
		{"Mood 2", "mood", -1},
		{"Mood 3", "mood", 6},
		{"Mood 4", "mood", 66},
		{"Mood 5", "mood", 0.3},
		{"Too long", "notes", strings.Repeat("Something", 600)},
		{"Too long ago", "date", "2019-12-31"},
		{"Invalid year", "date", "024-13-25"},
		{"Invalid month", "date", "2024-13-25"},
		{"Invalid date", "date", "2023-12-33"},
		{"Tomorrow", "date", time.Now().Add(time.Hour * 24).Format("2006-01-02")},
		{"Empty activities", "activities", []int{}},
	})
}
