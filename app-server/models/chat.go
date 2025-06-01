package models

import (
	"reflect"
	"time"
)

type Chat struct {
	Id         int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:'记录Id'"`
	UserIdFrom int64     `json:"user_id_from" gorm:"index;not null;foreignKey:Id;references:id;comment:'发送用户Id'"`
	UserFrom   User      `json:"user_from" gorm:"foreignKey:UserIdFrom;references:Id;comment:'发送用户'"`
	UserIdTo   int64     `json:"user_id_to" gorm:"index;not null;foreignKey:Id;references:id;comment:'接收用户Id'"`
	UserTo     User      `json:"user_to" gorm:"foreignKey:UserIdTo;references:Id;comment:'接收用户'"`
	Content    string    `json:"content" gorm:"type:text;not null;comment:'消息内容'"`
	CreateTime time.Time `json:"create_time" gorm:"not null;comment:'创建时间'"`
	StatusFrom int32     `json:"status_from" gorm:"not null;default:2;comment:'发送方状态（0: 删除, 1: 正常）'"`
	StatusTo   int32     `json:"status_to" gorm:"not null;default:2;comment:'接收方状态（0: 删除, 1: 正常）'"`
}

func (Chat) TableName() string {
	return "chat"
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
