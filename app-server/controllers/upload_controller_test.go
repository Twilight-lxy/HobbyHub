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

func TestAddFile(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	file := &models.File{FileName: "test.jpg", FileType: "image/jpeg", FileSize: 1024, FileHash: "abc123"}

	// Test successful file creation
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `files`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := AddFile(file)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test file creation failure
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("INSERT INTO `files`")).
		WillReturnError(errors.New("insert error"))
	mock2.ExpectRollback()

	err = AddFile(file)
	assert.EqualError(t, err, "insert error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestGetFileById(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	fileId := int64(1)
	createTime, _ := time.Parse("2006-01-02 15:04:05", "2023-06-15 10:00:00")
	expectedFile := &models.File{
		ID:         fileId,
		FileName:   "test.jpg",
		FileType:   "image/jpeg",
		FileSize:   1024,
		FileHash:   "abc123",
		CreateTime: createTime,
	}

	rows := sqlmock.NewRows([]string{"id", "file_name", "file_type", "file_size", "file_hash", "create_time"}).
		AddRow(expectedFile.ID, expectedFile.FileName, expectedFile.FileType, expectedFile.FileSize, expectedFile.FileHash, expectedFile.CreateTime)

	// Test successful file retrieval
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `files` WHERE id = ? ORDER BY `files`.`id` LIMIT ?")).
		WithArgs(fileId, 1).
		WillReturnRows(rows)

	file, err := GetFileById(fileId)
	assert.NoError(t, err)
	assert.Equal(t, expectedFile.ID, file.ID)
	assert.Equal(t, expectedFile.FileName, file.FileName)
	assert.Equal(t, expectedFile.FileHash, file.FileHash)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test file not found
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `files` WHERE id = ? ORDER BY `files`.`id` LIMIT ?")).
		WithArgs(fileId, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	file, err = GetFileById(fileId)
	assert.Nil(t, file)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestGetFileByHash(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	fileHash := "abc123"
	createTime, _ := time.Parse("2006-01-02 15:04:05", "2023-06-15 10:00:00")
	expectedFile := &models.File{
		ID:         1,
		FileName:   "test.jpg",
		FileType:   "image/jpeg",
		FileSize:   1024,
		FileHash:   fileHash,
		CreateTime: createTime,
	}

	rows := sqlmock.NewRows([]string{"id", "file_name", "file_type", "file_size", "file_hash", "create_time"}).
		AddRow(expectedFile.ID, expectedFile.FileName, expectedFile.FileType, expectedFile.FileSize, expectedFile.FileHash, expectedFile.CreateTime)

	// Test successful file retrieval by hash
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `files` WHERE file_hash = ? ORDER BY `files`.`id` LIMIT ?")).
		WithArgs(fileHash, 1).
		WillReturnRows(rows)

	file, err := GetFileByHash(fileHash)
	assert.NoError(t, err)
	assert.Equal(t, expectedFile.ID, file.ID)
	assert.Equal(t, expectedFile.FileHash, file.FileHash)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test file not found by hash
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `files` WHERE file_hash = ? ORDER BY `files`.`id` LIMIT ?")).
		WithArgs(fileHash, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	file, err = GetFileByHash(fileHash)
	assert.Nil(t, file)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestUpdateFile(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	file := models.File{ID: 1, FileName: "updated.jpg", FileType: "image/jpeg", FileSize: 2048, FileHash: "xyz789"}

	// Test successful file update
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `files` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := UpdateFile(&file)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test file update error
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("UPDATE `files` SET")).
		WillReturnError(errors.New("update error"))
	mock2.ExpectRollback()

	err = UpdateFile(&file)
	assert.EqualError(t, err, "update error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}

func TestDeleteFileById(t *testing.T) {
	mock, teardown := SetupMockDB(t)
	defer teardown()

	fileId := int64(1)

	// Test successful file deletion
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `files` WHERE `files`.`id` = ?")).
		WithArgs(fileId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := DeleteFileById(fileId)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test file deletion error
	mock2, teardown2 := SetupMockDB(t)
	defer teardown2()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta("DELETE FROM `files` WHERE `files`.`id` = ?")).
		WithArgs(fileId).
		WillReturnError(errors.New("delete error"))
	mock2.ExpectRollback()

	err = DeleteFileById(fileId)
	assert.EqualError(t, err, "delete error")
	assert.NoError(t, mock2.ExpectationsWereMet())
}
