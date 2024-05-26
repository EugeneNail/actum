package activities

type Activity struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Icon         string `json:"icon"`
	CollectionId int    `json:"collectionId"`
	UserId       int    `json:"userId"`
}

func New(name string, icon string, collectionId int, userId int) Activity {
	return Activity{0, name, icon, collectionId, userId}
}
