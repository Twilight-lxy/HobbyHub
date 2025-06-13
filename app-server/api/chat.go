package api

import (
	"net/http"
	"sort"
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
// @Param to_user_id query int true "对方用户ID"
// @Param starttime query string false "开始时间，格式为YYYY-MM-DD HH:MM:SS"
// @Param endtime query string false "结束时间，格式为YYYY-MM-DD HH:MM:SS"
// @Param Authorization header string true "JWT Token"
// @Success 200 {array} models.Chat
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/chat	 [get]
func GetChatHistory(c *gin.Context) {
	// 解析参数

	toUserIdStr := c.Query("to_user_id")
	startTime := c.Query("starttime")
	endTime := c.Query("endtime")
	jwtToken := c.GetHeader("Authorization")

	// 验证JWT并确保用户有权限访问这些聊天记录
	if jwtToken == "" {
		c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "jwt token is required"})
		return
	}

	jwtUser, err := utils.ParseJWT(jwtToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "unauthorized"})
		return
	}

	fromUserId := jwtUser.Id

	// 验证必需参数
	if toUserIdStr == "" {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "to_user_id are required"})
		return
	}

	toUserId, err := utils.StringToInt64(toUserIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "invalid to_user_id"})
		return
	}

	// 获取从fromUserId到toUserId的聊天记录
	chatsForward, err := controllers.GetAllChatByFromUserIdToUserId(fromUserId, toUserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to get chat history"})
		return
	}
	// 过滤已删除记录
	var filteredChatsForward []models.Chat
	for _, chat := range chatsForward {
		if chat.StatusFrom != 0 {
			filteredChatsForward = append(filteredChatsForward, chat)
		}
	}
	chatsForward = filteredChatsForward
	// 获取从toUserId到fromUserId的聊天记录
	chatsBackward, err := controllers.GetAllChatByFromUserIdToUserId(toUserId, fromUserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to get chat history"})
		return
	}
	// 过滤已删除记录
	var filteredChatsBackward []models.Chat
	for _, chat := range chatsBackward {
		if chat.StatusTo != 0 {
			filteredChatsBackward = append(filteredChatsBackward, chat)
		}
	}
	chatsBackward = filteredChatsBackward
	// 合并两个方向的聊天记录
	allChats := append(chatsForward, chatsBackward...)

	// 按创建时间降序排序
	if len(allChats) > 0 {
		sort.Slice(allChats, func(i, j int) bool {
			return allChats[i].CreateTime.After(allChats[j].CreateTime)
		})
	}

	// 处理时间过滤
	if startTime != "" || endTime != "" {
		var startTimeObj, endTimeObj time.Time

		if startTime != "" {
			startTimeObj = utils.ParseTimeFromString(startTime)
		}

		if endTime != "" {
			endTimeObj = utils.ParseTimeFromString(endTime)
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
// @Router /v1/chat [post]
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
		if friend.FriendId == req.UserIdTo && friend.Status == 1 { // 1表示好友关系已建立
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
		CreateTime: utils.GetCurrentTime(),
		StatusFrom: 1,
		StatusTo:   1,
	}

	if err := controllers.AddChat(chat); err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to create chat message"})
		return
	}

	c.JSON(http.StatusOK, *chat)
}

// @Summary 删除聊天消息
// @Description 删除聊天消息
// @Tags 聊天相关接口
// @Accept json
// @Produce json
// @Param id path integer true "需要删除的Chat 记录ID"
// @Param Authorization header string true "JWT Token"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/chat/{id} [delete]
func DeleteChat(c *gin.Context) {
	chatIdStr := c.Param("id")
	jwtToken := c.GetHeader("Authorization")

	if jwtToken == "" {
		c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "jwt token is required"})
		return
	}
	jwtUser, err := utils.ParseJWT(jwtToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "unauthorized"})
		return
	}
	chatId, err := utils.StringToInt64(chatIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "invalid chat id"})
		return
	}
	err = controllers.DeleteChatById(chatId, jwtUser.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to delete chat message"})
		return
	}
	c.JSON(http.StatusOK, &models.SuccessResponse{SuccessMessage: "chat message deleted successfully"})
}
