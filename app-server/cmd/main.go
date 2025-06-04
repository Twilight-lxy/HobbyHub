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
		fmt.Printf("加载配置失败: %v", err)
		fmt.Scanln()
		return
	}
	// 这里可以根据 cfg 初始化数据库等
	err := config.InitDatabase(config.GetConfig())
	if err != nil {
		fmt.Printf("初始化数据库失败: %v", err)
		fmt.Scanln()
		return
	}

	r := gin.Default()

	// 设置路由前缀 /api/v1
	apiV1 := r.Group("/api/v1")
	{
		//login
		apiV1.POST("/login", api.UserLogin)
		// User routes
		user := apiV1.Group("/user")
		{
			user.GET("/", api.GetUserInfo)
			user.PUT("/", api.UserRegister)
			user.POST("/", api.UpdateUserInfo)
		}
		// Chat routes
		chat := apiV1.Group("/chat")
		{
			chat.GET("/", api.GetChatHistory)   // 获取聊天记录
			chat.POST("/", api.SendChat)        // 发送聊天消息
			chat.DELETE("/:id", api.DeleteChat) // 删除聊天记录
		}
		// Friend routes
		friend := apiV1.Group("/friend")
		{
			friend.GET("/", api.GetFriendList)      // 获取好友列表
			friend.POST("/", api.SendFriendRequest) // 发送好友申请
			friend.PUT("/", api.UpdateFriendStatus) // 更新好友申请状态
			friend.DELETE("/:id", api.DeleteFriend) // 删除好友
		}
		// File routes
		file := apiV1.Group("/file")
		{
			file.POST("/", api.UploadFile)      // 上传文件
			file.GET("/:id", api.DownloadFile)  // 下载文件
			file.DELETE("/:id", api.DeleteFile) // 删除文件
		}
		// Activity routes
		activity := apiV1.Group("/activity")
		{
			activity.GET("/:id", api.GetActivitie)      // 获取活动列表
			activity.GET("/", api.GetAllActivitie)      // 获取活动详情
			activity.PUT("/", api.CreateActivity)       // 新建活动
			activity.POST("/:id", api.UpdateActivity)   // 更新活动信息
			activity.DELETE("/:id", api.DeleteActivity) // 软删除活动
		}
	}

	// Swagger 相关路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(fmt.Sprintf("%s:%d", config.GetConfig().Server.Host, config.GetConfig().Server.Port)) // 按配置文件端口启动
}
