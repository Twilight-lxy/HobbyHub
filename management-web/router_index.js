// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'
import Login from '@/views/auth/Login.vue'
import Layout from '@/layout/index.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { title: '登录' }
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
        meta: { title: '仪表盘', icon: 'dashboard' }
      },
      {
        path: 'user',
        name: 'UserManagement',
        component: () => import('@/views/admin/UserManagement.vue'),
        meta: { title: '用户管理', icon: 'user' }
      },
      {
        path: 'activity',
        name: 'ActivityManagement',
        component: () => import('@/views/admin/ActivityManagement.vue'),
        meta: { title: '活动管理', icon: 'activity' }
      },
      {
        path: 'comment',
        name: 'CommentManagement',
        component: () => import('@/views/admin/CommentManagement.vue'),
        meta: { title: '评论管理', icon: 'comment' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

// 路由守卫 - 权限验证
router.beforeEach((to, from, next) => {
  const isAuthenticated = localStorage.getItem('token')
  
  if (to.name !== 'Login' && !isAuthenticated) {
    next({ name: 'Login' })
  } else {
    next()
  }
})

export default router