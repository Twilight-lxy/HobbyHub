import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as loginApi } from '@/api/auth'
import { getUserInfo as getUserInfoApi } from '@/api/user'
import router from '@/router'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref({})
  const roles = ref([])

  // 登录
  const login = async (loginForm) => {
    try {
      const res = await loginApi(loginForm)
      token.value = res.data.token
      localStorage.setItem('token', res.data.token)
      await getUserInfo()
      return Promise.resolve()
    } catch (error) {
      return Promise.reject(error)
    }
  }

  // 获取用户信息
  const getUserInfo = async () => {
    try {
      const res = await getUserInfoApi()
      userInfo.value = res.data
      roles.value = res.data.roles || []
      return Promise.resolve(res.data)
    } catch (error) {
      return Promise.reject(error)
    }
  }

  // 退出登录
  const logout = () => {
    token.value = ''
    userInfo.value = {}
    roles.value = []
    localStorage.removeItem('token')
    router.push('/login')
  }

  // 重置状态
  const resetState = () => {
    token.value = ''
    userInfo.value = {}
    roles.value = []
  }

  return {
    token,
    userInfo,
    roles,
    login,
    getUserInfo,
    logout,
    resetState
  }
}) 