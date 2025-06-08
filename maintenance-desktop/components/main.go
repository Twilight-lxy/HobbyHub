package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MainUI struct {
	app        *tview.Application
	configPath string

	// 布局组件
	root       *tview.Flex
	navigation *tview.List
	content    *tview.Flex

	// 内容页面
	serverStatus *ServerStatus
	logViewer    *LogViewer
	database     *DatabaseViewer
	configViewer *ConfigViewer

	// 当前选中的页面
	currentPage int
}

func NewMainUI(app *tview.Application, configPath string) *MainUI {
	ui := &MainUI{
		app:        app,
		configPath: configPath,
	}

	ui.setupComponents()
	ui.setupLayout()
	ui.setupNavigation()

	return ui
}

func (ui *MainUI) setupComponents() {
	// 创建各个组件
	ui.serverStatus = NewServerStatus()
	ui.logViewer = NewLogViewer()
	ui.database = NewDatabaseViewer(ui.app)
	ui.configViewer = NewConfigViewer(ui.configPath)

	// 创建导航栏
	ui.navigation = tview.NewList().
		AddItem("🖥️  服务器状态", "查看CPU、内存等系统信息", '1', func() {
			ui.showPage(0)
		}).
		AddItem("📋 日志", "查看应用程序日志", '2', func() {
			ui.showPage(1)
		}).
		AddItem("🗄️  数据库", "查看数据库内容", '3', func() {
			ui.showPage(2)
		}).
		AddItem("⚙️  配置", "查看和编辑配置文件", '4', func() {
			ui.showPage(3)
		}).
		AddItem("❌ 退出", "退出应用程序", 'q', func() {
			ui.app.Stop()
		})

	ui.navigation.SetBorder(true).
		SetTitle("导航").
		SetTitleAlign(tview.AlignCenter)

	// 创建内容区域
	ui.content = tview.NewFlex().SetDirection(tview.FlexRow)
	ui.content.SetBorder(true).
		SetTitle("内容").
		SetTitleAlign(tview.AlignCenter)
}

func (ui *MainUI) setupLayout() {
	// 创建主布局：左侧导航，右侧内容
	ui.root = tview.NewFlex().
		AddItem(ui.navigation, 25, 1, true). // 导航栏固定宽度25
		AddItem(ui.content, 0, 3, false)     // 内容区域自适应
}

func (ui *MainUI) setupNavigation() {
	// 设置按键处理
	ui.root.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			ui.app.Stop()
			return nil
		case tcell.KeyTab:
			// Tab 键在导航栏和内容区域之间切换焦点
			if ui.app.GetFocus() == ui.navigation {
				ui.app.SetFocus(ui.content)
			} else {
				ui.app.SetFocus(ui.navigation)
			}
			return nil
		}
		return event
	})

	// 默认显示服务器状态页面
	ui.showPage(0)
}

func (ui *MainUI) showPage(pageIndex int) {
	ui.currentPage = pageIndex
	ui.content.Clear()

	switch pageIndex {
	case 0: // 服务器状态
		ui.content.AddItem(ui.serverStatus.GetView(), 0, 1, false)
		ui.serverStatus.Start() // 开始更新数据
	case 1: // 日志
		ui.content.AddItem(ui.logViewer.GetView(), 0, 1, false)
		ui.logViewer.Start()
	case 2: // 数据库
		ui.content.AddItem(ui.database.GetView(), 0, 1, false)
		ui.database.Refresh()
	case 3: // 配置
		ui.content.AddItem(ui.configViewer.GetView(), 0, 1, false)
		ui.configViewer.Refresh()
	}

	// 停止其他页面的更新
	ui.stopOtherPages(pageIndex)
}

func (ui *MainUI) stopOtherPages(currentPage int) {
	if currentPage != 0 {
		ui.serverStatus.Stop()
	}
	if currentPage != 1 {
		ui.logViewer.Stop()
	}
}

func (ui *MainUI) GetRoot() tview.Primitive {
	return ui.root
}
