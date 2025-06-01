package controllers

import (
	"errors"
	"regexp"
	"testing"

	"hobbyhub-server/config"
	"hobbyhub-server/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupMockDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)
	origin := config.DB
	config.DB = gdb
	return mock, func() {
		config.DB = origin
		db.Close()
	}
}

func TestAddUser(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	user := &models.User{Username: "testuser"}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := AddUser(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()
	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("INSERT INTO `user`")).
		WillReturnError(errors.New("insert error"))
	mock2.ExpectRollback()

	err = AddUser(user)
	assert.EqualError(t, err, "insert error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestGetUserByUserId(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	userId := int64(1)
	expectedUser := &models.User{Id: userId, Username: "testuser"}

	rows := sqlmock.NewRows([]string{"id", "username"}).
		AddRow(expectedUser.Id, expectedUser.Username)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user` WHERE id = ? ORDER BY `user`.`id` LIMIT ?")).
		WithArgs(userId, 1).
		WillReturnRows(rows)
	user, err := GetUserByUserId(userId)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Id, user.Id)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test not found
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user` WHERE id = ? ORDER BY `user`.`id` LIMIT ?")).
		WithArgs(userId, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err = GetUserByUserId(userId)
	assert.Nil(t, user)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestGetUserByUserName(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	userName := "testuser"
	expectedUser := &models.User{Id: 1, Username: userName}

	rows := sqlmock.NewRows([]string{"id", "username"}).
		AddRow(expectedUser.Id, expectedUser.Username)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user` WHERE username = ? ORDER BY `user`.`id` LIMIT ?")).
		WithArgs(userName, 1).
		WillReturnRows(rows)

	user, err := GetUserByUserName(userName)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Id, user.Id)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test not found
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user` WHERE username = ? ORDER BY `user`.`id` LIMIT ?")).
		WithArgs(userName, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err = GetUserByUserName(userName)
	assert.Nil(t, user)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestUpdateUser(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	user := models.User{Id: 1, Username: "updateduser"}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `user` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := UpdateUser(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test update error
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("UPDATE `user` SET")).
		WillReturnError(errors.New("update error"))
	mock2.ExpectRollback()

	err = UpdateUser(user)
	assert.EqualError(t, err, "update error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestDeleteUserByUserId(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	userId := int64(1)

	// 事务开始
	mock.ExpectBegin()

	// 删除用户的好友关系
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `friend` WHERE user_id = ? OR friend_id = ?")).
		WithArgs(userId, userId).
		WillReturnResult(sqlmock.NewResult(1, 2)) // 假设删除了两条好友记录

	// 删除用户的聊天记录
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `chat` WHERE user_id_from = ? OR user_id_to = ?")).
		WithArgs(userId, userId).
		WillReturnResult(sqlmock.NewResult(1, 3)) // 假设删除了三条聊天记录

	// 删除用户的活动评论
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `activity_comment` WHERE user_id = ?")).
		WithArgs(userId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// 删除用户的活动成员记录
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `activity_member` WHERE user_id = ?")).
		WithArgs(userId).
		WillReturnResult(sqlmock.NewResult(1, 2))

	// 删除用户创建的活动
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `activity` WHERE user_id = ?")).
		WithArgs(userId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// 删除用户本身
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `user` WHERE `user`.`id` = ?")).
		WithArgs(userId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// 提交事务
	mock.ExpectCommit()

	err := DeleteUserByUserId(userId)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试删除错误场景 - 在删除好友记录时失败
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("DELETE FROM `friend` WHERE user_id = ? OR friend_id = ?")).
		WithArgs(userId, userId).
		WillReturnError(errors.New("delete friend error"))
	mock2.ExpectRollback()

	err = DeleteUserByUserId(userId)
	assert.EqualError(t, err, "delete friend error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestCountUserRelations(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	userId := int64(1)

	// 模拟统计好友数查询
	friendCountRows := sqlmock.NewRows([]string{"count"}).AddRow(5)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `friend` WHERE user_id = ? AND status = ?")).
		WithArgs(userId, 1).
		WillReturnRows(friendCountRows)

	// 模拟统计创建的活动数查询 - 修改这里匹配硬编码的if_delete = 0
	createdActivityCountRows := sqlmock.NewRows([]string{"count"}).AddRow(3)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `activity` WHERE user_id = ? AND if_delete = 0")).
		WithArgs(userId).
		WillReturnRows(createdActivityCountRows)

	// 模拟统计参与的活动数查询
	joinedActivityCountRows := sqlmock.NewRows([]string{"count"}).AddRow(7)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `activity_member` WHERE user_id = ?")).
		WithArgs(userId).
		WillReturnRows(joinedActivityCountRows)

	// 模拟统计未读消息数查询
	unreadMessageCountRows := sqlmock.NewRows([]string{"count"}).AddRow(12)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `chat` WHERE user_id_to = ? AND status_to = ?")).
		WithArgs(userId, 2).
		WillReturnRows(unreadMessageCountRows)

	relations, err := CountUserRelations(userId)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), relations["friendCount"])
	assert.Equal(t, int64(3), relations["createdActivityCount"])
	assert.Equal(t, int64(7), relations["joinedActivityCount"])
	assert.Equal(t, int64(12), relations["unreadMessageCount"])
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试错误场景 - 在统计好友数时失败
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `friend` WHERE user_id = ? AND status = ?")).
		WithArgs(userId, 1).
		WillReturnError(errors.New("count friend error"))

	relations, err = CountUserRelations(userId)
	assert.EqualError(t, err, "count friend error")
	assert.Nil(t, relations)
	assert.NoError(t, mock2.ExpectationsWereMet())
}
