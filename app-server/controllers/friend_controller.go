package controllers

import (
	"hobbyhub-server/config"
	"hobbyhub-server/models"
)

func AddFriend(friend *models.Friend) error {
	if err := config.DB.Create(friend).Error; err != nil {
		return err
	}
	_friend := models.Friend{
		UserId:     friend.FriendId,
		FriendId:   friend.UserId,
		Status:     friend.Status,
		CreateTime: friend.CreateTime,
	}
	if err := config.DB.Create(&_friend).Error; err != nil {
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

func GetAllFriendsByUserId(userId int64) ([]models.Friend, error) {
	var friends []models.Friend
	if err := config.DB.Where("user_id = ?", userId).Find(&friends).Error; err != nil {
		return nil, err
	}
	return friends, nil
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
