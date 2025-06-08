package components

import (
	"fmt"
	"io/ioutil"

	"hobbyhub-server/config"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gopkg.in/yaml.v3"
)

type ConfigViewer struct {
	view       *tview.Flex
	configView *tview.TextView
	statusBar  *tview.TextView
	configPath string
	editMode   bool
}

func NewConfigViewer(configPath string) *ConfigViewer {
	cv := &ConfigViewer{
		configPath: configPath,
	}
	cv.setupViews()
	return cv
}

func (cv *ConfigViewer) setupViews() {
	// 创建配置显示区域
	cv.configView = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true).
		SetWordWrap(true)

	cv.configView.SetBorder(true).
		SetTitle("配置文件内容").
		SetTitleAlign(tview.AlignCenter)

	// 创建状态栏
	cv.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]按 'r' 刷新, 'e' 编辑模式, 's' 保存[-]")

	cv.statusBar.SetBorder(true).
		SetTitle("控制").
		SetTitleAlign(tview.AlignCenter)

	// 创建主布局
	cv.view = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(cv.configView, 0, 1, true).
		AddItem(cv.statusBar, 3, 1, false)

	// 设置按键处理
	cv.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'r':
			cv.Refresh()
			return nil
		case 'e':
			cv.toggleEditMode()
			return nil
		case 's':
			cv.saveConfig()
			return nil
		}
		return event
	})
}

func (cv *ConfigViewer) Refresh() {
	content := cv.loadConfigFile()
	cv.configView.SetText(content)
	cv.statusBar.SetText(fmt.Sprintf("[green]配置已刷新 - %s[-]", cv.configPath))
}

func (cv *ConfigViewer) loadConfigFile() string {
	// 首先尝试从配置对象获取结构化数据
	cfg := config.GetConfig()

	var content string
	content += "[yellow]HobbyHub 配置文件[-]\n"
	content += fmt.Sprintf("[gray]文件路径: %s[-]\n\n", cv.configPath)

	// 显示服务器配置
	content += "[green]🖥️ 服务器配置:[-]\n"
	content += fmt.Sprintf("  主机: %s\n", cfg.Server.Host)
	content += fmt.Sprintf("  端口: %d\n", cfg.Server.Port)
	content += "\n"

	// 显示数据库配置
	content += "[blue]🗄️ 数据库配置:[-]\n"
	content += fmt.Sprintf("  类型: %s\n", cfg.Database.Type)
	content += fmt.Sprintf("  主机: %s\n", cfg.Database.Host)
	content += fmt.Sprintf("  端口: %d\n", cfg.Database.Port)
	content += fmt.Sprintf("  数据库: %s\n", cfg.Database.Database)
	content += fmt.Sprintf("  用户名: %s\n", cfg.Database.Username)
	content += fmt.Sprintf("  密码: %s\n", cv.maskPassword(cfg.Database.Password))
	content += fmt.Sprintf("  字符集: %s\n", cfg.Database.Charset)
	content += "\n"

	// 显示认证配置
	content += "[purple]🔐 认证配置:[-]\n"
	content += fmt.Sprintf("  JWT密钥: %s\n", cv.maskPassword(cfg.Authentication.JwtSecret))
	content += "\n"

	// 显示文件配置
	content += "[cyan]📁 文件配置:[-]\n"
	content += fmt.Sprintf("  上传路径: %s\n", cfg.File.UploadPath)
	content += fmt.Sprintf("  最大大小: %d MB\n", cfg.File.MaxSize)
	content += fmt.Sprintf("  允许类型: %v\n", cfg.File.AllowedTypes)
	content += "\n"

	// 尝试读取原始文件内容
	if rawContent, err := cv.readRawConfigFile(); err == nil {
		content += "[yellow]📄 原始YAML内容:[-]\n"
		content += "[white]" + rawContent + "[-]\n"
	} else {
		content += fmt.Sprintf("[red]读取原始文件失败: %v[-]\n", err)
	}

	return content
}

func (cv *ConfigViewer) readRawConfigFile() (string, error) {
	data, err := ioutil.ReadFile(cv.configPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (cv *ConfigViewer) maskPassword(password string) string {
	if len(password) == 0 {
		return "[未设置]"
	}
	if len(password) <= 4 {
		return "****"
	}
	return password[:2] + "****" + password[len(password)-2:]
}

func (cv *ConfigViewer) toggleEditMode() {
	cv.editMode = !cv.editMode
	if cv.editMode {
		cv.statusBar.SetText("[red]编辑模式 - 按 's' 保存, 'e' 退出编辑[-]")
		cv.configView.SetTitle("配置文件编辑")

		// 在编辑模式下显示原始YAML内容
		if rawContent, err := cv.readRawConfigFile(); err == nil {
			cv.configView.SetText(rawContent)
		}
	} else {
		cv.statusBar.SetText("[yellow]按 'r' 刷新, 'e' 编辑模式, 's' 保存[-]")
		cv.configView.SetTitle("配置文件内容")
		cv.Refresh()
	}
}

func (cv *ConfigViewer) saveConfig() {
	if !cv.editMode {
		cv.statusBar.SetText("[red]请先进入编辑模式 (按 'e')[-]")
		return
	}

	// 获取编辑后的内容
	content := cv.configView.GetText(false)

	// 验证YAML格式
	var testConfig config.Config
	if err := yaml.Unmarshal([]byte(content), &testConfig); err != nil {
		cv.statusBar.SetText(fmt.Sprintf("[red]YAML格式错误: %v[-]", err))
		return
	}

	// 保存到文件
	if err := ioutil.WriteFile(cv.configPath, []byte(content), 0644); err != nil {
		cv.statusBar.SetText(fmt.Sprintf("[red]保存失败: %v[-]", err))
		return
	}

	// 重新加载配置
	if err := config.LoadConfig(cv.configPath); err != nil {
		cv.statusBar.SetText(fmt.Sprintf("[red]重新加载配置失败: %v[-]", err))
		return
	}

	cv.statusBar.SetText("[green]配置已保存并重新加载[-]")
	cv.editMode = false
	cv.configView.SetTitle("配置文件内容")
	cv.Refresh()
}

func (cv *ConfigViewer) GetView() tview.Primitive {
	return cv.view
}
