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

func DeleteUserByUserId(userId int64) error {
	if err := config.DB.Delete(&models.User{}, userId).Error; err != nil {
		return err
	}
	return nil
}
