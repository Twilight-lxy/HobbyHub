package models

import (
	"reflect"
	"time"
)

type Friend struct {
	Id         int64     `json:"id"`
	UserId     int64     `json:"user_id"`
	FriendId   int64     `json:"friend_id"`
	Status     int       `json:"status"` // 0: Pending, 1: Accepted, 2: Rejected
	CreateTime time.Time `json:"create_time"`
}

func (f *Friend) TableName() string {
	return "it_friend"
}

func (u *Friend) UpdateFriendFields(newu Friend) {
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
