<template>
  <div class="activity-detail-container">
    <el-page-header @back="goBack" :title="'返回活动列表'" :content="'活动详情'" />
    
    <!-- 活动基本信息 -->
    <el-card class="info-card mt-20" v-loading="loading">
      <div class="activity-header">
        <div class="activity-info">
          <h1 class="activity-title">{{ activity.name }}</h1>
          <div class="activity-meta">
            <el-tag :type="getStatusType(activity.status)" class="status-tag">
              {{ getStatusText(activity.status) }}
            </el-tag>
            <el-tag type="info" class="category-tag">{{ getCategoryName(activity.categoryId) }}</el-tag>
            <div v-for="tag in activity.tags" :key="tag" class="activity-tag">
              <el-tag size="small">{{ tag }}</el-tag>
            </div>
          </div>
        </div>
        <div class="activity-actions">
          <el-button type="primary" @click="handleEdit" v-if="isAdmin">
            <el-icon><Edit /></el-icon>
            编辑活动
          </el-button>
          <el-button type="success" @click="handleJoin" :disabled="activity.status !== 0 && activity.status !== 1">
            <el-icon><Plus /></el-icon>
            {{ isJoined ? '已参与' : '立即参与' }}
          </el-button>
          <el-button type="warning" @click="handleShare">
            <el-icon><Share /></el-icon>
            分享
          </el-button>
        </div>
      </div>
      
      <el-divider />
      
      <el-row :gutter="20">
        <el-col :span="16">
          <div class="activity-cover">
            <el-image :src="coverImage" fit="cover" />
          </div>
          
          <div class="activity-section">
            <h2 class="section-title">活动详情</h2>
            <div class="activity-content">{{ activity.description }}</div>
          </div>
          
          <div class="activity-section">
            <h2 class="section-title">参与须知</h2>
            <div class="notice-list">
              <div class="notice-item">
                <el-icon><Calendar /></el-icon>
                <span>活动时间：{{ formatDateTime(activity.createTime) }} 至 {{ formatDateTime(activity.updateTime) }}</span>
              </div>
              <div class="notice-item">
                <el-icon><Location /></el-icon>
                <span>活动地点：{{ activity.location || '线上活动' }}</span>
              </div>
              <div class="notice-item">
                <el-icon><User /></el-icon>
                <span>参与人数：{{ participantCount }}/{{ maxParticipants }}</span>
              </div>
              <div class="notice-item">
                <el-icon><Clock /></el-icon>
                <span>报名截止：{{ activity.registrationDeadline || '无限制' }}</span>
              </div>
              <div class="notice-item">
                <el-icon><InfoFilled /></el-icon>
                <span>注意事项：请准时参加，遵守活动规则，注意安全。</span>
              </div>
            </div>
          </div>
          
          <div class="activity-section">
            <h2 class="section-title">活动评论</h2>
            <div class="comment-list" v-if="comments.length > 0">
              <div class="comment-item" v-for="comment in comments" :key="comment.id">
                <div class="comment-user">
                  <el-avatar :size="40" :src="comment.author?.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'" />
                  <div class="comment-info">
                    <div class="comment-name">{{ comment.author?.nickname || '用户' }}</div>
                    <div class="comment-time">{{ formatDateTime(comment.createTime) }}</div>
                  </div>
                </div>
                <div class="comment-content">{{ comment.content }}</div>
                <div class="comment-actions">
                  <el-button type="text" @click="handleReply(comment)">回复</el-button>
                  <el-button type="text" @click="handleLike(comment)">
                    <el-icon><Star /></el-icon>
                    {{ comment.likes || 0 }}
                  </el-button>
                </div>
              </div>
            </div>
            <div class="comment-empty" v-else>
              暂无评论，快来发表第一条评论吧！
            </div>
            
            <div class="comment-form">
              <el-input
                v-model="commentForm.content"
                type="textarea"
                :rows="3"
                placeholder="请输入您的评论"
              />
              <div class="form-actions">
                <el-button type="primary" @click="submitComment">发表评论</el-button>
              </div>
            </div>
          </div>
        </el-col>
        
        <el-col :span="8">
          <div class="side-card">
            <h3>活动组织者</h3>
            <div class="organizer">
              <el-avatar :size="50" :src="organizer.avatar" />
              <div class="organizer-info">
                <div class="organizer-name">{{ organizer.name }}</div>
                <div class="organizer-desc">{{ organizer.description }}</div>
              </div>
            </div>
          </div>
          
          <div class="side-card">
            <h3>参与人员 ({{ participantCount }}/{{ maxParticipants }})</h3>
            <div class="participant-list">
              <el-avatar
                v-for="participant in participants"
                :key="participant.id"
                :size="40"
                :src="participant.avatar"
                class="participant-avatar"
              />
              <div class="more-participants" v-if="participants.length < participantCount">
                +{{ participantCount - participants.length }}
              </div>
            </div>
          </div>
          
          <div class="side-card">
            <h3>相关活动</h3>
            <div class="related-list">
              <div class="related-item" v-for="item in relatedActivities" :key="item.id" @click="viewActivity(item.id)">
                <el-image :src="getFirstImage(item.images)" fit="cover" class="related-image" />
                <div class="related-info">
                  <div class="related-title">{{ item.name }}</div>
                  <div class="related-time">{{ formatDateTime(item.createTime) }}</div>
                </div>
              </div>
            </div>
          </div>
          
        </el-col>
      </el-row>
    </el-card>
    
    <!-- 参与对话框 -->
    <el-dialog
      title="参与活动"
      v-model="joinDialogVisible"
      width="400px"
    >
      <el-form :model="joinForm" label-width="80px">
        <el-form-item label="姓名">
          <el-input v-model="joinForm.name" placeholder="请输入您的姓名" />
        </el-form-item>
        <el-form-item label="手机号码">
          <el-input v-model="joinForm.phone" placeholder="请输入您的手机号码" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-model="joinForm.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入备注信息"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="joinDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmJoin">确定参与</el-button>
      </template>
    </el-dialog>
    
    <!-- 分享对话框 -->
    <el-dialog
      title="分享活动"
      v-model="shareDialogVisible"
      width="400px"
    >
      <div class="share-container">
        <div class="qrcode-container">
          <el-image src="https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png" fit="cover" class="qrcode-image" />
        </div>
        <div class="share-platforms">
          <div class="share-item">
            <el-button circle>
              <el-icon><ChatDotSquare /></el-icon>
            </el-button>
            <span>微信</span>
          </div>
          <div class="share-item">
            <el-button circle>
              <el-icon><ChatLineSquare /></el-icon>
            </el-button>
            <span>朋友圈</span>
          </div>
          <div class="share-item">
            <el-button circle>
              <el-icon><Message /></el-icon>
            </el-button>
            <span>微博</span>
          </div>
          <div class="share-item">
            <el-button circle>
              <el-icon><Link /></el-icon>
            </el-button>
            <span>复制链接</span>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'
