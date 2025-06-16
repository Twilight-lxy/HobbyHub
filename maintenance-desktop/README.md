# 后端运维TUI（maintenance-desktop）

## 项目简介

本项目为 HobbyHub 后端服务的运维终端界面（TUI，Text-based User Interface），用于便捷地在命令行下对后端服务进行监控、管理和维护。

## 主要功能

- 查看服务器运行状态
- 查看和管理日志
- 数据库内容浏览
- 配置文件查看
- 支持多平台运行（Windows、Linux）

## 目录结构

```
maintenance-desktop/
├── components/           # 主要功能组件
│   ├── config_viewer.go      # 配置文件查看器
│   ├── database_viewer.go    # 数据库内容查看器
│   ├── log_viewer.go         # 日志查看器
│   ├── main.go               # TUI 主界面
│   └── server_status.go      # 服务器状态监控
├── go.mod                # Go 依赖管理
├── go.sum                # Go 依赖校验
├── hobbyhub-server-tui.exe   # Windows下编译产物
├── main.go               # 程序入口
```

## 快速开始

### 1. 环境准备

- Go 1.18 及以上版本

### 2. 安装依赖

```sh
go mod tidy
```

### 3. 运行

```sh
go run main.go
```

或编译后运行：

```sh
go build -o hobbyhub-server-tui.exe
./hobbyhub-server-tui.exe
```

## 主要依赖

- [tview](https://github.com/rivo/tview)：终端UI库
- [tcell](https://github.com/gdamore/tcell)：终端处理库

## 贡献指南

欢迎提交 issue 和 PR，完善功能或修复 bug。
