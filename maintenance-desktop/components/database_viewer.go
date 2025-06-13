package components

import (
	"fmt"
	"strconv"

	"hobbyhub-server/controllers"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type DatabaseViewer struct {
	view         *tview.Flex
	tableList    *tview.List
	tableView    *tview.Table
	statusBar    *tview.TextView
	currentTable string
	app          *tview.Application
}

func NewDatabaseViewer(app *tview.Application) *DatabaseViewer {
	dv := &DatabaseViewer{app: app}
	dv.setupViews()
	return dv
}

func (dv *DatabaseViewer) setupViews() {
	// 创建表列表
	dv.tableList = tview.NewList().
		AddItem("👥 用户 (user)", "查看用户信息", '1', func() {
			dv.showTable("user")
		}).
		AddItem("🎯 活动 (activity)", "查看活动信息", '2', func() {
			dv.showTable("activity")
		}).
		AddItem("💬 聊天 (chat)", "查看聊天记录", '3', func() {
			dv.showTable("chat")
		}).
		AddItem("👫 好友 (friend)", "查看好友关系", '4', func() {
			dv.showTable("friend")
		}).
		AddItem("📝 活动评论 (activity_comment)", "查看活动评论", '5', func() {
			dv.showTable("activity_comment")
		}).
		AddItem("👥 活动成员 (activity_member)", "查看活动成员", '6', func() {
			dv.showTable("activity_member")
		}).
		AddItem("📁 文件 (file)", "查看文件信息", '7', func() {
			dv.showTable("file")
		}).
		AddItem("👮 管理员 (admin)", "查看管理员信息", '8', func() {
			dv.showTable("admin")
		}).
		AddItem("🔄 刷新", "刷新当前表", 'r', func() {
			dv.Refresh()
		})

	dv.tableList.SetBorder(true).
		SetTitle("数据表").
		SetTitleAlign(tview.AlignCenter)

	// 创建表格显示区域
	dv.tableView = tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)

	dv.tableView.SetBorder(true).
		SetTitle("数据内容").
		SetTitleAlign(tview.AlignCenter)

	// 创建状态栏
	dv.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]选择左侧的表以查看数据[-]")

	dv.statusBar.SetBorder(true).
		SetTitle("状态").
		SetTitleAlign(tview.AlignCenter)

	// 创建主布局
	leftPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(dv.tableList, 0, 1, true).
		AddItem(dv.statusBar, 3, 1, false)

	dv.view = tview.NewFlex().
		AddItem(leftPanel, 30, 1, true).
		AddItem(dv.tableView, 0, 3, false)
	// 设置按键处理
	dv.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			if dv.tableList.HasFocus() {
				dv.app.SetFocus(dv.tableView)
			} else {
				dv.app.SetFocus(dv.tableList)
			}
			return nil
		}
		return event
	})

}

func (dv *DatabaseViewer) showTable(tableName string) {
	dv.currentTable = tableName
	dv.statusBar.SetText(fmt.Sprintf("[blue]正在加载表: %s[-]", tableName))

	// 清空表格
	dv.tableView.Clear()

	switch tableName {
	case "user":
		dv.showUserTable()
	case "activity":
		dv.showActivityTable()
	case "chat":
		dv.showChatTable()
	case "friend":
		dv.showFriendTable()
	case "activity_comment":
		dv.showActivityCommentTable()
	case "activity_member":
		dv.showActivityMemberTable()
	case "file":
		dv.showFileTable()
	case "admin":
		dv.showAdminTable()
	default:
		dv.statusBar.SetText(fmt.Sprintf("[red]未知表: %s[-]", tableName))
	}
}

func (dv *DatabaseViewer) showUserTable() {
	dv.tableView.SetCell(0, 0, tview.NewTableCell("用户表数据").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	dv.tableView.Clear()
	// 设置表头
	headers := []string{"ID", "用户名", "姓名", "性别", "地址", "创建时间"}
	for i, header := range headers {
		dv.tableView.SetCell(0, i, tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))
	}

	// 获取用户数据 - 这里需要实现一个获取所有用户的方法
	users, err := controllers.GetAllUsers()
	if err != nil {
		dv.statusBar.SetText(fmt.Sprintf("[red]加载用户数据失败: %v[-]", err))
		return
	}

	for i, user := range users {
		row := i + 1
		dv.tableView.SetCell(row, 0, tview.NewTableCell(strconv.FormatInt(user.Id, 10)))
		dv.tableView.SetCell(row, 1, tview.NewTableCell(user.Username))
		dv.tableView.SetCell(row, 2, tview.NewTableCell(user.Name))
		dv.tableView.SetCell(row, 3, tview.NewTableCell(user.Gender))
		dv.tableView.SetCell(row, 4, tview.NewTableCell(user.Addr))
		dv.tableView.SetCell(row, 5, tview.NewTableCell(user.CreateTime.Format("2006-01-02 15:04:05")))
	}

	dv.statusBar.SetText(fmt.Sprintf("[green]已加载 %d 条用户记录[-]", len(users)))
}

