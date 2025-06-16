# HobbyHub 鸿蒙手机端

![HarmonyOS](https://img.shields.io/badge/HarmonyOS-5.0.5-blue)
![ArkTS](https://img.shields.io/badge/ArkTS-TypeScript-green)
![DevEco Studio](https://img.shields.io/badge/IDE-DevEco%20Studio-orange)

## 📱 项目简介

HobbyHub 鸿蒙手机端是一款基于 HarmonyOS 开发的兴趣活动社交应用，旨在帮助用户发现和参与志同道合的兴趣活动。应用采用 ArkTS 语言开发，使用声明式 UI 框架 ArkUI 构建现代化的用户界面。

## ✨ 核心功能

### 🏠 主要模块

- **首页推荐** - 智能推荐感兴趣的活动
- **活动管理** - 创建、编辑、参与各类兴趣活动
- **社交互动** - 好友管理、实时聊天、社区交流
- **个人中心** - 用户信息管理、活动历史查看

### 🎯 特色功能

- **智能匹配** - 基于兴趣标签的活动推荐算法
- **实时通讯** - 支持文字、图片等多媒体消息
- **位置服务** - 基于地理位置的附近活动发现
- **文件分享** - 支持图片、文档等文件上传分享

## 🛠️ 技术架构

### 开发环境

- **开发语言**: ArkTS (TypeScript for HarmonyOS)
- **UI框架**: ArkUI 声明式开发范式
- **IDE**: DevEco Studio
- **构建工具**: Hvigor
- **测试框架**: Hypium

### 技术特点

- **响应式编程** - 使用 `@State` 状态管理实现数据驱动UI
- **组件化设计** - 模块化组件开发，提高代码复用性
- **原生性能** - 充分利用鸿蒙系统特性，保证应用流畅性
- **多设备适配** - 支持手机、平板、2合1设备

## 📁 项目结构

```
app-harmonyos/
├── AppScope/                    # 应用全局配置
│   ├── app.json5               # 应用配置文件
│   └── resources/              # 全局资源文件
├── entry/                      # 主入口模块
│   ├── src/main/
│   │   ├── ets/                # 主要源代码目录
│   │   │   ├── api/            # API接口层
│   │   │   │   ├── LoginApi.ets     # 登录认证API
│   │   │   │   ├── EventApi.ets     # 活动相关API
│   │   │   │   ├── UserApi.ets      # 用户信息API
│   │   │   │   └── ...
│   │   │   ├── common/         # 公共配置和常量
│   │   │   │   └── Constans.ets     # API常量定义
│   │   │   ├── component/      # 可复用组件
│   │   │   ├── entryability/   # 应用入口能力
│   │   │   │   └── EntryAbility.ets # 主入口能力
│   │   │   ├── pages/          # 页面视图
│   │   │   │   ├── Index.ets        # 主页面框架
│   │   │   │   ├── LoginPage.ets    # 登录页面
│   │   │   │   ├── EventDetailPage.ets # 活动详情页
│   │   │   │   └── ...
│   │   │   ├── utils/          # 工具类库
│   │   │   │   └── HttpRequest.ets  # HTTP请求封装
│   │   │   ├── view/           # 视图组件
│   │   │   │   ├── HomeViewComponent.ets    # 首页组件
│   │   │   │   ├── EventViewComponent.ets   # 活动组件
│   │   │   │   ├── FriendsViewController.ets # 好友组件
│   │   │   │   └── UserInfoViewComponent.ets # 个人中心组件
│   │   │   └── viewModel/      # 视图模型层
│   │   ├── module.json5        # 模块配置文件
│   │   └── resources/          # 资源文件
│   ├── src/ohosTest/          # 单元测试
│   └── src/test/              # 集成测试
├── oh_modules/                # 鸿蒙依赖模块
├── oh-package.json5           # 项目依赖配置
└── build-profile.json5        # 构建配置
```

## 🚀 快速开始

### 环境要求

- **DevEco Studio**: 5.0 或以上版本
- **HarmonyOS SDK**: API Level 12 或以上
- **Node.js**: 16.0 或以上版本（可选，用于构建工具）

### 安装步骤

1. **克隆项目**

```bash
git clone <repository-url>
cd HobbyHub/app-harmonyos
```

2. **打开项目**

   - 启动 DevEco Studio
   - 选择 "Open" 打开项目目录
   - 等待依赖自动下载和配置
3. **配置后端服务**

   - 修改 `entry/src/main/ets/common/Constans.ets` 中的 `BaseUrl`
   - 确保后端服务已启动并可访问
4. **运行项目**

   - 连接鸿蒙设备或启动模拟器
   - 点击 "Run" 按钮编译并安装应用

## ⚙️ 配置说明

### API配置

在 `entry/src/main/ets/common/Constans.ets` 中配置后端服务地址：

```typescript
export const BaseUrl = 'http://your-server-ip:8081'
```

### 应用权限

应用需要以下权限（已在 `module.json5` 中配置）：

- `ohos.permission.INTERNET` - 网络访问权限
- `ohos.permission.LOCATION` - 精确位置权限
- `ohos.permission.APPROXIMATELY_LOCATION` - 大概位置权限

### 构建配置

主要配置文件：

- `oh-package.json5` - 依赖管理
- `build-profile.json5` - 构建配置
- `module.json5` - 模块配置

## 🏗️ 开发指南

### 项目架构说明

#### 1. MVVM架构模式

```
View (页面/组件) ←→ ViewModel (业务逻辑) ←→ Model (数据层/API)
```

#### 2. 状态管理

使用 ArkUI 内置的状态管理机制：

```typescript
@State private data: DataType = initialValue;
@Prop readonly props: PropType;
@Link shared data: SharedType;
```

#### 3. 组件生命周期

```typescript
onPageShow() { /* 页面显示时 */ }
onPageHide() { /* 页面隐藏时 */ }
aboutToAppear() { /* 组件即将出现 */ }
aboutToDisappear() { /* 组件即将消失 */ }
```

### API调用示例

```typescript
import { login } from '../api/LoginApi';

// 调用登录API
try {
  const response = await login(username, password);
  if (response.success) {
    // 登录成功处理
    console.log('登录成功:', response.data);
  } else {
    // 登录失败处理
    console.error('登录失败:', response.message);
  }
} catch (error) {
  console.error('网络错误:', error);
}
```

### 自定义组件开发

```typescript
@Component
export struct CustomComponent {
  @State private value: string = '';
  
  build() {
    Column() {
      Text(this.value)
        .fontSize(16)
        .fontColor('#333333')
    }
    .width('100%')
    .padding(16)
  }
}
```

## 🧪 测试

### 运行单元测试

```bash
# 在DevEco Studio中
# 右键点击测试文件或测试目录
# 选择 "Run Tests"
```

### 测试覆盖的功能

- API接口调用测试
- 组件渲染测试
- 业务逻辑测试
- 用户交互测试

## 📦 构建和发布

### 开发版本构建

```bash
# 在DevEco Studio中选择构建类型为 "debug"
# 点击 Build -> Build Hap(s)/App(s)
```

### 生产版本构建

```bash
# 在DevEco Studio中选择构建类型为 "release"
# 配置签名证书
# 点击 Build -> Build Hap(s)/App(s)
```

### 签名配置

1. 生成或导入签名证书
2. 在 `build-profile.json5` 中配置签名信息
3. 确保证书配置正确

## 🔧 常见问题

### Q: 编译时出现依赖错误

**A**: 检查 `oh-package.json5` 中的依赖版本，执行清理重新构建。

### Q: 网络请求失败

**A**:

1. 检查设备网络连接
2. 确认后端服务地址和端口
3. 验证设备可以访问后端服务

### Q: 应用安装失败

**A**:

1. 检查设备开发者模式是否开启
2. 确认USB调试权限
3. 检查应用签名配置

### Q: 页面渲染异常

**A**:

1. 检查状态变量绑定
2. 验证组件生命周期方法
3. 查看DevEco Studio日志输出

## 🤝 贡献指南

1. Fork 项目仓库
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

### 代码规范

- 遵循 ArkTS 官方编码规范
- 使用有意义的变量和函数命名
- 添加必要的注释和文档
- 保持代码格式一致性

## 🙏 致谢

感谢以下开源项目和社区的支持：

- [HarmonyOS](https://developer.harmonyos.com/) - 华为鸿蒙操作系统
- [ArkUI](https://developer.harmonyos.com/cn/docs/documentation/doc-guides-V3/arkui-overview-0000001281001133-V3) - 鸿蒙声明式UI框架
- [DevEco Studio](https://developer.harmonyos.com/cn/develop/deveco-studio) - 鸿蒙开发IDE

---

⭐ **如果这个项目对您有帮助，请给我们一个 Star！**
