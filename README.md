# HobbyHub - 多平台兴趣活动社交系统

## 🎯 项目简介

HobbyHub 是一个多平台兴趣活动社交系统，旨在帮助用户发现和参与志同道合的兴趣活动。该系统由四个主要部分组成：鸿蒙手机APP、Web管理端、桌面运维端和后端服务。

### 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  HarmonyOS App  │    │  Web Management │    │Desktop Maintenance│
│   (用户端)       │    │   (管理端)       │    │    (运维端)       │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────┴───────────┐
                    │    Backend Server       │
                    │   (Go + Gin + MySQL)   │
                    └─────────────────────────┘
```

## 🚀 功能特性

### 📱 鸿蒙手机APP端 (app-harmonyos)

- **兴趣活动发起**：用户可以随时随地发起各类兴趣活动
- **智能匹配**：寻找具有相同或相似兴趣的人群
- **兴趣圈分享**：通过兴趣圈分享活动，吸引同好者
- **互动功能**：支持查看、报名、评论和点赞活动
- **实时聊天**：参与者之间可以实时交流

### 🌐 Web管理端 (management-web)

- **用户管理**：管理平台用户，查看用户详情和状态
- **活动管理**：创建、编辑、删除和上下架活动
- **评论审核**：管理用户评论，维护社区环境
- **数据统计**：仪表盘展示用户、活动、参与等核心数据
- **权限管理**：分级管理员权限控制

### 🖥️ 桌面运维端 (maintenance-desktop)

- **服务器监控**：实时查看服务器运行状态
- **日志管理**：查看和分析系统操作日志
- **数据库管理**：浏览和管理数据库内容
- **配置管理**：查看和修改系统配置
- **预警系统**：自定义预警规则，及时发现问题

### ⚙️ 后端服务 (app-server)

- **RESTful API**：提供完整的用户、活动、聊天等API接口
- **JWT认证**：安全的用户认证和授权机制
- **文件管理**：支持图片、文档等文件上传和管理
- **数据库ORM**：使用GORM进行数据库操作
- **API文档**：集成Swagger自动生成API文档

## 🛠️ 技术栈

### 鸿蒙手机APP端

- **前端**：ArkTS + HarmonyOS SDK
- **HTTP客户端**：axios
- **测试框架**：Hypium

### Web管理端

- **前端框架**：Vue 3 + Composition API
- **UI组件库**：Element Plus
- **状态管理**：Pinia
- **路由管理**：Vue Router
- **HTTP客户端**：Axios
- **图表库**：ECharts
- **构建工具**：Vite

### 桌面运维端

- **开发语言**：Go
- **UI框架**：tview (Terminal UI)
- **终端处理**：tcell

### 后端服务

- **开发语言**：Go
- **Web框架**：Gin
- **数据库ORM**：GORM
- **数据库**：MySQL
- **认证**：JWT (golang-jwt/jwt)
- **API文档**：Swagger (swaggo/gin-swagger)
- **测试框架**：testify + go-sqlmock

## 📁 项目结构

```
HobbyHub/
├── app-harmonyos/          # 鸿蒙手机APP
│   ├── AppScope/           # 应用配置
│   ├── entry/              # 主入口模块
│   ├── oh_modules/         # 依赖模块
│   └── ...
├── app-server/             # 后端服务
│   ├── api/                # 路由定义
│   ├── controllers/        # 控制器
│   ├── models/             # 数据模型
│   ├── utils/              # 工具函数
│   ├── cmd/                # 程序入口
│   └── ...
├── management-web/         # Web管理端
│   ├── src/
│   │   ├── views/          # 页面视图
│   │   ├── components/     # 组件
│   │   ├── api/            # API接口
│   │   └── ...
│   └── ...
├── maintenance-desktop/    # 桌面运维端
│   ├── components/         # TUI组件
│   └── main.go
└── docs/                   # 项目文档
```

## 🚀 快速开始

### 1. 后端服务启动

```bash
cd app-server
go mod tidy
go run ./cmd/main.go
```

访问 Swagger 文档：http://localhost:8081/swagger/index.html

### 2. Web管理端启动

```bash
cd management-web
npm install
npm run dev
```

访问管理界面：http://localhost:5173

### 3. 桌面运维端启动

```bash
cd maintenance-desktop
go mod tidy
go run main.go
```

### 4. 鸿蒙APP开发

使用 DevEco Studio 打开 `app-harmonyos` 目录进行开发和调试。

## 🔧 配置说明

### 后端服务配置

编辑 `app-server/cmd/config.yaml`：

```yaml
server:
  host: localhost
  port: 8081
database:
  type: mysql
  username: root
  password: "your_password"
  host: localhost
  port: 3306
  database: hobbyhub
authentication:
  jwtsecret: "your_jwt_secret"
file:
  upload_path: "./uploads"
  max_size: 10
  allowed_types: ["png", "jpg", "jpeg", "gif", "pdf"]
```

### Web管理端配置

在 `management-web/src/utils/request.js` 中修改后端API地址。

## 🧪 测试

### 后端测试

```bash
cd app-server
go test ./...
```

### 前端测试

```bash
cd management-web
npm run test
```

## 📊 系统数据流

```
[HarmonyOS App] ──HTTP──> [Backend API] ──SQL──> [MySQL Database]
       ↓                        ↑                       ↑
[Web Management] ──HTTP──> [Backend API]                │
       ↓                        ↑                       │
[Desktop TUI] ──────File/DB──────┴───────────────────────┘
```

⭐ 如果这个项目对您有帮助，请给我们一个 Star！
