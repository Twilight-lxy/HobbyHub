// src/api/user.js
import request from '@/utils/request'

export function getUserList(params) {
  return request({
    url: '/api/v1/admin/users',
    method: 'get',
    params
  })
}

export function getUserDetail(id) {
  return request({
    url: `/api/v1/admin/users/${id}`,
    method: 'get'
  })
}

export function createUser(data) {
  return request({
    url: '/api/v1/admin/users',
    method: 'post',
    data
  })
}

export function updateUser(id, data) {
  return request({
    url: `/api/v1/admin/users/${id}`,
    method: 'put',
    data
  })
}

export function deleteUser(id) {
  return request({
    url: `/api/v1/admin/users/${id}`,
    method: 'delete'
  })
}