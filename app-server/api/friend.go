package api

import (
	"hobbyhub-server/controllers"
	"hobbyhub-server/models"
	"hobbyhub-server/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary 获取好友列表
// @Description 通过用户ID获取好友列表
// @Tags 好友相关接口
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Success 200 {array} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/friend [get]
func GetFriendList(c *gin.Context) {
	jwtToken := c.GetHeader("Authorization")

	if jwtToken == "" {
		c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "jwt token is required"})
		return
	}

	// 验证JWT并确保用户有权限访问好友列表
	jwtUser, err := utils.ParseJWT(jwtToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "unauthorized access"})
		return
	}

	friends, err := controllers.GetAllFriendsByUserId(jwtUser.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, &models.ErrorResponse{ErrorMessage: "no friends"})
		return
	}

	c.JSON(http.StatusOK, friends)
}

type FriendRequest struct {
	UserId int64 `json:"user_id" binding:"required"` // 好友用户ID
}

// @Summary 发送好友申请
// @Description 通过用户ID发送好友申请
// @Tags 好友相关接口
// @Produce json
// @Param friendRequest body FriendRequest true "好友申请信息"
// @Param Authorization header string true "JWT Token"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/friend [post]
func SendFriendRequest(c *gin.Context) {
	var request FriendRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "invalid request"})
		return
	}

	jwtToken := c.GetHeader("Authorization")
	if jwtToken == "" {
		c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "jwt token is required"})
		return
	}

	jwtUser, err := utils.ParseJWT(jwtToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "unauthorized access"})
		return
	}

	friend := &models.Friend{
		UserId:     jwtUser.Id,
		FriendId:   request.UserId,
		Status:     2,
		CreateTime: time.Now(),
	}

	if err := controllers.AddFriend(friend); err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to send friend request"})
		return
	}

	c.JSON(http.StatusOK, &models.SuccessResponse{SuccessMessage: "friend request sent successfully"})
}

type UpdateFriendRequest struct {
	FriendId int64 `json:"friend_id" binding:"required"`        // 好友ID
	Status   int   `json:"status" binding:"required,oneof=0 1"` // 申请状态，0-拒绝，1-同意
}

// @Summary 更新好友状态
// @Description 通过ID更新好友状态
// @Tags 好友相关接口
// @Produce json
// @Param friendRequest body FriendRequest true "好友申请信息"
// @Param Authorization header string true "JWT Token"
// @Success 200 {object} models.Friend
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/friend [put]
func UpdateFriendStatus(c *gin.Context) {
	var request UpdateFriendRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "invalid request"})
		return
	}

	jwtToken := c.GetHeader("Authorization")
	if jwtToken == "" {
		c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "jwt token is required"})
		return
	}

	jwtUser, err := utils.ParseJWT(jwtToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "unauthorized access"})
		return
	}

	friend := &models.Friend{
		Status: request.Status,
	}

	dbFriend1, dbFriend2, err := controllers.GetFriendByUserIdAndFriendId(jwtUser.Id, request.FriendId)
	if err != nil {
		c.JSON(http.StatusNotFound, &models.ErrorResponse{ErrorMessage: "friend not found"})
		return
	}

	dbFriend1.UpdateFriendFields(*friend)
	dbFriend2.UpdateFriendFields(*friend)
	if err := controllers.UpdateFriendSynchronize(dbFriend1, dbFriend2); err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to update friend status"})
		return
	}
	c.JSON(http.StatusOK, dbFriend1)
}

// @Summary 删除好友
// @Description 通过ID删除好友状态
// @Tags 好友相关接口
// @Produce json
// @Param id path integer true "需要删除的Friend ID"
// @Param Authorization header string true "JWT Token"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/friend/:id [delete]
func DeleteFriend(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "invalid friend ID"})
		return
	}

	jwtToken := c.GetHeader("Authorization")
	if jwtToken == "" {
		c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "jwt token is required"})
		return
	}

	_, err := utils.ParseJWT(jwtToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "unauthorized access"})
		return
	}

	friendId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "invalid friend ID format"})
		return
	}

	dbFriend, err := controllers.GetFriendById(friendId)
	if err != nil {
		c.JSON(http.StatusNotFound, &models.ErrorResponse{ErrorMessage: "friend not found"})
		return
	}

	if err := controllers.DeleteFriendById(dbFriend.Id); err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to delete friend"})
		return
	}

	c.JSON(http.StatusOK, &models.SuccessResponse{SuccessMessage: "friend deleted successfully"})
}
