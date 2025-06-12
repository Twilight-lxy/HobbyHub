<template>
  <div class="activity-form-container">
    <el-page-header @back="goBack" :title="'返回活动列表'" :content="isEdit ? '编辑活动' : '创建活动'" />
    
    <el-card class="form-card mt-20" v-loading="loading">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        label-position="right"
      >
        <el-form-item label="活动名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入活动名称" />
        </el-form-item>
        
        <el-form-item label="活动分类" prop="categoryId">
          <el-select v-model="form.categoryId" placeholder="请选择活动分类" style="width: 100%">
            <el-option
              v-for="item in categoryOptions"
              :key="item.id"
              :label="item.name"
              :value="item.name"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="活动时间" prop="timeRange">
          <el-date-picker
            v-model="form.timeRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            style="width: 100%"
          />
        </el-form-item>
        
        <el-form-item label="活动地点" prop="location">
          <el-input v-model="form.location" placeholder="请输入活动地点" />
        </el-form-item>
        
        <el-form-item label="最大参与人数" prop="maxParticipants">
          <el-input-number v-model="form.maxParticipants" :min="1" :max="1000" />
        </el-form-item>
        
        <el-form-item label="报名截止时间" prop="registrationDeadline">
          <el-date-picker
            v-model="form.registrationDeadline"
            type="datetime"
            placeholder="请选择报名截止时间"
            style="width: 100%"
          />
        </el-form-item>
        
        <el-form-item label="活动状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="0">未开始</el-radio>
            <el-radio :label="1">进行中</el-radio>
            <el-radio :label="2">已结束</el-radio>
            <el-radio :label="3">已取消</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="活动标签" prop="tags">
          <el-tag
            v-for="tag in form.tags"
            :key="tag"
            closable
            @close="handleRemoveTag(tag)"
            style="margin-right: 10px; margin-bottom: 10px;"
          >
            {{ tag }}
          </el-tag>
          <el-input
            v-if="inputTagVisible"
            ref="tagInputRef"
            v-model="inputTagValue"
            class="input-new-tag"
            size="small"
            @keyup.enter="handleInputConfirm"
            @blur="handleInputConfirm"
          />
          <el-button v-else @click="showTagInput" size="small">+ 添加标签</el-button>
        </el-form-item>
        
        <el-form-item label="活动简介" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入活动简介"
          />
        </el-form-item>
        
        <el-form-item label="活动详情" prop="content">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="6"
            placeholder="请输入活动详情"
          />
        </el-form-item>
        
        <el-form-item label="封面图片" prop="coverImage">
          <el-upload
            class="avatar-uploader"
            action="/api/file/upload"
            :show-file-list="false"
            :on-success="handleCoverSuccess"
            :before-upload="beforeCoverUpload"
          >
            <el-image v-if="displayCoverImage" :src="displayCoverImage" class="cover-image" />
            <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
          </el-upload>
          <div class="upload-tip">建议上传尺寸比例为 16:9 的图片</div>
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="submitForm">保存</el-button>
          <el-button @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'
import { formatImageUrl } from '@/utils/image'

const route = useRoute()
const router = useRouter()

// 表单引用
const formRef = ref(null)
const tagInputRef = ref(null)

// 加载状态
const loading = ref(false)

// 是否是编辑模式
const isEdit = computed(() => !!route.params.id)

// 活动ID
const activityId = ref(route.params.id)

// 标签输入
const inputTagVisible = ref(false)
const inputTagValue = ref('')

// 分类选项
const categoryOptions = ref([])

// 表单数据
const form = reactive({
  name: '',
  categoryId: null,
  timeRange: [],
  location: '',
  maxParticipants: 20,
  registrationDeadline: '',
  status: 0,
  tags: [],
  description: '',
  content: '',
  coverImage: ''
})

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入活动名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  categoryId: [
    { required: true, message: '请选择活动分类', trigger: 'change' }
  ],
  timeRange: [
    { required: true, message: '请选择活动时间', trigger: 'change' },
    { 
      validator: (rule, value, callback) => {
        if (value && value.length === 2) {
          if (value[0] >= value[1]) {
            callback(new Error('开始时间必须早于结束时间'))
          } else {
            callback()
          }
        } else {
          callback()
        }
      }, 
      trigger: 'change' 
    }
  ],
  location: [
    { required: true, message: '请输入活动地点', trigger: 'blur' }
  ],
  maxParticipants: [
    { required: true, message: '请输入最大参与人数', trigger: 'blur' },
    { type: 'number', min: 1, message: '参与人数必须大于0', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入活动简介', trigger: 'blur' },
    { max: 200, message: '简介不能超过200个字符', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入活动详情', trigger: 'blur' }
  ],
  coverImage: [
    { required: true, message: '请上传封面图片', trigger: 'change' }
  ]
}

// 返回活动列表
const goBack = () => {
  router.push('/activity/list')
}

// 添加标签
const showTagInput = () => {
  inputTagVisible.value = true
  nextTick(() => {
    tagInputRef.value?.focus()
  })
}

// 确认添加标签
const handleInputConfirm = () => {
  if (inputTagValue.value) {
    if (form.tags.indexOf(inputTagValue.value) === -1) {
      form.tags.push(inputTagValue.value)
    }
  }
  inputTagVisible.value = false
  inputTagValue.value = ''
}

// 移除标签
const handleRemoveTag = (tag) => {
  const index = form.tags.indexOf(tag)
  if (index !== -1) {
    form.tags.splice(index, 1)
  }
}

// 处理封面上传成功
const handleCoverSuccess = (res) => {
  if (res && res.code === 200 && res.data) {
    form.coverImage = res.data.url
    ElMessage.success('上传成功')
  } else {
    ElMessage.error(res?.msg || '上传失败')
  }
}

// 上传前检查
const beforeCoverUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isImage) {
    ElMessage.error('上传封面图片只能是图片格式!')
  }
  if (!isLt2M) {
    ElMessage.error('上传封面图片大小不能超过 2MB!')
  }
  return isImage && isLt2M
}

