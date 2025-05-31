package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"hobbyhub-server/controllers"
	"hobbyhub-server/models"
	utils "hobbyhub-server/utils"

	"github.com/gin-gonic/gin"
)

// @Summary 获取用户信息
// @Description 通过用户ID获取用户信息；可选用户id或用户名查询，优先使用用户id；不填写id或用户名，使用jwt token获取
// @Tags 用户相关接口
// @Produce json
// @Param id query int false "用户ID"
// @Param username query string false "用户名"
// @Param Authorization header string false "JWT Token"
// @Success 200 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/user [get]
func GetUserInfo(c *gin.Context) {
	idStr := c.Query("id")
	usernameStr := c.Query("username")
	jwtToken := c.GetHeader("Authorization")
	var user *models.User
	var err error
	var tokenIsValid bool = false
	if idStr == "" && usernameStr == "" && jwtToken == "" {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "id or username or Authorization parameter is required"})
		return
	} else if idStr == "" && usernameStr == "" && jwtToken != "" {
		// 如果没有提供 id 和 username，但有 JWT Token，则使用 JWT Token 获取用户信息
		log.Println("GetUserInfo called with JWT token")
		jwtUser, err := utils.ParseJWT(jwtToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "invalid jwt token"})
			return
		}
		tokenIsValid = true
		user, err = controllers.GetUserByUserId(jwtUser.Id)
		if err != nil {
			c.JSON(http.StatusNotFound, &models.ErrorResponse{ErrorMessage: "user not found"})
			return
		}
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

	if jwtToken != "" && !tokenIsValid {
		// 如果有 JWT Token，验证用户身份
		jwtUser, err := utils.ParseJWT(jwtToken)
		if err == nil {
			if jwtUser.Id == user.Id {
				// 如果 JWT Token 验证通过，设置 tokenIsValid 为 true
				tokenIsValid = true
			}
		} else {
			log.Println("Invalid JWT token:", err)
		}
	}
	if !tokenIsValid {
		// 如果没有 JWT Token，返回部分用户信息
		user.Username = ""
		user.Addr = ""
		user.CreateTime = time.Time{} // 不返回创建时间
		user.Lat = 0
		user.Lon = 0
	}
	user.Password = "" // 不返回密码
	c.JSON(http.StatusOK, user)
}

type UsernameAndPassword struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary 用户登录
// @Description 用户名密码登录
// @Tags 用户相关接口
// @Accept json
// @Produce json
// @Param loginRequest body UsernameAndPassword true "登录请求体，包含用户名和密码"
// @Success 200 {string} string "JWT Token"
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/login [post]
func UserLogin(c *gin.Context) {
	var req UsernameAndPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "username and password are required"})
		return
	}
	dbUser, err := controllers.GetUserByUserName(req.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, &models.ErrorResponse{ErrorMessage: "user not found or invalid credentials"})
		return
	}
	if dbUser.Password != req.Password {
		c.JSON(http.StatusNotFound, &models.ErrorResponse{ErrorMessage: "username or password is incorrect"})
		return
	}
	// 设置 JWT Token
	jwtToken, err := utils.GenerateJWT(dbUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to generate JWT token"})
		return
	}
	c.JSON(http.StatusOK, jwtToken)
}

// @Summary 用户注册
// @Description 用户名密码注册
// @Tags 用户相关接口
// @Accept json
// @Produce json
// @Param loginRequest body UsernameAndPassword true "注册请求体，包含用户名和密码"
// @Success 200 {string} string "JWT Token"
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/user [put]
func UserRegister(c *gin.Context) {
	var req UsernameAndPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "username and password are required"})
		return
	}
	// 检查用户名是否已存在
	existingUser, err := controllers.GetUserByUserName(req.Username)
	if err == nil && existingUser != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "username already exists"})
		return
	}
	// 创建新用户
	newUser := &models.User{
		Username:   req.Username,
		Password:   req.Password,
		CreateTime: time.Now(),
	}
	if err := controllers.AddUser(newUser); err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to create user"})
		return
	}
	// 设置 JWT Token
	jwtToken, err := utils.GenerateJWT(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to generate JWT token"})
		return
	}
	c.JSON(http.StatusOK, jwtToken)
}

// @Summary 更新用户信息
// @Description 更新用户信息
// @Tags 用户相关接口
// @Accept json
// @Produce json
// @Param userInfo body models.User true "用户数据，必须包含id"
// @Param Authorization header string true "JWT Token"
// @Success 200 {object} models.User "修改后的用户信息"
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/user [post]
func UpdateUserInfo(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "invalid user data"})
		return
	}
	jwtToken := c.GetHeader("Authorization")
	if jwtToken == "" {
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{ErrorMessage: "jwt token is required"})
		return
	}
	jwtUser, err := utils.ParseJWT(jwtToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &models.ErrorResponse{ErrorMessage: "unauthorized"})
		return
	}
	jwtUser.UpdateUserFields(user)
	if err := controllers.UpdateUser(*jwtUser); err != nil {
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{ErrorMessage: "failed to update user"})
		return
	}
	jwtUser.Password = "" // 不返回密码
	c.JSON(http.StatusOK, jwtUser)
}
