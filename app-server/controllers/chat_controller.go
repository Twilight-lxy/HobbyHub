package controllers

import (
	"errors"
	"hobbyhub-server/config"
	"hobbyhub-server/models"
)

// AddFriend 添加聊天记录
func AddChat(chat *models.Chat) error {
	return config.DB.Create(chat).Error
}

// GetChatHistory 获取聊天记录
func GetChatById(chatId int64) (*models.Chat, error) {
	var chat models.Chat
	if err := config.DB.Where("id = ?", chatId).First(&chat).Error; err != nil {
		return nil, err
	}
	return &chat, nil
}

// GetAllChatByFromUserIdToUserId 获取指定用户之间的所有聊天记录
func GetAllChatByFromUserIdToUserId(fromUserId, toUserId int64) ([]models.Chat, error) {
	var chats []models.Chat
	if err := config.DB.Where("user_id_from = ? AND user_id_to = ?", fromUserId, toUserId).Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}

// UpdateChat 更新聊天记录
func UpdateChat(chat *models.Chat) error {
	return config.DB.Save(chat).Error
}

func DeleteChatById(chatId int64, userId int64) error {
	var chat models.Chat
	if err := config.DB.First(&chat, chatId).Error; err != nil {
		return err
	}

	// 根据用户是发送方还是接收方，更新相应的状态
	if chat.UserIdFrom == userId {
		chat.StatusFrom = 0 // 删除
	} else if chat.UserIdTo == userId {
		chat.StatusTo = 0 // 删除
	} else {
		return errors.New("用户与聊天Id无法匹配") // 用户不是发送方或接收方，无法删除
	}

	return config.DB.Save(&chat).Error
}
