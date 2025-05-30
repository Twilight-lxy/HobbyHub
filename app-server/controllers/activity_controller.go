package controllers

import (
	"hobbyhub-server/config"
	"hobbyhub-server/models"
)

func AddActivity(activity *models.Activity) error {
	if err := config.DB.Create(activity).Error; err != nil {
		return err
	}
	return nil
}
func GetActivityById(activityId int64) (*models.Activity, error) {
	var activity models.Activity
	if err := config.DB.Where("id = ?", activityId).First(&activity).Error; err != nil {
		return nil, err
	}
	return &activity, nil
}
func UpdateActivity(activity *models.Activity) error {
	if err := config.DB.Save(activity).Error; err != nil {
		return err
	}
	return nil
}
func DeleteActivityById(activityId int64) error {
	if err := config.DB.Delete(&models.Activity{}, activityId).Error; err != nil {
		return err
	}
	return nil
}

func AddActivityMember(activityMember *models.ActivityMember) error {
	if err := config.DB.Create(activityMember).Error; err != nil {
		return err
	}
	return nil
}

func GetActivityMembersByActivityId(activityId int64) ([]models.ActivityMember, error) {
	var members []models.ActivityMember
	if err := config.DB.Where("activity_id = ?", activityId).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func GetActivityMembersByUserId(userId int64) ([]models.ActivityMember, error) {
	var members []models.ActivityMember
	if err := config.DB.Where("user_id = ?", userId).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func UpdateActivityMember(activityMember *models.ActivityMember) error {
	if err := config.DB.Save(activityMember).Error; err != nil {
		return err
	}
	return nil
}

func DeleteActivityMember(activityId, userId int64) error {
	if err := config.DB.Where("activity_id = ? AND user_id = ?", activityId, userId).Delete(&models.ActivityMember{}).Error; err != nil {
		return err
	}
	return nil
}

func AddActivityComment(activityComment *models.ActivityComment) error {
	if err := config.DB.Create(activityComment).Error; err != nil {
		return err
	}
	return nil
}

func GetActivityCommentsByActivityId(activityId int64) ([]models.ActivityComment, error) {
	var comments []models.ActivityComment
	if err := config.DB.Where("event_id = ?", activityId).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func GetActivityCommentsByUserId(userId int64) ([]models.ActivityComment, error) {
	var comments []models.ActivityComment
	if err := config.DB.Where("user_id = ?", userId).Find(&comments).Error; err != nil {
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

func UpdateActivityComment(activityComment *models.ActivityComment) error {
	if err := config.DB.Save(activityComment).Error; err != nil {
		return err
	}
	return nil
}

func DeleteActivityComment(commentId int64) error {
	if err := config.DB.Delete(&models.ActivityComment{}, commentId).Error; err != nil {
		return err
	}
	return nil
}
