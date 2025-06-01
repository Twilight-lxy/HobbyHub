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
		Status:     2,
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

// GetAllFriendsByUserId 获取用户的所有好友关系，并预加载好友用户信息
func GetAllFriendsByUserId(userId int64) ([]models.Friend, error) {
	var friends []models.Friend

	// 使用Preload预加载关联的User和FriendUser数据
	if err := config.DB.Where("user_id = ?", userId).
		Preload("User").
		Preload("FriendUser").
		Find(&friends).Error; err != nil {
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
func UpdateFriendSynchronize(friend1 *models.Friend, friend2 *models.Friend) error {
	// 开启事务
	tx := config.DB.Begin()

	// 事务操作出错时进行回滚的函数
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 在事务中更新第一个好友关系
	if err := tx.Save(friend1).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}

	// 在事务中更新第二个好友关系
	if err := tx.Save(friend2).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

func DeleteFriendById(friendId int64) error {
	if err := config.DB.Delete(&models.Friend{}, friendId).Error; err != nil {
		return err
	}
	return nil
}
func GetFriendByUserIdAndFriendId(userId, friendId int64) (*models.Friend, *models.Friend, error) {
	var friends1 *models.Friend
	var friends2 *models.Friend

	// First query
	if err := config.DB.Where("user_id = ? AND friend_id = ?", userId, friendId).Find(&friends1).Error; err != nil {
		return nil, nil, err
	}

	// Second query
	if err := config.DB.Where("user_id = ? AND friend_id = ?", friendId, userId).Find(&friends2).Error; err != nil {
		return nil, nil, err
	}
	return friends1, friends2, nil
}
