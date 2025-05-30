package api

import (
	"net/http"
	"strconv"
	"time"

	"hobbyhub-server/controllers"
	"hobbyhub-server/models"
	"hobbyhub-server/utils"

	"github.com/gin-gonic/gin"
)

// @Summary 获取聊天记录
// @Description 根据两个用户的id获取聊天记录列表
// @Tags 聊天相关接口
// @Accept json
// @Produce json
// @Param from_user_id query int true "发送者用户ID"
// @Param to_user_id query int true "接收者用户ID"
// @Param starttime query string false "开始时间，格式为YYYY-MM-DD HH:MM:SS"
// @Param endtime query string false "结束时间，格式为YYYY-MM-DD HH:MM:SS"
// @Param Authorization header string true "JWT Token"
// @Success 200 {array} models.Chat
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/chat/history [get]
func GetChatHistory(c *gin.Context) {
	// 解析参数
	fromUserIDStr := c.Query("from_user_id")
	toUserIDStr := c.Query("to_user_id")
	startTime := c.Query("starttime")
	endTime := c.Query("endtime")
	jwtToken := c.GetHeader("Authorization")

	// 验证必需参数
	if fromUserIDStr == "" || toUserIDStr == "" {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "from_user_id and to_user_id are required"})
		return
	}

	// 将用户ID从字符串转换为int64
	fromUserID, err := strconv.ParseInt(fromUserIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "invalid from_user_id"})
		return
	}

	toUserID, err := strconv.ParseInt(toUserIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "invalid to_user_id"})
		return
	}

	// 验证JWT并确保用户有权限访问这些聊天记录
	if jwtToken == "" {
		c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "jwt token is required"})
		return
	}

	jwtUser, err := utils.ParseJWT(jwtToken)
	if err != nil || (jwtUser.ID != fromUserID && jwtUser.ID != toUserID) {
		c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "unauthorized"})
		return
	}

	// 获取从fromUserID到toUserID的聊天记录
	chatsForward, err := controllers.GetAllChatByFromUserIDToUserID(fromUserID, toUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to get chat history"})
		return
	}

	// 获取从toUserID到fromUserID的聊天记录
	chatsBackward, err := controllers.GetAllChatByFromUserIDToUserID(toUserID, fromUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to get chat history"})
		return
	}

	// 合并两个方向的聊天记录
	allChats := append(chatsForward, chatsBackward...)

	// 处理时间过滤
	if startTime != "" || endTime != "" {
		var startTimeObj, endTimeObj time.Time
		var err error

		if startTime != "" {
			startTimeObj, err = time.Parse("2006-01-02 15:04:05", startTime)
			if err != nil {
				c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "invalid starttime format"})
				return
			}
		}

		if endTime != "" {
			endTimeObj, err = time.Parse("2006-01-02 15:04:05", endTime)
			if err != nil {
				c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "invalid endtime format"})
				return
			}
		}

		var filteredChats []models.Chat
		for _, chat := range allChats {
			if startTime != "" && chat.CreateTime.Before(startTimeObj) {
				continue
			}

			if endTime != "" && chat.CreateTime.After(endTimeObj) {
				continue
			}

			filteredChats = append(filteredChats, chat)
		}

		c.JSON(http.StatusOK, filteredChats)
	} else {
		c.JSON(http.StatusOK, allChats)
	}
}

type newChatRequest struct {
	UserIDTo int64  `json:"user_id_to" binding:"required"` // 接收者用户ID
	Content  string `json:"content" binding:"required"`    // 聊天内容
}

// @Summary 发送聊天消息
// @Description 发送聊天消息
// @Tags 聊天相关接口
// @Accept json
// @Produce json
// @Param chat body newChatRequest true "聊天记录内容"
// @Param Authorization header string true "JWT Token"
// @Success 200 {array} models.Chat
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/chat/send [post]
func SendChat(c *gin.Context) {
	var req newChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "invalid request format"})
		return
	}

	// 验证JWT并获取发送者ID
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
	// 验证接收者用户ID是否是发送者用户的好友

	// 创建聊天记录
	chat := &models.Chat{
		UserIDFrom: jwtUser.ID,
		UserIDTo:   req.UserIDTo,
		Content:    req.Content,
		CreateTime: time.Now(),
		StatusFrom: 1, // 已发送状态
		StatusTo:   0, // 未读状态
	}

	if err := controllers.AddChat(chat); err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to create chat message"})
		return
	}

	c.JSON(http.StatusOK, chat)
}
