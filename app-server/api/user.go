package api

import (
	"log"
	"net/http"
	"strconv"

	"hobbyhub-server/controllers"
	"hobbyhub-server/models"

	"github.com/gin-gonic/gin"
)

// @Summary 获取用户信息
// @Description 通过用户ID获取用户信息，可选用户id或用户名查询，优先使用用户id
// @Tags 用户相关接口
// @Produce json
// @Param id query int false "用户ID"
// @Param username query string false "用户名"
// @Success 200 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/user/info [get]
func GetUserInfo(c *gin.Context) {
	idStr := c.Query("id")
	usernameStr := c.Query("username")
	jwtToken, err := c.Cookie("jwt_token")
	if err != nil {
		jwtToken = ""
	}
	var user *models.User
	if idStr == "" && usernameStr == "" {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "id or username parameter is required"})
		return
	} else if idStr != "" {
		log.Println("GetUserInfo called with id:", idStr)
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "invalid user id"})
			return
		}
		user, err = controllers.GetUserByUserId(id)
		if err != nil {
			c.JSON(http.StatusNotFound, &models.ErrorResponse{ErrorMessage: "user not found"})
			return
		}
	} else {
		log.Println("GetUserInfo called with username:", usernameStr)
		user, err = controllers.GetUserByUserName(usernameStr)
		if err != nil {
			c.JSON(http.StatusNotFound, &models.ErrorResponse{ErrorMessage: "user not found"})
			return
		}
	}

	if jwtToken != "" {
		// 如果有 JWT Token，验证用户身份
		log.Println("JWT Token is", jwtToken)
	} else {
		// 如果没有 JWT Token，返回部分用户信息
		user.Password = "" // 不返回密码
		c.SetCookie("jwt_token", "token_test", 3600, "/", "", false, true)
	}
	c.JSON(http.StatusOK, user)
}
