<template>
  <div class="login-container">
    <!-- 背景装饰元素 -->
    <div class="bg-decorator top-left"></div>
    <div class="bg-decorator bottom-right"></div>
    
    <div class="login-card-wrapper">
      <el-card class="login-card">
        <div class="login-header">
          <div class="logo">
            <img src="@/assets/logo.svg" alt="Logo" class="logo-image" />
          </div>
          <h2 class="title">兴趣小队管理系统</h2>
          <p class="subtitle">请登录您的账户</p>
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
              class="form-input"
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
              class="form-input"
            />
          </el-form-item>
          
          <el-form-item class="remember-me">
            <el-checkbox v-model="loginForm.rememberMe" class="custom-checkbox">记住我</el-checkbox>
          </el-form-item>
          
          <el-form-item>
            <el-button
              type="primary"
              :loading="loading"
              class="login-button"
              @click="handleLogin"
            >
              <span v-if="!loading">登录</span>
              <span v-else>正在登录...</span>
            </el-button>
          </el-form-item>
        </el-form>
        
        <div class="login-tips">
          <el-alert
            title="HobbyHub的管理账号如下：账号：admin，密码：admin"
            type="info"
            :closable="false"
            show-icon
            class="custom-alert"
          />
        </div>

        
      </el-card>
      
      <!-- 页脚 -->
      <div class="login-footer">
        <p>© 2025 HobbyHub管理系统</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'
import axios from 'axios'
import { setToken, setAdmin,getToken } from '@/utils/auth'

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
    { min: 5, max: 20, message: '长度在 5 到 20 个字符', trigger: 'blur' }
  ]
}

// 初始化管理员账号
const initAdmin = async () => {
  try {
    // 使用axios直接发送请求，不使用request工具
    const response = await axios.post('/api/v1/admin/login')
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
        const response = await axios.post('/api/v1/admin/login', {
          username: loginForm.username,
          password: loginForm.password
        })
        if (response.status === 200) {
          const { token } = response.data.token
          // 调试信息
          console.log('登录成功，获取到的令牌:', token)
          // 保存登录状态，确保token完整保存
          if (token!=0) {
            // 使用auth工具类保存令牌和管理员信息
            setToken(token)
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
/* 背景样式 */
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  position: relative;
  overflow: hidden;
  
  /* 渐变背景 */
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  
  /* 背景装饰元素 */
  .bg-decorator {
    position: absolute;
    width: 500px;
    height: 500px;
    border-radius: 50%;
    filter: blur(100px);
    opacity: 0.5;
    z-index: 0;
  }
  
  .top-left {
    top: -250px;
    left: -250px;
    background: linear-gradient(135deg, #8a2be2 0%, #9370db 100%);
  }
  
  .bottom-right {
    bottom: -250px;
    right: -250px;
    background: linear-gradient(135deg, #4169e1 0%, #00bfff 100%);
  }
}

/* 卡片包装器 */
.login-card-wrapper {
  position: relative;
  z-index: 1;
  width: 400px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  
  /* 卡片悬浮动画 */
  .login-card {
    width: 100%;
    padding: 40px 40px;
    border-radius: 20px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.15);
    background-color: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.2);
    transition: transform 0.3s ease, box-shadow 0.3s ease;
    
    &:hover {
      transform: translateY(-5px);
      box-shadow: 0 15px 40px rgba(0, 0, 0, 0.2);
    }
  }
  
  /* 页脚样式 */
  .login-footer {
    margin-top: 20px;
    color: rgba(255, 255, 255, 0.7);
    font-size: 12px;
  }
}

/* 头部样式 */
.login-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 40px;
  
  .logo {
    width: 90px;
    height: 90px;
    margin-bottom: 20px;
    border-radius: 50%;
    box-shadow: 0 8px 20px rgba(102, 126, 234, 0.25);
    overflow: hidden;
    
    .logo-image {
      width: 100%;
      height: 100%;
      object-fit: contain;
      padding: 10px;
      background-color: white;
    }
  }
  
  .title {
    font-size: 28px;
    font-weight: bold;
    color: #303133;
    margin: 0;
    margin-bottom: 8px;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  }
  
  .subtitle {
    font-size: 14px;
    color: #606266;
    margin: 0;
  }
}

/* 表单样式 */
.login-form {
  .form-input {
    /* 自定义输入框样式 */
    .el-input__wrapper {
      border-radius: 12px;
      padding: 12px 15px;
      border: 1px solid #dcdfe6;
      box-shadow: 0 2px 5px rgba(0, 0, 0, 0.05);
      transition: border-color 0.3s ease, box-shadow 0.3s ease;
      
      &:hover {
        border-color: #c0c4cc;
      }
      
      &:focus-within {
        border-color: #409eff;
        box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.1);
      }
    }
    
    .el-input__prefix {
      margin-right: 8px;
      color: #909399;
    }
  }
  
  .remember-me {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 25px;
    margin-top: 5px;
    
    .custom-checkbox {
      .el-checkbox__label {
        color: #606266;
      }
      
      .el-checkbox__input.is-checked .el-checkbox__inner {
        background-color: #409eff;
        border-color: #409eff;
      }
    }
  }
  
  .login-button {
    width: 100%;
    padding: 14px 0;
    font-size: 16px;
    border-radius: 12px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border: none;
    box-shadow: 0 8px 20px rgba(102, 126, 234, 0.3);
    transition: all 0.3s ease;
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 10px 25px rgba(102, 126, 234, 0.4);
      background: linear-gradient(135deg, #5a6ee2 0%, #6c3fb2 100%);
    }
    
    &:active {
      transform: translateY(0);
      box-shadow: 0 4px 10px rgba(102, 126, 234, 0.3);
    }
  }
}

/* 提示框样式 */
.login-tips {
  margin-top: 25px;
  
  .custom-alert {
    border-radius: 10px;
    padding: 12px 15px;
    background-color: rgba(64, 158, 255, 0.05);
    border-color: rgba(64, 158, 255, 0.1);
    
    .el-alert__title {
      font-size: 13px;
      color: #606266;
    }
    
    .el-alert__icon {
      color: #409eff;
    }
  }
}

/* 表单验证错误样式 */
.el-form-item.is-error .el-input__wrapper {
  border-color: #f56c6c;
}

.el-form-item__error {
  padding-top: 5px;
  font-size: 12px;
}
</style>
