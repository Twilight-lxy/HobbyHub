package models

import "time"

type ActivityComment struct {
	ID         int64     `json:"id"`
	EventID    int64     `json:"eventId"`
	UserID     int64     `json:"userId"`
	Content    string    `json:"content"`
	CreateTime time.Time `json:"createTime"`
}

func (ActivityComment) TableName() string {
	return "it_event_comment"
}
