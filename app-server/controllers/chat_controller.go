package controllers

import (
	"hobbyhub-server/config"
	"hobbyhub-server/models"
)

// AddChat 添加聊天记录
func AddChat(chat *models.Chat) error {
	return config.DB.Create(chat).Error
}

// GetChatById 通过ID获取聊天记录，并预加载用户信息
func GetChatById(chatId int64) (*models.Chat, error) {
	var chat models.Chat
	if err := config.DB.Where("id = ?", chatId).
		Preload("UserFrom").
		Preload("UserTo").
		First(&chat).Error; err != nil {
		return nil, err
	}
	return &chat, nil
}
func GetAllChatByFromUserIdToUserId(fromUserId int64, toUserId int64) ([]models.Chat, error) {
	var chats []models.Chat
	if err := config.DB.Where("user_id_from = ? AND user_id_to = ?", fromUserId, toUserId).
		Preload("UserFrom").
		Preload("UserTo").
		Order("create_time DESC").
		Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}

// UpdateChat 更新聊天记录
func UpdateChat(chat *models.Chat) error {
	return config.DB.Save(chat).Error
}

// DeleteChatById 删除聊天记录（实际上是标记为删除状态）
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
		return config.DB.Error // 用户不是发送方或接收方，无法删除
	}

	return config.DB.Save(&chat).Error
}
