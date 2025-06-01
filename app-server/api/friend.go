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

type FriendResponse struct {
	Id         int64     `json:"id" binding:"required"`          // 记录ID
	FriendId   int64     `json:"friend_id" binding:"required"`   // 好友ID
	Status     int       `json:"status" binding:"required"`      // 好友状态（0: 拒绝, 1: 接受, 2: 等待接受, 3：已发出申请）
	CreateTime time.Time `json:"create_time" binding:"required"` // 创建时间
}

func infoFromModelsFriend(friend models.Friend) FriendResponse {
	// 将 models.Friend 转换为 FriendResponse
	return FriendResponse{
		Id:         friend.Id,
		FriendId:   friend.FriendId,
		Status:     friend.Status,
		CreateTime: friend.CreateTime,
	}
}
func infoFromModelsFriends(friends []models.Friend) []FriendResponse {
	// 将 []models.Friend 转换为 []FriendResponse
	var friendResponses []FriendResponse
	for _, friend := range friends {
		friendResponses = append(friendResponses, infoFromModelsFriend(friend))
	}
	return friendResponses
}

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

	c.JSON(http.StatusOK, infoFromModelsFriends(friends))
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

	friend1, friend2, _ := controllers.GetFriendByUserIdAndFriendId(jwtUser.Id, request.UserId)
	if friend1 != nil && friend2 != nil {
		if friend1.Status == 1 && friend2.Status == 1 {
			c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "friend already exists"})
		} else if friend1.Status == 3 && friend2.Status == 2 {
			c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "friend request already sent"})
		} else if friend1.Status == 2 && friend2.Status == 3 {
			friend1.Status = 1 // 更新状态为接受
			friend2.Status = 1 // 更新状态为接受
			if err := controllers.UpdateFriendSynchronize(friend1, friend2); err != nil {
				c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to update friend status"})
				return
			}
			c.JSON(http.StatusOK, &models.SuccessResponse{SuccessMessage: "friend request updated successfully"})
			return
		} else if friend1.Status == 0 || friend2.Status == 0 {
			c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "friend request already rejected"})
			friend1.Status = 3 // 更新状态为已发出申请
			friend2.Status = 2 // 更新状态为等待接受
			if err := controllers.UpdateFriendSynchronize(friend1, friend2); err != nil {
				c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to update friend status"})
				return
			}
			c.JSON(http.StatusOK, &models.SuccessResponse{SuccessMessage: "friend request sent successfully"})
			return
		} else {
			c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "friend request error"})
			return
		}
	} else {
		friend := &models.Friend{
			UserId:     jwtUser.Id,
			FriendId:   request.UserId,
			Status:     3,
			CreateTime: time.Now(),
		}

		if err := controllers.AddFriend(friend); err != nil {
			c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to send friend request"})
			return
		}

		c.JSON(http.StatusOK, &models.SuccessResponse{SuccessMessage: "friend request sent successfully"})
		return
	}
}

type UpdateFriendRequest struct {
	FriendId int64 `json:"friend_id" binding:"required"`        // 好友ID
	Status   int   `json:"status" binding:"required,oneof=0 1"` // 申请状态，0-拒绝，1-同意
}

// @Summary 更新好友状态
// @Description 通过ID更新好友状态
// @Tags 好友相关接口
// @Produce json
// @Param updatefriendRequest body UpdateFriendRequest true "好友申请信息"
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
	if dbFriend1.Status != 2 && dbFriend1.Status != 3 {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "friend request not in pending state"})
		return
	}
	dbFriend1.UpdateFriendFields(*friend)
	dbFriend2.UpdateFriendFields(*friend)
	if err := controllers.UpdateFriendSynchronize(dbFriend1, dbFriend2); err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to update friend status"})
		return
	}
	c.JSON(http.StatusOK, infoFromModelsFriend(*dbFriend1))
}

// @Summary 删除好友
// @Description 通过ID删除好友状态
// @Tags 好友相关接口
// @Produce json
// @Param id path integer true "需要删除的Friend 记录ID"
// @Param Authorization header string true "JWT Token"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/friend/{id} [delete]
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

	dbFriend.Status = 0 // 设置状态为0，表示删除好友

	if err := controllers.UpdateFriend(dbFriend); err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to delete friend"})
		return
	}

	c.JSON(http.StatusOK, &models.SuccessResponse{SuccessMessage: "friend deleted successfully"})
}
