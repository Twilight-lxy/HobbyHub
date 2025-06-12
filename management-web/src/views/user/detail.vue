<template>
  <div class="user-detail-container">
    <el-page-header @back="goBack" :title="'返回用户列表'" :content="'用户详情'" />
    
    <el-row :gutter="20" class="mt-20">
      <!-- 用户基本信息 -->
      <el-col :span="20">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>基本信息</span>
              <el-button type="primary" link @click="handleEdit">
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
            </div>
          </template>
          <div class="user-info">
            <div class="avatar-container">
              <el-avatar :size="100" :src="formatImage(userInfo.profilePicture)" />
            </div>
            <div class="info-list">
              <div class="info-item">
                <span class="label">用户名：</span>
                <span class="value">{{ userInfo.name }}</span>
              </div>
              <div class="info-item">
                <span class="label">昵称：</span>
                <span class="value">{{ userInfo.nickname }}</span>
              </div>
              <div class="info-item">
                <span class="label">小红书ID：</span>
                <span class="value">{{ userInfo.xhsId }}</span>
              </div>
              <div class="info-item">
                <span class="label">个人简介：</span>
                <span class="value">{{ userInfo.bio }}</span>
              </div>
              <div class="info-item">
                <span class="label">性别：</span>
                <span class="value">{{ userInfo.sex }}</span>
              </div>
              <div class="info-item">
                <span class="label">出生日期：</span>
                <span class="value">{{ formatDate(userInfo.birthdate) }}</span>
              </div>
              <div class="info-item">
                <span class="label">手机号码：</span>
                <span class="value">{{ userInfo.phone }}</span>
              </div>
              <div class="info-item">
                <span class="label">邮箱：</span>
                <span class="value">{{ userInfo.email }}</span>
              </div>
              <div class="info-item">
                <span class="label">地区：</span>
                <span class="value">{{ userInfo.region }}</span>
              </div>
              <div class="info-item">
                <span class="label">职业：</span>
                <span class="value">{{ userInfo.profession }}</span>
              </div>
              <div class="info-item">
                <span class="label">学校：</span>
                <span class="value">{{ userInfo.school }}</span>
              </div>
              <div class="info-item">
                <span class="label">角色：</span>
                <span class="value">
                  <el-tag v-if="userInfo.roleId === 1" type="danger">超级管理员</el-tag>
                  <el-tag v-else-if="userInfo.roleId === 2" type="warning">管理员</el-tag>
                  <el-tag v-else-if="userInfo.roleId === 3" type="success">普通用户</el-tag>
                  <el-tag v-else type="info">访客</el-tag>
                </span>
              </div>
              <div class="info-item">
                <span class="label">状态：</span>
                <span class="value">
                  <el-tag v-if="userInfo.status === 1" type="success">启用</el-tag>
                  <el-tag v-else type="danger">禁用</el-tag>
                </span>
              </div>
              <div class="info-item">
                <span class="label">注册时间：</span>
                <span class="value">{{ userInfo.createTime }}</span>
              </div>
            </div>
          </div>
          
          <!-- 背景图片 -->
          <div class="background-image-section" v-if="userInfo.backgroundImage">
            <h3>用户背景图片</h3>
            <el-image 
              :src="formatImage(userInfo.backgroundImage)" 
              fit="cover"
              style="width: 100%; max-height: 200px; border-radius: 8px;"
            />
          </div>
        </el-card>
      </el-col>
      
    </el-row>
    
    <!-- 编辑用户对话框 -->
    <el-dialog
      title="编辑用户"
      v-model="dialogVisible"
      width="600px"
    >
      <el-form :model="editForm" :rules="rules" ref="editFormRef" label-width="100px">
        <el-form-item label="用户名" prop="name">
          <el-input v-model="editForm.name" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="editForm.nickname" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item label="小红书ID" prop="xhsId">
          <el-input v-model="editForm.xhsId" placeholder="请输入小红书ID" />
        </el-form-item>
        <el-form-item label="个人简介" prop="bio">
          <el-input v-model="editForm.bio" type="textarea" :rows="3" placeholder="请输入个人简介" />
        </el-form-item>
        <el-form-item label="性别" prop="sex">
          <el-radio-group v-model="editForm.sex">
            <el-radio label="男">男</el-radio>
            <el-radio label="女">女</el-radio>
            <el-radio label="未知">未知</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="出生日期" prop="birthdate">
          <el-date-picker
            v-model="editForm.birthdate"
            type="date"
            placeholder="请选择出生日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="手机号码" prop="phone">
          <el-input v-model="editForm.phone" placeholder="请输入手机号码" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="editForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="地区" prop="region">
          <el-input v-model="editForm.region" placeholder="请输入地区" />
        </el-form-item>
        <el-form-item label="职业" prop="profession">
          <el-input v-model="editForm.profession" placeholder="请输入职业" />
        </el-form-item>
        <el-form-item label="学校" prop="school">
          <el-input v-model="editForm.school" placeholder="请输入学校" />
        </el-form-item>
        <el-form-item label="角色" prop="roleId">
          <el-select v-model="editForm.roleId" placeholder="请选择角色">
            <el-option v-for="role in roleOptions" :key="role.id" :label="role.name" :value="role.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="editForm.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="头像" prop="profilePicture">
          <el-input v-model="editForm.profilePicture" placeholder="请输入头像图片路径" />
          <div class="preview-image" v-if="editForm.profilePicture">
            <el-image :src="formatImage(editForm.profilePicture)" style="width: 100px; height: 100px;" />
          </div>
        </el-form-item>
        <el-form-item label="背景图片" prop="backgroundImage">
          <el-input v-model="editForm.backgroundImage" placeholder="请输入背景图片路径" />
          <div class="preview-image" v-if="editForm.backgroundImage">
            <el-image :src="formatImage(editForm.backgroundImage)" style="width: 200px; height: 100px;" />
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import request from '@/utils/request'

