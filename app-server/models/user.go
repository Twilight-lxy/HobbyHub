package models

import "reflect"

type User struct {
	ID         int64   `json:"id" gorm:"primaryKey"`
	Username   string  `json:"username" gorm:"unique"`
	Password   string  `json:"password"`
	Name       string  `json:"name"`
	Gender     string  `json:"gender"`
	Addr       string  `json:"addr"`
	HeadImg    string  `json:"headImg"`
	CreateTime string  `json:"createTime"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
}

func (User) TableName() string {
	return "it_user"
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