func (dv *DatabaseViewer) showActivityTable() {
	dv.tableView.SetCell(0, 0, tview.NewTableCell("活动表数据").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	dv.tableView.Clear()
	// 设置表头
	headers := []string{"ID", "名称", "介绍", "地址", "创建者ID", "状态", "创建时间"}
	for i, header := range headers {
		dv.tableView.SetCell(0, i, tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))
	}

	// 获取活动数据
	activities, err := controllers.GetAllActivities()
	if err != nil {
		dv.statusBar.SetText(fmt.Sprintf("[red]加载活动数据失败: %v[-]", err))
		return
	}

	for i, activity := range activities {
		row := i + 1
		dv.tableView.SetCell(row, 0, tview.NewTableCell(strconv.FormatInt(activity.Id, 10)))
		dv.tableView.SetCell(row, 1, tview.NewTableCell(activity.Name))
		dv.tableView.SetCell(row, 2, tview.NewTableCell(activity.Intro))
		dv.tableView.SetCell(row, 3, tview.NewTableCell(activity.Addr))
		dv.tableView.SetCell(row, 4, tview.NewTableCell(strconv.FormatInt(activity.UserId, 10)))
		dv.tableView.SetCell(row, 5, tview.NewTableCell(strconv.Itoa(int(activity.State))))
		dv.tableView.SetCell(row, 6, tview.NewTableCell(activity.CreateTime.Format("2006-01-02 15:04:05")))
	}

	dv.statusBar.SetText(fmt.Sprintf("[green]已加载 %d 条活动记录[-]", len(activities)))
}

func (dv *DatabaseViewer) showChatTable() {
	dv.tableView.SetCell(0, 0, tview.NewTableCell("聊天表数据").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	dv.tableView.Clear()
	// 设置表头
	headers := []string{"ID", "发送者ID", "接收者ID", "内容", "创建时间", "发送者状态", "接收者状态"}
	for i, header := range headers {
		dv.tableView.SetCell(0, i, tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))
	}

	// 获取聊天数据
	chats, err := controllers.GetAllChats()
	if err != nil {
		dv.statusBar.SetText(fmt.Sprintf("[red]加载聊天数据失败: %v[-]", err))
		return
	}

	for i, chat := range chats {
		row := i + 1
		content := chat.Content
		if len(content) > 30 {
			content = content[:30] + "..."
		}

		dv.tableView.SetCell(row, 0, tview.NewTableCell(strconv.FormatInt(chat.Id, 10)))
		dv.tableView.SetCell(row, 1, tview.NewTableCell(strconv.FormatInt(chat.UserIdFrom, 10)))
		dv.tableView.SetCell(row, 2, tview.NewTableCell(strconv.FormatInt(chat.UserIdTo, 10)))
		dv.tableView.SetCell(row, 3, tview.NewTableCell(content))
		dv.tableView.SetCell(row, 4, tview.NewTableCell(chat.CreateTime.Format("2006-01-02 15:04:05")))
		dv.tableView.SetCell(row, 5, tview.NewTableCell(strconv.Itoa(int(chat.StatusFrom))))
		dv.tableView.SetCell(row, 6, tview.NewTableCell(strconv.Itoa(int(chat.StatusTo))))
	}

	dv.statusBar.SetText(fmt.Sprintf("[green]已加载 %d 条聊天记录[-]", len(chats)))
}

func (dv *DatabaseViewer) showFriendTable() {
	dv.tableView.SetCell(0, 0, tview.NewTableCell("好友表数据").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	dv.tableView.Clear()
	// 设置表头
	headers := []string{"ID", "用户ID", "好友ID", "状态", "创建时间"}
	for i, header := range headers {
		dv.tableView.SetCell(0, i, tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))
	}

	friens, err := controllers.GetAllFriends()
	if err != nil {
		dv.statusBar.SetText(fmt.Sprintf("[red]加载好友数据失败: %v[-]", err))
		return
	}
	for i, friend := range friens {
		row := i + 1
		dv.tableView.SetCell(row, 0, tview.NewTableCell(strconv.FormatInt(friend.Id, 10)))
		dv.tableView.SetCell(row, 1, tview.NewTableCell(strconv.FormatInt(friend.UserId, 10)))
		dv.tableView.SetCell(row, 2, tview.NewTableCell(strconv.FormatInt(friend.FriendId, 10)))
		dv.tableView.SetCell(row, 3, tview.NewTableCell(strconv.Itoa(int(friend.Status))))
		dv.tableView.SetCell(row, 4, tview.NewTableCell(friend.CreateTime.Format("2006-01-02 15:04:05")))
	}
	dv.statusBar.SetText(fmt.Sprintf("[green]已加载 %d 条好友记录[-]", len(friens)))
}

