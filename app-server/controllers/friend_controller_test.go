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

	friend := &models.Friend{UserID: 1, FriendID: 2}

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
		ID:       friendId,
		UserID:   1,
		FriendID: 2,
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "friend_id"}).
		AddRow(expectedFriend.ID, expectedFriend.UserID, expectedFriend.FriendID)

	// Test successful friend retrieval
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `friends` WHERE id = ? ORDER BY `friends`.`id` LIMIT ?")).
		WithArgs(friendId, 1).
		WillReturnRows(rows)

	friend, err := GetFriendById(friendId)
	assert.NoError(t, err)
	assert.Equal(t, expectedFriend.ID, friend.ID)
	assert.Equal(t, expectedFriend.UserID, friend.UserID)
	assert.Equal(t, expectedFriend.FriendID, friend.FriendID)
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

func TestGetAllFriendsIdByUserID(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	userId := int64(1)
	expectedIds := []int64{1, 2, 3}

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(expectedIds[0]).
		AddRow(expectedIds[1]).
		AddRow(expectedIds[2])

	// Test successful retrieval of all friend IDs
	mock.ExpectQuery(regexp.QuoteMeta("SELECT `id` FROM `friends` WHERE user_id = ?")).
		WithArgs(userId).
		WillReturnRows(rows)

	friendIds, err := GetAllFriendsIdByUserId(userId)
	assert.NoError(t, err)
	assert.Equal(t, expectedIds, friendIds)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test error case
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT `id` FROM `friends` WHERE user_id = ?")).
		WithArgs(userId).
		WillReturnError(errors.New("query error"))

	friendIds, err = GetAllFriendsIdByUserId(userId)
	assert.Nil(t, friendIds)
	assert.Error(t, err)
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestUpdateFriend(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	friend := models.Friend{ID: 1, UserID: 1, FriendID: 3}

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

func TestGetFriendByUserIDAndFriendID(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	userId := int64(1)
	friendId := int64(2)
	expectedFriend := &models.Friend{
		ID:       1,
		UserID:   userId,
		FriendID: friendId,
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "friend_id"}).
		AddRow(expectedFriend.ID, expectedFriend.UserID, expectedFriend.FriendID)

	// Test successful friend retrieval
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `friends` WHERE user_id = ? AND friend_id = ? ORDER BY `friends`.`id` LIMIT ?")).
		WithArgs(userId, friendId, 1).
		WillReturnRows(rows)

	friend, err := GetFriendByUserIdAndFriendId(userId, friendId)
	assert.NoError(t, err)
	assert.Equal(t, expectedFriend.ID, friend.ID)
	assert.Equal(t, expectedFriend.UserID, friend.UserID)
	assert.Equal(t, expectedFriend.FriendID, friend.FriendID)
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
