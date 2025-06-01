package controllers

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"hobbyhub-server/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAddFriend(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	friend := &models.Friend{UserId: 1, FriendId: 2, CreateTime: time.Now()}

	// 测试成功场景：两个INSERT都成功
	// 第一次插入原始朋友关系
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `friend`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// 第二次插入反向朋友关系
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `friend`")).
		WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectCommit()

	err := AddFriend(friend)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试第一次插入失败场景
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("INSERT INTO `friend`")).
		WillReturnError(errors.New("insert error"))
	mock2.ExpectRollback()

	err = AddFriend(friend)
	assert.EqualError(t, err, "insert error")
	assert.NoError(t, mock2.ExpectationsWereMet())

	// 测试第二次插入失败场景
	mock3, teardown3 := SetupMockDB(t)
	defer teardown3()

	// 第一次插入成功
	mock3.ExpectBegin()
	mock3.ExpectExec(regexp.QuoteMeta("INSERT INTO `friend`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock3.ExpectCommit()

	// 第二次插入失败
	mock3.ExpectBegin()
	mock3.ExpectExec(regexp.QuoteMeta("INSERT INTO `friend`")).
		WillReturnError(errors.New("second insert error"))
	mock3.ExpectRollback()

	err = AddFriend(friend)
	assert.EqualError(t, err, "second insert error")
	assert.NoError(t, mock3.ExpectationsWereMet())
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
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `friend` WHERE id = ? ORDER BY `friend`.`id` LIMIT ?")).
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

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `friend` WHERE id = ? ORDER BY `friend`.`id` LIMIT ?")).
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
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `friend` WHERE user_id = ?")).
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

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `friend` WHERE user_id = ?")).
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
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `friend` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := UpdateFriend(&friend)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test friend update error
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("UPDATE `friend` SET")).
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
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `friend` WHERE `friend`.`id` = ?")).
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
	mock2.ExpectExec(regexp.QuoteMeta("DELETE FROM `friend` WHERE `friend`.`id` = ?")).
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

	// 创建第一个方向的好友关系（userId -> friendId）
	friend1 := &models.Friend{
		Id:       1,
		UserId:   userId,
		FriendId: friendId,
	}

	// 创建第二个方向的好友关系（friendId -> userId）
	friend2 := &models.Friend{
		Id:       2,
		UserId:   friendId,
		FriendId: userId,
	}

	// 第一个查询的结果行
	rows1 := sqlmock.NewRows([]string{"id", "user_id", "friend_id"}).
		AddRow(friend1.Id, friend1.UserId, friend1.FriendId)

	// 第二个查询的结果行
	rows2 := sqlmock.NewRows([]string{"id", "user_id", "friend_id"}).
		AddRow(friend2.Id, friend2.UserId, friend2.FriendId)

	// 测试成功获取两个方向的好友关系
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `friend` WHERE user_id = ? AND friend_id = ?")).
		WithArgs(userId, friendId).
		WillReturnRows(rows1)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `friend` WHERE user_id = ? AND friend_id = ?")).
		WithArgs(friendId, userId).
		WillReturnRows(rows2)

	friendA, friendB, err := GetFriendByUserIdAndFriendId(userId, friendId)
	assert.NoError(t, err)
	assert.NotNil(t, friendA) // 应该返回第一个好友关系
	assert.NotNil(t, friendB) // 应该返回第二个好友关系

	// 验证第一个好友关系
	assert.Equal(t, friend1.Id, friendA.Id)
	assert.Equal(t, friend1.UserId, friendA.UserId)
	assert.Equal(t, friend1.FriendId, friendA.FriendId)

	// 验证第二个好友关系
	assert.Equal(t, friend2.Id, friendB.Id)
	assert.Equal(t, friend2.UserId, friendB.UserId)
	assert.Equal(t, friend2.FriendId, friendB.FriendId)

	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试查询失败的情况
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `friend` WHERE user_id = ? AND friend_id = ?")).
		WithArgs(userId, friendId).
		WillReturnError(errors.New("query error"))

	friendA, friendB, err = GetFriendByUserIdAndFriendId(userId, friendId)
	assert.Nil(t, friendA)
	assert.Nil(t, friendB)
	assert.EqualError(t, err, "query error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestUpdateFriendSynchronize(t *testing.T) {
	// 测试场景1：成功同步更新两个好友关系
	mock, teardown := SetupMockDB(t)
	defer teardown()

	friend1 := &models.Friend{
		Id:       1,
		UserId:   1,
		FriendId: 2,
		Status:   1, // 假设要更新的状态
	}

	friend2 := &models.Friend{
		Id:       2,
		UserId:   2,
		FriendId: 1,
		Status:   1, // 假设要更新的状态
	}

	// 期望开启事务
	mock.ExpectBegin()

	// 期望更新第一个好友关系成功
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `friend` SET")).
		WithArgs(
			friend1.UserId,
			friend1.FriendId,
			friend1.Status,
			sqlmock.AnyArg(), // create_time
			1,                // id
		).WillReturnResult(sqlmock.NewResult(1, 1))

	// 期望更新第二个好友关系成功
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `friend` SET")).
		WithArgs(
			friend2.UserId,
			friend2.FriendId,
			friend2.Status,
			sqlmock.AnyArg(), // create_time
			2,                // id
		).WillReturnResult(sqlmock.NewResult(2, 1))

	// 期望提交事务
	mock.ExpectCommit()

	// 执行测试
	err := UpdateFriendSynchronize(friend1, friend2)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试场景2：第二次更新失败应回滚
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	// 期望开启事务
	mock2.ExpectBegin()

	// 期望更新第一个好友关系成功
	mock2.ExpectExec(regexp.QuoteMeta("UPDATE `friend` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// 期望更新第二个好友关系失败
	mock2.ExpectExec(regexp.QuoteMeta("UPDATE `friend` SET")).
		WillReturnError(errors.New("update error"))

	// 期望回滚事务
	mock2.ExpectRollback()

	// 执行测试
	err = UpdateFriendSynchronize(friend1, friend2)
	assert.Error(t, err)
	assert.EqualError(t, err, "update error")
	assert.NoError(t, mock2.ExpectationsWereMet())

	// 测试场景3：事务开始失败
	mock3, teardown3 := SetupMockDB(t)
	defer teardown3()

	// 期望开启事务失败
	mock3.ExpectBegin().WillReturnError(errors.New("transaction error"))

	// 执行测试
	err = UpdateFriendSynchronize(friend1, friend2)
	assert.Error(t, err)
	assert.EqualError(t, err, "transaction error")
	assert.NoError(t, mock3.ExpectationsWereMet())
}
