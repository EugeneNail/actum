package records

import (
	"github.com/EugeneNail/actum/internal/database/resource/activities"
	"github.com/EugeneNail/actum/internal/database/resource/collections"
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"testing"
)

func TestShow(t *testing.T) {
	client, database := startup.Records(t)

	collections.NewFactory(1).Make(1).Insert()
	activities.NewFactory(1, 1).Make(4).Insert()

	client.
		Post("/api/records", `{
			"date": "2024-01-01",
			"mood": 2,
			"activities": [1, 2, 3, 4],
			"notes": "My cat has died"
		}`).
		AssertStatus(http.StatusCreated)

	database.
		AssertCount("records", 1).
		AssertCount("records_activities", 4).
		AssertHas("records", map[string]any{
			"id":      1,
			"date":    "2024-01-01",
			"mood":    2,
			"notes":   "My cat has died",
			"user_id": 1,
		})

	for i := 1; i <= 4; i++ {
		database.AssertHas("records_activities", map[string]any{
			"record_id":   1,
			"activity_id": i,
		})
	}

	var record struct {
		Date       string `json:"date"`
		Mood       int    `json:"mood"`
		Notes      string `json:"notes"`
		Activities []int  `json:"activities"`
	}

	client.
		Get("/api/records/1").
		AssertStatus(http.StatusOK).
		ReadData(&record)

	if record.Date != "2024-01-01" ||
		record.Mood != 2 ||
		record.Notes != "My cat has died" ||
		len(record.Activities) != 4 {
		t.Errorf("Incorrect data")
	}
}

func TestShowNotFound(t *testing.T) {
	client, database := startup.Records(t)

	database.AssertCount("records", 0)
	client.Get("/api/records/1").AssertStatus(http.StatusNotFound)
}

func TestShowSomeoneElsesRecord(t *testing.T) {
	client, database := startup.Records(t)

	collections.NewFactory(1).Make(1).Insert()
	activities.NewFactory(1, 1).Make(3).Insert()

	database.AssertCount("activities", 3)

	client.
		Post("/api/records", `{
			"mood": 1,
			"date": "2020-01-01",
			"notes": "Test",
			"activities": [1]
		}`).
		AssertStatus(http.StatusCreated)

	client.ChangeUser()
	client.
		Get("/api/records/1").
		AssertStatus(http.StatusForbidden)
}
