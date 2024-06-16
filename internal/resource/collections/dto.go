package collections

import (
	"github.com/EugeneNail/actum/internal/resource/activities"
)

type Collection struct {
	Id         int                   `json:"id"`
	Name       string                `json:"name"`
	Color      int                   `json:"color"`
	UserId     int                   `json:"userId"`
	Activities []activities.Activity `json:"activities"`
}

func New(name string, color int, userId int) Collection {
	return Collection{0, name, color, userId, []activities.Activity{}}
}