import { formatImageUrl, formatFirstImage } from '@/utils/image'

const route = useRoute()
const router = useRouter()

// 活动ID
const activityId = ref(route.params.id)

// 加载状态
const loading = ref(false)

// 是否是管理员
const isAdmin = ref(false)

// 是否已参与
const isJoined = ref(false)

// 对话框可见性
const joinDialogVisible = ref(false)
const shareDialogVisible = ref(false)

// 参与表单
const joinForm = reactive({
  name: '',
  phone: '',
  remark: ''
})

// 评论表单
const commentForm = reactive({
  content: '',
  postId: null
})

// 分类选项
const categoryOptions = ref([])

// 活动信息
const activity = ref({})
const coverImage = computed(() => {
  if (activity.value.images && activity.value.images.length > 0) {
    return formatImageUrl(activity.value.images[0], 'https://via.placeholder.com/800x400')
  }
  return 'https://via.placeholder.com/800x400'
})

// 活动组织者
const organizer = reactive({
  id: 1,
  name: '活动组织者',
  avatar: '',
  description: '活动组织者描述'
})

// 参与人数和最大人数
const participantCount = computed(() => {
  return activity.value.participantCount || 0
})
const maxParticipants = computed(() => {
  // 在后端，price字段实际上表示活动人数上限
  if (activity.value.price) {
    return activity.value.price * 10
  }
  return 0
})

// 参与者列表
const participants = ref([])

// 评论列表
const comments = ref([])

// 相关活动
const relatedActivities = ref([])

// 获取分类名称
const getCategoryName = (categoryId) => {
  if (!categoryId || !activity.value.category) return '未知分类'
  return activity.value.category.join(', ')
}

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

// 格式化日期时间
const formatDateTime = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString()
}

// 获取图片数组的第一张图片
const getFirstImage = (images) => {
  return formatFirstImage(images, 'https://via.placeholder.com/100x60')
}

// 返回活动列表
const goBack = () => {
  router.push('/activity/list')
}

