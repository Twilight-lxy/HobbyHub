package models

import "time"

type Activity struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Addr       string    `json:"addr"`
	Intro      string    `json:"intro"`
	HeadImg    string    `json:"headImg"`
	UserID     int64     `json:"userId"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
	StartTime  time.Time `json:"startTime"`
	State      int       `json:"state"`
	IfDelete   int       `json:"ifDelete"`
	Lat        float64   `json:"lat"`
	Lon        float64   `json:"lon"`
}

func (Activity) TableName() string {
	return "it_event"
}

type ActivityMember struct {
	ID         int64  `json:"id"`
	EventID    int64  `json:"eventId"`
	UserID     int64  `json:"userId"`
	CreateTime string `json:"createTime"`
}

func (ActivityMember) TableName() string {
	return "it_event_member"
}

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
