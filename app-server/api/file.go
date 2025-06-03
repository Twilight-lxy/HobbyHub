package api

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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

	_, err := utils.ParseJWT(jwtToken)
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

	// 写入文件
	// 使用TeeReader同时进行写入和哈希计算
	fileExt := filepath.Ext(filename)
	destFilename := filehash + fileExt
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
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "File hash mismatch"})
		return
	}

	fileInfo := &models.File{
		FileName:   filename,
		FileType:   fileExt,
		FileSize:   header.Size,
		FileHash:   filehash,
		CreateTime: time.Now(),
	}
	if err := controllers.AddFile(fileInfo); err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "Failed to save file info"})
		return
	}
	c.JSON(http.StatusOK, fileInfo)
}
