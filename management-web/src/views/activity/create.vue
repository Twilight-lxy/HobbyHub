<template>
  <div class="activity-create-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <el-row :gutter="20">
        <el-col :span="20">
          <h1 class="page-title">创建新活动</h1>
        </el-col>
        <el-col :span="4" class="text-right">
          <el-button @click="goBack" size="medium">
            <el-icon><ArrowLeft /></el-icon> 返回
          </el-button>
        </el-col>
      </el-row>
    </div>

    <!-- 主内容区 -->
    <el-card class="create-card" shadow="hover">
      <el-form
        ref="activityFormRef"
        :model="activityForm"
        :rules="activityRules"
        label-width="120px"
        class="activity-form"
        @submit.prevent="submitForm"
      >
        <!-- 活动ID (自动生成) -->
        <el-form-item label="活动ID">
          <el-input
            v-model="activityForm.id"
            placeholder="自动生成"
            disabled
            readonly
            class="form-input"
          />
        </el-form-item>

        <!-- 活动名称 -->
        <el-form-item label="活动名称" prop="name">
          <el-input
            v-model="activityForm.name"
            placeholder="请输入活动名称"
            class="form-input"
            @input="handleInput"
          />
        </el-form-item>

        <!-- 活动地点 -->
        <el-form-item label="活动地点" prop="addr">
          <el-input
            v-model="activityForm.addr"
            placeholder="请输入活动地点"
            class="form-input"
          />
        </el-form-item>

        <!-- 创建者ID (用户输入) -->
        <el-form-item label="创建者ID" prop="userId">
          <el-input
            v-model="activityForm.userId"
            placeholder="请输入创建者ID"
            class="form-input"
          />
        </el-form-item>

        <!-- 创建时间 (自动生成) -->
        <el-form-item label="创建时间">
          <el-input
            v-model="activityForm.createTime"
            placeholder="自动生成"
            disabled
            readonly
            class="form-input"
          />
        </el-form-item>

        <!-- 开始时间 -->
        <el-form-item label="开始时间" prop="startTime">
          <el-date-picker
            v-model="activityForm.startTime"
            type="datetime"
            placeholder="请选择开始日期和时间"
            value-format="YYYY-MM-DD HH:mm:ss"
            class="form-input"
          />
        </el-form-item>

        <!-- 活动简介 -->
        <el-form-item label="活动简介" prop="intro">
          <el-input
            v-model="activityForm.intro"
            type="textarea"
            :rows="4"
            placeholder="请输入活动简介"
            class="form-textarea"
          />
        </el-form-item>

        <!-- 操作按钮 -->
        <el-form-item style="text-align: right; padding-right: 40px;">
          <el-button type="primary" :loading="loading" @click="submitForm">
            <span v-if="!loading">创建活动</span>
            <span v-else>创建中...</span>
          </el-button>
          <el-button @click="resetForm">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { getToken } from '@/utils/auth'

const router = useRouter()
const activityFormRef = ref(null)
const loading = ref(false)

// 活动表单数据
const activityForm = reactive({
  id: '',                 // 活动ID
  name: '',               // 活动名称
  addr: '',               // 活动地点
  userId: '',             // 创建者ID
  createTime: '',         // 创建时间
  startTime: '',          // 开始时间
  intro: ''               // 活动简介
})

// 表单验证规则（新增创建者ID数字验证）
const activityRules = {
  name: [
    { required: true, message: '请输入活动名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  addr: [
    { required: true, message: '请输入活动地点', trigger: 'blur' },
    { min: 2, max: 100, message: '长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  userId: [
    { required: true, message: '请输入创建者ID', trigger: 'blur' },
    { pattern: /^\d+$/, message: '创建者ID必须为数字', trigger: 'blur' } // 新增数字验证
  ],
  startTime: [
    { required: true, message: '请选择开始时间', trigger: 'change' }
  ],
  intro: [
    { max: 500, message: '长度不能超过 500 个字符', trigger: 'blur' }
  ]
}

// 处理返回
const goBack = () => {
  router.push('/activity')
}

// 处理输入
const handleInput = (val) => {
  // 可以在这里添加实时输入处理逻辑
}

// 重置表单
const resetForm = () => {
  activityFormRef.value?.resetFields()
  activityForm.id = ''
  activityForm.createTime = ''
}

// 提交表单
const submitForm = () => {
  activityFormRef.value?.validate(async (valid) => {
    if (valid) {
      loading.value = true
      
      try {
        // 生成活动ID（实际项目中应从后端获取）
        activityForm.id = generateActivityId()
        
        // 填充创建时间
        activityForm.createTime = formatDate(new Date())
        
        // 将创建者ID转换为数字类型
        const userId = parseInt(activityForm.userId, 10)
        
        // 发送创建活动请求
        const response = await axios.put('/api/v1/admin/activity', 
          {
            id:1,
            name:activityForm.name,
            addr:activityForm.addr,
            userId:userId, // 使用转换后的数字类型
            intro:activityForm.intro
            
          },
          {
            headers: {
              Authorization: getToken()
            }
          }
        );

        if (response.status === 200) {
          ElMessage.success('活动创建成功')
          router.push('/activity')
        } else {
          ElMessage.error(response.data.msg || '活动创建失败')
        }
      } catch (error) {
        console.error('活动创建失败', error)
        ElMessage.error('活动创建失败，请稍后重试')
      } finally {
        loading.value = false
      }
    }
  })
}

// 生成活动ID（示例方法）
const generateActivityId = () => {
  return 'ACT' + Date.now().toString().substring(6)
}

// 原生日期格式化函数
const formatDate = (date) => {
  if (!date) return ''
  
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}:${String(d.getSeconds()).padStart(2, '0')}`
}

// 组件挂载时初始化
onMounted(() => {
  // 自动填充创建时间
  activityForm.createTime = formatDate(new Date())
})
</script>

<style lang="scss" scoped>
// 样式保持不变...
</style>