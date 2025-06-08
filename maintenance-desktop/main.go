package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rivo/tview"

	"hobbyhub-maintenance-desktop/components"

	"hobbyhub-server/config"
)

func main() {
	// 支持命令行指定配置文件路径
	serverExePath := flag.String("server", "./hobbyhub-server.exe", "服务器可执行文件路径")
	configPath := flag.String("config", "", "配置文件路径 (默认为服务器可执行文件同级目录的config.yaml)")
	flag.Parse()
	// 检查服务器可执行文件是否存在
	if _, err := os.Stat(*serverExePath); os.IsNotExist(err) {
		log.Fatalf("服务器可执行文件不存在: %s", *serverExePath)
	}
	// 如果指定了服务器可执行文件路径，则获取其目录
	serverDir := *serverExePath
	if stat, err := os.Stat(serverDir); err == nil && stat.IsDir() {
		// 如果是目录，则使用该目录
		serverDir = serverDir
	} else {
		// 否则获取文件所在目录
		serverDir = serverDir[:len(serverDir)-len(stat.Name())]
	}
	// 如果未指定配置文件路径，则使用服务器可执行文件同级目录的config.yaml
	if *configPath == "" {
		*configPath = fmt.Sprintf("%sconfig.yaml", serverDir)
	}

	// 加载配置
	if err := config.LoadConfig(*configPath); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库连接
	if err := config.InitDatabase(config.GetConfig()); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 创建 TUI 应用
	app := tview.NewApplication()

	// 创建主界面
	ui := components.NewMainUI(app, *configPath)

	// 设置根组件并运行
	if err := app.SetRoot(ui.GetRoot(), true).EnableMouse(true).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "应用程序运行错误: %v\n", err)
		os.Exit(1)
	}
}
