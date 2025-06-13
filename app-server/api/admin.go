package api

import (
	"hobbyhub-server/controllers"
	"hobbyhub-server/models"
	"hobbyhub-server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 管理员登录
// @Description 用于管理员的用户名密码登录
// @Tags 管理相关接口
// @Accept json
// @Produce json
// @Param loginRequest body UsernameAndPassword true "登录请求体，包含用户名和密码"
// @Success 200 {object} JWTResponse "JWT Token"
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/admin/login [post]
func AdminLogin(c *gin.Context) {
	var loginRequest UsernameAndPassword
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Invalid request"})
		return
	}
	admin, err := controllers.GetAdminByUserName(loginRequest.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{ErrorMessage: "Admin not found"})
		return
	}
	if !utils.CheckPasswordHash(loginRequest.Password, admin.Password) {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{ErrorMessage: "Invalid password"})
		return
	}
	token, err := utils.GenerateAdminJWT(admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, JWTResponse{Token: token})
}

// @Summary 获取所有用户
// @Description 获取所有用户信息
// @Tags 管理相关接口
// @Accept json
// @Produce json
// @Param Authorization header string true "Admin JWT Token"
// @Param page query int false "页码，默认为1"
// @Param pageSize query int false "每页数量，默认为10"
// @Success 200 {array} models.User "用户列表"
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /v1/admin/users [get]
func GetAllUsers(c *gin.Context) {
	jwtToken := c.GetHeader("Authorization")
	if jwtToken == "" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{ErrorMessage: "Authorization header is required"})
		return
	}
	admin, err := utils.ParseAdminJWT(jwtToken)
	if err != nil || admin == nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{ErrorMessage: "Invalid token"})
		return
	}

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")

	users, err := controllers.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: "Failed to retrieve users"})
		return
	}
	if len(users) == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{ErrorMessage: "No users found"})
		return
	}
	// 分页处理
	pageInt, err := utils.StringToInt(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Invalid page number"})
		return
	}
	pageSizeInt, err := utils.StringToInt(pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Invalid page size"})
		return
	}
	if pageInt < 1 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Page number must be greater than 0"})
		return
	}
	if pageSizeInt < 1 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Page size must be greater than 0"})
		return
	}
	start := (pageInt - 1) * pageSizeInt
	if start >= len(users) {
		c.JSON(http.StatusOK, []models.User{})
		return
	}
	end := start + pageSizeInt
	if end > len(users) {
		end = len(users)
	}
	c.JSON(http.StatusOK, users[start:end])
}

// @Summary 获取所有活动
// @Description 获取所有活动信息
// @Tags 管理相关接口
// @Accept json
// @Produce json
// @Param Authorization header string true "Admin JWT Token"
// @Param page query int false "页码，默认为1"
// @Param pageSize query int false "每页数量，默认为10"
// @Success 200 {array} models.Activity "活动列表"
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func GetAllActivities(c *gin.Context) {
	jwtToken := c.GetHeader("Authorization")
	if jwtToken == "" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{ErrorMessage: "Authorization header is required"})
		return
	}
	admin, err := utils.ParseAdminJWT(jwtToken)
	if err != nil || admin == nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{ErrorMessage: "Invalid token"})
		return
	}

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")

	activities, err := controllers.GetAllActivities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: "Failed to retrieve activities"})
		return
	}
	if len(activities) == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{ErrorMessage: "No activities found"})
		return
	}
	// 分页处理
	pageInt, err := utils.StringToInt(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Invalid page number"})
		return
	}
	pageSizeInt, err := utils.StringToInt(pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Invalid page size"})
		return
	}
	if pageInt < 1 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Page number must be greater than 0"})
		return
	}
	if pageSizeInt < 1 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Page size must be greater than 0"})
		return
	}
	start := (pageInt - 1) * pageSizeInt
	if start >= len(activities) {
		c.JSON(http.StatusOK, []models.Activity{})
		return
	}
	end := start + pageSizeInt
	if end > len(activities) {
		end = len(activities)
	}
	c.JSON(http.StatusOK, activities[start:end])
}

// @Summary 创建活动
// @Description 创建一个新的活动
// @Tags 管理相关接口
// @Accept json
// @Produce json
// @Param Authorization header string true "Admin JWT Token"
// @Param activity body models.Activity true "活动内容"
// @Success 200 {object} models.Activity
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /v1/admin/activity [put]
func AdminCreateActivity(c *gin.Context) {
	jwtToken := c.GetHeader("Authorization")
	if jwtToken == "" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{ErrorMessage: "Authorization header is required"})
		return
	}
	admin, err := utils.ParseAdminJWT(jwtToken)
	if err != nil || admin == nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{ErrorMessage: "Invalid token"})
		return
	}

	var activity models.Activity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{ErrorMessage: "Invalid request format"})
		return
	}
	activity.Id = 0                              // 确保ID为0，表示新建活动
	activity.CreateTime = utils.GetCurrentTime() // 设置创建时间

	if err := controllers.AddActivity(&activity); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{ErrorMessage: "Failed to create activity"})
		return
	}

	c.JSON(http.StatusOK, activity)
}
