package records

import (
	"fmt"
	"time"
)

type Record struct {
	Id      int       `json:"id"`
	Mood    int       `json:"mood"`
	Date    time.Time `json:"date"`
	Weather int       `json:"weather"`
	Notes   string    `json:"notes"`
	UserId  int       `json:"userId"`
}

type IndexRecord struct {
	Id          int               `json:"id"`
	Date        string            `json:"date"`
	Weather     int               `json:"weather"`
	Mood        int               `json:"mood"`
	Notes       string            `json:"notes"`
	Collections []IndexCollection `json:"collections"`
	Photos      []string          `json:"photos"`
}

type IndexCollection struct {
	Id         int             `json:"-"`
	Name       string          `json:"name"`
	Color      int             `json:"color"`
	Activities []IndexActivity `json:"activities"`
}

type IndexActivity struct {
	RecordId     int    `json:"-"`
	CollectionId int    `json:"-"`
	Icon         int    `json:"icon"`
	Name         string `json:"name"`
}

func New(mood int, weather int, date string, notes string, userId int) (Record, error) {
	time, err := time.Parse("2006-01-02", date)
	if err != nil {
		return Record{}, fmt.Errorf("records.New(): %w", err)
	}

	return Record{0, mood, time, weather, notes, userId}, nil
}
