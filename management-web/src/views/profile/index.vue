<template>
  <div class="profile-container">
    <el-card class="profile-card">
      <template #header>
        <div class="card-header">
          <h3>个人信息</h3>
          <el-button type="primary" @click="handleEdit">编辑信息</el-button>
        </div>
      </template>
      
      <el-row :gutter="20">
        <el-col :span="8" class="avatar-container">
          <el-avatar :size="120" :src="adminInfo.avatar || defaultAvatar" />
          <div class="upload-avatar" v-if="isEditing">
            <el-upload
              class="avatar-uploader"
              action="/api/file/upload"
              :show-file-list="false"
              :on-success="handleAvatarSuccess"
              :before-upload="beforeAvatarUpload"
            >
              <el-button type="primary" size="small">更换头像</el-button>
            </el-upload>
          </div>
        </el-col>
        
        <el-col :span="16">
          <div v-if="!isEditing" class="info-display">
            <div class="info-item">
              <span class="label">用户名：</span>
              <span class="value">{{ adminInfo.username }}</span>
            </div>
            <div class="info-item">
              <span class="label">昵称：</span>
              <span class="value">{{ adminInfo.nickname }}</span>
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
            label-width="80px"
          >
            <el-form-item label="用户名" prop="username">
              <el-input v-model="form.username" disabled />
            </el-form-item>
            <el-form-item label="昵称" prop="nickname">
              <el-input v-model="form.nickname" />
            </el-form-item>
            <el-form-item label="手机号" prop="phone">
              <el-input v-model="form.phone" />
            </el-form-item>
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="form.email" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="submitForm">保存</el-button>
              <el-button @click="cancelEdit">取消</el-button>
            </el-form-item>
          </el-form>
        </el-col>
      </el-row>
    </el-card>
    
    <el-card class="profile-card">
      <template #header>
        <div class="card-header">
          <h3>修改密码</h3>
        </div>
      </template>
      
      <el-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-width="100px"
      >
        <el-form-item label="原密码" prop="oldPassword">
          <el-input
            v-model="passwordForm.oldPassword"
            type="password"
            placeholder="请输入原密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="passwordForm.newPassword"
            type="password"
            placeholder="请输入新密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirmPassword">
          <el-input
            v-model="passwordForm.confirmPassword"
            type="password"
            placeholder="请确认新密码"
            show-password
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
import { ElMessage } from 'element-plus'
import axios from 'axios'

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
const passwordLoading = ref(false)
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
    const response = await axios.get('/api/admin/info', {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('token')}`
      }
    })
    
    if (response.data.code === 200) {
      adminInfo.value = response.data.data
      // 初始化表单数据
      Object.assign(form, {
        username: adminInfo.value.username,
        nickname: adminInfo.value.nickname,
        phone: adminInfo.value.phone,
        email: adminInfo.value.email,
        avatar: adminInfo.value.avatar
      })
    } else {
      ElMessage.error('获取管理员信息失败')
    }
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
      try {
        const response = await axios.put('/api/admin/update', form, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`
          }
        })
        
        if (response.data.code === 200) {
          ElMessage.success('更新成功')
          // 更新本地存储的管理员信息
          const admin = JSON.parse(localStorage.getItem('adminInfo') || '{}')
          admin.nickname = form.nickname
          admin.avatar = form.avatar
          localStorage.setItem('adminInfo', JSON.stringify(admin))
          // 更新显示的管理员信息
          getAdminInfo()
          isEditing.value = false
        } else {
          ElMessage.error(response.data.msg || '更新失败')
        }
      } catch (error) {
        console.error('更新失败', error)
        ElMessage.error('更新失败，请检查网络连接')
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
        const response = await axios.put('/api/admin/updatePassword', {
          oldPassword: passwordForm.oldPassword,
          newPassword: passwordForm.newPassword
        }, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`
          }
        })
        
        if (response.data.code === 200) {
          ElMessage.success('密码更新成功，请重新登录')
          // 清除登录状态
          localStorage.removeItem('token')
          localStorage.removeItem('adminInfo')
          // 跳转到登录页
          setTimeout(() => {
            window.location.href = '/login'
          }, 1500)
        } else {
          ElMessage.error(response.data.msg || '密码更新失败')
        }
      } catch (error) {
        console.error('密码更新失败', error)
        ElMessage.error('密码更新失败，请检查网络连接')
      } finally {
        passwordLoading.value = false
      }
    }
  })
}

// 头像上传前的校验
const beforeAvatarUpload = (file) => {
  const isJPG = file.type === 'image/jpeg'
  const isPNG = file.type === 'image/png'
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isJPG && !isPNG) {
    ElMessage.error('上传头像图片只能是 JPG 或 PNG 格式!')
  }
  if (!isLt2M) {
    ElMessage.error('上传头像图片大小不能超过 2MB!')
  }
  return (isJPG || isPNG) && isLt2M
}

// 头像上传成功回调
const handleAvatarSuccess = (res, file) => {
  if (res.code === 200) {
    form.avatar = res.data
    ElMessage.success('头像上传成功')
  } else {
    ElMessage.error(res.msg || '头像上传失败')
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
.profile-container {
  padding: 20px;
  
  .profile-card {
    margin-bottom: 20px;
    
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      h3 {
        margin: 0;
        font-size: 18px;
        font-weight: 500;
      }
    }
    
    .avatar-container {
      display: flex;
      flex-direction: column;
      align-items: center;
      
      .el-avatar {
        margin-bottom: 15px;
      }
      
      .upload-avatar {
        margin-top: 10px;
      }
    }
    
    .info-display {
      .info-item {
        margin-bottom: 15px;
        
        .label {
          font-weight: bold;
          margin-right: 10px;
          color: #606266;
        }
        
        .value {
          color: #303133;
        }
      }
    }
  }
}
</style> 