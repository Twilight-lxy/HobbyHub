# 兴趣小队管理系统

这是一个基于Vue 3和Element Plus的兴趣小队管理系统，用于管理兴趣小队平台的活动、用户、小队和参与记录。

## 功能特性

- 📊 **仪表盘**：展示系统概览数据，包括活动数量、用户数量、小队数量和参与记录数量
- 👥 **用户管理**：管理平台用户，查看用户详情
- 🎯 **活动管理**：创建、编辑、删除和上下架活动
- 👨‍👩‍👧‍👦 **小队管理**：管理兴趣小队，查看小队详情
- 📝 **参与记录**：查看用户参与活动的记录
- 🔑 **系统管理**：角色管理和操作日志

## 技术栈

- **前端框架**：Vue 3
- **状态管理**：Pinia
- **UI组件库**：Element Plus
- **路由**：Vue Router
- **HTTP请求**：Axios
- **图表**：ECharts
- **构建工具**：Vite

## 开发指南

### 安装依赖

```bash
npm install
```

### 开发环境运行

```bash
npm run dev
```

### 构建生产版本

```bash
npm run build
```

### 预览生产版本

```bash
npm run preview
```

## 项目结构

```
src/
├── api/            # API请求模块
├── assets/         # 静态资源
├── components/     # 公共组件
├── layout/         # 布局组件
├── router/         # 路由配置
├── stores/         # 状态管理
├── utils/          # 工具函数
└── views/          # 页面视图
```

## 接口文档

系统默认连接到`http://localhost:8080/api`作为后端API地址，可以在`src/utils/request.js`中修改。

## 许可证

[MIT](LICENSE)
