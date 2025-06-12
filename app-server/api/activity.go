package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"hobbyhub-server/controllers"
	"hobbyhub-server/models"
	"hobbyhub-server/utils"
)

// @Summary 获取活动信息
// @Description 获取指定活动的完整信息
// @Tags 活动相关接口
// @Accept json
// @Produce json
// @Param id path int true "活动ID"
// @Success 200 {array} models.Activity
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/activity/{id} [get]
func GetActivitie(c *gin.Context) {
	// 获取活动ID
	activityIdStr := c.Param("id")
	if activityIdStr == "" {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "activity id is required"})
		return
	}
	activityId, err := strconv.ParseInt(activityIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "Invalid activity id format"})
		return
	}

	// 获取活动信息
	activity, err := controllers.GetActivityById(activityId)
	if err != nil {
		c.JSON(http.StatusNotFound, &models.ErrorResponse{ErrorMessage: "activity not found"})
		return
	}

	c.JSON(http.StatusOK, activity)
}

type simpleActivity struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (s *simpleActivity) LoadFromModelActivity(activity models.Activity) {
	s.Id = activity.Id
	s.Name = activity.Name
}

// @Summary 获取所有活动Id
// @Description 获取所有活动的id与名称
// @Tags 活动相关接口
// @Accept json
// @Produce json
// @Success 200 {array} models.Activity
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/activity [get]
func GetAllActivitie(c *gin.Context) {
	// 获取所有活动ID
	activitys, err := controllers.GetAllActivities()
	if err != nil {
		c.JSON(http.StatusNotFound, &models.ErrorResponse{ErrorMessage: "activities not found"})
		return
	}
	var activityResponse []simpleActivity
	for _, activity := range activitys {
		var temp simpleActivity
		temp.LoadFromModelActivity(activity)
		activityResponse = append(activityResponse, temp)
	}
	c.JSON(http.StatusOK, activitys)
}

// @Summary 修改活动
// @Description 修改指定活动
// @Tags 活动相关接口
// @Accept json
// @Produce json
// @Param id path integer true "活动id"
// @Param activity body models.Activity true "活动内容"
// @Param Authorization header string true "JWT Token"
// @Success 200 {array} simpleActivity
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/activity/{id} [post]
func UpdateActivity(c *gin.Context) {
	// 获取活动ID
	activityIdStr := c.Param("id")
	if activityIdStr == "" {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "activity id is required"})
		return
	}
	activityId, err := strconv.ParseInt(activityIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "Invalid activity id format"})
		return
	}
	// 验证JWT并获取用户Id
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

	dbActivity, err := controllers.GetActivityById(activityId)
	if err != nil {
		c.JSON(http.StatusNotFound, &models.ErrorResponse{ErrorMessage: "activity not found"})
		return
	}
	// 检查用户是否有权限修改活动
	if dbActivity.UserId != jwtUser.Id {
		c.JSON(http.StatusForbidden, &models.ErrorResponse{ErrorMessage: "you do not have permission to update this activity"})
		return
	}

	var activity models.Activity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "invalid request format"})
		return
	}

	activity.Id = 0
	activity.IfDelete = 0
	dbActivity.UpdateActivityFields(activity)

	// 更新活动信息
	if err := controllers.UpdateActivity(dbActivity); err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to update activity"})
		return
	}

	c.JSON(http.StatusOK, &simpleActivity{Id: activity.Id, Name: activity.Name})
}

// @Summary 创建活动
// @Description 创建新的活动
// @Tags 活动相关接口
// @Accept json
// @Produce json
// @Param activity body models.Activity true "活动内容"
// @Param Authorization header string true "JWT Token"
// @Success 200 {array} models.Activity
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/activity [put]
func CreateActivity(c *gin.Context) {
	// 验证JWT并获取用户Id
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
	var activity models.Activity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "invalid request format"})
		return
	}
	activity.UserId = jwtUser.Id
	activity.IfDelete = 0
	activity.Id = 0 // 确保新创建的活动没有ID
	// 创建活动
	if err := controllers.AddActivity(&activity); err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to create activity"})
		return
	}
	c.JSON(http.StatusOK, activity)
}

