package controllers

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"hobbyhub-server/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
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
		UserId:     1,
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
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `activity`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := AddActivity(activity)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试添加活动失败
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("INSERT INTO `activity`")).
		WillReturnError(errors.New("insert error"))
	mock2.ExpectRollback()

	err = AddActivity(activity)
	assert.EqualError(t, err, "insert error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestUpdateActivity(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	now := time.Now()
	activity := &models.Activity{
		Id:         1,
		Name:       "更新后的活动",
		Intro:      "这是更新后的活动描述",
		Addr:       "更新后的地点",
		HeadImg:    "updated.jpg",
		UserId:     1,
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
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `activity` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := UpdateActivity(activity)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试更新活动失败
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("UPDATE `activity` SET")).
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

	// 测试成功软删除活动（设置if_delete = 1）
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `activity` SET `if_delete`=? WHERE id = ?")).
		WithArgs(1, activityId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := DeleteActivityById(activityId)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试删除活动失败
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("UPDATE `activity` SET `if_delete`=? WHERE id = ?")).
		WithArgs(1, activityId).
		WillReturnError(errors.New("update error"))
	mock2.ExpectRollback()

	err = DeleteActivityById(activityId)
	assert.EqualError(t, err, "update error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

// TestAddActivityMember 测试添加活动成员功能
func TestAddActivityMember(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	createTime, err := time.Parse("2006-01-02 15:04:05", "2025-05-30 10:00:00")
	assert.NoError(t, err)
	activityMember := &models.ActivityMember{
		EventId:    1,
		UserId:     2,
		CreateTime: createTime,
	}

	// 测试成功添加活动成员
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `activity_member`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = AddActivityMember(activityMember)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试添加活动成员失败
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("INSERT INTO `activity_member`")).
		WillReturnError(errors.New("insert error"))
	mock2.ExpectRollback()

	err = AddActivityMember(activityMember)
	assert.EqualError(t, err, "insert error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

// TestGetActivityMembersByActivityId 测试根据活动Id获取成员列表
func TestGetActivityMembersByActivityId(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	activityId := int64(1)
	layout := "2006-01-02 15:04:05"
	createTime1, _ := time.Parse(layout, "2025-05-30 10:00:00")
	createTime2, _ := time.Parse(layout, "2025-05-30 11:00:00")
	createTime3, _ := time.Parse(layout, "2025-05-30 12:00:00")
	expectedMembers := []models.ActivityMember{
		{Id: 1, EventId: activityId, UserId: 101, CreateTime: createTime1},
		{Id: 2, EventId: activityId, UserId: 102, CreateTime: createTime2},
		{Id: 3, EventId: activityId, UserId: 103, CreateTime: createTime3},
	}

	// 创建活动成员查询结果
	rows := sqlmock.NewRows([]string{"id", "event_id", "user_id", "create_time"})
	for _, m := range expectedMembers {
		rows.AddRow(m.Id, m.EventId, m.UserId, m.CreateTime)
	}

	// 创建用户查询结果
	userRows := sqlmock.NewRows([]string{"id", "username", "name"}).
		AddRow(101, "user101", "User 101").
		AddRow(102, "user102", "User 102").
		AddRow(103, "user103", "User 103")

	// 测试成功获取活动成员并预加载用户
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `activity_member` WHERE event_id = ?")).
		WithArgs(activityId).
		WillReturnRows(rows)

	// 预加载用户
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user` WHERE `user`.`id` IN (?,?,?)")).
		WithArgs(101, 102, 103).
		WillReturnRows(userRows)

	members, err := GetActivityMembersByActivityId(activityId)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedMembers), len(members))
	for i, m := range members {
		assert.Equal(t, expectedMembers[i].Id, m.Id)
		assert.Equal(t, expectedMembers[i].EventId, m.EventId)
		assert.Equal(t, expectedMembers[i].UserId, m.UserId)

		// 验证预加载用户信息
		assert.NotNil(t, m.User)
		assert.Equal(t, m.UserId, m.User.Id)
		assert.Equal(t, fmt.Sprintf("user%d", m.UserId), m.User.Username)
	}
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试查询错误
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `activity_member` WHERE event_id = ?")).
		WithArgs(activityId).
		WillReturnError(errors.New("query error"))

	members, err = GetActivityMembersByActivityId(activityId)
	assert.Nil(t, members)
	assert.EqualError(t, err, "query error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

// TestGetActivityMembersByUserId 测试根据用户Id获取活动成员列表
func TestGetActivityMembersByUserId(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	userId := int64(1)
	layout := "2006-01-02 15:04:05"
	createTime1, _ := time.Parse(layout, "2025-05-30 10:00:00")
	createTime2, _ := time.Parse(layout, "2025-05-30 11:00:00")
	createTime3, _ := time.Parse(layout, "2025-05-30 12:00:00")
	expectedMembers := []models.ActivityMember{
		{Id: 1, EventId: 101, UserId: userId, CreateTime: createTime1},
		{Id: 2, EventId: 102, UserId: userId, CreateTime: createTime2},
		{Id: 3, EventId: 103, UserId: userId, CreateTime: createTime3},
	}

	// 成员查询结果
	rows := sqlmock.NewRows([]string{"id", "event_id", "user_id", "create_time"})
	for _, m := range expectedMembers {
		rows.AddRow(m.Id, m.EventId, m.UserId, m.CreateTime)
	}

	// 活动查询结果
	activityRows := sqlmock.NewRows([]string{"id", "name", "user_id"}).
		AddRow(101, "活动101", 201).
		AddRow(102, "活动102", 202).
		AddRow(103, "活动103", 203)

	// 活动创建者查询结果
	activityCreatorRows := sqlmock.NewRows([]string{"id", "username", "name"}).
		AddRow(201, "creator201", "Creator 201").
		AddRow(202, "creator202", "Creator 202").
		AddRow(203, "creator203", "Creator 203")

	// 测试成功获取用户参与的活动
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `activity_member` WHERE user_id = ?")).
		WithArgs(userId).
		WillReturnRows(rows)

	// 预加载Activity
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `activity` WHERE `activity`.`id` IN (?,?,?)")).
		WithArgs(101, 102, 103).
		WillReturnRows(activityRows)

	// 预加载Activity.User
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user` WHERE `user`.`id` IN (?,?,?)")).
		WithArgs(201, 202, 203).
		WillReturnRows(activityCreatorRows)

	members, err := GetActivityMembersByUserId(userId)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedMembers), len(members))
	for i, m := range members {
		assert.Equal(t, expectedMembers[i].Id, m.Id)
		assert.Equal(t, expectedMembers[i].EventId, m.EventId)
		assert.Equal(t, expectedMembers[i].UserId, m.UserId)

		// 验证预加载活动
		assert.NotNil(t, m.Activity)
		assert.Equal(t, m.EventId, m.Activity.Id)
		assert.Equal(t, fmt.Sprintf("活动%d", m.EventId), m.Activity.Name)

		// 验证预加载活动创建者
		assert.NotNil(t, m.Activity.User)
		creatorId := m.Activity.UserId
		assert.Equal(t, fmt.Sprintf("creator%d", creatorId), m.Activity.User.Username)
	}
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试查询错误
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `activity_member` WHERE user_id = ?")).
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
		Id:         1,
		EventId:    1,
		UserId:     2,
		CreateTime: createTime,
	}

	// 测试成功更新活动成员
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `activity_member` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = UpdateActivityMember(activityMember)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试更新活动成员失败
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("UPDATE `activity_member` SET")).
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
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `activity_member` WHERE event_id = ? AND user_id = ?")).
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
	mock2.ExpectExec(regexp.QuoteMeta("DELETE FROM `activity_member` WHERE event_id = ? AND user_id = ?")).
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
		EventId:    1,
		UserId:     2,
		Content:    "这是一条测试评论",
		CreateTime: now,
	}

	// 测试成功添加活动评论
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `activity_comment`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := AddActivityComment(activityComment)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试添加活动评论失败
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("INSERT INTO `activity_comment`")).
		WillReturnError(errors.New("insert error"))
	mock2.ExpectRollback()

	err = AddActivityComment(activityComment)
	assert.EqualError(t, err, "insert error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

// TestGetActivityCommentsByActivityId 测试通过活动Id获取评论
func TestGetActivityCommentsByActivityId(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	activityId := int64(1)
	now := time.Now()
	expectedComments := []models.ActivityComment{
		{Id: 1, EventId: activityId, UserId: 101, Content: "评论1", CreateTime: now.Add(-time.Hour * 2)},
		{Id: 2, EventId: activityId, UserId: 102, Content: "评论2", CreateTime: now.Add(-time.Hour)},
		{Id: 3, EventId: activityId, UserId: 103, Content: "评论3", CreateTime: now},
	}

	// 评论查询结果
	rows := sqlmock.NewRows([]string{"id", "event_id", "user_id", "content", "create_time"})
	for _, c := range expectedComments {
		rows.AddRow(c.Id, c.EventId, c.UserId, c.Content, c.CreateTime)
	}

	// 用户查询结果
	userRows := sqlmock.NewRows([]string{"id", "username", "name"}).
		AddRow(101, "user101", "User 101").
		AddRow(102, "user102", "User 102").
		AddRow(103, "user103", "User 103")

	// 测试成功获取评论并预加载用户信息
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `activity_comment` WHERE event_id = ? ORDER BY create_time DESC")).
		WithArgs(activityId).
		WillReturnRows(rows)

	// 预加载用户
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user` WHERE `user`.`id` IN (?,?,?)")).
		WithArgs(101, 102, 103).
		WillReturnRows(userRows)

	comments, err := GetActivityCommentsByActivityId(activityId)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedComments), len(comments))
	for i, c := range comments {
		assert.Equal(t, expectedComments[i].Id, c.Id)
		assert.Equal(t, expectedComments[i].EventId, c.EventId)
		assert.Equal(t, expectedComments[i].UserId, c.UserId)
		assert.Equal(t, expectedComments[i].Content, c.Content)

		// 验证预加载用户信息
		assert.NotNil(t, c.User)
		assert.Equal(t, c.UserId, c.User.Id)
		assert.Equal(t, fmt.Sprintf("user%d", c.UserId), c.User.Username)
	}
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试查询错误
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `activity_comment` WHERE event_id = ? ORDER BY create_time DESC")).
		WithArgs(activityId).
		WillReturnError(errors.New("query error"))

	comments, err = GetActivityCommentsByActivityId(activityId)
	assert.Nil(t, comments)
	assert.EqualError(t, err, "query error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

// TestGetActivityCommentsByUserId 测试通过用户Id获取评论
func TestGetActivityCommentsByUserId(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	userId := int64(1)
	now := time.Now()
	expectedComments := []models.ActivityComment{
		{Id: 1, EventId: 101, UserId: userId, Content: "评论1", CreateTime: now.Add(-time.Hour * 2)},
		{Id: 2, EventId: 102, UserId: userId, Content: "评论2", CreateTime: now.Add(-time.Hour)},
		{Id: 3, EventId: 103, UserId: userId, Content: "评论3", CreateTime: now},
	}

	// 评论查询结果
	rows := sqlmock.NewRows([]string{"id", "event_id", "user_id", "content", "create_time"})
	for _, c := range expectedComments {
		rows.AddRow(c.Id, c.EventId, c.UserId, c.Content, c.CreateTime)
	}

	// 活动查询结果
	activityRows := sqlmock.NewRows([]string{"id", "name", "user_id"}).
		AddRow(101, "活动101", 201).
		AddRow(102, "活动102", 202).
		AddRow(103, "活动103", 203)

	// 活动创建者查询结果
	activityCreatorRows := sqlmock.NewRows([]string{"id", "username", "name"}).
		AddRow(201, "creator201", "Creator 201").
		AddRow(202, "creator202", "Creator 202").
		AddRow(203, "creator203", "Creator 203")

	// 测试成功获取评论
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `activity_comment` WHERE user_id = ? ORDER BY create_time DESC")).
		WithArgs(userId).
		WillReturnRows(rows)

	// 预加载Activity
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `activity` WHERE `activity`.`id` IN (?,?,?)")).
		WithArgs(101, 102, 103).
		WillReturnRows(activityRows)

	// 预加载Activity.User
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user` WHERE `user`.`id` IN (?,?,?)")).
		WithArgs(201, 202, 203).
		WillReturnRows(activityCreatorRows)

	comments, err := GetActivityCommentsByUserId(userId)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedComments), len(comments))
	for i, c := range comments {
		assert.Equal(t, expectedComments[i].Id, c.Id)
		assert.Equal(t, expectedComments[i].EventId, c.EventId)
		assert.Equal(t, expectedComments[i].UserId, c.UserId)
		assert.Equal(t, expectedComments[i].Content, c.Content)

		// 验证预加载Activity
		assert.NotNil(t, c.Activity)
		assert.Equal(t, c.EventId, c.Activity.Id)
		assert.Equal(t, fmt.Sprintf("活动%d", c.EventId), c.Activity.Name)

		// 验证预加载Activity.User
		assert.NotNil(t, c.Activity.User)
		creatorId := c.Activity.UserId
		assert.Equal(t, fmt.Sprintf("creator%d", creatorId), c.Activity.User.Username)
	}
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试查询错误
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `activity_comment` WHERE user_id = ? ORDER BY create_time DESC")).
		WithArgs(userId).
		WillReturnError(errors.New("query error"))

	comments, err = GetActivityCommentsByUserId(userId)
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
		Id:         1,
		EventId:    1,
		UserId:     2,
		Content:    "这是更新后的评论",
		CreateTime: now,
	}

	// 测试成功更新评论
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `activity_comment` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := UpdateActivityComment(activityComment)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// 测试更新评论失败
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("UPDATE `activity_comment` SET")).
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
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `activity_comment` WHERE `activity_comment`.`id` = ?")).
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
	mock2.ExpectExec(regexp.QuoteMeta("DELETE FROM `activity_comment` WHERE `activity_comment`.`id` = ?")).
		WithArgs(commentId).
		WillReturnError(errors.New("delete error"))
	mock2.ExpectRollback()

	err = DeleteActivityComment(commentId)
	assert.EqualError(t, err, "delete error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}
