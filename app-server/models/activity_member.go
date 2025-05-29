package models

type ActivityMember struct {
	ID         int64  `json:"id"`
	EventID    int64  `json:"eventId"`
	UserID     int64  `json:"userId"`
	CreateTime string `json:"createTime"`
}

func (ActivityMember) TableName() string {
	return "it_event_member"
}
