package controllers

import (
	"hobbyhub-server/config"
	"hobbyhub-server/models"
)

func AddFriend(friend *models.Friend) error {
	if err := config.DB.Create(friend).Error; err != nil {
		return err
	}
	return nil
}

func GetFriendById(friendId int64) (*models.Friend, error) {
	var friend models.Friend
	if err := config.DB.Where("id = ?", friendId).First(&friend).Error; err != nil {
		return nil, err
	}
	return &friend, nil
}

func GetAllFriendsIdByUserId(userId int64) ([]int64, error) {
	var friendIds []int64
	if err := config.DB.Model(&models.Friend{}).
		Where("user_id = ?", userId).
		Pluck("id", &friendIds).Error; err != nil {
		return nil, err
	}
	return friendIds, nil
}
func UpdateFriend(friend *models.Friend) error {
	if err := config.DB.Save(friend).Error; err != nil {
		return err
	}
	return nil
}
func DeleteFriendById(friendId int64) error {
	if err := config.DB.Delete(&models.Friend{}, friendId).Error; err != nil {
		return err
	}
	return nil
}
func GetFriendByUserIdAndFriendId(userId, friendId int64) (*models.Friend, error) {
	var friend models.Friend
	if err := config.DB.Where("user_id = ? AND friend_id = ?", userId, friendId).First(&friend).Error; err != nil {
		return nil, err
	}
	return &friend, nil
}
