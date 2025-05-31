package controllers

import (
	"hobbyhub-server/config"
	"hobbyhub-server/models"
)

func AddFile(file *models.File) error {
	if err := config.DB.Create(file).Error; err != nil {
		return err
	}
	return nil
}
func GetFileById(fileId int64) (*models.File, error) {
	var file models.File
	if err := config.DB.Where("id = ?", fileId).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}
func GetFileByHash(fileHash string) (*models.File, error) {
	var file models.File
	if err := config.DB.Where("file_hash = ?", fileHash).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}
func UpdateFile(file *models.File) error {
	if err := config.DB.Save(file).Error; err != nil {
		return err
	}
	return nil
}
func DeleteFileById(fileId int64) error {
	if err := config.DB.Delete(&models.File{}, fileId).Error; err != nil {
		return err
	}
	return nil
}
