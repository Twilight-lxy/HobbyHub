package controllers

import (
	"errors"
	"hobbyhub-server/config"
	"hobbyhub-server/models"
)

// AddUser 添加新用户
func AddUser(user *models.User) error {
	return config.DB.Create(user).Error
}
func GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByUserId 通过ID获取用户详情
func GetUserByUserId(userId int64) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUserName 通过用户名获取用户详情
func GetUserByUserName(userName string) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("username = ?", userName).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func UpdateUser(user models.User) error {
	if user.Id == 0 {
		return errors.New("用户ID不能为空")
	}
	return config.DB.Save(&user).Error
}

// DeleteUserByUserId 删除用户（需谨慎，应考虑关联数据）
func DeleteUserByUserId(userId int64) error {
	// 开启事务处理关联数据
	tx := config.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除用户的好友关系
	if err := tx.Where("user_id = ? OR friend_id = ?", userId, userId).
		Delete(&models.Friend{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除用户的聊天记录
	if err := tx.Where("user_id_from = ? OR user_id_to = ?", userId, userId).
		Delete(&models.Chat{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除用户的活动评论
	if err := tx.Where("user_id = ?", userId).
		Delete(&models.ActivityComment{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除用户的活动成员记录
	if err := tx.Where("user_id = ?", userId).
		Delete(&models.ActivityMember{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除用户创建的活动（或者考虑转移所有权）
	if err := tx.Where("user_id = ?", userId).
		Delete(&models.Activity{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 最后删除用户本身
	if err := tx.Delete(&models.User{}, userId).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// CountUserRelations 获取用户的关系统计
func CountUserRelations(userId int64) (map[string]int64, error) {
	result := make(map[string]int64)

	// 统计好友数
	var friendCount int64
	if err := config.DB.Model(&models.Friend{}).
		Where("user_id = ? AND status = ?", userId, 1).
		Count(&friendCount).Error; err != nil {
		return nil, err
	}
	result["friendCount"] = friendCount

	// 统计创建的活动数
	var createdActivityCount int64
	if err := config.DB.Model(&models.Activity{}).
		Where("user_id = ? AND if_delete = 0", userId).
		Count(&createdActivityCount).Error; err != nil {
		return nil, err
	}
	result["createdActivityCount"] = createdActivityCount

	// 统计参与的活动数
	var joinedActivityCount int64
	if err := config.DB.Model(&models.ActivityMember{}).
		Where("user_id = ?", userId).
		Count(&joinedActivityCount).Error; err != nil {
		return nil, err
	}
	result["joinedActivityCount"] = joinedActivityCount

	// 统计未读消息数
	var unreadMessageCount int64
	if err := config.DB.Model(&models.Chat{}).
		Where("user_id_to = ? AND status_to = ?", userId, 2).
		Count(&unreadMessageCount).Error; err != nil {
		return nil, err
	}
	result["unreadMessageCount"] = unreadMessageCount

	return result, nil
}
