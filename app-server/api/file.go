package api

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"hobbyhub-server/config"
	"hobbyhub-server/controllers"
	"hobbyhub-server/models"
	"hobbyhub-server/utils"

	"github.com/gin-gonic/gin"
)

// @Summary 上传文件
// @Description 上传文件
// @Tags 文件相关接口
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "JWT token"
// @Param filename formData string true "文件名"
// @Param filehash formData string true "文件Hash"
// @Param file formData file true "文件数据"
// @Success 200 {object} models.File
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /v1/file [post]
func UploadFile(c *gin.Context) {
	// 检查用户是否登录
	jwtToken := c.GetHeader("Authorization")
	if jwtToken == "" {
		c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "jwt token is required"})
		return
	}

	jwtUser, err := utils.ParseJWT(jwtToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "invalid jwt token"})
		return
	}

	//限制文件大小
	fileConfig := config.GetConfig().File
	if err := c.Request.ParseMultipartForm(fileConfig.MaxSize << 20); err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "File size exceeds limit"})
		return
	}
	//获取文件名和文件hash
	filename := c.PostForm("filename")
	filehash := c.PostForm("filehash")
	if filename == "" || filehash == "" {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "filename and filehash are required"})
		return
	}

	// 获取文件数据
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "file is required"})
		return
	}
	defer file.Close()

	//验证文件大小
	if header.Size > fileConfig.MaxSize<<20 {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "File size exceeds limit"})
		return
	}
	//验证文件类型
	fileExt := filepath.Ext(filename)
	flag := false
	for _, allowType := range fileConfig.AllowedTypes {
		if strings.EqualFold(fileExt, "."+allowType) {
			flag = true
			break
		}
	}
	if !flag {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "File type not allowed"})
		return
	}

	fileInfo := &models.File{
		FileName:     filename,
		FileType:     fileExt,
		FileSize:     header.Size,
		FileHash:     filehash,
		CreateTime:   utils.GetCurrentTime(),
		UpLoadUserId: jwtUser.Id,
	}
	if err := controllers.AddFile(fileInfo); err != nil {
		linkFileInfo, err := controllers.GetFileByHash(filehash)
		if err != nil && linkFileInfo != nil {
			c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "Failed to check existing file"})
			return
		}
		fileInfo.FileHash = linkFileInfo.FileHash + "_" + utils.GenerateRandomString(9) // 避免哈希冲突
		fileInfo.LinkFileId = linkFileInfo.Id                                           // 关联已存在的文件ID
		err = controllers.AddFile(fileInfo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "upload file failed"})
		}
		linkFileInfo.LinkFileId -= 1
		controllers.UpdateFile(linkFileInfo)
		c.JSON(http.StatusOK, fileInfo) // 返回已存在的文件信息
		return
	}

	// 写入文件
	// 使用TeeReader同时进行写入和哈希计算

	destFilename := fmt.Sprintf("%d%s", fileInfo.Id, fileExt)
	// 确保上传路径存在
	if err := os.MkdirAll(fileConfig.UploadPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "Failed to create upload directory"})
		return
	}
	destPath := filepath.Join(fileConfig.UploadPath, destFilename)
	out, err := os.Create(destPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "Failed to create file"})
		return
	}
	defer out.Close()

	hash := sha256.New()
	teeReader := io.TeeReader(file, hash)

	if _, err := io.Copy(out, teeReader); err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "Failed to save file"})
		return
	}

	calculatedHash := fmt.Sprintf("%x", hash.Sum(nil))
	if !strings.EqualFold(calculatedHash, filehash) {
		// 删除已写入的文件
		os.Remove(destPath)
		// 删除数据库中的记录
		controllers.DeleteFileById(fileInfo.Id)
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "File hash mismatch"})
		return
	}

	c.JSON(http.StatusOK, fileInfo)
}

