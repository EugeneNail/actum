package activities

type Activity struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Icon         int    `json:"icon"`
	CollectionId int    `json:"collectionId"`
	UserId       int    `json:"userId"`
}

type IndexActivity struct {
	RecordId     int    `json:"recordId"`
	CollectionId int    `json:"collectionId"`
	Icon         int    `json:"icon"`
	Name         string `json:"name"`
}

func New(name string, icon int, collectionId int, userId int) Activity {
	return Activity{0, name, icon, collectionId, userId}
}
