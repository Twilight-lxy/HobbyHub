<template>
  <div class="team-detail-container">
    <el-page-header @back="goBack" :title="'返回小队列表'" :content="'小队详情'" />
    
    <el-card class="info-card mt-20">
      <div class="team-header">
        <div class="team-info">
          <h1 class="team-title">{{ team.name }}</h1>
          <div class="team-meta">
            <el-tag :type="getStatusType(team.status)" class="status-tag">
              {{ getStatusText(team.status) }}
            </el-tag>
            <el-tag type="info" class="category-tag">{{ getCategoryName(team.categoryId) }}</el-tag>
          </div>
        </div>
        <div class="team-actions">
          <el-button type="primary" @click="handleEdit" v-if="isAdmin || isLeader">
            <el-icon><Edit /></el-icon>
            编辑小队
          </el-button>
          <el-button type="success" @click="handleJoin" :disabled="team.status !== 0 || isJoined">
            <el-icon><Plus /></el-icon>
            {{ isJoined ? '已加入' : '加入小队' }}
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
          <div class="team-cover">
            <el-image :src="coverImage" fit="cover" />
          </div>
          
          <div class="team-section">
            <h2 class="section-title">小队简介</h2>
            <div class="team-content">{{ team.description }}</div>
          </div>
          
          <div class="team-section">
            <h2 class="section-title">小队成员 ({{ team.memberCount }}/{{ team.maxMembers }})</h2>
            <el-table :data="members" style="width: 100%">
              <el-table-column label="头像" width="80">
                <template #default="{ row }">
                  <el-avatar :size="40" :src="row.avatar" />
                </template>
              </el-table-column>
              <el-table-column prop="name" label="用户名" />
              <el-table-column prop="joinTime" label="加入时间" />
              <el-table-column label="角色" width="100">
                <template #default="{ row }">
                  <el-tag :type="row.isLeader ? 'danger' : ''">
                    {{ row.isLeader ? '队长' : '成员' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="150" v-if="isAdmin || isLeader">
                <template #default="{ row }">
                  <el-button 
                    type="danger" 
                    link 
                    size="small" 
                    :disabled="row.isLeader"
                    @click="handleRemoveMember(row)"
                  >
                    移除
                  </el-button>
                  <el-button 
                    type="primary" 
                    link 
                    size="small" 
                    v-if="isLeader && !row.isLeader"
                    @click="handleSetLeader(row)"
                  >
                    设为队长
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
          
          <div class="team-section">
            <h2 class="section-title">小队动态</h2>
            <div class="activity-list" v-if="activities.length > 0">
              <div class="activity-item" v-for="activity in activities" :key="activity.id">
                <div class="activity-time">{{ activity.time }}</div>
                <div class="activity-content">
                  <div class="activity-user">
                    <el-avatar :size="32" :src="activity.user.avatar" />
                    <span class="user-name">{{ activity.user.name }}</span>
                  </div>
                  <div class="activity-message">{{ activity.message }}</div>
                </div>
              </div>
            </div>
            <div class="empty-block" v-else>
              <el-empty description="暂无小队动态" />
            </div>
          </div>
        </el-col>
        
        <el-col :span="8">
          <div class="side-card">
            <h3>队长信息</h3>
            <div class="leader-info">
              <el-avatar :size="60" :src="team.leader?.avatar" />
              <div class="leader-detail">
                <div class="leader-name">{{ team.leader?.name }}</div>
                <div class="leader-contact">
                  <el-button type="primary" link size="small">
                    <el-icon><Message /></el-icon>
                    发送消息
                  </el-button>
                </div>
              </div>
            </div>
          </div>
          
          <div class="side-card">
            <h3>相关活动</h3>
            <div class="related-list">
              <div class="related-item" v-for="item in relatedActivities" :key="item.id" @click="viewActivity(item.id)">
                <el-image :src="formatActivityImage(item.coverImage)" fit="cover" class="related-image" />
                <div class="related-info">
                  <div class="related-title">{{ item.name }}</div>
                  <div class="related-time">{{ item.startTime }}</div>
                </div>
              </div>
            </div>
            <div class="empty-block" v-if="relatedActivities.length === 0">
              <el-empty description="暂无相关活动" />
            </div>
          </div>
          
          <div class="side-card">
            <h3>小队二维码</h3>
            <div class="qrcode-container">
              <el-image src="https://via.placeholder.com/200" fit="cover" class="qrcode-image" />
              <div class="qrcode-tip">扫码加入小队</div>
            </div>
          </div>
        </el-col>
      </el-row>
    </el-card>
    
    <!-- 加入小队对话框 -->
    <el-dialog
      title="加入小队"
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
        <el-form-item label="加入理由">
          <el-input
            v-model="joinForm.reason"
            type="textarea"
            :rows="3"
            placeholder="请简单介绍一下自己，以及为什么想加入小队"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="joinDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmJoin">确定加入</el-button>
      </template>
    </el-dialog>
    
    <!-- 分享对话框 -->
    <el-dialog
      title="分享小队"
      v-model="shareDialogVisible"
      width="400px"
    >
      <div class="share-container">
        <div class="qrcode-container">
          <el-image src="https://via.placeholder.com/200" fit="cover" class="qrcode-image" />
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
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatImageUrl, formatFirstImage } from '@/utils/image'
import request from '@/utils/request'

const route = useRoute()
const router = useRouter()

// 小队ID
const teamId = ref(route.params.id)

// 当前用户ID
const currentUserId = ref(1)

// 是否是管理员
const isAdmin = ref(true)

// 是否是队长
const isLeader = computed(() => {
  return team.value.leader && team.value.leader.id === currentUserId.value
})

// 是否已加入
const isJoined = ref(false)

// 对话框可见性
const joinDialogVisible = ref(false)
const shareDialogVisible = ref(false)

// 加入表单
const joinForm = reactive({
  name: '',
  phone: '',
  reason: ''
})

// 分类选项
const categoryOptions = [
  { id: 1, name: '户外活动' },
  { id: 2, name: '文化活动' },
  { id: 3, name: '体育活动' },
  { id: 4, name: '技术交流' },
  { id: 5, name: '兴趣培训' }
]

// 小队信息
const team = reactive({
  id: teamId.value,
  name: '摄影爱好者小队',
  categoryId: 2,
  coverImage: 'https://via.placeholder.com/800x400',
  description: '这是一个摄影爱好者的小队，我们会定期组织外出拍摄活动，分享摄影技巧和心得。欢迎所有对摄影感兴趣的朋友加入我们！\n\n小队成立于2023年3月，目前已经组织了多次外出拍摄活动，涵盖风景、人像、街拍等多种题材。我们有经验丰富的摄影师提供指导，也欢迎摄影新手加入学习。',
  memberCount: 15,
  maxMembers: 30,
  status: 0,
  leader: {
    id: 2,
    name: '摄影达人',
    avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
  }
})

// 小队成员
const members = ref([
  {
    id: 2,
    name: '摄影达人',
    avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
    joinTime: '2023-03-01 10:00:00',
    isLeader: true
  },
  {
    id: 3,
    name: '风景控',
    avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
    joinTime: '2023-03-02 14:30:00',
    isLeader: false
  },
  {
    id: 4,
    name: '人像摄影师',
    avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
    joinTime: '2023-03-05 09:15:00',
    isLeader: false
  },
  {
    id: 5,
    name: '街拍爱好者',
    avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
    joinTime: '2023-03-10 16:45:00',
    isLeader: false
  }
])

// 小队动态
const activities = ref([
  {
    id: 1,
    time: '2023-06-15 14:30:00',
    user: {
      id: 2,
      name: '摄影达人',
      avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
    },
    message: '创建了小队'
  },
  {
    id: 2,
    time: '2023-06-16 10:20:00',
    user: {
      id: 3,
      name: '风景控',
      avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
    },
    message: '加入了小队'
  },
  {
    id: 3,
    time: '2023-06-18 09:45:00',
    user: {
      id: 2,
      name: '摄影达人',
      avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
    },
    message: '发布了活动：周末西湖外拍'
  }
])

// 相关活动
const relatedActivities = ref([
  {
    id: 101,
    name: '周末西湖外拍',
    coverImage: 'https://via.placeholder.com/100x60',
    startTime: '2023-06-24 09:00:00'
  },
  {
    id: 102,
    name: '人像摄影技巧分享',
    coverImage: 'https://via.placeholder.com/100x60',
    startTime: '2023-07-01 14:00:00'
  },
  {
    id: 103,
    name: '后期修图培训',
    coverImage: 'https://via.placeholder.com/100x60',
    startTime: '2023-07-08 10:00:00'
  }
])

// 计算属性：处理封面图片
const coverImage = computed(() => {
  return formatImageUrl(team.value.coverImage, 'https://via.placeholder.com/800x400')
})

// 获取分类名称
const getCategoryName = (categoryId) => {
  const category = categoryOptions.find(item => item.id === categoryId)
  return category ? category.name : '未知分类'
}

// 获取状态类型
const getStatusType = (status) => {
  const map = {
    0: 'success',
    1: 'warning',
    2: 'info'
  }
  return map[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const map = {
    0: '招募中',
    1: '已满员',
    2: '已关闭'
  }
  return map[status] || '未知'
}

// 返回小队列表
const goBack = () => {
  router.push('/team/list')
}

// 查看活动详情
const viewActivity = (id) => {
  router.push(`/activity/detail/${id}`)
}

// 编辑小队
const handleEdit = () => {
  router.push(`/team/edit/${teamId.value}`)
}

// 加入小队
const handleJoin = () => {
  if (isJoined.value) {
    ElMessage.info('您已加入该小队')
    return
  }
  joinDialogVisible.value = true
}

// 确认加入
const confirmJoin = () => {
  // 模拟API请求
  setTimeout(() => {
    isJoined.value = true
    team.memberCount++
    members.value.push({
      id: currentUserId.value,
      name: '当前用户',
      avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
      joinTime: new Date().toLocaleString(),
      isLeader: false
    })
    activities.value.unshift({
      id: activities.value.length + 1,
      time: new Date().toLocaleString(),
      user: {
        id: currentUserId.value,
        name: '当前用户',
        avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
      },
      message: '加入了小队'
    })
    ElMessage.success('加入成功')
    joinDialogVisible.value = false
  }, 500)
}

// 分享小队
const handleShare = () => {
  shareDialogVisible.value = true
}

// 移除成员
const handleRemoveMember = (member) => {
  ElMessageBox.confirm(`确定要移除成员 ${member.name} 吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    // 模拟API请求
    setTimeout(() => {
      members.value = members.value.filter(item => item.id !== member.id)
      team.memberCount--
      activities.value.unshift({
        id: activities.value.length + 1,
        time: new Date().toLocaleString(),
        user: {
          id: currentUserId.value,
          name: '当前用户',
          avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
        },
        message: `将成员 ${member.name} 移出了小队`
      })
      ElMessage.success('移除成功')
    }, 500)
  }).catch(() => {})
}

// 设置队长
const handleSetLeader = (member) => {
  ElMessageBox.confirm(`确定要将 ${member.name} 设为队长吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    // 模拟API请求
    setTimeout(() => {
      // 更新原队长
      const oldLeader = members.value.find(item => item.isLeader)
      if (oldLeader) {
        oldLeader.isLeader = false
      }
      
      // 设置新队长
      const newLeader = members.value.find(item => item.id === member.id)
      if (newLeader) {
        newLeader.isLeader = true
        team.leader = {
          id: newLeader.id,
          name: newLeader.name,
          avatar: newLeader.avatar
        }
      }
      
      activities.value.unshift({
        id: activities.value.length + 1,
        time: new Date().toLocaleString(),
        user: {
          id: currentUserId.value,
          name: '当前用户',
          avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
        },
        message: `将 ${member.name} 设为了队长`
      })
      ElMessage.success('设置成功')
    }, 500)
  }).catch(() => {})
}

// 获取小队详情
const getTeamDetail = async (id) => {
  try {
    // 实际调用API获取小队详情
    const res = await request.get(`/api/product_cart/${id}`)
    if (res && res.code === 200) {
      const data = res.data
      
      // 将后端数据映射到前端模型
      if (data) {
        // 基本信息
        team.id = data.id
        team.name = data.product?.name || '未命名小队'
        team.categoryId = data.product?.category ? parseInt(data.product.category.split(',')[0]) : 1
        team.coverImage = data.product?.images ? data.product.images.split(',')[0] : 'https://via.placeholder.com/800x400'
        team.description = data.product?.description || '暂无描述'
        team.memberCount = data.count || 0
        team.maxMembers = data.product?.price ? data.product.price * 10 : 30
        team.status = data.status === 'selected' ? 1 : (data.status === 'unselected' ? 0 : 2)
        
        // 队长信息
        if (data.user) {
          team.leader = {
            id: data.user.id,
            name: data.user.nickname || data.user.username || '未知用户',
            avatar: data.user.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
          }
        }
        
        // 检查当前用户是否是队长
        isJoined.value = data.user?.id === currentUserId.value
      }
    } else {
      ElMessage.error(res?.msg || '获取小队详情失败')
    }
  } catch (error) {
    console.error('获取小队详情失败', error)
    ElMessage.error('获取小队详情失败')
  }
}

// 格式化活动图片
const formatActivityImage = (url) => {
  return formatImageUrl(url, 'https://via.placeholder.com/100x60')
}

onMounted(() => {
  getTeamDetail(teamId.value)
  
  // 检查当前用户是否已加入
  const isMember = members.value.some(item => item.id === currentUserId.value)
  isJoined.value = isMember
})
</script>

<style lang="scss" scoped>
.team-detail-container {
  .mt-20 {
    margin-top: 20px;
  }
  
  .info-card {
    margin-bottom: 20px;
  }
  
  .team-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 20px;
    
    .team-info {
      .team-title {
        font-size: 24px;
        margin-bottom: 15px;
      }
      
      .team-meta {
        display: flex;
        align-items: center;
        
        .status-tag,
        .category-tag {
          margin-right: 10px;
        }
      }
    }
    
    .team-actions {
      display: flex;
      gap: 10px;
    }
  }
  
  .team-cover {
    margin-bottom: 30px;
    
    .el-image {
      width: 100%;
      border-radius: 8px;
    }
  }
  
  .team-section {
    margin-bottom: 30px;
    
    .section-title {
      font-size: 20px;
      margin-bottom: 20px;
      padding-bottom: 10px;
      border-bottom: 1px solid #ebeef5;
    }
    
    .team-content {
      white-space: pre-line;
      line-height: 1.8;
    }
  }
  
  .activity-list {
    .activity-item {
      display: flex;
      margin-bottom: 20px;
      
      .activity-time {
        width: 150px;
        color: #909399;
        font-size: 14px;
      }
      
      .activity-content {
        flex: 1;
        
        .activity-user {
          display: flex;
          align-items: center;
          margin-bottom: 10px;
          
          .user-name {
            margin-left: 10px;
            font-weight: bold;
          }
        }
        
        .activity-message {
          padding: 10px 15px;
          background-color: #f5f7fa;
          border-radius: 4px;
        }
      }
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
    
    .leader-info {
      display: flex;
      align-items: center;
      
      .leader-detail {
        margin-left: 15px;
        
        .leader-name {
          font-weight: bold;
          margin-bottom: 5px;
        }
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
  
  .empty-block {
    padding: 20px 0;
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