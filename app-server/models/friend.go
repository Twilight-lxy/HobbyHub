package models

import (
	"reflect"
	"time"
)

type Friend struct {
	Id         int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:'记录Id'"`
	UserId     int64     `json:"user_id" gorm:"index;not null;foreignKey:Id;references:id;comment:'用户Id'"`
	User       User      `json:"user" gorm:"foreignKey:UserId;references:Id;comment:'用户'"`
	FriendId   int64     `json:"friend_id" gorm:"index;not null;foreignKey:Id;references:id;comment:'好友Id'"`
	FriendUser User      `json:"friend_user" gorm:"foreignKey:FriendId;references:Id;comment:'好友用户'"`
	Status     int       `json:"status" gorm:"not null default:0;comment:'状态（0: 拒绝, 1: 接受, 2: 等待接受, 3：已发出申请）"`
	CreateTime time.Time `json:"create_time" gorm:"not null;comment:'创建时间'"`
}

func (f *Friend) TableName() string {
	return "friend"
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
