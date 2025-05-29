package models

type Chat struct {
	ID         int64  `json:"id"`
	UserIDFrom int64  `json:"user_id_from"`
	UserIDTo   int64  `json:"user_id_to"`
	Content    string `json:"content"`
	CreateTime string `json:"create_time"`
	StatusFrom int32  `json:"status_from"`
	StatusTo   int32  `json:"status_to"`
}

func (Chat) TableName() string {
	return "it_chat"
}
