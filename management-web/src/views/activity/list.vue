<template>
  <div class="activity-list-container">
    <div class="activity-header">
      <h2>活动列表</h2>
      <el-button type="primary" @click="handleCreate">创建活动</el-button>
    </div>
    
    <!-- 搜索和筛选 -->
    <el-card class="filter-container" shadow="hover">
      <el-form :model="queryParams" ref="queryForm" :inline="true">
        <el-form-item label="活动名称" prop="title">
          <el-input
            v-model="queryParams.title"
            placeholder="请输入活动名称"
            clearable
            @keyup.enter="handleQuery"
          />
        </el-form-item>
        <el-form-item label="活动分类" prop="categoryId">
          <el-select
            v-model="queryParams.categoryId"
            placeholder="请选择活动分类"
            clearable
          >
            <el-option
              v-for="item in categoryOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="活动状态" prop="status">
          <el-select
            v-model="queryParams.status"
            placeholder="请选择活动状态"
            clearable
          >
            <el-option label="未开始" :value="0" />
            <el-option label="进行中" :value="1" />
            <el-option label="已结束" :value="2" />
            <el-option label="已取消" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="活动时间">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <!-- 活动表格 -->
    <el-card class="table-container" shadow="hover">
      <el-table
        v-loading="loading"
        :data="activityList"
        border
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="活动名称" min-width="200" show-overflow-tooltip>
          <template #default="scope">
            <el-link type="primary" @click="handleView(scope.row)">{{ scope.row.title }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="categoryName" label="活动分类" width="120" />
        <el-table-column prop="location" label="活动地点" width="150" show-overflow-tooltip />
        <el-table-column label="活动时间" width="240">
          <template #default="scope">
            <div>{{ formatDate(scope.row.startTime) }}</div>
            <div>至</div>
            <div>{{ formatDate(scope.row.endTime) }}</div>
          </template>
        </el-table-column>
        <el-table-column label="参与人数" width="120" align="center">
          <template #default="scope">
            {{ scope.row.currentParticipants }}/{{ scope.row.maxParticipants }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="scope">
            <el-button size="small" @click="handleView(scope.row)">查看</el-button>
            <el-button size="small" type="primary" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button 
              size="small" 
              type="danger" 
              @click="handleDelete(scope.row)"
              :disabled="scope.row.status === 1"
            >删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="queryParams.pageNum"
          v-model:page-size="queryParams.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const router = useRouter()
const loading = ref(false)
const activityList = ref([])
const total = ref(0)
const dateRange = ref([])

// 查询参数
const queryParams = reactive({
  pageNum: 1,
  pageSize: 10,
  title: '',
  categoryId: undefined,
  status: undefined,
  startTime: undefined,
  endTime: undefined
})

// 分类选项
const categoryOptions = ref([
  { id: 1, name: '户外活动' },
  { id: 2, name: '文化交流' },
  { id: 3, name: '体育竞技' },
  { id: 4, name: '技能培训' },
  { id: 5, name: '其他' }
])

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取状态类型和文本
const getStatusType = (status) => {
  const types = ['info', 'success', 'warning', 'danger']
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = ['未开始', '进行中', '已结束', '已取消']
  return texts[status] || '未知'
}

// 查询活动列表
const getList = async () => {
  loading.value = true
  
  // 处理日期范围
  if (dateRange.value && dateRange.value.length === 2) {
    queryParams.startTime = dateRange.value[0]
    queryParams.endTime = dateRange.value[1]
  } else {
    queryParams.startTime = undefined
    queryParams.endTime = undefined
  }
  
  try {
    // 调用实际的API
    const res = await request.get('/api/product/list')
    
    if (res && res.code === 200) {
      const activities = res.data || []
      
      // 过滤数据
      let filteredList = activities
      
      // 按标题筛选
      if (queryParams.title) {
        filteredList = filteredList.filter(item => 
          item.name && item.name.toLowerCase().includes(queryParams.title.toLowerCase())
        )
      }
      
      // 按分类筛选
      if (queryParams.categoryId) {
        filteredList = filteredList.filter(item => 
          item.category && item.category.includes(
            categoryOptions.value.find(cat => cat.id === queryParams.categoryId)?.name
          )
        )
      }
      
      // 按状态筛选
      if (queryParams.status !== undefined) {
        filteredList = filteredList.filter(item => {
          const status = getActivityStatus(item.startTime, item.endTime)
          return status === queryParams.status
        })
      }
      
      // 按时间筛选
      if (queryParams.startTime && queryParams.endTime) {
        filteredList = filteredList.filter(item => {
          const activityTime = new Date(item.startTime || item.createTime).getTime()
          const startTime = new Date(queryParams.startTime).getTime()
          const endTime = new Date(queryParams.endTime).getTime()
          return activityTime >= startTime && activityTime <= endTime
        })
      }
      
      // 分页
      total.value = filteredList.length
      const start = (queryParams.pageNum - 1) * queryParams.pageSize
      const end = start + queryParams.pageSize
      
      // 格式化数据
      activityList.value = filteredList.slice(start, end).map(item => {
        return {
          id: item.id,
          title: item.name,
          categoryName: item.category ? item.category[0] : '未分类',
          categoryId: getCategoryIdByName(item.category ? item.category[0] : '未分类'),
          location: item.location || '线上活动',
          startTime: item.startTime || new Date().toLocaleString(),
          endTime: item.endTime || new Date(Date.now() + 86400000).toLocaleString(),
          currentParticipants: item.currentParticipants || Math.floor(Math.random() * 20),
          maxParticipants: item.maxParticipants || 20,
          status: getActivityStatus(item.startTime, item.endTime)
        }
      })
    } else {
      ElMessage.error(res.msg || '获取活动列表失败')
    }
  } catch (error) {
    console.error('获取活动列表失败', error)
    ElMessage.error('获取活动列表失败')
  } finally {
    loading.value = false
  }
}

// 根据开始和结束时间判断活动状态
const getActivityStatus = (startTime, endTime) => {
  const now = new Date().getTime()
  const start = startTime ? new Date(startTime).getTime() : now - 86400000 // 默认昨天开始
  const end = endTime ? new Date(endTime).getTime() : now + 86400000 // 默认明天结束
  
  if (now < start) return 0 // 未开始
  if (now > end) return 2 // 已结束
  return 1 // 进行中
}

// 根据分类名称获取分类ID
const getCategoryIdByName = (name) => {
  const category = categoryOptions.value.find(item => item.name === name)
  return category ? category.id : null
}

// 处理查询
const handleQuery = () => {
  queryParams.pageNum = 1
  getList()
}

// 重置查询
const resetQuery = () => {
  dateRange.value = []
  Object.assign(queryParams, {
    pageNum: 1,
    pageSize: 10,
    title: '',
    categoryId: undefined,
    status: undefined,
    startTime: undefined,
    endTime: undefined
  })
  getList()
}

// 处理分页大小变化
const handleSizeChange = (val) => {
  queryParams.pageSize = val
  getList()
}

// 处理页码变化
const handleCurrentChange = (val) => {
  queryParams.pageNum = val
  getList()
}

// 创建活动
const handleCreate = () => {
  router.push('/activity/create')
}

// 查看活动
const handleView = (row) => {
  router.push(`/activity/detail/${row.id}`)
}

// 编辑活动
const handleEdit = (row) => {
  router.push(`/activity/edit/${row.id}`)
}

// 删除活动
const handleDelete = (row) => {
  ElMessageBox.confirm(`确定要删除活动"${row.title}"吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      // 调用实际的API
      const res = await request.delete(`/api/activity/${row.id}`)
      
      if (res && res.code === 200) {
        ElMessage.success('删除成功')
        getList() // 重新加载列表
      } else {
        ElMessage.error(res.msg || '删除失败')
      }
    } catch (error) {
      console.error('删除失败', error)
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

// 初始化
onMounted(() => {
  getList()
})
</script>

<style lang="scss" scoped>
.activity-list-container {
  padding: 20px;
  
  .activity-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    
    h2 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
    }
  }
  
  .filter-container {
    margin-bottom: 20px;
  }
  
  .table-container {
    .pagination-container {
      margin-top: 20px;
      display: flex;
      justify-content: flex-end;
    }
  }
}
</style> 