// 查看活动详情
const viewActivity = (id) => {
  router.push(`/activity/detail/${id}`)
}

// 编辑活动
const handleEdit = () => {
  router.push(`/activity/edit/${activityId.value}`)
}

// 参与活动
const handleJoin = () => {
  if (isJoined.value) {
    ElMessage.info('您已参与该活动')
    return
  }
  joinDialogVisible.value = true
}

// 确认参与
const confirmJoin = async () => {
  try {
    const productOrder = {
      productId: activityId.value,
      count: 1, // 参与人数
      address: joinForm.remark // 使用备注作为地址
    }
    
    const res = await request.post('/api/product_order/save', productOrder)
    if (res && res.code === 200) {
      isJoined.value = true
      ElMessage.success('参与成功')
      joinDialogVisible.value = false
      // 刷新活动详情
      getActivityDetail(activityId.value)
    } else {
      ElMessage.error(res.msg || '参与失败')
    }
  } catch (error) {
    console.error('参与活动失败', error)
    ElMessage.error('参与失败，请重试')
  }
}

// 分享活动
const handleShare = () => {
  shareDialogVisible.value = true
}

// 回复评论
const handleReply = (comment) => {
  commentForm.content = `@${comment.author?.nickname || '用户'} `
}

// 点赞评论
const handleLike = (comment) => {
  // 实际项目中应该调用点赞API
  comment.likes = (comment.likes || 0) + 1
  ElMessage.success('点赞成功')
}

// 提交评论
const submitComment = async () => {
  if (!commentForm.content.trim()) {
    ElMessage.warning('请输入评论内容')
    return
  }
  
  try {
    const postComment = {
      postId: activityId.value,
      content: commentForm.content,
      createTime: new Date()
    }
    
    const res = await request.post('/api/post_comment', postComment)
    if (res && res.code === 200) {
      ElMessage.success('评论成功')
      commentForm.content = ''
      // 重新获取评论列表
      getComments(activityId.value)
    } else {
      ElMessage.error(res.msg || '评论失败')
    }
  } catch (error) {
    console.error('评论失败', error)
    ElMessage.error('评论失败，请重试')
  }
}

// 获取活动详情
const getActivityDetail = async (id) => {
  loading.value = true
  try {
    const res = await request.get(`/api/product/${id}`)
    if (res && res.code === 200) {
      activity.value = res.data
      // 设置评论表单的postId
      commentForm.postId = id
      
      // 获取评论列表
      getComments(id)
      
      // 获取相关活动
      getRelatedActivities()
      
      // 检查是否已参与
      checkIfJoined(id)
    } else {
      ElMessage.error(res.msg || '获取活动详情失败')
    }
  } catch (error) {
    console.error('获取活动详情失败', error)
    ElMessage.error('获取活动详情失败')
  } finally {
    loading.value = false
  }
}

// 获取评论列表
const getComments = async (id) => {
  try {
    const res = await request.get(`/api/post_comment/list/${id}`)
    if (res && res.code === 200) {
      comments.value = res.data || []
    }
  } catch (error) {
    console.error('获取评论列表失败', error)
  }
}

// 获取相关活动
const getRelatedActivities = async () => {
  try {
    const res = await request.get('/api/product/list')
    if (res && res.code === 200) {
      // 过滤掉当前活动，并只取3个相关活动
      relatedActivities.value = res.data
        .filter(item => item.id !== Number(activityId.value))
        .slice(0, 3)
    }
  } catch (error) {
    console.error('获取相关活动失败', error)
  }
}

// 检查是否已参与
const checkIfJoined = async (id) => {
  try {
    const res = await request.get('/api/product_order/list')
    if (res && res.code === 200) {
      const orders = res.data || []
      isJoined.value = orders.some(order => order.productId === Number(id))
    }
  } catch (error) {
    console.error('检查参与状态失败', error)
  }
}

// 获取分类列表
const getCategories = async () => {
  try {
    const res = await request.get('/api/tabs/list')
    if (res && res.code === 200) {
      categoryOptions.value = res.data || []
    }
  } catch (error) {
    console.error('获取分类列表失败', error)
  }
}

onMounted(() => {
  getActivityDetail(activityId.value)
  getCategories()
})
</script>

