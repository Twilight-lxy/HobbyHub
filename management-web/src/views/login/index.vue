<template>
  <div class="login-container">
    <el-card class="login-card">
      <div class="login-header">
        <div class="logo">
          <img src="@/assets/logo.svg" alt="Logo" class="logo-image" />
        </div>
        <h2 class="title">兴趣小队管理系统</h2>
      </div>
      
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="用户名"
            prefix-icon="User"
            size="large"
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            placeholder="密码"
            prefix-icon="Lock"
            size="large"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        
        <el-form-item class="remember-me">
          <el-checkbox v-model="loginForm.rememberMe">记住我</el-checkbox>
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            :loading="loading"
            class="login-button"
            @click="handleLogin"
          >
            登录
          </el-button>
        </el-form-item>
      </el-form>
      
      <div class="login-tips">
        <el-alert
          title="管理员账号：admin，密码：123456"
          type="info"
          :closable="false"
          show-icon
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'
import axios from 'axios'
import { setToken, setAdmin } from '@/utils/auth'

const router = useRouter()
const loginFormRef = ref(null)
const loading = ref(false)

// 登录表单数据
const loginForm = reactive({
  username: '',
  password: '',
  rememberMe: false
})

// 表单验证规则
const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ]
}

// 初始化管理员账号
const initAdmin = async () => {
  try {
    // 使用axios直接发送请求，不使用request工具
    const response = await axios.get('/api/admin/init')
    if (response.data.code === 200) {
      console.log('管理员账号初始化成功:', response.data.msg)
    }
  } catch (error) {
    console.error('初始化管理员账号失败', error)
    if (error.response) {
      console.error('错误状态码:', error.response.status)
      console.error('错误信息:', error.response.data)
    }
  }
}

// 处理登录
const handleLogin = () => {
  loginFormRef.value?.validate(async (valid) => {
    if (valid) {
      loading.value = true
      
      try {
        // 使用axios直接发送登录请求，不使用拦截器
        const response = await axios.post('/api/admin/login', {
          username: loginForm.username,
          password: loginForm.password
        })
        
        if (response.data.code === 200) {
          const { token, admin } = response.data.data
          
          // 调试信息
          console.log('登录成功，获取到的令牌:', token)
          console.log('令牌类型:', typeof token)
          console.log('令牌长度:', token ? token.length : 0)
          console.log('令牌点数量:', token ? token.split('.').length - 1 : 0)
          
          // 保存登录状态，确保token完整保存
          if (token) {
            // 使用auth工具类保存令牌和管理员信息
            setToken(token)
            setAdmin(admin)
            
            if (loginForm.rememberMe) {
              localStorage.setItem('remember', 'true')
              localStorage.setItem('username', loginForm.username)
            } else {
              localStorage.removeItem('remember')
              localStorage.removeItem('username')
            }
            
            ElMessage.success('登录成功')
            router.push('/')
          } else {
            ElMessage.error('登录失败：未获取到有效的令牌')
          }
        } else {
          ElMessage.error(response.data.msg || '登录失败')
        }
      } catch (error) {
        console.error('登录失败', error)
        if (error.response) {
          console.error('错误状态码:', error.response.status)
          console.error('错误信息:', error.response.data)
          ElMessage.error(`登录失败: ${error.response.data.msg || error.message}`)
        } else {
          ElMessage.error('登录失败，请检查网络连接')
        }
      } finally {
        loading.value = false
      }
    }
  })
}

// 检查是否记住了用户名
onMounted(() => {
  // 初始化管理员账号
  initAdmin()
  
  // 检查是否记住了用户名
  const remember = localStorage.getItem('remember')
  if (remember === 'true') {
    const username = localStorage.getItem('username')
    if (username) {
      loginForm.username = username
      loginForm.rememberMe = true
    }
  }
})
</script>

<style lang="scss" scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f5f7fa;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  
  .login-card {
    width: 400px;
    padding: 20px 30px;
    border-radius: 8px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    background-color: rgba(255, 255, 255, 0.9);
    backdrop-filter: blur(10px);
    
    .login-header {
      display: flex;
      flex-direction: column;
      align-items: center;
      margin-bottom: 30px;
      
      .logo {
        width: 80px;
        height: 80px;
        margin-bottom: 15px;
        
        .logo-image {
          width: 100%;
          height: 100%;
          object-fit: contain;
        }
      }
      
      .title {
        font-size: 24px;
        font-weight: bold;
        color: #303133;
        margin: 0;
      }
    }
    
    .login-form {
      .remember-me {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 20px;
      }
      
      .login-button {
        width: 100%;
        padding: 12px 0;
        font-size: 16px;
      }
    }
    
    .login-tips {
      margin-top: 20px;
    }
  }
}
</style> 