package collections

import "github.com/EugeneNail/actum/internal/resource/activities"

type Collection struct {
	Id         int                   `json:"id"`
	Name       string                `json:"name"`
	UserId     int                   `json:"userId"`
	Activities []activities.Activity `json:"activities"`
}

func New(name string, userId int) Collection {
	return Collection{0, name, userId, []activities.Activity{}}
}
