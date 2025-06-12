# HobbyHub 后端系统详细设计报告

## 1. 项目概述

### 1.1 项目简介
HobbyHub 是一个多平台兴趣活动管理系统的后端服务，基于 Go 语言和 Gin 框架开发。系统支持用户管理、好友系统、聊天功能、活动管理、文件上传等核心功能，并集成了 JWT 认证和 Swagger API 文档。

### 1.2 技术栈
- **编程语言**: Go 1.19+
- **Web框架**: Gin
- **数据库**: MySQL / SQLite
- **ORM框架**: GORM
- **认证**: JWT (golang-jwt/jwt)
- **API文档**: Swagger (swaggo/gin-swagger)
- **测试框架**: Testify, go-sqlmock
- **配置管理**: YAML

### 1.3 系统特点
- RESTful API 设计
- JWT 认证机制
- 文件上传与管理
- 软删除机制
- 完善的单元测试
- 在线 API 文档

## 2. 系统架构设计

### 2.1 总体架构
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   客户端应用     │    │   Web 前端      │    │   移动端应用     │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────┴─────────────┐
                    │      Gin Web Server       │
                    │     (API Gateway)         │
                    └─────────────┬─────────────┘
                                 │
                    ┌─────────────┴─────────────┐
                    │     Business Logic        │
                    │    (Controllers)          │
                    └─────────────┬─────────────┘
                                 │
                    ┌─────────────┴─────────────┐
                    │      数据访问层           │
                    │      (GORM ORM)           │
                    └─────────────┬─────────────┘
                                 │
                    ┌─────────────┴─────────────┐
                    │      数据存储层           │
                    │    (MySQL/SQLite)         │
                    └───────────────────────────┘
