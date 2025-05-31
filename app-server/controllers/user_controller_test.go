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
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `it_user`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := AddUser(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()
	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("INSERT INTO `it_user`")).
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

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `it_user` WHERE id = ? ORDER BY `it_user`.`id` LIMIT ?")).
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

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `it_user` WHERE id = ? ORDER BY `it_user`.`id` LIMIT ?")).
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

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `it_user` WHERE username = ? ORDER BY `it_user`.`id` LIMIT ?")).
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

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `it_user` WHERE username = ? ORDER BY `it_user`.`id` LIMIT ?")).
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
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `it_user` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := UpdateUser(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test update error
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("UPDATE `it_user` SET")).
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

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `it_user` WHERE `it_user`.`id` = ?")).
		WithArgs(userId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := DeleteUserByUserId(userId)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test delete error
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("DELETE FROM `it_user` WHERE `it_user`.`id` = ?")).
		WithArgs(userId).
		WillReturnError(errors.New("delete error"))
	mock2.ExpectRollback()

	err = DeleteUserByUserId(userId)
	assert.EqualError(t, err, "delete error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}
