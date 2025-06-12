<template>
  <div class="order-detail-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>参与记录详情</span>
          <el-button type="primary" @click="goBack">返回列表</el-button>
        </div>
      </template>
      
      <el-descriptions :column="2" border>
        <el-descriptions-item label="记录ID">{{ orderDetail.id }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ orderDetail.user?.username || '未知' }}</el-descriptions-item>
        <el-descriptions-item label="活动名称">{{ orderDetail.product?.name || '未知' }}</el-descriptions-item>
        <el-descriptions-item label="参与人数">{{ orderDetail.count || 0 }}</el-descriptions-item>
        <el-descriptions-item label="总人数">{{ orderDetail.totalPrice || 0 }}</el-descriptions-item>
        <el-descriptions-item label="参与时间">{{ orderDetail.createTime || '未知' }}</el-descriptions-item>
      </el-descriptions>
      
      <div class="section-title">活动详情</div>
      <el-descriptions :column="2" border v-if="orderDetail.product">
        <el-descriptions-item label="活动ID">{{ orderDetail.product.id }}</el-descriptions-item>
        <el-descriptions-item label="活动分类">{{ orderDetail.product.category }}</el-descriptions-item>
        <el-descriptions-item label="人数上限">{{ orderDetail.product.price * 10 }}</el-descriptions-item>
        <el-descriptions-item label="剩余名额">{{ orderDetail.product.stock }}</el-descriptions-item>
        <el-descriptions-item label="活动描述" :span="2">{{ orderDetail.product.description }}</el-descriptions-item>
        <el-descriptions-item label="活动图片" :span="2">
          <div class="image-preview">
            <el-image
              v-for="(img, index) in orderDetail.product.images?.split(',')"
              :key="index"
              :src="img"
              :preview-src-list="orderDetail.product.images?.split(',')"
              fit="cover"
              class="preview-image"
            />
          </div>
        </el-descriptions-item>
      </el-descriptions>
      <div v-else class="empty-tip">暂无活动详情</div>
      
      <div class="section-title">用户信息</div>
      <el-descriptions :column="2" border v-if="orderDetail.user">
        <el-descriptions-item label="用户ID">{{ orderDetail.user.id }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ orderDetail.user.username }}</el-descriptions-item>
        <el-descriptions-item label="昵称">{{ orderDetail.user.nickname || '未设置' }}</el-descriptions-item>
        <el-descriptions-item label="手机号">{{ orderDetail.user.phone || '未设置' }}</el-descriptions-item>
        <el-descriptions-item label="邮箱">{{ orderDetail.user.email || '未设置' }}</el-descriptions-item>
        <el-descriptions-item label="注册时间">{{ orderDetail.user.createTime }}</el-descriptions-item>
      </el-descriptions>
      <div v-else class="empty-tip">暂无用户详情</div>
      
      <div class="button-group">
        <el-button type="primary" @click="handleEdit">编辑记录</el-button>
        <el-button type="danger" @click="handleDelete">删除记录</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getOrderDetail, deleteOrder } from '@/api/order'

const route = useRoute()
const router = useRouter()
const orderId = route.params.id

// 参与记录详情
const orderDetail = ref({})

// 获取参与记录详情
const fetchOrderDetail = async () => {
  try {
    const res = await getOrderDetail(orderId)
    orderDetail.value = res.data
  } catch (error) {
    console.error('获取参与记录详情失败', error)
  }
}

// 返回列表
const goBack = () => {
  router.push('/order/list')
}

// 处理编辑
const handleEdit = () => {
  ElMessage.info('编辑功能开发中')
}

// 处理删除
const handleDelete = () => {
  ElMessageBox.confirm('确定要删除该参与记录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteOrder(orderId)
      ElMessage.success('删除成功')
      router.push('/order/list')
    } catch (error) {
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

onMounted(() => {
  fetchOrderDetail()
})
</script>

<style lang="scss" scoped>
.order-detail-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .section-title {
    font-size: 16px;
    font-weight: bold;
    margin: 20px 0 10px;
    padding-bottom: 10px;
    border-bottom: 1px solid #ebeef5;
  }
  
  .empty-tip {
    color: #909399;
    text-align: center;
    padding: 20px 0;
  }
  
  .image-preview {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
    
    .preview-image {
      width: 100px;
      height: 100px;
      border-radius: 4px;
    }
  }
  
  .button-group {
    margin-top: 20px;
    display: flex;
    justify-content: center;
    gap: 20px;
  }
}
</style> 