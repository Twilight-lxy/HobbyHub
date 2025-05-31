package controllers

import (
	"hobbyhub-server/config"
	"hobbyhub-server/models"
)

func AddChat(chat *models.Chat) error {
	if err := config.DB.Create(chat).Error; err != nil {
		return err
	}
	return nil
}

func GetChatById(chatId int64) (*models.Chat, error) {
	var chat models.Chat
	if err := config.DB.Where("id = ?", chatId).First(&chat).Error; err != nil {
		return nil, err
	}
	return &chat, nil
}

func GetAllChatByFromUserIdToUserId(fromUserId, toUserId int64) ([]models.Chat, error) {
	var chats []models.Chat
	if err := config.DB.Where("user_id_from = ? AND user_id_to = ?", fromUserId, toUserId).Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}

func UpdateChat(chat *models.Chat) error {
	if err := config.DB.Save(chat).Error; err != nil {
		return err
	}
	return nil
}

func DeleteChatById(chatId int64) error {
	if err := config.DB.Delete(&models.Chat{}, chatId).Error; err != nil {
		return err
	}
	return nil
}
