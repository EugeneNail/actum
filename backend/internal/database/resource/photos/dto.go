package photos

type Photo struct {
	Id       int    `json:"id"`
	Name     string `json:"uuid"`
	RecordId *int   `json:"recordId"`
	UserId   int    `json:"userId"`
}

func New(name string, recordId *int, userId int) Photo {
	return Photo{0, name, recordId, userId}
}
