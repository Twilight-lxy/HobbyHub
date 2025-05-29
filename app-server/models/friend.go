package models

type Friend struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	FriendID   int64  `json:"friend_id"`
	Status     int    `json:"apply"` // 0: Pending, 1: Accepted, 2: Rejected
	CreateTime string `json:"create_time"`
}
