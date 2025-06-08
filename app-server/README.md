# HobbyHub App Server

## 简介

HobbyHub App Server 是 HobbyHub 多平台兴趣活动系统的后端服务，基于 Go 语言和 Gin 框架开发，支持 MySQL 数据库，提供用户、活动等相关 API 接口，并集成 Swagger 在线文档。

## 目录结构

```
app-server/
├── api/            # 路由与接口定义
├── cmd/            # 程序入口与配置文件
├── config/         # 配置加载与数据库初始化
├── controllers/    # 业务逻辑控制器
├── docs/           # Swagger 文档
├── models/         # 数据模型
├── utils/          # 工具函数
├── go.mod
├── go.sum
└── README.md
```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置数据库

编辑 `cmd/config.yaml`，填写你的 MySQL 数据库信息：

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
    database: iteam01
    charset: utf8mb4
authentication:
    jwtsecret: "your_jwt_secret"
file:
    upload_path: "./uploads"
    max_size: 10
    allowed_types: ["png", "jpg", "jpeg", "gif", "pdf", "doc", "docx"]
```

### 3. 运行服务

```bash
go run ./cmd/main.go
```

支持通过命令行参数指定配置文件路径：

```bash
go run ./cmd/main.go -config=cmd/config.yaml
```

### 4. 访问 API

- 用户信息接口示例：
  `GET http://localhost:8081/api/v1/user/info/{id}`

### 5. Swagger 文档

启动服务后访问：
[http://localhost:8081/swagger/index.html](http://localhost:8081/swagger/index.html)

## 测试

运行所有测试：

```bash
go test ./...
```

运行特定模块测试：

```bash
go test ./controllers
go test ./utils
```

## 主要依赖

- [Gin](https://github.com/gin-gonic/gin) - HTTP Web框架
- [GORM](https://gorm.io/) - ORM框架
- [Swaggo/gin-swagger](https://github.com/swaggo/gin-swagger) - API文档生成
- [golang-jwt/jwt](https://github.com/golang-jwt/jwt) - JWT认证
- [testify](https://github.com/stretchr/testify) - 测试框架
- [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) - SQL测试Mock

## 贡献

欢迎提交 issue 和 PR！
