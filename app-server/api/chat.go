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
	fromUserIdStr := c.Query("from_user_id")
	toUserIdStr := c.Query("to_user_id")
	startTime := c.Query("starttime")
	endTime := c.Query("endtime")
	jwtToken := c.GetHeader("Authorization")

	// 验证必需参数
	if fromUserIdStr == "" || toUserIdStr == "" {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "from_user_id and to_user_id are required"})
		return
	}

	// 将用户ID从字符串转换为int64
	fromUserId, err := strconv.ParseInt(fromUserIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "invalid from_user_id"})
		return
	}

	toUserId, err := strconv.ParseInt(toUserIdStr, 10, 64)
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
	if err != nil || (jwtUser.Id != fromUserId && jwtUser.Id != toUserId) {
		c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "unauthorized"})
		return
	}

	// 获取从fromUserId到toUserId的聊天记录
	chatsForward, err := controllers.GetAllChatByFromUserIdToUserId(fromUserId, toUserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to get chat history"})
		return
	}

	// 获取从toUserId到fromUserId的聊天记录
	chatsBackward, err := controllers.GetAllChatByFromUserIdToUserId(toUserId, fromUserId)
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
	UserIdTo int64  `json:"user_id_to" binding:"required"` // 接收者用户Id
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

	// 验证JWT并获取发送者Id
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
	// 验证接收者用户Id是否是发送者用户的好友
	friends, err := controllers.GetAllFriendsByUserId(jwtUser.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to get friends list"})
		return
	}
	isFriend := false
	for _, friend := range friends {
		if friend.Id == req.UserIdTo {
			isFriend = true
			break
		}
	}
	if !isFriend {
		c.JSON(http.StatusForbidden, &models.ErrorResponse{ErrorMessage: "you can only send messages to your friends"})
		return
	}
	// 创建聊天记录
	chat := &models.Chat{
		UserIdFrom: jwtUser.Id,
		UserIdTo:   req.UserIdTo,
		Content:    req.Content,
		CreateTime: time.Now(),
		StatusFrom: 0,
		StatusTo:   0,
	}

	if err := controllers.AddChat(chat); err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to create chat message"})
		return
	}

	c.JSON(http.StatusOK, chat)
}