const route = useRoute()
const router = useRouter()

// 用户ID
const userId = ref(route.params.id)

// 用户信息
const userInfo = ref({
  id: userId.value,
  name: '',
  nickname: '',
  xhsId: '',
  bio: '',
  sex: '',
  birthdate: '',
  profilePicture: '',
  backgroundImage: '',
  phone: '',
  email: '',
  region: '',
  profession: '',
  school: '',
  roleId: 3,
  status: 1,
  createTime: ''
})

// 统计信息
const stats = reactive({
  activityCount: 0,
  teamCount: 0,
  loginCount: 0
})

// 角色选项
const roleOptions = [
  { id: 1, name: '超级管理员' },
  { id: 2, name: '管理员' },
  { id: 3, name: '普通用户' },
  { id: 4, name: '访客' }
]

// 参与的活动
const activities = ref([])

// 加入的小队
const teams = ref([])

// 对话框可见性
const dialogVisible = ref(false)

// 编辑表单
const editForm = reactive({
  id: '',
  name: '',
  nickname: '',
  xhsId: '',
  bio: '',
  sex: '',
  birthdate: '',
  phone: '',
  email: '',
  region: '',
  profession: '',
  school: '',
  profilePicture: '',
  backgroundImage: '',
  roleId: '',
  status: 1
})

// 表单引用
const editFormRef = ref(null)

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  roleId: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
}

// 图表引用
const loginChartRef = ref(null)
let loginChart = null

// 获取状态类型
const getStatusType = (status) => {
  const map = {
    0: 'info',
    1: 'success',
    2: 'warning',
    3: 'danger'
  }
  return map[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const map = {
    0: '未开始',
    1: '进行中',
    2: '已结束',
    3: '已取消'
  }
  return map[status] || '未知'
}

// 返回用户列表
const goBack = () => {
  router.push('/user/list')
}

// 查看活动详情
const viewActivity = (id) => {
  router.push(`/activity/detail/${id}`)
}

// 查看小队详情
const viewTeam = (id) => {
  router.push(`/team/detail/${id}`)
}

// 处理编辑
const handleEdit = () => {
  Object.assign(editForm, { ...userInfo.value })
  dialogVisible.value = true
}

// 提交表单
const submitForm = () => {
  editFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        // 调用后端API更新用户信息
        const res = await request.put('/api/user/admin/update', editForm)
        
        if (res && res.code === 200) {
          // 更新成功，刷新用户信息
          ElMessage.success('编辑成功')
          // 刷新用户详情
          getUserDetail()
          // 关闭对话框
          dialogVisible.value = false
        } else {
          ElMessage.error(res?.msg || '编辑失败')
        }
      } catch (error) {
        console.error('编辑用户失败', error)
        ElMessage.error('编辑用户失败')
      }
    }
  })
}

// 初始化登录图表
const initLoginChart = () => {
  if (!loginChartRef.value) return
  
  loginChart = echarts.init(loginChartRef.value)
  
  const option = {
    title: {
      text: '近30天登录记录',
      left: 'center'
    },
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: Array.from({ length: 30 }, (_, i) => {
        const date = new Date()
        date.setDate(date.getDate() - 29 + i)
        return `${date.getMonth() + 1}/${date.getDate()}`
      })
    },
    yAxis: {
      type: 'value',
      minInterval: 1
    },
    series: [
      {
        name: '登录次数',
        type: 'line',
        data: Array.from({ length: 30 }, () => Math.floor(Math.random() * 2)),
        markPoint: {
          data: [
            { type: 'max', name: '最大值' },
            { type: 'min', name: '最小值' }
          ]
        }
      }
    ]
  }
  
  loginChart.setOption(option)
}

// 监听窗口大小变化
const handleResize = () => {
  loginChart?.resize()
}