func (dv *DatabaseViewer) showActivityCommentTable() {
	dv.tableView.SetCell(0, 0, tview.NewTableCell("活动评论表数据").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	dv.tableView.Clear()
	// 设置表头
	headers := []string{"ID", "活动ID", "用户ID", "评论内容", "创建时间"}
	for i, header := range headers {
		dv.tableView.SetCell(0, i, tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))
	}
	comments, err := controllers.GetAllComments()
	if err != nil {
		dv.statusBar.SetText(fmt.Sprintf("[red]加载活动评论数据失败: %v[-]", err))
		return
	}
	for i, comment := range comments {
		row := i + 1
		dv.tableView.SetCell(row, 0, tview.NewTableCell(strconv.FormatInt(comment.Id, 10)))
		dv.tableView.SetCell(row, 1, tview.NewTableCell(strconv.FormatInt(comment.ActivityId, 10)))
		dv.tableView.SetCell(row, 2, tview.NewTableCell(strconv.FormatInt(comment.UserId, 10)))
		dv.tableView.SetCell(row, 3, tview.NewTableCell(comment.Content))
		dv.tableView.SetCell(row, 4, tview.NewTableCell(comment.CreateTime.Format("2006-01-02 15:04:05")))
	}
	dv.statusBar.SetText(fmt.Sprintf("[green]已加载 %d 条活动评论记录[-]", len(comments)))
	dv.tableView.ScrollToBeginning()
}

func (dv *DatabaseViewer) showActivityMemberTable() {
	dv.tableView.SetCell(0, 0, tview.NewTableCell("活动成员表数据").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	dv.tableView.Clear()
	// 设置表头
	headers := []string{"ID", "活动ID", "用户ID", "创建时间"}
	for i, header := range headers {
		dv.tableView.SetCell(0, i, tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))
	}
	members, err := controllers.GetAllActivityMembers()
	if err != nil {
		dv.statusBar.SetText(fmt.Sprintf("[red]加载活动成员数据失败: %v[-]", err))
		return
	}
	for i, member := range members {
		row := i + 1
		dv.tableView.SetCell(row, 0, tview.NewTableCell(strconv.FormatInt(member.Id, 10)))
		dv.tableView.SetCell(row, 1, tview.NewTableCell(strconv.FormatInt(member.ActivityId, 10)))
		dv.tableView.SetCell(row, 2, tview.NewTableCell(strconv.FormatInt(member.UserId, 10)))
		dv.tableView.SetCell(row, 3, tview.NewTableCell(member.CreateTime.Format("2006-01-02 15:04:05")))
	}
	dv.statusBar.SetText(fmt.Sprintf("[green]已加载 %d 条活动成员记录[-]", len(members)))
	dv.tableView.ScrollToBeginning()
}

func (dv *DatabaseViewer) showFileTable() {
	dv.tableView.SetCell(0, 0, tview.NewTableCell("文件表数据").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	dv.tableView.Clear()
	// 设置表头
	headers := []string{"ID", "文件名", "文件hash", "上传时间"}
	for i, header := range headers {
		dv.tableView.SetCell(0, i, tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))
	}
	files, err := controllers.GetAllFiles()
	if err != nil {
		dv.statusBar.SetText(fmt.Sprintf("[red]加载文件数据失败: %v[-]", err))
		return
	}
	for i, file := range files {
		row := i + 1
		dv.tableView.SetCell(row, 0, tview.NewTableCell(strconv.FormatInt(file.Id, 10)))
		dv.tableView.SetCell(row, 1, tview.NewTableCell(file.FileName))
		dv.tableView.SetCell(row, 2, tview.NewTableCell(file.FileHash))
		dv.tableView.SetCell(row, 3, tview.NewTableCell(file.CreateTime.Format("2006-01-02 15:04:05")))
	}
	dv.statusBar.SetText(fmt.Sprintf("[green]已加载 %d 条文件记录[-]", len(files)))
	dv.tableView.ScrollToBeginning()
}

func (dv *DatabaseViewer) showAdminTable() {
	dv.tableView.SetCell(0, 0, tview.NewTableCell("管理员表数据").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	dv.tableView.Clear()
	// 设置表头
	headers := []string{"ID", "用户名", "名称"}
	for i, header := range headers {
		dv.tableView.SetCell(0, i, tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))
	}
	admins, err := controllers.GetAllAdmins()
	if err != nil {
		dv.statusBar.SetText(fmt.Sprintf("[red]加载管理员数据失败: %v[-]", err))
		return
	}
	for i, admin := range admins {
		row := i + 1
		dv.tableView.SetCell(row, 0, tview.NewTableCell(strconv.FormatInt(admin.Id, 10)))
		dv.tableView.SetCell(row, 1, tview.NewTableCell(admin.Username))
		dv.tableView.SetCell(row, 2, tview.NewTableCell(admin.Name))
	}
	dv.statusBar.SetText(fmt.Sprintf("[green]已加载 %d 条管理员记录[-]", len(admins)))
	dv.tableView.ScrollToBeginning()
}

func (dv *DatabaseViewer) Refresh() {
	if dv.currentTable != "" {
		dv.showTable(dv.currentTable)
	}
}

func (dv *DatabaseViewer) GetView() tview.Primitive {
	return dv.view
}
