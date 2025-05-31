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

func TestAddFriend(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	friend := &models.Friend{UserId: 1, FriendId: 2}

	// Test successful friend creation
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `friends`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := AddFriend(friend)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test friend creation failure
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("INSERT INTO `friends`")).
		WillReturnError(errors.New("insert error"))
	mock2.ExpectRollback()

	err = AddFriend(friend)
	assert.EqualError(t, err, "insert error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestGetFriendById(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	friendId := int64(1)
	expectedFriend := &models.Friend{
		Id:       friendId,
		UserId:   1,
		FriendId: 2,
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "friend_id"}).
		AddRow(expectedFriend.Id, expectedFriend.UserId, expectedFriend.FriendId)

	// Test successful friend retrieval
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `friends` WHERE id = ? ORDER BY `friends`.`id` LIMIT ?")).
		WithArgs(friendId, 1).
		WillReturnRows(rows)

	friend, err := GetFriendById(friendId)
	assert.NoError(t, err)
	assert.Equal(t, expectedFriend.Id, friend.Id)
	assert.Equal(t, expectedFriend.UserId, friend.UserId)
	assert.Equal(t, expectedFriend.FriendId, friend.FriendId)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test friend not found
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `friends` WHERE id = ? ORDER BY `friends`.`id` LIMIT ?")).
		WithArgs(friendId, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	friend, err = GetFriendById(friendId)
	assert.Nil(t, friend)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestGetAllFriendsByUserId(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	userId := int64(1)
	expectedFriends := []models.Friend{
		{Id: 1, UserId: userId, FriendId: 10},
		{Id: 2, UserId: userId, FriendId: 11},
		{Id: 3, UserId: userId, FriendId: 12},
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "friend_id"})
	for _, f := range expectedFriends {
		rows.AddRow(f.Id, f.UserId, f.FriendId)
	}

	// 测试成功获取所有好友
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `friends` WHERE user_id = ?")).
		WithArgs(userId).
		WillReturnRows(rows)

	friends, err := GetAllFriendsByUserId(userId)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedFriends), len(friends))
	for i, f := range friends {
		assert.Equal(t, expectedFriends[i].Id, f.Id)
		assert.Equal(t, expectedFriends[i].UserId, f.UserId)
		assert.Equal(t, expectedFriends[i].FriendId, f.FriendId)
	}
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试查询错误情况
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `friends` WHERE user_id = ?")).
		WithArgs(userId).
		WillReturnError(errors.New("query error"))

	friends, err = GetAllFriendsByUserId(userId)
	assert.Nil(t, friends)
	assert.EqualError(t, err, "query error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestUpdateFriend(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	friend := models.Friend{Id: 1, UserId: 1, FriendId: 3}

	// Test successful friend update
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `friends` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := UpdateFriend(&friend)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test friend update error
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("UPDATE `friends` SET")).
		WillReturnError(errors.New("update error"))
	mock2.ExpectRollback()

	err = UpdateFriend(&friend)
	assert.EqualError(t, err, "update error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestDeleteFriendById(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	friendId := int64(1)

	// Test successful friend deletion
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `friends` WHERE `friends`.`id` = ?")).
		WithArgs(friendId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := DeleteFriendById(friendId)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test friend deletion error
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("DELETE FROM `friends` WHERE `friends`.`id` = ?")).
		WithArgs(friendId).
		WillReturnError(errors.New("delete error"))
	mock2.ExpectRollback()

	err = DeleteFriendById(friendId)
	assert.EqualError(t, err, "delete error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestGetFriendByUserIdAndFriendId(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	userId := int64(1)
	friendId := int64(2)
	expectedFriend := &models.Friend{
		Id:       1,
		UserId:   userId,
		FriendId: friendId,
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "friend_id"}).
		AddRow(expectedFriend.Id, expectedFriend.UserId, expectedFriend.FriendId)

	// Test successful friend retrieval
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `friends` WHERE user_id = ? AND friend_id = ? ORDER BY `friends`.`id` LIMIT ?")).
		WithArgs(userId, friendId, 1).
		WillReturnRows(rows)

	friend, err := GetFriendByUserIdAndFriendId(userId, friendId)
	assert.NoError(t, err)
	assert.Equal(t, expectedFriend.Id, friend.Id)
	assert.Equal(t, expectedFriend.UserId, friend.UserId)
	assert.Equal(t, expectedFriend.FriendId, friend.FriendId)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test friend not found
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `friends` WHERE user_id = ? AND friend_id = ? ORDER BY `friends`.`id` LIMIT ?")).
		WithArgs(userId, friendId, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	friend, err = GetFriendByUserIdAndFriendId(userId, friendId)
	assert.Nil(t, friend)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoError(t, mock2.ExpectationsWereMet())
}
