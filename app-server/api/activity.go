package api

import (
	"net/http"
	"strconv"

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
// @Success 200 {array} simpleActivity
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
	c.JSON(http.StatusOK, activityResponse)
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
