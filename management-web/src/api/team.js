import request from '@/utils/request'

// 获取小队列表
export function getTeamList(params) {
  return request({
    url: '/product_cart/admin/list',
    method: 'get',
    params
  })
}

// 获取小队详情
export function getTeamDetail(id) {
  return request({
    url: `/product_cart/${id}`,
    method: 'get'
  })
}

// 创建小队
export function createTeam(data) {
  return request({
    url: '/product_cart',
    method: 'post',
    data
  })
}

// 更新小队
export function updateTeam(id, data) {
  return request({
    url: `/product_cart/${id}`,
    method: 'put',
    data
  })
}

// 删除小队
export function deleteTeam(ids) {
  return request({
    url: `/product_cart/remove/ids?ids=${ids.join(',')}`,
    method: 'delete'
  })
} 