package models

import (
	"reflect"
	"time"
)

type Activity struct {
	Id         int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:'活动Id'"`
	Name       string    `json:"name" gorm:"type:varchar(255);not null;comment:'活动名称'"`
	Addr       string    `json:"addr" gorm:"type:varchar(255);comment:'活动地址'"`
	Intro      string    `json:"intro" gorm:"type:text;comment:'活动简介'"`
	HeadImg    string    `json:"headImg" gorm:"type:varchar(255);comment:'活动封面图片'"`
	UserId     int64     `json:"userId" gorm:"index;not null;comment:'创建者Id'"`
	CreateTime time.Time `json:"createTime" gorm:"not null;comment:'创建时间'"`
	UpdateTime time.Time `json:"updateTime" gorm:"not null;comment:'更新时间'"`
	StartTime  time.Time `json:"startTime" gorm:"not null;comment:'活动开始时间'"`
	State      int       `json:"state" gorm:"not null;default:0;comment:'活动状态（0: 未开始, 1: 进行中, 2: 已结束, 3: 已取消）'"`
	IfDelete   int       `json:"ifDelete" gorm:"not null;default:0;comment:'删除状态（0: 正常, 1: 已删除）'"`
	Lat        float64   `json:"lat" gorm:"comment:'纬度'"`
	Lon        float64   `json:"lon" gorm:"comment:'经度'"`
}

func (Activity) TableName() string {
	return "activity"
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
	Id         int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:'记录Id'"`
	EventId    int64     `json:"eventId" gorm:"index;not null;comment:'活动Id'"`
	Activity   Activity  `json:"-" gorm:"foreignKey:EventId;references:Id;comment:'活动'"`
	UserId     int64     `json:"userId" gorm:"index;not null;comment:'用户Id'"`
	User       User      `json:"user" gorm:"foreignKey:UserId;references:Id;comment:'用户'"`
	CreateTime time.Time `json:"createTime" gorm:"not null;comment:'创建时间'"`
}

func (ActivityMember) TableName() string {
	return "activity_member"
}

func (u *ActivityMember) UpdateActivityMemberFields(newu ActivityMember) {
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

type ActivityComment struct {
	Id         int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:'记录Id'"`
	EventId    int64     `json:"eventId" gorm:"index;not null;comment:'活动Id'"`
	Activity   Activity  `json:"-" gorm:"foreignKey:EventId;references:Id;comment:'活动'"`
	UserId     int64     `json:"userId" gorm:"index;not null;comment:'用户Id'"`
	User       User      `json:"user" gorm:"foreignKey:UserId;references:Id;comment:'用户'"`
	Content    string    `json:"content" gorm:"type:text;not null;comment:'评论内容'"`
	CreateTime time.Time `json:"createTime" gorm:"not null;comment:'创建时间'"`
}

func (ActivityComment) TableName() string {
	return "activity_comment"
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
