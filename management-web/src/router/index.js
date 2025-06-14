import { createRouter, createWebHashHistory } from 'vue-router'

// 布局组件
const Layout = () => import('@/layout/index.vue')

// 路由配置
const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录', hidden: true }
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '主页面', icon: 'Odometer' }
      }
    ]
  },
  {
    path: '/user',
    component: Layout,
    redirect: '/user/list',
    meta: { title: '用户管理', icon: 'User' },
    children: [
      {
        path: 'list',
        name: 'UserList',
        component: () => import('@/views/user/list.vue'),
        meta: { title: '查看用户列表' }
      },
      {
        path: 'message',
        name: 'UserMessage',
        component: () => import('@/views/user/message.vue'),
        meta: { title: '用户评论管理' }
      },
    ]
  },
  {
    path: '/activity',
    component: Layout,
    redirect: '/activity/list',
    meta: { title: '活动管理', icon: 'Calendar' },
    children: [
      {
        path: 'list',
        name: 'ActivityList',
        component: () => import('@/views/activity/list.vue'),
        meta: { title: '活动列表' }
      },
      {
        path: 'category',
        name: 'ActivityCategory',
        component: () => import('@/views/activity/category.vue'),
        meta: { title: '活动分类' }
      },
      {
        path: 'create',
        name: 'ActivityCreate',
        component: () => import('@/views/activity/form.vue'),
        meta: { title: '创建活动' }
      },
      {
        path: 'edit/:id',
        name: 'ActivityEdit',
        component: () => import('@/views/activity/form.vue'),
        meta: { title: '编辑活动', hidden: true }
      },
      {
        path: 'detail/:id',
        name: 'ActivityDetail',
        component: () => import('@/views/activity/detail.vue'),
        meta: { title: '活动详情', hidden: true }
      }
    ]
  },
  {
    path: '/profile',
    component: Layout,
    redirect: '/profile/index',
    children: [
      {
        path: 'index',
        name: 'Profile',
        component: () => import('@/views/profile/index.vue'),
        meta: { title: '个人中心', icon: 'UserFilled', hidden: true }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - 兴趣小队管理系统` : '兴趣小队管理系统'
  
  // 判断是否登录
  const token = localStorage.getItem('token')
  if (to.path !== '/login' && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router 