<style lang="scss" scoped>
.activity-detail-container {
  .mt-20 {
    margin-top: 20px;
  }
  
  .info-card {
    margin-bottom: 20px;
  }
  
  .activity-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 20px;
    
    .activity-info {
      .activity-title {
        font-size: 24px;
        margin-bottom: 15px;
      }
      
      .activity-meta {
        display: flex;
        align-items: center;
        flex-wrap: wrap;
        
        .status-tag,
        .category-tag,
        .activity-tag {
          margin-right: 10px;
          margin-bottom: 10px;
        }
      }
    }
    
    .activity-actions {
      display: flex;
      gap: 10px;
    }
  }
  
  .activity-cover {
    margin-bottom: 30px;
    
    .el-image {
      width: 100%;
      border-radius: 8px;
    }
  }
  
  .activity-section {
    margin-bottom: 30px;
    
    .section-title {
      font-size: 20px;
      margin-bottom: 20px;
      padding-bottom: 10px;
      border-bottom: 1px solid #ebeef5;
    }
    
    .activity-content {
      white-space: pre-line;
      line-height: 1.8;
    }
  }
  
  .notice-list {
    .notice-item {
      display: flex;
      align-items: center;
      margin-bottom: 15px;
      
      .el-icon {
        margin-right: 10px;
        color: #409EFF;
      }
    }
  }
  
  .comment-list {
    .comment-item {
      padding: 15px 0;
      border-bottom: 1px solid #ebeef5;
      
      &:last-child {
        border-bottom: none;
      }
      
      .comment-user {
        display: flex;
        align-items: center;
        margin-bottom: 10px;
        
        .comment-info {
          margin-left: 10px;
          
          .comment-name {
            font-weight: bold;
          }
          
          .comment-time {
            font-size: 12px;
            color: #909399;
          }
        }
      }
      
      .comment-content {
        margin-bottom: 10px;
        line-height: 1.6;
      }
      
      .comment-actions {
        display: flex;
        justify-content: flex-end;
      }
    }
  }
  
  .comment-empty {
    text-align: center;
    color: #909399;
    padding: 30px 0;
  }
  
  .comment-form {
    margin-top: 20px;
    
    .form-actions {
      margin-top: 10px;
      display: flex;
      justify-content: flex-end;
    }
  }
  
  .side-card {
    background-color: #f5f7fa;
    border-radius: 8px;
    padding: 20px;
    margin-bottom: 20px;
    
    h3 {
      margin-top: 0;
      margin-bottom: 15px;
      font-size: 18px;
      border-bottom: 1px solid #ebeef5;
      padding-bottom: 10px;
    }
    
    .organizer {
      display: flex;
      align-items: center;
      
      .organizer-info {
        margin-left: 15px;
        
        .organizer-name {
          font-weight: bold;
          margin-bottom: 5px;
        }
        
        .organizer-desc {
          font-size: 12px;
          color: #606266;
        }
      }
    }
    
    .participant-list {
      display: flex;
      flex-wrap: wrap;
      gap: 10px;
      
      .participant-avatar {
        cursor: pointer;
      }
      
      .more-participants {
        width: 40px;
        height: 40px;
        border-radius: 50%;
        background-color: #f0f0f0;
        display: flex;
        justify-content: center;
        align-items: center;
        font-size: 12px;
        color: #606266;
      }
    }
    
    .related-list {
      .related-item {
        display: flex;
        align-items: center;
        padding: 10px 0;
        border-bottom: 1px solid #ebeef5;
        cursor: pointer;
        
        &:last-child {
          border-bottom: none;
        }
        
        &:hover {
          .related-title {
            color: #409EFF;
          }
        }
        
        .related-image {
          width: 60px;
          height: 40px;
          border-radius: 4px;
          margin-right: 10px;
        }
        
        .related-info {
          flex: 1;
          
          .related-title {
            font-weight: bold;
            margin-bottom: 5px;
          }
          
          .related-time {
            font-size: 12px;
            color: #909399;
          }
        }
      }
    }
    
    .qrcode-container {
      display: flex;
      flex-direction: column;
      align-items: center;
      
      .qrcode-image {
        width: 200px;
        height: 200px;
      }
      
      .qrcode-tip {
        margin-top: 10px;
        font-size: 14px;
        color: #606266;
      }
    }
  }
  
  .share-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    
    .qrcode-container {
      margin-bottom: 20px;
    }
    
    .share-platforms {
      display: flex;
      justify-content: center;
      gap: 20px;
      
      .share-item {
        display: flex;
        flex-direction: column;
        align-items: center;
        
        span {
          margin-top: 8px;
          font-size: 12px;
        }
      }
    }
  }
}
</style> 