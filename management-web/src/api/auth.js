import request from '@/utils/request'

// 管理员登录
export function login(data) {
  return request({
    url: '/api/v1/admin/login',
    method: 'post',
    data
  })
}

// 管理员注销
export function logout() {
  return request({
    url: '/auth/logout',
    method: 'post'
  })
} 