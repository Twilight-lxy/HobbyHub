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
