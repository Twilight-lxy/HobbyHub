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

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}
	// 这里可以根据 cfg 初始化数据库等
	config.InitDatabase(cfg)

	r := gin.Default()

	// 设置路由前缀 /api/v1
	apiV1 := r.Group("/api/v1")
	{
		// User routes
		apiV1.GET("/user/info/", api.GetUserInfo)
	}

	// Swagger 相关路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)) // 按配置文件端口启动
}
