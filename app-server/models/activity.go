package models

import (
	"reflect"
	"time"
)

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

func (u *Activity) UpdateActivityFields(newu Activity) {
	// 使用反射来检查字段是否为零值，避免硬编码每个字段
	v := reflect.ValueOf(newu)
	t := reflect.TypeOf(newu)

	dbv := reflect.ValueOf(u).Elem()

	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := t.Field(i)

		// 检查字段是否为零值
		if !reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			// 跳过 ID 字段，通常不希望更新主键
			if fieldType.Name == "ID" {
				continue
			}
			// 更新非零值字段
			dbField := dbv.FieldByName(fieldType.Name)
			if dbField.IsValid() && dbField.CanSet() {
				dbField.Set(field)
			}
		}
	}
}

type ActivityMember struct {
	ID         int64     `json:"id"`
	EventID    int64     `json:"eventId"`
	UserID     int64     `json:"userId"`
	CreateTime time.Time `json:"createTime"`
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

func (u *ActivityComment) UpdateActivityCommentFields(newu ActivityComment) {
	// 使用反射来检查字段是否为零值，避免硬编码每个字段
	v := reflect.ValueOf(newu)
	t := reflect.TypeOf(newu)

	dbv := reflect.ValueOf(u).Elem()

	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := t.Field(i)

		// 检查字段是否为零值
		if !reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			// 跳过 ID 字段，通常不希望更新主键
			if fieldType.Name == "ID" {
				continue
			}
			// 更新非零值字段
			dbField := dbv.FieldByName(fieldType.Name)
			if dbField.IsValid() && dbField.CanSet() {
				dbField.Set(field)
			}
		}
	}
}
