package models

import (
	"reflect"
	"time"
)

type Chat struct {
	Id         int64     `json:"id"`
	UserIdFrom int64     `json:"user_id_from"`
	UserIdTo   int64     `json:"user_id_to"`
	Content    string    `json:"content"`
	CreateTime time.Time `json:"create_time"`
	StatusFrom int32     `json:"status_from"`
	StatusTo   int32     `json:"status_to"`
}

func (Chat) TableName() string {
	return "it_chat"
}

func (u *Chat) UpdateChatFields(newu Chat) {
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
