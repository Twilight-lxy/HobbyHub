package controllers

import (
	"hobbyhub-server/config"
	"hobbyhub-server/models"
)

// AddActivity 添加新活动
func AddActivity(activity *models.Activity) error {
	return config.DB.Create(activity).Error
}

// GetAllActivitById 获取指定活动
func GetActivityById(activityId int64) (*models.Activity, error) {
	var activity models.Activity
	if err := config.DB.Where("id = ? AND if_delete = 0", activityId).
		Preload("User").
		Preload("Members").
		Preload("Members.User").
		Preload("Comments").
		Preload("Comments.User").
		First(&activity).Error; err != nil {
		return nil, err
	}
	return &activity, nil
}

// GetAllActivities 获取所有活动
func GetAllActivities() ([]models.Activity, error) {
	var activities []models.Activity
	if err := config.DB.Where("if_delete = 0").
		Order("create_time DESC").
		Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

// UpdateActivity 更新活动信息
func UpdateActivity(activity *models.Activity) error {
	if err := config.DB.Save(activity).Error; err != nil {
		return err
	}
	return nil
}

// DeleteActivityById 软删除活动(标记删除状态)
func DeleteActivityById(activityId int64) error {
	// 使用事务进行软删除
	tx := config.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&models.Activity{}).
		Where("id = ?", activityId).
		Update("if_delete", 1).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetActivitiesByUserId 获取用户创建的活动
func GetActivitiesByUserId(userId int64) ([]models.Activity, error) {
	var activities []models.Activity
	if err := config.DB.Where("user_id = ? AND if_delete = 0", userId).
		Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

// AddActivityMember 添加活动成员
func AddActivityMember(activityMember *models.ActivityMember) error {
	if err := config.DB.Create(activityMember).Error; err != nil {
		return err
	}
	return nil
}

// GetActivityMembersByActivityId 获取活动的所有成员
func GetActivityMembersByActivityId(activityId int64) ([]models.ActivityMember, error) {
	var members []models.ActivityMember
	if err := config.DB.Where("event_id = ?", activityId).
		Preload("User").
		Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

// GetActivityMembersByUserId 获取用户参加的所有活动的成员记录
func GetActivityMembersByUserId(userId int64) ([]models.ActivityMember, error) {
	var members []models.ActivityMember
	if err := config.DB.Where("user_id = ?", userId).
		Preload("Activity").
		Preload("Activity.User").
		Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

// UpdateActivityMember 更新活动成员
func UpdateActivityMember(activityMember *models.ActivityMember) error {
	if err := config.DB.Save(activityMember).Error; err != nil {
		return err
	}
	return nil
}

// DeleteActivityMember 删除活动成员
func DeleteActivityMember(activityId, userId int64) error {
	if err := config.DB.Where("event_id = ? AND user_id = ?", activityId, userId).
		Delete(&models.ActivityMember{}).Error; err != nil {
		return err
	}
	return nil
}

// AddActivityComment 添加活动评论
func AddActivityComment(activityComment *models.ActivityComment) error {
	if err := config.DB.Create(activityComment).Error; err != nil {
		return err
	}
	return nil
}

// GetActivityCommentsByActivityId 获取活动的所有评论
func GetActivityCommentsByActivityId(activityId int64) ([]models.ActivityComment, error) {
	var comments []models.ActivityComment
	if err := config.DB.Where("event_id = ?", activityId).
		Preload("User").
		Order("create_time DESC").
		Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// GetActivityCommentsByUserId 获取用户的所有活动评论
func GetActivityCommentsByUserId(userId int64) ([]models.ActivityComment, error) {
	var comments []models.ActivityComment
	if err := config.DB.Where("user_id = ?", userId).
		Preload("Activity").
		Preload("Activity.User").
		Order("create_time DESC").
		Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func GetActivityCommentsByActivityIdAndUserId(activityId, userId int64) ([]models.ActivityComment, error) {
	var comments []models.ActivityComment
	if err := config.DB.Where("event_id = ? AND user_id = ?", activityId, userId).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// UpdateActivityComment 更新活动评论
func UpdateActivityComment(activityComment *models.ActivityComment) error {
	if err := config.DB.Save(activityComment).Error; err != nil {
		return err
	}
	return nil
}

// DeleteActivityComment 删除活动评论
func DeleteActivityComment(commentId int64) error {
	if err := config.DB.Delete(&models.ActivityComment{}, commentId).Error; err != nil {
		return err
	}
	return nil
}
