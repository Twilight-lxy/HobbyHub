package models

import "time"

type File struct {
	Id           int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:'文件Id'"`
	FileName     string    `json:"file_name" gorm:"type:varchar(255);not null;comment:'文件名'"`
	FileType     string    `json:"file_type" gorm:"type:varchar(50);not null;index;comment:'文件类型'"`
	FileSize     int64     `json:"file_size" gorm:"not null;comment:'文件大小(字节)'"`
	FileHash     string    `json:"file_hash" gorm:"type:varchar(74);index;unique;comment:'文件哈希值'"`
	CreateTime   time.Time `json:"create_time" gorm:"not null;comment:'创建时间'"`
	LinkFileId   int64     `json:"link_file_id" gorm:"not null;default:0;index;comment:'关联文件Id'"`
	UpLoadUserId int64     `json:"upload_user_id" gorm:"not null;index;comment:'上传用户Id'"`
}

// 定义表名
func (File) TableName() string {
	return "file"
}
