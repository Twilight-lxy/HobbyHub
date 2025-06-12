import request from '@/utils/request'

// 获取参与记录列表
export function getOrderList(params) {
  return request({
    url: '/product_order/admin/list',
    method: 'get',
    params
  })
}

// 获取参与记录详情
export function getOrderDetail(id) {
  return request({
    url: `/product_order/${id}`,
    method: 'get'
  })
}

// 创建参与记录
export function createOrder(data) {
  return request({
    url: '/product_order/save',
    method: 'post',
    data
  })
}

// 批量创建参与记录
export function createOrderBatch(data) {
  return request({
    url: '/product_order/save/batch',
    method: 'post',
    data
  })
}

// 更新参与记录
export function updateOrder(id, data) {
  return request({
    url: `/product_order/${id}`,
    method: 'put',
    data
  })
}

// 删除参与记录
export function deleteOrder(id) {
  return request({
    url: `/product_order/${id}`,
    method: 'delete'
  })
} 