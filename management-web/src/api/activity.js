import request from '@/utils/request'

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
    url: '/api/v1/admin/activity',
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