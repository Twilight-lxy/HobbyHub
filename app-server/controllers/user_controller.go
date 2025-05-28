package controllers

import (
	"hobbyhub-server/config"
	"hobbyhub-server/models"
)

func AddUser(user *models.User) error {
	if err := config.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserByUserId(userId int64) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByUserName(userName string) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("username = ?", userName).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(user models.User) error {
	if err := config.DB.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func DeleteUserByUserId(userId int64) error {
	if err := config.DB.Delete(&models.User{}, userId).Error; err != nil {
		return err
	}
	return nil
}
