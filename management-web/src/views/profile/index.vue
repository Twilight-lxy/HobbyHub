<template>
  <div class="profile-container">
    <el-card class="profile-card">
      <template #header>
        <div class="card-header">
          <h3>个人信息</h3>
          <el-button 
            type="primary" 
            @click="handleEdit"
            :class="{ 'edit-button-active': isEditing }"
          >
            {{ isEditing ? '保存修改' : '编辑信息' }}
          </el-button>
        </div>
      </template>
      
      <el-row :gutter="30" class="profile-row">
        <el-col :span="8" class="avatar-column">
          <div class="avatar-wrapper">
            <el-avatar 
              :size="140" 
              :src="adminInfo.avatar || defaultAvatar" 
              class="profile-avatar"
              @click="openAvatarUpload"
            />
            <div class="avatar-upload-overlay" v-if="isEditing">
              <el-icon class="avatar-edit-icon">
                <Edit />
              </el-icon>
              <div class="avatar-upload-text">更换头像</div>
            </div>
          </div>
        </el-col>
        
        <el-col :span="16" class="info-column">
          <div v-if="!isEditing" class="info-display">
            <div class="info-item">
              <span class="label">用户名：</span>
              <span class="value">{{ adminInfo.username }}</span>
            </div>
            <div class="info-item">
              <span class="label">昵称：</span>
              <span class="value">{{ adminInfo.nickname || '未设置' }}</span>
            </div>
            <div class="info-item">
              <span class="label">手机号：</span>
              <span class="value">{{ adminInfo.phone || '未设置' }}</span>
            </div>
            <div class="info-item">
              <span class="label">邮箱：</span>
              <span class="value">{{ adminInfo.email || '未设置' }}</span>
            </div>
            <div class="info-item">
              <span class="label">角色：</span>
              <span class="value">系统管理员</span>
            </div>
            <div class="info-item">
              <span class="label">创建时间：</span>
              <span class="value">{{ formatDate(adminInfo.createTime) }}</span>
            </div>
            <div class="info-item">
              <span class="label">最后登录：</span>
              <span class="value">{{ formatDate(adminInfo.lastLoginTime) }}</span>
            </div>
          </div>
          
          <el-form
            v-else
            ref="formRef"
            :model="form"
            :rules="rules"
            label-width="100px"
            class="edit-form"
          >
            <el-form-item label="用户名" prop="username">
              <el-input v-model="form.username" disabled class="disabled-input" />
            </el-form-item>
            <el-form-item label="昵称" prop="nickname">
              <el-input v-model="form.nickname" placeholder="请输入昵称" />
            </el-form-item>
            <el-form-item label="手机号" prop="phone">
              <el-input v-model="form.phone" placeholder="请输入手机号" />
            </el-form-item>
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="form.email" placeholder="请输入邮箱" />
            </el-form-item>
            <el-form-item class="form-actions">
              <el-button @click="cancelEdit">取消</el-button>
              <el-button type="primary" @click="submitForm" :loading="formLoading">
                保存
              </el-button>
            </el-form-item>
          </el-form>
        </el-col>
      </el-row>
    </el-card>
    
    <el-card class="profile-card password-card">
      <template #header>
        <div class="card-header">
          <h3>修改密码</h3>
        </div>
      </template>
      
      <el-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-width="120px"
        class="password-form"
      >
        <el-form-item label="原密码" prop="oldPassword">
          <el-input
            v-model="passwordForm.oldPassword"
            type="password"
            placeholder="请输入原密码"
            show-password
            class="password-input"
          />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="passwordForm.newPassword"
            type="password"
            placeholder="请输入新密码"
            show-password
            class="password-input"
          />
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirmPassword">
          <el-input
            v-model="passwordForm.confirmPassword"
            type="password"
            placeholder="请确认新密码"
            show-password
            class="password-input"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="updatePassword" :loading="passwordLoading">
            更新密码
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'
import { Edit } from '@element-plus/icons-vue'

// 默认头像
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

// 管理员信息
const adminInfo = ref({
  id: '',
  username: '',
  nickname: '',
  avatar: '',
  phone: '',
  email: '',
  createTime: '',
  lastLoginTime: ''
})

// 编辑状态
const isEditing = ref(false)
const formRef = ref(null)
const formLoading = ref(false)
const passwordLoading = ref(false)

const form = reactive({
  username: '',
  nickname: '',
  phone: '',
  email: '',
  avatar: ''
})

