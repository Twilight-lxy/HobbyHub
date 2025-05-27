package models

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
