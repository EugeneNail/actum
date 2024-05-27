package records

import (
	"fmt"
	"time"
)

type RecordActivity struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	RecordId int    `json:"recordId"`
}

type Record struct {
	Id         int              `json:"id"`
	Mood       string           `json:"mood"`
	Date       time.Time        `json:"date"`
	Notes      string           `json:"notes"`
	UserId     int              `json:"userId"`
	Activities []RecordActivity `json:"activities"`
}

func New(mood string, date string, notes string, userId int) (Record, error) {
	time, err := time.Parse("2006-01-02", date)
	if err != nil {
		return Record{}, fmt.Errorf("records.New(): %w", err)
	}

	return Record{0, mood, time, notes, userId, []RecordActivity{}}, nil
}
