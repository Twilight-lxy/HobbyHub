hobbyhub-web/
├── public/                 # 静态资源
├── src/
│   ├── api/                # API 接口封装
│   ├── assets/             # 静态资源
│   ├── components/         # 公共组件
│   ├── router/             # 路由配置
│   ├── store/              # Vuex 状态管理
│   ├── styles/             # 全局样式
│   ├── utils/              # 工具函数
│   ├── views/              # 页面组件
│   │   ├── admin/          # 管理员页面
│   │   │   ├── UserManagement.vue  # 用户管理
│   │   │   ├── ActivityManagement.vue # 活动管理
│   │   │   └── CommentManagement.vue # 评论管理
│   │   ├── auth/           # 认证相关
│   │   └── dashboard/      # 仪表盘
│   ├── App.vue             # 根组件
│   └── main.js             # 入口文件
├── package.json
└── vue.config.js