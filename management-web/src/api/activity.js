import request from '@/utils/request'

// 获取活动列表
export function getActivityList(params) {
  return request({
    url: '/product/list',
    method: 'get',
    params
  })
}

// 获取活动详情
export function getActivityDetail(id) {
  return request({
    url: `/product/${id}`,
    method: 'get'
  })
}

// 创建活动
export function createActivity(data) {
  return request({
    url: '/product',
    method: 'post',
    data
  })
}

// 更新活动
export function updateActivity(id, data) {
  return request({
    url: `/product/${id}`,
    method: 'put',
    data
  })
}

// 删除活动
export function deleteActivity(id) {
  return request({
    url: `/product/${id}`,
    method: 'delete'
  })
}

// 上架/下架活动
export function toggleActivityStatus(id, isActive) {
  return request({
    url: `/product/status/${id}`,
    method: 'put',
    data: { isActive }
  })
} 