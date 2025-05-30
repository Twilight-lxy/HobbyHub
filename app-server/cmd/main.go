package main

import (
	"flag"
	"fmt"
	"hobbyhub-server/api"
	"hobbyhub-server/config"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "hobbyhub-server/docs"
)

// @title HobbyHub API
// @version 1.0
// @description HobbyHub 后端 API 文档
// @host localhost:8081
// @BasePath /api

func main() {
	// 支持命令行指定配置文件路径，默认 config.yaml
	configPath := flag.String("config", "config.yaml", "配置文件路径")
	flag.Parse()

	if err := config.LoadConfig(*configPath); err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}
	// 这里可以根据 cfg 初始化数据库等
	config.InitDatabase(config.GetConfig())

	r := gin.Default()

	// 设置路由前缀 /api/v1
	apiV1 := r.Group("/api/v1")
	{
		// User routes
		user := apiV1.Group("/user")
		{
			user.GET("/info/", api.GetUserInfo)
			user.POST("/login/", api.UserLogin)
			user.POST("/register/", api.UserRegister)
			user.POST("/update/", api.UpdateUserInfo)
		}
		// Chat routes
		chat := apiV1.Group("/chat")
		{
			chat.GET("/history", api.GetChatHistory) // 获取聊天记录
			chat.POST("/send", api.SendChat)         // 发送聊天消息
		}
	}

	// Swagger 相关路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(fmt.Sprintf("%s:%d", config.GetConfig().Server.Host, config.GetConfig().Server.Port)) // 按配置文件端口启动
}
