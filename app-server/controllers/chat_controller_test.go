package controllers

import (
	"errors"
	"regexp"
	"testing"

	"hobbyhub-server/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAddChat(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	chat := &models.Chat{
		UserIdFrom: 1,
		UserIdTo:   2,
		Content:    "Hello there!",
	}

	// 测试成功添加聊天记录
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `chat`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := AddChat(chat)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试添加聊天记录失败
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("INSERT INTO `chat`")).
		WillReturnError(errors.New("insert error"))
	mock2.ExpectRollback()

	err = AddChat(chat)
	assert.EqualError(t, err, "insert error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestGetChatById(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	chatId := int64(1)
	expectedChat := &models.Chat{
		Id:         chatId,
		UserIdFrom: 1,
		UserIdTo:   2,
		Content:    "Hello there!",
	}

	rows := sqlmock.NewRows([]string{"id", "user_id_from", "user_id_to", "content"}).
		AddRow(expectedChat.Id, expectedChat.UserIdFrom, expectedChat.UserIdTo, expectedChat.Content)

	// 测试成功获取聊天记录
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `chat` WHERE id = ? ORDER BY `chat`.`id` LIMIT ?")).
		WithArgs(chatId, 1).
		WillReturnRows(rows)

	chat, err := GetChatById(chatId)
	assert.NoError(t, err)
	assert.Equal(t, expectedChat.Id, chat.Id)
	assert.Equal(t, expectedChat.UserIdFrom, chat.UserIdFrom)
	assert.Equal(t, expectedChat.UserIdTo, chat.UserIdTo)
	assert.Equal(t, expectedChat.Content, chat.Content)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试找不到聊天记录
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `chat` WHERE id = ? ORDER BY `chat`.`id` LIMIT ?")).
		WithArgs(chatId, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	chat, err = GetChatById(chatId)
	assert.Nil(t, chat)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestGetAllChatByFromUserIdToUserId(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	fromUserId := int64(1)
	toUserId := int64(2)
	expectedChats := []models.Chat{
		{Id: 1, UserIdFrom: fromUserId, UserIdTo: toUserId, Content: "Message 1"},
		{Id: 2, UserIdFrom: fromUserId, UserIdTo: toUserId, Content: "Message 2"},
		{Id: 3, UserIdFrom: fromUserId, UserIdTo: toUserId, Content: "Message 3"},
	}

	rows := sqlmock.NewRows([]string{"id", "user_id_from", "user_id_to", "content"})
	for _, chat := range expectedChats {
		rows.AddRow(chat.Id, chat.UserIdFrom, chat.UserIdTo, chat.Content)
	}

	// 测试成功获取所有聊天记录
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `chat` WHERE user_id_from = ? AND user_id_to = ?")).
		WithArgs(fromUserId, toUserId).
		WillReturnRows(rows)

	chats, err := GetAllChatByFromUserIdToUserId(fromUserId, toUserId)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedChats), len(chats))
	for i, chat := range chats {
		assert.Equal(t, expectedChats[i].Id, chat.Id)
		assert.Equal(t, expectedChats[i].UserIdFrom, chat.UserIdFrom)
		assert.Equal(t, expectedChats[i].UserIdTo, chat.UserIdTo)
		assert.Equal(t, expectedChats[i].Content, chat.Content)
	}
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试获取聊天记录出错
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `chat` WHERE user_id_from = ? AND user_id_to = ?")).
		WithArgs(fromUserId, toUserId).
		WillReturnError(errors.New("query error"))

	chats, err = GetAllChatByFromUserIdToUserId(fromUserId, toUserId)
	assert.Nil(t, chats)
	assert.EqualError(t, err, "query error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestUpdateChat(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	chat := &models.Chat{
		Id:         1,
		UserIdFrom: 1,
		UserIdTo:   2,
		Content:    "Updated message",
	}

	// 测试成功更新聊天记录
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `chat` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := UpdateChat(chat)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试更新聊天记录失败
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("UPDATE `chat` SET")).
		WillReturnError(errors.New("update error"))
	mock2.ExpectRollback()

	err = UpdateChat(chat)
	assert.EqualError(t, err, "update error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestDeleteChatById(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	chatId := int64(1)
	userId := int64(2) // 用户是接收者

	// 原始聊天记录
	chatRows := sqlmock.NewRows([]string{"id", "user_id_from", "user_id_to", "content", "status_from", "status_to"}).
		AddRow(chatId, 1, userId, "Hello", 2, 2)

	// 测试标记删除成功 - 用户是接收者
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `chat` WHERE `chat`.`id` = ? ORDER BY `chat`.`id` LIMIT ?")).
		WithArgs(chatId, 1).
		WillReturnRows(chatRows)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `chat` SET")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 0, chatId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := DeleteChatById(chatId, userId)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试标记删除 - 用户是发送者
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	senderId := int64(1)

	// 原始聊天记录
	chatRows2 := sqlmock.NewRows([]string{"id", "user_id_from", "user_id_to", "content", "status_from", "status_to"}).
		AddRow(chatId, senderId, 2, "Hello", 2, 2)

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `chat` WHERE `chat`.`id` = ? ORDER BY `chat`.`id` LIMIT ?")).
		WithArgs(chatId, 1).
		WillReturnRows(chatRows2)

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("UPDATE `chat` SET")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 0, sqlmock.AnyArg(), chatId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock2.ExpectCommit()

	err = DeleteChatById(chatId, senderId)
	assert.NoError(t, err)
	assert.NoError(t, mock2.ExpectationsWereMet())

	// 测试查询消息失败
	mock3, teardown3 := SetupMockDB(t)
	defer teardown3()

	// 修复：匹配GORM的查询格式并返回错误
	mock3.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `chat` WHERE `chat`.`id` = ? ORDER BY `chat`.`id` LIMIT ?")).
		WithArgs(chatId, 1).
		WillReturnError(errors.New("query error"))

	err = DeleteChatById(chatId, userId)
	assert.EqualError(t, err, "query error")
	assert.NoError(t, mock3.ExpectationsWereMet())
}