// 获取用户详情
const getUserDetail = async () => {
  try {
    // 从后端获取用户基本信息
    const res = await request.get(`/api/user/${userId.value}`)
    if (res && res.code === 200) {
      const userData = res.data
      if (userData) {
        // 映射后端返回的数据到前端显示模型
        userInfo.value = {
          id: userData.id,
          name: userData.name || '',
          nickname: userData.nickname || userData.name || '',
          xhsId: userData.xhsId || '',
          bio: userData.bio || '',
          sex: userData.sex || '未知',
          birthdate: userData.birthdate || '',
          phone: userData.phone || '',
          email: userData.email || '',
          region: userData.region || '',
          profession: userData.profession || '',
          school: userData.school || '',
          profilePicture: userData.profilePicture || '',
          backgroundImage: userData.backgroundImage || '',
          roleId: userData.roleId || 3,
          status: userData.status || 1,
          createTime: userData.createTime || ''
        }
      }
    } else {
      ElMessage.error(res?.msg || '获取用户详情失败')
    }
    
  } catch (error) {
    console.error('获取用户详情失败', error)
    ElMessage.error('获取用户详情失败')
  }
}

// 获取用户参与的活动
const getActivities = async () => {
  try {
    // 调用获取用户参与活动的API
    const res = await request.get(`/api/product_order/list/user`, {
      params: { userId: userId.value }
    })
    
    if (res && res.code === 200) {
      // 处理返回的活动数据
      activities.value = (res.data || []).map(item => {
        return {
          id: item.productId,
          name: item.product?.name || '未知活动',
          categoryName: item.product?.category ? item.product.category.split(',')[0] : '未分类',
          coverImage: item.product?.images ? item.product.images.split(',')[0] : 'https://via.placeholder.com/100x60',
          startTime: item.product?.startTime || '',
          endTime: item.product?.endTime || '',
          location: item.product?.address || '未知地点',
          status: item.status === 'paid' ? 1 : (item.status === 'created' ? 0 : 2),
          joinTime: item.createTime
        }
      })
    }
  } catch (error) {
    console.error('获取用户参与活动失败', error)
  }
}

// 获取用户加入的小队
const getTeams = async () => {
  try {
    // 调用获取用户小队的API
    const res = await request.get(`/api/product_cart/list/user`, {
      params: { userId: userId.value }
    })
    
    if (res && res.code === 200) {
      // 处理返回的小队数据
      teams.value = (res.data || []).map(item => {
        return {
          id: item.id,
          name: item.product?.name || '未命名小队',
          avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
          activityName: item.product?.name || '未知活动',
          memberCount: item.count || 0,
          maxCount: item.product?.price ? item.product.price * 10 : 20,
          createTime: item.createTime || '',
          joinTime: item.createTime || ''
        }
      })
    }
  } catch (error) {
    console.error('获取用户小队失败', error)
  }
}

// 获取用户统计信息
const getStats = async () => {
  stats.activityCount = activities.value.length
  stats.teamCount = teams.value.length
  
  try {
    // 这里可以调用获取登录次数的API，目前使用默认值
    stats.loginCount = 10
  } catch (error) {
    console.error('获取统计信息失败', error)
  }
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '未设置';
  
  // 如果日期字符串包含T，表示是ISO格式
  if (dateString.includes('T')) {
    const date = new Date(dateString);
    return date.getFullYear() + '-' + 
           String(date.getMonth() + 1).padStart(2, '0') + '-' + 
           String(date.getDate()).padStart(2, '0');
  }
  
  return dateString;
}

// 格式化图片路径
const formatImage = (imagePath) => {
  if (!imagePath) return 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png';
  
  // 如果是完整URL，直接返回
  if (imagePath.startsWith('http')) {
    return imagePath;
  }
  
  // 否则拼接后端服务器地址
  // 这里假设您的后端API地址是/api开头，实际环境可能需要调整
  return `http://localhost:8081${imagePath}`;
}

onMounted(() => {
  getUserDetail()
})
</script>

<style lang="scss" scoped>
.user-detail-container {
  .mt-20 {
    margin-top: 20px;
  }
  
  .box-card {
    margin-bottom: 20px;
    
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
  }
  
  .user-info {
    .avatar-container {
      display: flex;
      justify-content: center;
      margin-bottom: 20px;
    }
    
    .info-list {
      .info-item {
        display: flex;
        margin-bottom: 15px;
        
        .label {
          width: 100px;
          color: #606266;
        }
        
        .value {
          flex: 1;
          color: #303133;
        }
      }
    }
  }
  
  .background-image-section {
    margin-top: 30px;
    
    h3 {
      font-size: 16px;
      margin-bottom: 15px;
      color: #303133;
    }
  }
  
  .preview-image {
    margin-top: 10px;
  }
}
</style>