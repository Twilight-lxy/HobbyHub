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

func TestAddActivity(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	now := time.Now()
	activity := &models.Activity{
		Name:       "测试活动",
		Intro:      "这是一个测试活动",
		Addr:       "测试地点",
		HeadImg:    "test.jpg",
		UserID:     1,
		CreateTime: now,
		UpdateTime: now,
		StartTime:  now.Add(time.Hour),
		State:      1,
		IfDelete:   0,
		Lat:        23.456,
		Lon:        113.567,
	}

	// 测试成功添加活动
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `it_event`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := AddActivity(activity)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试添加活动失败
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("INSERT INTO `it_event`")).
		WillReturnError(errors.New("insert error"))
	mock2.ExpectRollback()

	err = AddActivity(activity)
	assert.EqualError(t, err, "insert error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestGetActivityById(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	activityId := int64(1)
	now := time.Now()
	expectedActivity := &models.Activity{
		ID:         activityId,
		Name:       "测试活动",
		Intro:      "这是一个测试活动",
		Addr:       "测试地点",
		HeadImg:    "test.jpg",
		UserID:     1,
		CreateTime: now,
		UpdateTime: now,
		StartTime:  now.Add(time.Hour),
		State:      1,
		IfDelete:   0,
		Lat:        23.456,
		Lon:        113.567,
	}

	rows := sqlmock.NewRows([]string{
		"id", "name", "addr", "intro", "head_img", "user_id", "create_time",
		"update_time", "start_time", "state", "if_delete", "lat", "lon",
	}).AddRow(
		expectedActivity.ID, expectedActivity.Name, expectedActivity.Addr,
		expectedActivity.Intro, expectedActivity.HeadImg, expectedActivity.UserID,
		expectedActivity.CreateTime, expectedActivity.UpdateTime, expectedActivity.StartTime,
		expectedActivity.State, expectedActivity.IfDelete, expectedActivity.Lat, expectedActivity.Lon,
	)

	// 测试成功获取活动
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `it_event` WHERE id = ? ORDER BY `it_event`.`id` LIMIT ?")).
		WithArgs(activityId, 1).
		WillReturnRows(rows)

	activity, err := GetActivityById(activityId)
	assert.NoError(t, err)
	assert.Equal(t, expectedActivity.ID, activity.ID)
	assert.Equal(t, expectedActivity.Name, activity.Name)
	assert.Equal(t, expectedActivity.Intro, activity.Intro)
	assert.Equal(t, expectedActivity.Addr, activity.Addr)
	assert.Equal(t, expectedActivity.UserID, activity.UserID)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试找不到活动
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `it_event` WHERE id = ? ORDER BY `it_event`.`id` LIMIT ?")).
		WithArgs(activityId, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	activity, err = GetActivityById(activityId)
	assert.Nil(t, activity)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestUpdateActivity(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	now := time.Now()
	activity := &models.Activity{
		ID:         1,
		Name:       "更新后的活动",
		Intro:      "这是更新后的活动描述",
		Addr:       "更新后的地点",
		HeadImg:    "updated.jpg",
		UserID:     1,
		CreateTime: now.Add(-time.Hour),
		UpdateTime: now,
		StartTime:  now.Add(time.Hour * 3),
		State:      2,
		IfDelete:   0,
		Lat:        24.456,
		Lon:        114.567,
	}

	// 测试成功更新活动
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `it_event` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := UpdateActivity(activity)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试更新活动失败
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("UPDATE `it_event` SET")).
		WillReturnError(errors.New("update error"))
	mock2.ExpectRollback()

	err = UpdateActivity(activity)
	assert.EqualError(t, err, "update error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestDeleteActivityById(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	activityId := int64(1)

	// 测试成功删除活动
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `it_event` WHERE `it_event`.`id` = ?")).
		WithArgs(activityId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := DeleteActivityById(activityId)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试删除活动失败
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("DELETE FROM `it_event` WHERE `it_event`.`id` = ?")).
		WithArgs(activityId).
		WillReturnError(errors.New("delete error"))
	mock2.ExpectRollback()

	err = DeleteActivityById(activityId)
	assert.EqualError(t, err, "delete error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}