// 获取分类列表
const getCategoryOptions = async () => {
  try {
    const res = await request.get('/api/tabs/list')
    if (res && res.code === 200) {
      categoryOptions.value = res.data || []
    } else {
      ElMessage.error(res?.msg || '获取分类列表失败')
    }
  } catch (error) {
    console.error('获取分类列表失败', error)
    ElMessage.error('获取分类列表失败')
  }
}

// 提交表单
const submitForm = () => {
  formRef.value?.validate(async (valid) => {
    if (valid) {
      try {
        loading.value = true
        
        // 处理时间范围
        const formData = { ...form }
        if (form.timeRange && form.timeRange.length === 2) {
          formData.startTime = form.timeRange[0]
          formData.endTime = form.timeRange[1]
        }
        delete formData.timeRange
        
        // 处理分类和标签
        formData.category = form.categoryId ? [form.categoryId] : []
        
        // 处理图片
        formData.images = form.coverImage ? [form.coverImage] : []
        
        // 发送请求
        let res
        if (isEdit.value) {
          res = await request.put(`/api/product/${activityId.value}`, formData)
        } else {
          res = await request.post('/api/product', formData)
        }
        
        if (res && res.code === 200) {
          ElMessage.success(isEdit.value ? '编辑成功' : '创建成功')
          router.push('/activity/list')
        } else {
          ElMessage.error(res?.msg || '操作失败')
        }
      } catch (error) {
        console.error('提交失败', error)
        ElMessage.error('操作失败，请重试')
      } finally {
        loading.value = false
      }
    } else {
      ElMessage.warning('请完善表单信息')
      return false
    }
  })
}

// 获取活动详情
const getActivityDetail = async (id) => {
  loading.value = true
  try {
    const res = await request.get(`/api/product/${id}`)
    if (res && res.code === 200) {
      const activityData = res.data
      if (activityData) {
        // 填充表单
        form.name = activityData.name || ''
        form.description = activityData.description || ''
        form.content = activityData.content || ''
        form.status = activityData.status || 0
        
        // 处理分类
        if (activityData.category && activityData.category.length > 0) {
          form.categoryId = activityData.category[0]
        }
        
        // 处理标签
        form.tags = activityData.tags || []
        
        // 处理图片
        if (activityData.images && activityData.images.length > 0) {
          form.coverImage = activityData.images[0]
        }
        
        // 处理时间范围
        if (activityData.startTime && activityData.endTime) {
          form.timeRange = [
            new Date(activityData.startTime),
            new Date(activityData.endTime)
          ]
        }
        
        // 处理地点
        form.location = activityData.location || ''
        
        // 处理最大参与人数 (price字段实际是活动人数上限)
        form.maxParticipants = activityData.price ? activityData.price * 10 : 20
        
        // 处理报名截止时间
        if (activityData.registrationDeadline) {
          form.registrationDeadline = new Date(activityData.registrationDeadline)
        }
      }
    } else {
      ElMessage.error(res?.msg || '获取活动详情失败')
    }
  } catch (error) {
    console.error('获取活动详情失败', error)
    ElMessage.error('获取活动详情失败')
  } finally {
    loading.value = false
  }
}

// 计算属性：处理封面图片显示
const displayCoverImage = computed(() => {
  return formatImageUrl(form.coverImage)
})

onMounted(() => {
  getCategoryOptions()
  if (isEdit.value) {
    getActivityDetail(activityId.value)
  }
})
</script>

<style lang="scss" scoped>
.activity-form-container {
  .mt-20 {
    margin-top: 20px;
  }
  
  .form-card {
    margin-bottom: 20px;
  }
  
  .input-new-tag {
    width: 90px;
    margin-right: 10px;
    vertical-align: bottom;
  }
  
  .avatar-uploader {
    .el-upload {
      border: 1px dashed #d9d9d9;
      border-radius: 6px;
      cursor: pointer;
      position: relative;
      overflow: hidden;
      transition: var(--el-transition-duration-fast);
      
      &:hover {
        border-color: var(--el-color-primary);
      }
    }
  }
  
  .avatar-uploader-icon {
    font-size: 28px;
    color: #8c939d;
    width: 178px;
    height: 178px;
    text-align: center;
    line-height: 178px;
  }
  
  .cover-image {
    width: 300px;
    height: 169px;
    display: block;
  }
  
  .upload-tip {
    font-size: 12px;
    color: #606266;
    margin-top: 10px;
  }
}
</style> 