```

### 2.2 目录结构
```
app-server/
├── api/                    # HTTP 路由处理层
│   ├── user.go            # 用户相关接口
│   ├── chat.go            # 聊天相关接口
│   ├── activity.go        # 活动相关接口
│   └── file.go            # 文件相关接口
├── cmd/                   # 程序入口
│   ├── main.go            # 主程序入口
│   └── config.yaml        # 配置文件
├── config/                # 配置管理
│   ├── read_config.go     # 配置读取
│   └── mysql.go           # 数据库初始化
├── controllers/           # 业务逻辑控制器
│   ├── user_controller.go
│   └── activity_controller.go
├── models/                # 数据模型
│   ├── user.go
│   ├── chat.go
│   └── activity.go
├── utils/                 # 工具函数
│   ├── jwt.go             # JWT 处理
│   ├── jwt_test.go        # JWT 测试
│   ├── random_string_test.go
│   └── password.go        # 密码加密
└── docs/                  # Swagger 文档
```

## 3. 数据库设计

### 3.1 数据库选型
系统支持 MySQL 和 SQLite 两种数据库：
- **MySQL**: 生产环境推荐，支持高并发
- **SQLite**: 开发测试环境，轻量级部署

### 3.2 核心数据表

#### 3.2.1 用户表 (user)
```sql
CREATE TABLE `user` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT COMMENT '用户Id',
  `username` varchar(50) UNIQUE COMMENT '用户名',
  `password` varchar(255) COMMENT '密码',
  `name` varchar(50) COMMENT '姓名',
  `gender` varchar(10) COMMENT '性别',
  `addr` varchar(255) COMMENT '地址',
  `head_img` varchar(255) COMMENT '头像图片',
  `create_time` datetime COMMENT '创建时间',
  `lat` double COMMENT '纬度',
  `lon` double COMMENT '经度'
);
```

#### 3.2.2 聊天记录表 (chat)
```sql
CREATE TABLE `chat` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT COMMENT '记录Id',
  `user_id_from` bigint NOT NULL COMMENT '发送用户Id',
  `user_id_to` bigint NOT NULL COMMENT '接收用户Id',
  `content` text NOT NULL COMMENT '消息内容',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `status_from` int NOT NULL DEFAULT 2 COMMENT '发送方状态（0: 删除, 1: 正常）',
  `status_to` int NOT NULL DEFAULT 2 COMMENT '接收方状态（0: 删除, 1: 正常）',
  INDEX `idx_user_from` (`user_id_from`),
  INDEX `idx_user_to` (`user_id_to`)
);
```

#### 3.2.3 活动表 (activity)
```sql
CREATE TABLE `activity` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT COMMENT '活动Id',
  `name` varchar(100) NOT NULL COMMENT '活动名称',
  `description` text COMMENT '活动描述',
  `user_id` bigint NOT NULL COMMENT '创建者Id',
  `create_time` datetime COMMENT '创建时间',
  `start_time` datetime COMMENT '开始时间',
  `end_time` datetime COMMENT '结束时间',
  `location` varchar(255) COMMENT '活动地点',
  `max_participants` int COMMENT '最大参与人数',
  `if_delete` tinyint DEFAULT 0 COMMENT '是否删除（0: 正常, 1: 已删除）'
);
```

#### 3.2.4 活动成员表 (activity_member)
```sql
CREATE TABLE `activity_member` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `activity_id` bigint NOT NULL COMMENT '活动Id',
  `user_id` bigint NOT NULL COMMENT '用户Id',
  `create_time` datetime COMMENT '加入时间',
  UNIQUE KEY `uk_activity_user` (`activity_id`, `user_id`)
);
```

### 3.3 数据库特性

#### 3.3.1 软删除机制
- 活动表使用 `if_delete` 字段实现软删除
- 聊天记录使用 `status_from` 和 `status_to` 实现双向软删除
- 保证数据完整性和可恢复性

#### 3.3.2 索引设计
- 用户表：`username` 唯一索引
- 聊天表：`user_id_from`, `user_id_to` 索引
- 活动成员表：`(activity_id, user_id)` 复合唯一索引

## 4. API 设计

### 4.1 API 规范
- **协议**: HTTP/HTTPS
- **格式**: JSON
- **认证**: JWT Bearer Token
- **版本**: /api/v1/

### 4.2 核心 API 接口

#### 4.2.1 用户认证
```
POST /api/v1/login          # 用户登录
PUT  /api/v1/user           # 用户注册
GET  /api/v1/user           # 获取用户信息
POST /api/v1/user           # 更新用户信息
```

#### 4.2.2 聊天功能
```
GET    /api/v1/chat         # 获取聊天记录
POST   /api/v1/chat         # 发送聊天消息
DELETE /api/v1/chat/:id     # 删除聊天记录
```

#### 4.2.3 活动管理
```
GET    /api/v1/activity           # 获取所有活动
GET    /api/v1/activity/:id       # 获取活动详情
PUT    /api/v1/activity           # 创建活动
POST   /api/v1/activity/:id       # 更新活动
DELETE /api/v1/activity/:id       # 删除活动
GET    /api/v1/activity/:id/member    # 获取活动成员
PUT    /api/v1/activity/:id/member    # 加入活动
DELETE /api/v1/activity/:id/member    # 退出活动
```

#### 4.2.4 文件管理
```
POST   /api/v1/file         # 上传文件
GET    /api/v1/file/:id     # 下载文件
DELETE /api/v1/file/:id     # 删除文件
```

### 4.3 统一响应格式

#### 4.3.1 成功响应
```json
{
  "data": {...},
  "code": 200,
  "message": "success"
}
```

#### 4.3.2 错误响应
```json
{
  "errorMessage": "具体错误信息",
  "code": 400,
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 5. 核心功能模块

### 5.1 认证与授权模块

#### 5.1.1 JWT 认证
```go
// JWT Token 生成
func GenerateJWT(u *models.User) (string, error) {
    claims := jwt.MapClaims{
        "id":  u.Id,
        "exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(config.GetConfig().Authentication.JwtSecret))
}
```

#### 5.1.2 密码加密
- 使用 bcrypt 算法进行密码哈希
- 盐值自动生成，提高安全性

### 5.2 文件管理模块

#### 5.2.1 文件上传特性
- 文件大小限制（可配置）
- 文件类型白名单验证
- SHA256 哈希校验
- 重复文件检测和链接

#### 5.2.2 存储策略
```go
// 文件存储路径规则
destFilename := fmt.Sprintf("%d%s", fileInfo.Id, fileExt)
destPath := filepath.Join(fileConfig.UploadPath, destFilename)
```

### 5.3 聊天系统

#### 5.3.1 好友验证机制
- 只能向好友发送消息
- 消息发送前验证好友关系
- 支持消息状态管理

#### 5.3.2 消息查询优化
- 双向查询合并
- 时间范围过滤
- 软删除状态过滤

### 5.4 活动管理

#### 5.4.1 权限控制
- 只有活动创建者可以修改/删除活动
- 成员管理权限验证
- 评论权限验证

#### 5.4.2 软删除机制
- 活动软删除，保留数据完整性
- 相关数据级联处理
- 删除权限验证

## 6. 配置管理

### 6.1 配置文件结构
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
  charset: utf8mb4

authentication:
  jwtsecret: "your_jwt_secret"

file:
  upload_path: "./uploads"
  max_size: 10
  allowed_types: ["png", "jpg", "jpeg", "gif", "pdf", "doc", "docx"]
```

### 6.2 配置特性
- 支持命令行参数指定配置文件
- 默认配置自动生成
- 环境变量支持
- 热重载配置（部分配置）

## 7. 错误处理与日志

### 7.1 错误处理策略
- 统一错误响应格式
- 错误码规范化
- 敏感信息过滤
- 详细错误日志记录

### 7.2 日志系统
- 同时输出到控制台和文件
- 结构化日志格式
- 日志级别控制
- 日志轮转支持

## 8. 测试设计

### 8.1 测试覆盖
- 单元测试：Utils 模块
- 集成测试：Controllers 模块
- API 测试：HTTP 接口测试
- 数据库测试：Mock 数据库测试

### 8.2 测试工具
```go
// JWT 测试示例
func TestGenerateJWTAndParseJWT(t *testing.T) {
    restore := mockJwtSecret("testsecret")
    defer restore()
    
    dbRestore := setupDBMock()
    defer dbRestore()
    
    user := &models.User{Id: 123}
    tokenString, err := GenerateJWT(user)
    assert.NoError(t, err)
    
    parsedUser, err := ParseJWT(tokenString)
    assert.NoError(t, err)
    assert.Equal(t, user.Id, parsedUser.Id)
}
```

## 9. 安全设计

### 9.1 认证安全
- JWT Token 有效期控制（72小时）
- 密码 bcrypt 哈希存储
- 敏感信息不返回（密码字段）

### 9.2 数据安全
- SQL 注入防护（GORM ORM）
- XSS 防护（输入验证）
- 文件上传安全验证
- 权限验证机制

### 9.3 API 安全
- CORS 配置
- 请求频率限制
- 输入参数验证
- 错误信息安全过滤

## 10. 性能优化

### 10.1 数据库优化
- 合适的索引设计
- 查询语句优化
- 连接池配置
- 慢查询监控

### 10.2 缓存策略
- JWT Token 缓存
- 用户信息缓存
- 活动数据缓存
- 静态文件缓存

### 10.3 并发处理
- Gin 框架天然支持并发
- 数据库连接池管理
- 文件上传并发控制
- 内存使用优化

## 11. 部署与运维

### 11.1 部署方式
- 二进制文件部署
- Docker 容器化部署
- 系统服务部署
- 负载均衡部署

### 11.2 监控与告警
- 健康检查接口
- 系统指标监控
- 错误率监控
- 性能指标追踪

### 11.3 备份策略
- 数据库定期备份
- 文件存储备份
- 配置文件备份
- 日志文件归档

## 12. 扩展性设计

### 12.1 水平扩展
- 无状态服务设计
- 数据库读写分离
- 文件存储分布式
- 缓存集群支持

### 12.2 功能扩展
- 插件化架构预留
- API 版本管理
- 第三方服务集成
- 微服务拆分准备

## 13. 总结

HobbyHub 后端系统采用了现代化的 Go 语言技术栈，具有良好的架构设计和完善的功能模块。系统在安全性、性能、可维护性等方面都有充分考虑，能够满足中小型兴趣活动管理平台的需求。

### 13.1 技术亮点
- 完整的 RESTful API 设计
- JWT 认证机制
- 软删除数据保护
- 完善的单元测试
- 在线 API 文档

### 13.2 发展建议
- 引入缓存机制提升性能
- 添加消息队列支持
- 实现实时通讯功能
- 增加数据分析模块
- 支持多租户架构

---

*本文档基于 HobbyHub v1.0 系统代码分析生成，如有疑问请联系开发团队。*