// @Summary 删除活动
// @Description 软删除指定活动
// @Tags 活动相关接口
// @Accept json
// @Produce json
// @Param id path integer true "活动id"
// @Param Authorization header string true "JWT Token"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/activity/{id} [delete]
func DeleteActivity(c *gin.Context) {
	// 获取活动ID
	activityIdStr := c.Param("id")
	if activityIdStr == "" {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "activity id is required"})
		return
	}
	activityId, err := strconv.ParseInt(activityIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "Invalid activity id format"})
		return
	}

	// 验证JWT并获取用户Id
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

	// 获取活动信息
	dbActivity, err := controllers.GetActivityById(activityId)
	if err != nil {
		c.JSON(http.StatusNotFound, &models.ErrorResponse{ErrorMessage: "activity not found"})
		return
	}

	// 检查用户是否有权限删除该活动
	if dbActivity.UserId != jwtUser.Id {
		c.JSON(http.StatusForbidden, &models.ErrorResponse{ErrorMessage: "you do not have permission to delete this activity"})
		return
	}

	if err = controllers.DeleteActivityById(activityId); err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "cannot delete activity"})
	}

	c.JSON(http.StatusOK, &models.SuccessResponse{SuccessMessage: "activity deleted successfully"})
}

// @Summary 获取活动成员
// @Description 获取指定活动所有成员
// @Tags 活动相关接口
// @Accept json
// @Produce json
// @Param id path integer true "活动id"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/activity/{id}/member [get]
func GetActivityMembers(c *gin.Context) {
	// 获取活动ID
	activityIdStr := c.Param("id")
	if activityIdStr == "" {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "activity id is required"})
		return
	}
	activityId, err := strconv.ParseInt(activityIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "Invalid activity id format"})
		return
	}

	// 获取活动成员
	members, err := controllers.GetActivityMembersByActivityId(activityId)
	if err != nil {
		c.JSON(http.StatusNotFound, &models.ErrorResponse{ErrorMessage: "activity members not found"})
		return
	}

	c.JSON(http.StatusOK, members)
}

// @Summary 加入活动
// @Description 加入指定活动
// @Tags 活动相关接口
// @Accept json
// @Produce json
// @Param id path integer true "活动id"
// @Param Authorization header string true "JWT Token"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/activity/{id}/member [put]
func JoinActivity(c *gin.Context) {
	// 获取活动ID
	activityIdStr := c.Param("id")
	if activityIdStr == "" {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "activity id is required"})
		return
	}
	activityId, err := strconv.ParseInt(activityIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "Invalid activity id format"})
		return
	}

	// 验证JWT并获取用户Id
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

	// 创建活动成员记录
	member := &models.ActivityMember{
		UserId:     jwtUser.Id,
		ActivityId: activityId,
		CreateTime: time.Now(),
	}

	if err := controllers.AddActivityMember(member); err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to join activity"})
		return
	}

	c.JSON(http.StatusOK, &models.SuccessResponse{SuccessMessage: "joined activity successfully"})
}

// @Summary 退出活动
// @Description 退出指定活动
// @Tags 活动相关接口
// @Accept json
// @Produce json
// @Param id path integer true "活动id"
// @Param Authorization header string true "JWT Token"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/activity/{id}/member [delete]
func LeaveActivity(c *gin.Context) {
	// 获取活动ID
	activityIdStr := c.Param("id")
	if activityIdStr == "" {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "activity id is required"})
		return
	}
	activityId, err := strconv.ParseInt(activityIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "Invalid activity id format"})
		return
	}

	// 验证JWT并获取用户Id
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

	if err := controllers.DeleteActivityMember(activityId, jwtUser.Id); err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to leave activity"})
		return
	}

	c.JSON(http.StatusOK, &models.SuccessResponse{SuccessMessage: "left activity successfully"})
}