// @Summary 下载文件
// @Description 下载文件
// @Tags 文件相关接口
// @Produce octet-stream
// @Produce json
// @Param Authorization header string true "JWT token"
// @Param id path int true "文件ID"
// @Success 200 {file} binary
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /v1/file/{id} [get]
func DownloadFile(c *gin.Context) {
	// 检查用户是否登录
	jwtToken := c.GetHeader("Authorization")
	if jwtToken == "" {
		c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "jwt token is required"})
		return
	}

	_, err := utils.ParseJWT(jwtToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "invalid jwt token"})
		return
	}
	fileIdStr := c.Param("id")
	if fileIdStr == "" {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "fileId is required"})
		return
	}

	fileId, err := utils.StringToInt64(fileIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "Invalid fileId format"})
		return
	}

	fileInfo, err := controllers.GetFileById(fileId)
	if err != nil {
		c.JSON(http.StatusNotFound, &models.ErrorResponse{ErrorMessage: "file not found"})
		return
	}

	if fileInfo.LinkFileId > 0 {
		// 如果是关联文件，获取原始文件信息
		originalFile, err := controllers.GetFileById(fileInfo.LinkFileId)
		if err != nil {
			c.JSON(http.StatusNotFound, &models.ErrorResponse{ErrorMessage: "original file not found"})
			return
		}
		fileInfo = originalFile
	}
	fileConfig := config.GetConfig().File
	destPath := filepath.Join(fileConfig.UploadPath, fmt.Sprintf("%d%s", fileInfo.Id, fileInfo.FileType))
	c.FileAttachment(destPath, fileInfo.FileName)
}

// @Summary 删除文件
// @Description 删除文件
// @Tags 文件相关接口
// @Produce json
// @Param Authorization header string true "JWT token"
// @Param id path int true "文件ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /v1/file/{id} [delete]
func DeleteFile(c *gin.Context) {
	// 检查用户是否登录
	jwtToken := c.GetHeader("Authorization")
	if jwtToken == "" {
		c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "jwt token is required"})
		return
	}

	jwtUser, err := utils.ParseJWT(jwtToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "invalid jwt token"})
		return
	}

	fileIdStr := c.Param("id")
	if fileIdStr == "" {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "fileId is required"})
		return
	}

	fileId, err := utils.StringToInt64(fileIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "Invalid fileId format"})
		return
	}

	fileInfo, err := controllers.GetFileById(fileId)
	if err != nil {
		c.JSON(http.StatusNotFound, &models.ErrorResponse{ErrorMessage: "file not found"})
		return
	}
	if fileInfo.UpLoadUserId != jwtUser.Id {
		c.JSON(http.StatusForbidden, &models.ErrorResponse{ErrorMessage: "You do not have permission to delete this file"})
		return
	}
	if fileInfo.LinkFileId > 0 {
		linkFile, err := controllers.GetFileById(fileInfo.LinkFileId)
		if err != nil {
			c.JSON(http.StatusNotFound, &models.ErrorResponse{ErrorMessage: "linked file not found"})
			return
		}
		linkFile.LinkFileId += 1 // 减少关联计数
		if err := controllers.UpdateFile(linkFile); err != nil {
			c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "Failed to update linked file record"})
			return
		}
		if err := controllers.DeleteFileById(fileId); err != nil {
			c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "Failed to delete file record from database"})
			return
		}
	} else if fileInfo.LinkFileId == 0 {
		destPath := filepath.Join(config.GetConfig().File.UploadPath, fmt.Sprintf("%d%s", fileInfo.Id, fileInfo.FileType))
		if err := os.Remove(destPath); err != nil {
			c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "Failed to delete file from disk"})
			return
		}
		if err := controllers.DeleteFileById(fileId); err != nil {
			c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "Failed to delete file record from database"})
			return
		}
	} else {
		fileInfo.UpLoadUserId = 0 // 清除上传用户ID
		if err := controllers.UpdateFile(fileInfo); err != nil {
			c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "Failed to update file record"})
			return
		}
	}
	c.JSON(http.StatusOK, &models.SuccessResponse{SuccessMessage: "File deleted successfully"})
}

func ServeBasicFile(c *gin.Context) {
	filename := c.Param("filename")
	log.Printf("Serving file: %s", filename)
	fileConfig := config.GetConfig().File
	filePath := filepath.Join(fileConfig.UploadPath, filename)
	c.File(filePath)
}
