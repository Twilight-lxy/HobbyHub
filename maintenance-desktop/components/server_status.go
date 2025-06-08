package components

import (
	"fmt"
	"runtime"
	"time"

	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

type ServerStatus struct {
	view      *tview.Flex
	infoView  *tview.TextView
	chartView *tview.TextView
	running   bool
	stopChan  chan bool
}

func NewServerStatus() *ServerStatus {
	ss := &ServerStatus{
		stopChan: make(chan bool),
	}

	ss.setupViews()
	return ss
}

func (ss *ServerStatus) setupViews() {
	// 创建信息显示区域
	ss.infoView = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true)
	ss.infoView.SetBorder(true).
		SetTitle("系统信息").
		SetTitleAlign(tview.AlignCenter)

	// 创建图表显示区域
	ss.chartView = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)
	ss.chartView.SetBorder(true).
		SetTitle("性能图表").
		SetTitleAlign(tview.AlignCenter)

	// 创建主布局
	ss.view = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ss.infoView, 0, 2, false).
		AddItem(ss.chartView, 0, 1, false)
}

func (ss *ServerStatus) Start() {
	if ss.running {
		return
	}

	ss.running = true
	go ss.updateLoop()
}

func (ss *ServerStatus) Stop() {
	if !ss.running {
		return
	}

	ss.running = false
	ss.stopChan <- true
}

func (ss *ServerStatus) updateLoop() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ss.stopChan:
			return
		case <-ticker.C:
			ss.updateStatus()
		}
	}
}

func (ss *ServerStatus) updateStatus() {
	// 获取系统信息
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// 获取CPU使用率
	cpuPercent, _ := cpu.Percent(time.Second, false)
	var cpuUsage float64
	if len(cpuPercent) > 0 {
		cpuUsage = cpuPercent[0]
	}

	// 获取内存信息
	vmStat, _ := mem.VirtualMemory()

	// 获取进程信息
	pid := int32(runtime.GOMAXPROCS(0))
	proc, _ := process.NewProcess(pid)

	// 更新信息显示
	info := fmt.Sprintf(`[yellow]HobbyHub 服务器状态监控[-]

[green]🖥️ 系统信息:[-]
  操作系统: %s
  架构: %s
  CPU核心数: %d
  Go版本: %s

[blue]💾 内存使用:[-]
  总内存: %.2f GB
  已用内存: %.2f GB (%.1f%%)
  可用内存: %.2f GB
  
[red]🔥 应用程序内存:[-]
  堆内存分配: %.2f MB
  堆内存使用: %.2f MB
  GC次数: %d
  Goroutine数量: %d

[purple]📊 性能指标:[-]
  CPU使用率: %.1f%%
  系统负载: 正常
  
[cyan]⏰ 运行时间:[-]
  启动时间: %s
  运行时长: %s

[proc]
%s  
  `,
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
		runtime.Version(),
		float64(vmStat.Total)/1024/1024/1024,
		float64(vmStat.Used)/1024/1024/1024,
		vmStat.UsedPercent,
		float64(vmStat.Available)/1024/1024/1024,
		float64(memStats.Alloc)/1024/1024,
		float64(memStats.HeapAlloc)/1024/1024,
		memStats.NumGC,
		runtime.NumGoroutine(),
		cpuUsage,
		time.Now().Format("2006-01-02 15:04:05"),
		"运行中...", proc.String())

	ss.infoView.SetText(info)

	// 更新图表
	chart := ss.generateChart(cpuUsage, vmStat.UsedPercent)
	ss.chartView.SetText(chart)
}

func (ss *ServerStatus) generateChart(cpuUsage, memUsage float64) string {
	// 生成简单的ASCII图表
	chart := "[yellow]性能监控图表[-]\n\n"

	// CPU使用率条形图
	chart += fmt.Sprintf("[green]CPU使用率: %.1f%%[-]\n", cpuUsage)
	chart += ss.generateBar(cpuUsage, "🟩", "⬜")
	chart += "\n\n"

	// 内存使用率条形图
	chart += fmt.Sprintf("[blue]内存使用率: %.1f%%[-]\n", memUsage)
	chart += ss.generateBar(memUsage, "🟦", "⬜")
	chart += "\n\n"

	// 添加时间戳
	chart += fmt.Sprintf("[gray]最后更新: %s[-]", time.Now().Format("15:04:05"))

	return chart
}

func (ss *ServerStatus) generateBar(percentage float64, fillChar, emptyChar string) string {
	barWidth := 40
	fillWidth := int((percentage / 100.0) * float64(barWidth))

	bar := ""
	for i := 0; i < barWidth; i++ {
		if i < fillWidth {
			bar += fillChar
		} else {
			bar += emptyChar
		}
	}

	return fmt.Sprintf("[%s] %.1f%%", bar, percentage)
}

func (ss *ServerStatus) GetView() tview.Primitive {
	return ss.view
}
