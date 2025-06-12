import axios from 'axios'
import { ElMessage } from 'element-plus'
import { getToken, clearAuth, isAdmin } from '@/utils/auth'

// 创建axios实例
const service = axios.create({
  baseURL: '', // 基础URL，使用代理时留空
  timeout: 10000 // 请求超时时间
})

// 请求拦截器
service.interceptors.request.use(
  config => {
    // 使用auth工具类获取token
    const token = getToken()
    if (token && token.trim() !== '') {
      // 添加Bearer前缀 - 这是后端期望的格式
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  error => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    const res = response.data
    
    // 如果状态码不是200，则判断为错误
    if (res.code !== 200) {
      ElMessage.error(res.msg || '系统错误')
      
      // 401: 未登录或token过期
      if (res.code === 401) {
        // 清除登录状态
        clearAuth()
        
        // 跳转到登录页
        setTimeout(() => {
          window.location.href = '/login'
        }, 1500)
      }
      
      return Promise.reject(new Error(res.msg || '系统错误'))
    } else {
      return res
    }
  },
  error => {
    console.error('响应错误:', error)
    
    // 处理网络错误
    let message = '网络错误，请检查您的网络连接'
    if (error.response) {
      switch (error.response.status) {
        case 401:
          message = '登录过期，请重新登录'
          // 清除登录状态
          clearAuth()
          
          // 跳转到登录页
          setTimeout(() => {
            window.location.href = '/login'
          }, 1500)
          break
        case 404:
          message = '请求的资源不存在'
          break
        case 500:
          message = '服务器内部错误'
          break
        default:
          message = `请求失败(${error.response.status})`
      }
    } else if (error.message.includes('timeout')) {
      message = '请求超时，请稍后再试'
    }
    
    ElMessage.error(message)
    return Promise.reject(error)
  }
)

export default service 