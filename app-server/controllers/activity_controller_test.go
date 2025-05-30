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

// TestAddActivityMember 测试添加活动成员功能
func TestAddActivityMember(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	createTime, err := time.Parse("2006-01-02 15:04:05", "2025-05-30 10:00:00")
	assert.NoError(t, err)
	activityMember := &models.ActivityMember{
		EventID:    1,
		UserID:     2,
		CreateTime: createTime,
	}

	// 测试成功添加活动成员
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `it_event_member`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = AddActivityMember(activityMember)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试添加活动成员失败
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("INSERT INTO `it_event_member`")).
		WillReturnError(errors.New("insert error"))
	mock2.ExpectRollback()

	err = AddActivityMember(activityMember)
	assert.EqualError(t, err, "insert error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

// TestGetActivityMembersByActivityId 测试根据活动ID获取成员列表
func TestGetActivityMembersByActivityId(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	activityId := int64(1)
	layout := "2006-01-02 15:04:05"
	createTime1, _ := time.Parse(layout, "2025-05-30 10:00:00")
	createTime2, _ := time.Parse(layout, "2025-05-30 11:00:00")
	createTime3, _ := time.Parse(layout, "2025-05-30 12:00:00")
	expectedMembers := []models.ActivityMember{
		{ID: 1, EventID: activityId, UserID: 101, CreateTime: createTime1},
		{ID: 2, EventID: activityId, UserID: 102, CreateTime: createTime2},
		{ID: 3, EventID: activityId, UserID: 103, CreateTime: createTime3},
	}

	rows := sqlmock.NewRows([]string{"id", "event_id", "user_id", "create_time"})
	for _, m := range expectedMembers {
		rows.AddRow(m.ID, m.EventID, m.UserID, m.CreateTime)
	}

	// 测试成功获取活动成员
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `it_event_member` WHERE activity_id = ?")).
		WithArgs(activityId).
		WillReturnRows(rows)

	members, err := GetActivityMembersByActivityId(activityId)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedMembers), len(members))
	for i, m := range members {
		assert.Equal(t, expectedMembers[i].ID, m.ID)
		assert.Equal(t, expectedMembers[i].EventID, m.EventID)
		assert.Equal(t, expectedMembers[i].UserID, m.UserID)
	}
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试查询错误
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `it_event_member` WHERE activity_id = ?")).
		WithArgs(activityId).
		WillReturnError(errors.New("query error"))

	members, err = GetActivityMembersByActivityId(activityId)
	assert.Nil(t, members)
	assert.EqualError(t, err, "query error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

// TestGetActivityMembersByUserId 测试根据用户ID获取活动成员列表
func TestGetActivityMembersByUserId(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	userId := int64(1)
	layout := "2006-01-02 15:04:05"
	createTime1, _ := time.Parse(layout, "2025-05-30 10:00:00")
	createTime2, _ := time.Parse(layout, "2025-05-30 11:00:00")
	createTime3, _ := time.Parse(layout, "2025-05-30 12:00:00")
	expectedMembers := []models.ActivityMember{
		{ID: 1, EventID: 101, UserID: userId, CreateTime: createTime1},
		{ID: 2, EventID: 102, UserID: userId, CreateTime: createTime2},
		{ID: 3, EventID: 103, UserID: userId, CreateTime: createTime3},
	}

	rows := sqlmock.NewRows([]string{"id", "event_id", "user_id", "create_time"})
	for _, m := range expectedMembers {
		rows.AddRow(m.ID, m.EventID, m.UserID, m.CreateTime)
	}

	// 测试成功获取用户参与的活动
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `it_event_member` WHERE user_id = ?")).
		WithArgs(userId).
		WillReturnRows(rows)

	members, err := GetActivityMembersByUserId(userId)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedMembers), len(members))
	for i, m := range members {
		assert.Equal(t, expectedMembers[i].ID, m.ID)
		assert.Equal(t, expectedMembers[i].EventID, m.EventID)
		assert.Equal(t, expectedMembers[i].UserID, m.UserID)
	}
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试查询错误
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `it_event_member` WHERE user_id = ?")).
		WithArgs(userId).
		WillReturnError(errors.New("query error"))

	members, err = GetActivityMembersByUserId(userId)
	assert.Nil(t, members)
	assert.EqualError(t, err, "query error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

// TestUpdateActivityMember 测试更新活动成员
func TestUpdateActivityMember(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	layout := "2006-01-02 15:04:05"
	createTime, err := time.Parse(layout, "2025-05-30 15:00:00")
	assert.NoError(t, err)

	activityMember := &models.ActivityMember{
		ID:         1,
		EventID:    1,
		UserID:     2,
		CreateTime: createTime,
	}

	// 测试成功更新活动成员
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `it_event_member` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = UpdateActivityMember(activityMember)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试更新活动成员失败
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("UPDATE `it_event_member` SET")).
		WillReturnError(errors.New("update error"))
	mock2.ExpectRollback()

	err = UpdateActivityMember(activityMember)
	assert.EqualError(t, err, "update error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

// TestDeleteActivityMember 测试删除活动成员
func TestDeleteActivityMember(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	activityId := int64(1)
	userId := int64(2)

	// 测试成功删除活动成员
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `it_event_member` WHERE activity_id = ? AND user_id = ?")).
		WithArgs(activityId, userId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := DeleteActivityMember(activityId, userId)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试删除活动成员失败
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("DELETE FROM `it_event_member` WHERE activity_id = ? AND user_id = ?")).
		WithArgs(activityId, userId).
		WillReturnError(errors.New("delete error"))
	mock2.ExpectRollback()

	err = DeleteActivityMember(activityId, userId)
	assert.EqualError(t, err, "delete error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

// TestAddActivityComment 测试添加活动评论
func TestAddActivityComment(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	now := time.Now()
	activityComment := &models.ActivityComment{
		EventID:    1,
		UserID:     2,
		Content:    "这是一条测试评论",
		CreateTime: now,
	}

	// 测试成功添加活动评论
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `it_event_comment`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := AddActivityComment(activityComment)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试添加活动评论失败
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("INSERT INTO `it_event_comment`")).
		WillReturnError(errors.New("insert error"))
	mock2.ExpectRollback()

	err = AddActivityComment(activityComment)
	assert.EqualError(t, err, "insert error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

// TestGetActivityCommentsByActivityId 测试通过活动ID获取评论
func TestGetActivityCommentsByActivityId(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	activityId := int64(1)
	now := time.Now()
	expectedComments := []models.ActivityComment{
		{ID: 1, EventID: activityId, UserID: 101, Content: "评论1", CreateTime: now.Add(-time.Hour * 2)},
		{ID: 2, EventID: activityId, UserID: 102, Content: "评论2", CreateTime: now.Add(-time.Hour)},
		{ID: 3, EventID: activityId, UserID: 103, Content: "评论3", CreateTime: now},
	}

	rows := sqlmock.NewRows([]string{"id", "event_id", "user_id", "content", "create_time"})
	for _, c := range expectedComments {
		rows.AddRow(c.ID, c.EventID, c.UserID, c.Content, c.CreateTime)
	}

	// 测试成功获取评论
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `it_event_comment` WHERE event_id = ?")).
		WithArgs(activityId).
		WillReturnRows(rows)

	comments, err := GetActivityCommentsByActivityId(activityId)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedComments), len(comments))
	for i, c := range comments {
		assert.Equal(t, expectedComments[i].ID, c.ID)
		assert.Equal(t, expectedComments[i].EventID, c.EventID)
		assert.Equal(t, expectedComments[i].UserID, c.UserID)
		assert.Equal(t, expectedComments[i].Content, c.Content)
	}
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试查询错误
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `it_event_comment` WHERE event_id = ?")).
		WithArgs(activityId).
		WillReturnError(errors.New("query error"))

	comments, err = GetActivityCommentsByActivityId(activityId)
	assert.Nil(t, comments)
	assert.EqualError(t, err, "query error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

// TestGetActivityCommentsByUserId 测试通过用户ID获取评论
func TestGetActivityCommentsByUserId(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	userId := int64(1)
	now := time.Now()
	expectedComments := []models.ActivityComment{
		{ID: 1, EventID: 101, UserID: userId, Content: "评论1", CreateTime: now.Add(-time.Hour * 2)},
		{ID: 2, EventID: 102, UserID: userId, Content: "评论2", CreateTime: now.Add(-time.Hour)},
		{ID: 3, EventID: 103, UserID: userId, Content: "评论3", CreateTime: now},
	}

	rows := sqlmock.NewRows([]string{"id", "event_id", "user_id", "content", "create_time"})
	for _, c := range expectedComments {
		rows.AddRow(c.ID, c.EventID, c.UserID, c.Content, c.CreateTime)
	}

	// 测试成功获取评论
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `it_event_comment` WHERE user_id = ?")).
		WithArgs(userId).
		WillReturnRows(rows)

	comments, err := GetActivityCommentsByUserId(userId)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedComments), len(comments))
	for i, c := range comments {
		assert.Equal(t, expectedComments[i].ID, c.ID)
		assert.Equal(t, expectedComments[i].EventID, c.EventID)
		assert.Equal(t, expectedComments[i].UserID, c.UserID)
		assert.Equal(t, expectedComments[i].Content, c.Content)
	}
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试查询错误
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `it_event_comment` WHERE user_id = ?")).
		WithArgs(userId).
		WillReturnError(errors.New("query error"))

	comments, err = GetActivityCommentsByUserId(userId)
	assert.Nil(t, comments)
	assert.EqualError(t, err, "query error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

// TestGetActivityCommentsByActivityIdAndUserId 测试通过活动ID和用户ID获取评论
func TestGetActivityCommentsByActivityIdAndUserId(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	activityId := int64(1)
	userId := int64(2)
	now := time.Now()
	expectedComments := []models.ActivityComment{
		{ID: 1, EventID: activityId, UserID: userId, Content: "评论1", CreateTime: now.Add(-time.Hour)},
		{ID: 2, EventID: activityId, UserID: userId, Content: "评论2", CreateTime: now},
	}

	rows := sqlmock.NewRows([]string{"id", "event_id", "user_id", "content", "create_time"})
	for _, c := range expectedComments {
		rows.AddRow(c.ID, c.EventID, c.UserID, c.Content, c.CreateTime)
	}

	// 测试成功获取评论
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `it_event_comment` WHERE event_id = ? AND user_id = ?")).
		WithArgs(activityId, userId).
		WillReturnRows(rows)

	comments, err := GetActivityCommentsByActivityIdAndUserId(activityId, userId)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedComments), len(comments))
	for i, c := range comments {
		assert.Equal(t, expectedComments[i].ID, c.ID)
		assert.Equal(t, expectedComments[i].EventID, c.EventID)
		assert.Equal(t, expectedComments[i].UserID, c.UserID)
		assert.Equal(t, expectedComments[i].Content, c.Content)
	}
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试查询错误
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `it_event_comment` WHERE event_id = ? AND user_id = ?")).
		WithArgs(activityId, userId).
		WillReturnError(errors.New("query error"))

	comments, err = GetActivityCommentsByActivityIdAndUserId(activityId, userId)
	assert.Nil(t, comments)
	assert.EqualError(t, err, "query error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

// TestUpdateActivityComment 测试更新活动评论
func TestUpdateActivityComment(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	now := time.Now()
	activityComment := &models.ActivityComment{
		ID:         1,
		EventID:    1,
		UserID:     2,
		Content:    "这是更新后的评论",
		CreateTime: now,
	}

	// 测试成功更新评论
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `it_event_comment` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := UpdateActivityComment(activityComment)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试更新评论失败
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("UPDATE `it_event_comment` SET")).
		WillReturnError(errors.New("update error"))
	mock2.ExpectRollback()

	err = UpdateActivityComment(activityComment)
	assert.EqualError(t, err, "update error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

// TestDeleteActivityComment 测试删除活动评论
func TestDeleteActivityComment(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	commentId := int64(1)

	// 测试成功删除评论
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `it_event_comment` WHERE `it_event_comment`.`id` = ?")).
		WithArgs(commentId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := DeleteActivityComment(commentId)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试删除评论失败
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("DELETE FROM `it_event_comment` WHERE `it_event_comment`.`id` = ?")).
		WithArgs(commentId).
		WillReturnError(errors.New("delete error"))
	mock2.ExpectRollback()

	err = DeleteActivityComment(commentId)
	assert.EqualError(t, err, "delete error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}
