package components

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type LogViewer struct {
	view      *tview.Flex
	logView   *tview.TextView
	statusBar *tview.TextView
	running   bool
	stopChan  chan bool
	logFile   string
}

func NewLogViewer() *LogViewer {
	lv := &LogViewer{
		stopChan: make(chan bool),
		logFile:  "app.log", // 默认日志文件
	}

	lv.setupViews()
	return lv
}

func (lv *LogViewer) setupViews() {
	// 创建日志显示区域
	lv.logView = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			// 自动滚动到底部
			lv.logView.ScrollToEnd()
		})

	lv.logView.SetBorder(true).
		SetTitle("应用程序日志").
		SetTitleAlign(tview.AlignCenter)

	// 创建状态栏
	lv.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]按 'c' 清空日志, 'r' 刷新, 'f' 跟踪日志[-]")

	lv.statusBar.SetBorder(true).
		SetTitle("控制").
		SetTitleAlign(tview.AlignCenter)

	// 创建主布局
	lv.view = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(lv.logView, 0, 1, true).
		AddItem(lv.statusBar, 3, 1, false)

	// 设置按键处理
	lv.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'c':
			lv.clearLogs()
			return nil
		case 'r':
			lv.refreshLogs()
			return nil
		case 'f':
			lv.toggleFollow()
			return nil
		}
		return event
	})
}

func (lv *LogViewer) Start() {
	if lv.running {
		return
	}

	lv.running = true
	lv.refreshLogs()
	go lv.followLogs()
}

func (lv *LogViewer) Stop() {
	if !lv.running {
		return
	}

	lv.running = false
	lv.stopChan <- true
}

func (lv *LogViewer) refreshLogs() {
	logs := lv.readLogFile()
	lv.logView.SetText(logs)
}

func (lv *LogViewer) readLogFile() string {
	// 尝试读取多个可能的日志文件位置
	logFiles := []string{
		"app.log",
		"../app.log",
		"logs/app.log",
		"../logs/app.log",
	}

	var content strings.Builder
	found := false

	for _, logFile := range logFiles {
		if file, err := os.Open(logFile); err == nil {
			defer file.Close()
			found = true

			scanner := bufio.NewScanner(file)
			lineCount := 0
			maxLines := 1000 // 只显示最后1000行

			var lines []string
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
				lineCount++
			}

			// 只显示最后的maxLines行
			start := 0
			if lineCount > maxLines {
				start = lineCount - maxLines
			}

			for i := start; i < len(lines); i++ {
				line := lines[i]
				coloredLine := lv.colorizeLogLine(line)
				content.WriteString(coloredLine + "\n")
			}
			break
		}
	}

	if !found {
		content.WriteString("[red]未找到日志文件[-]\n")
		content.WriteString("[gray]尝试查找的位置:[-]\n")
		for _, logFile := range logFiles {
			content.WriteString(fmt.Sprintf("  - %s\n", logFile))
		}
		content.WriteString("\n[yellow]提示: 启动应用程序后会生成日志文件[-]\n")

		// 添加一些模拟日志用于演示
		content.WriteString(lv.generateSampleLogs())
	}

	return content.String()
}

func (lv *LogViewer) colorizeLogLine(line string) string {
	line = strings.TrimSpace(line)
	if line == "" {
		return line
	}

	// 根据日志级别添加颜色
	if strings.Contains(line, "[ERROR]") || strings.Contains(line, "ERROR") {
		return "[red]" + line + "[-]"
	} else if strings.Contains(line, "[WARN]") || strings.Contains(line, "WARN") {
		return "[yellow]" + line + "[-]"
	} else if strings.Contains(line, "[INFO]") || strings.Contains(line, "INFO") {
		return "[green]" + line + "[-]"
	} else if strings.Contains(line, "[DEBUG]") || strings.Contains(line, "DEBUG") {
		return "[cyan]" + line + "[-]"
	} else if strings.Contains(line, "GET") || strings.Contains(line, "POST") || strings.Contains(line, "PUT") || strings.Contains(line, "DELETE") {
		return "[blue]" + line + "[-]"
	}

	return "[white]" + line + "[-]"
}

func (lv *LogViewer) generateSampleLogs() string {
	now := time.Now()
	logs := fmt.Sprintf(`
[green][INFO] %s HobbyHub 服务器启动[-]
[blue][INFO] %s 数据库连接成功[-]
[cyan][DEBUG] %s 加载配置文件: config.yaml[-]
[blue][INFO] %s GET /api/v1/user - 200 OK[-]
[blue][INFO] %s POST /api/v1/login - 200 OK[-]
[yellow][WARN] %s 用户登录失败: 密码错误[-]
[blue][INFO] %s GET /api/v1/activity - 200 OK[-]
[green][INFO] %s 新用户注册: user123[-]
[blue][INFO] %s PUT /api/v1/activity - 201 Created[-]
[cyan][DEBUG] %s JWT token 验证成功[-]
[blue][INFO] %s GET /swagger/index.html - 200 OK[-]
[green][INFO] %s 服务器运行在 localhost:8081[-]`,
		now.Add(-10*time.Minute).Format("2006/01/02 15:04:05"),
		now.Add(-9*time.Minute).Format("2006/01/02 15:04:05"),
		now.Add(-8*time.Minute).Format("2006/01/02 15:04:05"),
		now.Add(-7*time.Minute).Format("2006/01/02 15:04:05"),
		now.Add(-6*time.Minute).Format("2006/01/02 15:04:05"),
		now.Add(-5*time.Minute).Format("2006/01/02 15:04:05"),
		now.Add(-4*time.Minute).Format("2006/01/02 15:04:05"),
		now.Add(-3*time.Minute).Format("2006/01/02 15:04:05"),
		now.Add(-2*time.Minute).Format("2006/01/02 15:04:05"),
		now.Add(-1*time.Minute).Format("2006/01/02 15:04:05"),
		now.Format("2006/01/02 15:04:05"),
		now.Format("2006/01/02 15:04:05"))

	return logs
}

func (lv *LogViewer) followLogs() {
	// 实现日志跟踪功能
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-lv.stopChan:
			return
		case <-ticker.C:
			if lv.running {
				lv.refreshLogs()
			}
		}
	}
}

func (lv *LogViewer) clearLogs() {
	lv.logView.SetText("[yellow]日志已清空[-]\n")
}

func (lv *LogViewer) toggleFollow() {
	// 切换日志跟踪状态
	if lv.running {
		lv.statusBar.SetText("[red]日志跟踪已停止 - 按 'f' 重新开始[-]")
	} else {
		lv.statusBar.SetText("[green]日志跟踪已开始 - 按 'f' 停止[-]")
	}
}

func (lv *LogViewer) GetView() tview.Primitive {
	return lv.view
}
