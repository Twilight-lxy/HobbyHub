package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"hobbyhub-server/controllers"
	"hobbyhub-server/models"
)

// @Summary 获取活动信息
// @Description 获取聊天消息
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

// @Summary 获取所有活动Id
// @Description 获取聊天消息
// @Tags 活动相关接口
// @Accept json
// @Produce json
// @Success 200 {array} number
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/activity [get]
func GetAllActivitie(c *gin.Context) {}
