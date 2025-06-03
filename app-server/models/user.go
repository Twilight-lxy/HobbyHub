package models

import (
	"reflect"
	"time"
)

type User struct {
	Id         int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:'用户Id'"`
	Username   string    `json:"username" gorm:"unique;comment:'用户名'"`
	Password   string    `json:"password" gorm:"comment:'密码'"`
	Name       string    `json:"name" gorm:"comment:'姓名'"`
	Gender     string    `json:"gender" gorm:"comment:'性别'"`
	Addr       string    `json:"addr" gorm:"comment:'地址'"`
	HeadImg    string    `json:"headImg" gorm:"comment:'头像图片'"`
	CreateTime time.Time `json:"createTime" gorm:"comment:'创建时间'"`
	Lat        float64   `json:"lat" gorm:"comment:'纬度'"`
	Lon        float64   `json:"lon" gorm:"comment:'经度'"`
}

func (User) TableName() string {
	return "user"
}

func (u *User) UpdateUserFields(newu User) {
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