// 表单验证规则
const rules = {
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

// 密码表单
const passwordFormRef = ref(null)
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 密码验证规则
const passwordRules = {
  oldPassword: [
    { required: true, message: '请输入原密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 获取管理员信息
const getAdminInfo = async () => {
  try {
    // 模拟数据请求
    await new Promise(resolve => setTimeout(resolve, 500))
    
    // 模拟数据
    adminInfo.value = {
      id: 1,
      username: 'admin',
      nickname: '系统管理员',
      avatar: '',
      phone: '13800138000',
      email: 'admin@example.com',
      createTime: '2023-01-15 08:30:00',
      lastLoginTime: '2025-06-16 11:20:00'
    }
    
    // 初始化表单数据
    Object.assign(form, {
      username: adminInfo.value.username,
      nickname: adminInfo.value.nickname,
      phone: adminInfo.value.phone,
      email: adminInfo.value.email,
      avatar: adminInfo.value.avatar
    })
  } catch (error) {
    console.error('获取管理员信息失败', error)
    ElMessage.error('获取管理员信息失败，请检查网络连接')
  }
}

// 编辑信息
const handleEdit = () => {
  isEditing.value = true
}

// 取消编辑
const cancelEdit = () => {
  isEditing.value = false
  // 重置表单
  Object.assign(form, {
    username: adminInfo.value.username,
    nickname: adminInfo.value.nickname,
    phone: adminInfo.value.phone,
    email: adminInfo.value.email,
    avatar: adminInfo.value.avatar
  })
}

// 提交表单
const submitForm = () => {
  formRef.value?.validate(async (valid) => {
    if (valid) {
      formLoading.value = true
      
      try {
        // 模拟提交请求
        await new Promise(resolve => setTimeout(resolve, 800))
        
        ElMessage.success('更新成功')
        // 更新本地管理员信息
        adminInfo.value.nickname = form.nickname
        adminInfo.value.phone = form.phone
        adminInfo.value.email = form.email
        adminInfo.value.avatar = form.avatar
        isEditing.value = false
      } catch (error) {
        console.error('更新失败', error)
        ElMessage.error('更新失败，请检查网络连接')
      } finally {
        formLoading.value = false
      }
    }
  })
}

// 更新密码
const updatePassword = () => {
  passwordFormRef.value?.validate(async (valid) => {
    if (valid) {
      passwordLoading.value = true
      
      try {
        // 模拟密码更新请求
        await new Promise(resolve => setTimeout(resolve, 800))
        
        ElMessageBox.confirm(
          '密码更新成功，请重新登录以应用新密码',
          '提示',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'success'
          }
        ).then(() => {
          // 清除登录状态
          localStorage.removeItem('token')
          localStorage.removeItem('adminInfo')
          // 跳转到登录页
          setTimeout(() => {
            window.location.href = '/login'
          }, 1000)
        })
      } catch (error) {
        console.error('密码更新失败', error)
        ElMessage.error('密码更新失败，请检查网络连接')
      } finally {
        passwordLoading.value = false
      }
    }
  })
}

// 打开头像上传
const openAvatarUpload = () => {
  if (isEditing.value) {
    ElMessage.info('点击更换头像')
    // 这里可以添加头像上传逻辑
  }
}

// 格式化日期
const formatDate = (date) => {
  if (!date) return '未知'
  return new Date(date).toLocaleString()
}

// 生命周期钩子
onMounted(() => {
  getAdminInfo()
})
</script>

<style lang="scss" scoped>
:root {
  --primary-color: #722ED1;
  --primary-light: #8B5CF6;
  --secondary-color: #36CFC9;
  --text-color: #303133;
  --text-secondary: #606266;
  --bg-color: #f9fafc;
  --card-bg: #ffffff;
  --border-color: #ebeef5;
  --success-color: #67c23a;
}

.profile-container {
  padding: 30px;
  background-color: #f5f7fa;
  min-height: 100vh;
  
  .profile-card {
    border-radius: 16px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.08);
    transition: all 0.3s ease;
    margin-bottom: 30px;
    background-color: var(--card-bg);
    
    &:hover {
      box-shadow: 0 15px 35px rgba(0, 0, 0, 0.12);
    }
    
    .password-card {
      margin-top: 20px;
    }
    
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding-bottom: 0;
      
      h3 {
        margin: 0;
        font-size: 20px;
        font-weight: 600;
        color: var(--text-color);
      }
      
      .edit-button-active {
        background-color: var(--primary-light);
        border-color: var(--primary-light);
        
        &:hover {
          background-color: var(--primary-color);
          border-color: var(--primary-color);
        }
      }
    }
    
    .profile-row {
      padding: 30px 0;
    }
    
    .avatar-column {
      display: flex;
      justify-content: center;
      align-items: center;
      
      .avatar-wrapper {
        position: relative;
        cursor: pointer;
        
        .profile-avatar {
          transition: all 0.3s ease;
          box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
          
          &:hover {
            transform: scale(1.05);
            box-shadow: 0 10px 25px rgba(114, 46, 209, 0.2);
          }
        }
        
        .avatar-upload-overlay {
          position: absolute;
          bottom: 0;
          left: 0;
          width: 100%;
          height: 40%;
          background: linear-gradient(to top, rgba(0, 0, 0, 0.7), transparent);
          border-radius: 0 0 10px 10px;
          display: flex;
          flex-direction: column;
          justify-content: center;
          align-items: center;
          opacity: 0;
          transition: opacity 0.3s ease;
          
          .avatar-edit-icon {
            font-size: 24px;
            color: white;
            margin-bottom: 5px;
          }
          
          .avatar-upload-text {
            color: white;
            font-size: 14px;
          }
        }
        
        &:hover .avatar-upload-overlay {
          opacity: 1;
        }
      }
    }
    
    .info-column {
      .info-display {
        .info-item {
          display: flex;
          align-items: center;
          margin-bottom: 20px;
          padding-bottom: 20px;
          border-bottom: 1px solid #f0f2f5;
          
          &:last-child {
            margin-bottom: 0;
            padding-bottom: 0;
            border-bottom: none;
          }
          
          .label {
            font-weight: 500;
            color: var(--text-secondary);
            min-width: 80px;
            font-size: 15px;
          }
          
          .value {
            color: var(--text-color);
            font-size: 15px;
            flex: 1;
          }
        }
      }
      
      .edit-form {
        .disabled-input {
          background-color: #f5f7fa;
          border-color: #dcdfe6;
          cursor: not-allowed;
        }
        
        .form-actions {
          display: flex;
          justify-content: flex-end;
          gap: 10px;
        }
      }
    }
  }
  
  .password-form {
    .password-input {
      .el-input__inner {
        padding: 12px 15px;
      }
    }
  }
}
</style>