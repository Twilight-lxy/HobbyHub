/**
 * 认证工具类
 */

/**
 * 保存令牌到本地存储
 * @param {string} token 令牌
 */
export function setToken(token) {
  localStorage.setItem('token', token)
}

/**
 * 从本地存储获取令牌
 * @returns {string} 令牌
 */
export function getToken() {
  return localStorage.getItem('token')
}

/**
 * 从本地存储中移除令牌
 */
export function removeToken() {
  localStorage.removeItem('token')
}

/**
 * 保存管理员信息到本地存储
 * @param {Object} admin 管理员信息
 */
export function setAdmin(admin) {
  localStorage.setItem('adminInfo', JSON.stringify(admin))
}

/**
 * 从本地存储获取管理员信息
 * @returns {Object} 管理员信息
 */
export function getAdmin() {
  const adminInfo = localStorage.getItem('adminInfo')
  return adminInfo ? JSON.parse(adminInfo) : null
}

/**
 * 从本地存储中移除管理员信息
 */
export function removeAdmin() {
  localStorage.removeItem('adminInfo')
}

/**
 * 清除所有认证信息
 */
export function clearAuth() {
  removeToken()
  removeAdmin()
}

/**
 * 检查当前登录的是否为管理员
 * @returns {boolean} 是否为管理员
 */
export function isAdmin() {
  return !!getAdmin()
} 