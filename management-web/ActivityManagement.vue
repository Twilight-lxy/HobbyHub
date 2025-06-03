<template>
  <div class="activity-management">
    <el-card>
      <div class="header">
        <h2>活动管理</h2>
        <div class="filters">
          <el-input
            v-model="searchQuery"
            placeholder="搜索活动名称"
            style="width: 200px"
            @keyup.enter="fetchActivities" />
          <el-button type="primary" @click="fetchActivities">搜索</el-button>
        </div>
      </div>
      
      <el-table :data="activityList" border style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="活动名称" />
        <el-table-column prop="creator" label="创建者" />
        <el-table-column prop="startTime" label="开始时间" />
        <el-table-column prop="endTime" label="结束时间" />
        <el-table-column prop="participants" label="参与人数" />
        <el-table-column label="状态">
          <template #default="{row}">
            <el-tag :type="getStatusTagType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180">
          <template #default="{row}">
            <el-button size="small" @click="viewDetails(row)">详情</el-button>
            <el-button size="small" type="danger" @click="deleteActivity(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="pagination.current"
        :page-sizes="[10, 20, 50, 100]"
        :page-size="pagination.size"
        layout="total, sizes, prev, pager, next, jumper"
        :total="pagination.total">
      </el-pagination>
    </el-card>
    
    <!-- 活动详情对话框 -->
    <activity-detail 
      v-model="detailVisible"
      :activity="currentActivity" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getActivityList, deleteActivity } from '@/api/activity'
import ActivityDetail from './components/ActivityDetail.vue'

const activityList = ref([])
const searchQuery = ref('')
const pagination = ref({
  current: 1,
  size: 10,
  total: 0
})
const detailVisible = ref(false)
const currentActivity = ref({})

const fetchActivities = async () => {
  const params = {
    page: pagination.value.current,
    size: pagination.value.size,
    query: searchQuery.value
  }
  const res = await getActivityList(params)
  activityList.value = res.data.list
  pagination.value.total = res.data.total
}

// 状态显示处理
const getStatusTagType = (status) => {
  const map = {
    0: 'info',    // 未开始
    1: 'success', // 进行中
    2: 'warning', // 即将结束
    3: 'danger'   // 已结束
  }
  return map[status] || 'info'
}

const getStatusText = (status) => {
  const map = {
    0: '未开始',
    1: '进行中',
    2: '即将结束',
    3: '已结束'
  }
  return map[status] || '未知'
}

// 初始化加载数据
onMounted(() => {
  fetchActivities()
})

// 分页处理
const handleSizeChange = (val) => {
  pagination.value.size = val
  fetchActivities()
}

const handleCurrentChange = (val) => {
  pagination.value.current = val
  fetchActivities()
}

// 活动操作
const viewDetails = (activity) => {
  currentActivity.value = activity
  detailVisible.value = true
}

const deleteActivity = async (id) => {
  await ElMessageBox.confirm('确定删除该活动吗？', '提示', {
    type: 'warning'
  })
  await deleteActivity(id)
  ElMessage.success('删除成功')
  fetchActivities()
}
</script>

<style lang="scss" scoped>
.activity-management {
  padding: 20px;
  
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    
    .filters {
      display: flex;
      gap: 10px;
    }
  }
}
</style>