// @Summary 获取用户参加的所有活动
// @Description 获取用户参加的所有活动
// @Tags 活动相关接口
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Success 200 {array} simpleActivity
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/activity/member [get]
func GetUserActivities(c *gin.Context) {
	// 验证JWT并获取用户Id
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
	// 获取用户参加的所有活动成员记录
	members, err := controllers.GetActivityMembersByUserId(jwtUser.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, &models.ErrorResponse{ErrorMessage: "activity members not found"})
		return
	}
	// 获取活动ID列表
	var activitys []simpleActivity
	for _, member := range members {
		tempActivitie, err := controllers.GetActivityById(member.ActivityId)
		if err != nil {
			c.JSON(http.StatusNotFound, &models.ErrorResponse{ErrorMessage: "activity not found"})
			return
		}
		temp := simpleActivity{}
		temp.LoadFromModelActivity(*tempActivitie)
		activitys = append(activitys, temp)
	}
	// 获取活动信息

	if err != nil {
		c.JSON(http.StatusNotFound, &models.ErrorResponse{ErrorMessage: "activities not found"})
		return
	}
	c.JSON(http.StatusOK, activitys)
}

// @Summary 获取活动评论
// @Description 获取指定活动的所有评论
// @Tags 活动相关接口
// @Accept json
// @Produce json
// @Param id path integer true "活动id"
// @Success 200 {array} models.ActivityComment
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/activity/{id}/comment [get]
func GetActivityComments(c *gin.Context) {
	// 获取活动ID
	activityIdStr := c.Param("id")
	if activityIdStr == "" {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "activity id is required"})
		return
	}
	activityId, err := strconv.ParseInt(activityIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "Invalid activity id format"})
		return
	}

	// 获取活动评论
	comments, err := controllers.GetActivityCommentsByActivityId(activityId)
	if err != nil {
		c.JSON(http.StatusNotFound, &models.ErrorResponse{ErrorMessage: "activity comments not found"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

type activityCommentResponse struct {
	Comment string `json:"comment"`
}

// @Summary 添加活动评论
// @Description 添加指定活动的评论
// @Tags 活动相关接口
// @Accept json
// @Produce json
// @Param id path integer true "活动id"
// @Param comment body activityCommentResponse true "评论内容"
// @Param Authorization header string true "JWT Token"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/activity/{id}/comment [put]
func AddActivityComment(c *gin.Context) {
	// 获取活动ID
	activityIdStr := c.Param("id")
	if activityIdStr == "" {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "activity id is required"})
		return
	}
	activityId, err := strconv.ParseInt(activityIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "Invalid activity id format"})
		return
	}

	// 验证JWT并获取用户Id
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

	var comment activityCommentResponse
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "invalid request format"})
		return
	}

	activityComment := models.ActivityComment{
		UserId:     jwtUser.Id,
		ActivityId: activityId,
		Content:    comment.Comment,
		CreateTime: time.Now(),
	}

	if err := controllers.AddActivityComment(&activityComment); err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to add comment"})
		return
	}

	c.JSON(http.StatusOK, &models.SuccessResponse{SuccessMessage: "comment added successfully"})
}

// @Summary 删除活动评论
// @Description 删除指定活动的指定评论
// @Tags 活动相关接口
// @Accept json
// @Produce json
// @Param commentId path integer true "评论id"
// @Param Authorization header string true "JWT Token"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/activity/comment/{commentId} [delete]
func DeleteActivityComment(c *gin.Context) {
	// 获取评论ID
	commentIdStr := c.Param("commentId")
	if commentIdStr == "" {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "comment id is required"})
		return
	}
	commentId, err := strconv.ParseInt(commentIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "Invalid comment id format"})
		return
	}

	// 验证JWT并获取用户Id
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
	// 获取评论信息
	_, err = controllers.GetActivityCommentsByActivityIdAndUserId(commentId, jwtUser.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, &models.ErrorResponse{ErrorMessage: "comment not found or you do not have permission to delete it"})
		return
	}
	if err = controllers.DeleteActivityComment(commentId); err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, &models.SuccessResponse{SuccessMessage: "comment deleted successfully"})
}
