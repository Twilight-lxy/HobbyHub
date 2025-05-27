package api

import (
	"log"
	"net/http"
	"strconv"

	"hobbyhub-server/controllers"

	"github.com/gin-gonic/gin"
)

// @Summary 获取用户信息
// @Description 通过用户ID获取用户信息
// @Tags 用户相关接口
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} models.User
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /v1/user/info/{id} [get]
func GetUserInfo(c *gin.Context) {
	idStr := c.Param("id")
	jwtToken, err := c.Cookie("jwt_token")
	if err != nil {
		jwtToken = ""
	}
	log.Println("GetUserInfo called with id:", idStr)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	user, err := controllers.GetUserByUserId(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
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
