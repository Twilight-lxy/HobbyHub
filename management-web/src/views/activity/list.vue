<template>
  <div class="activity-management-container">
    <!-- 搜索区域 -->
    <el-card class="search-card">
      <el-form :inline="true" :model="queryParams" class="search-form">
        <el-form-item label="活动名称">
          <el-input v-model="queryParams.activityName" placeholder="请输入活动名称" clearable />
        </el-form-item>
        <el-form-item label="活动地点">
          <el-input v-model="queryParams.addr" placeholder="请输入活动地点" clearable />
        </el-form-item>
        <el-form-item label="创建时间">
          <el-date-picker
            v-model="createDateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="开始时间">
          <el-date-picker
            v-model="startDateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="活动状态">
          <el-select v-model="queryParams.state" placeholder="请选择活动状态">
            <el-option label="全部" value="" />
            <el-option label="待审核" value="0" />
            <el-option label="已审核" value="1" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="resetQuery">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据统计卡片 -->
    <el-row :gutter="20" class="stats-card-container">
      <el-col :span="12">
        <el-card class="stats-card" shadow="hover">
          <div class="stats-content">
            <div class="stats-title">总活动数</div>
            <div class="stats-value">{{ totalActivities }}</div>
            <div class="stats-trend">
              <i class="el-icon-caret-top" :class="trendClass"></i>
              <span>{{ trendText }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="stats-card" shadow="hover">
          <div class="stats-content">
            <div class="stats-title">待审核活动</div>
            <div class="stats-value">{{ pendingActivities }}</div>
            <div class="stats-trend">
              <i class="el-icon-caret-top" :class="pendingTrendClass"></i>
              <span>{{ pendingTrendText }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 表格区域 -->
    <el-card class="table-card">
      <template #header>
        <div class="card-header">
          <span>活动管理列表</span>
          <div class="right-actions">
            <el-button type="primary" @click="handleAdd">
              <el-icon><Plus /></el-icon>
              新增活动
            </el-button>
            <el-button type="success" v-if="selectedIds.length && selectedIds.length === 1" @click="handleApprove">
              <el-icon><Check /></el-icon>
              审核通过
            </el-button>
            <el-button type="danger" :disabled="!selectedIds.length" @click="handleBatchDelete">
              <el-icon><Delete /></el-icon>
              批量删除
            </el-button>
            <el-button type="success" @click="handleExport">
              <el-icon><Download /></el-icon>
              导出数据
            </el-button>
          </div>
        </div>
      </template>
      
      <!-- 表格 -->
      <el-table
        v-loading="loading"
        :data="activityList"
        @selection-change="handleSelectionChange"
        stripe
        fit
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="活动ID" prop="id" width="80" />
        <el-table-column label="活动名称" prop="name" width="180" />
        <el-table-column label="活动地点" prop="addr" width="150" />
        <el-table-column label="创建者ID" prop="userId" width="120" />
        <el-table-column label="创建时间" prop="createTime" width="180">
          <template #default="{ row }">
            {{ formatDate(row.createTime) }}
          </template>
        </el-table-column>
        <el-table-column label="开始时间" prop="startTime" width="180">
          <template #default="{ row }">
            {{ formatDate(row.startTime) }}
          </template>
        </el-table-column>
        <el-table-column label="活动简介" width="200">
          <template #default="{ row }">
            <el-tooltip content="{{ row.intro }}" placement="top">
              <span class="text-ellipsis">{{ row.intro }}</span>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="活动状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.state == 1 ? 'success' : 'warning'">
              {{ row.state == 1 ? '已审核' : '待审核' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleView(row)">
              <el-icon><View /></el-icon> 查看
            </el-button>
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              <el-icon><Edit /></el-icon> 编辑
            </el-button>
            <el-button 
              v-if="row.state == 0" 
              type="success" 
              link 
              size="small" 
              @click="handleApproveSingle(row.id)"
            >
              <el-icon><Check /></el-icon> 审核
            </el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">
              <el-icon><Delete /></el-icon> 删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <el-pagination
        v-if="total > 0"
        v-model:current-page="queryParams.pageNum"
        v-model:page-size="queryParams.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
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

// 查询参数（与后端一致的字段名）
const queryParams = reactive({
  activityName: '',
  addr: '',                // 活动地点
  createStartTime: '',
  createEndTime: '',
  startStartTime: '',
  startEndTime: '',
  state: '',               // 活动状态：0待审核，1已审核
  pageNum: 1,
  pageSize: 10
})

// 日期范围
const createDateRange = ref([])
const startDateRange = ref([])

// 活动列表
const activityList = ref([])

// 原始活动数据（从API获取后存储）
const activityData = ref([])

// 选中的活动ID
const selectedIds = ref([])

// 总数
const total = ref(0)

// 加载状态
const loading = ref(false)

// 统计数据
const totalActivities = ref(0)       // 总活动数
const pendingActivities = ref(0)     // 待审核活动数
const trendClass = ref('positive')   // 趋势样式
const trendText = ref('12.5% 较上月') // 趋势文本
const pendingTrendClass = ref('warning') // 待审核趋势样式
const pendingTrendText = ref('待处理中') // 待审核趋势文本

// 监听创建日期范围变化，更新查询参数
const watchCreateDateRange = () => {
  if (createDateRange.value && createDateRange.value.length === 2) {
    queryParams.createStartTime = createDateRange.value[0]
    queryParams.createEndTime = createDateRange.value[1]
  } else {
    queryParams.createStartTime = ''
    queryParams.createEndTime = ''
  }
}

// 监听开始日期范围变化，更新查询参数
const watchStartDateRange = () => {
  if (startDateRange.value && startDateRange.value.length === 2) {
    queryParams.startStartTime = startDateRange.value[0]
    queryParams.startEndTime = startDateRange.value[1]
  } else {
    queryParams.startStartTime = ''
    queryParams.startEndTime = ''
  }
}

// 过滤活动数据（包含状态、时间、名称、地点等条件）
const filterActivities = () => {
  let filtered = [...activityData.value]
  
  // 按活动名称过滤
  if (queryParams.activityName) {
    filtered = filtered.filter(activity => 
      activity.name.includes(queryParams.activityName)
    )
  }
  
  // 按活动地点过滤
  if (queryParams.addr) {
    filtered = filtered.filter(activity => 
      activity.addr.includes(queryParams.addr)
    )
  }
  
  // 按创建时间过滤
  if (queryParams.createStartTime && queryParams.createEndTime) {
    const start = new Date(queryParams.createStartTime)
    const end = new Date(queryParams.createEndTime)
    end.setDate(end.getDate() + 1) // 包含结束日期
    
    filtered = filtered.filter(activity => {
      const createTime = new Date(activity.createTime)
      return createTime >= start && createTime < end
    })
  }
  
  // 按开始时间过滤
  if (queryParams.startStartTime && queryParams.startEndTime) {
    const start = new Date(queryParams.startStartTime)
    const end = new Date(queryParams.startEndTime)
    end.setDate(end.getDate() + 1) // 包含结束日期
    
    filtered = filtered.filter(activity => {
      const startTime = new Date(activity.startTime)
      return startTime >= start && startTime < end
    })
  }
  
  // 按活动状态过滤（0待审核，1已审核）
  if (queryParams.state !== '') {
    filtered = filtered.filter(activity => activity.state == queryParams.state)
  }
  
  // 更新统计数据
  totalActivities.value = filtered.length
  pendingActivities.value = filtered.filter(activity => activity.state == 0).length
  
  // 分页处理
  const startIndex = (queryParams.pageNum - 1) * queryParams.pageSize
  const endIndex = startIndex + queryParams.pageSize
  
  activityList.value = filtered.slice(startIndex, endIndex)
  total.value = filtered.length
}

// 获取活动列表（从API请求数据）
const getList = async () => {
  loading.value = true
  try {
    const response = await axios.get('api/v1/activity', {
      headers: {
        Authorization: getToken()
      }
    })
    
    activityData.value = response.data || []
    console.log('活动数据:', activityData.value)
    
    // 更新统计数据
    totalActivities.value = activityData.value.length
    pendingActivities.value = activityData.value.filter(activity => activity.state == 0).length
    
    watchCreateDateRange()
    watchStartDateRange()
    filterActivities()
  } catch (error) {
    console.error('获取活动列表失败', error)
    ElMessage.error('获取活动列表失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

// 处理查询（触发数据获取和过滤）
const handleQuery = () => {
  queryParams.pageNum = 1
  getList()
}

// 重置查询条件
const resetQuery = () => {
  Object.assign(queryParams, {
    activityName: '',
    addr: '',                
    createStartTime: '',
    createEndTime: '',
    startStartTime: '',
    startEndTime: '',
    state: '',              
    pageNum: 1,
    pageSize: 10
  })
  createDateRange.value = []
  startDateRange.value = []
  getList()
}

// 处理表格选择变化
const handleSelectionChange = (selection) => {
  selectedIds.value = selection.map(item => item.id)
}

// 处理分页大小变化
const handleSizeChange = (size) => {
  queryParams.pageSize = size
  getList()
}

// 处理页码变化
const handleCurrentChange = (page) => {
  queryParams.pageNum = page
  getList()
}

// 格式化日期（从时间戳转换为可读格式）
const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 处理查看活动详情
const handleView = (row) => {
  router.push(`/activity/detail/${row.id}`)
}

// 处理新增活动
const handleAdd = () => {
  router.push('/activity/add')
}

// 处理编辑活动
const handleEdit = (row) => {
  router.push(`/activity/edit/${row.id}`)
}

// 处理删除单个活动
const handleDelete = (row) => {
  ElMessageBox.confirm(`确定要删除活动 "${row.name}" 吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await axios.delete(`/api/v1/activity/${row.id}`, {
        headers: {
          Authorization: getToken()
        }
      })
      
      // 更新本地数据
      activityData.value = activityData.value.filter(activity => activity.id !== row.id)
      
      ElMessage.success('删除成功')
      getList()
    } catch (error) {
      console.error('删除活动失败', error)
      ElMessage.error('删除活动失败，请稍后重试')
    }
  }).catch(() => {})
}

// 处理批量删除活动
const handleBatchDelete = () => {
  if (selectedIds.value.length === 0) {
    ElMessage.warning('请选择要删除的活动')
    return
  }
  
  ElMessageBox.confirm(`确定要删除选中的${selectedIds.value.length}个活动吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      // 批量删除API请求
      await Promise.all(selectedIds.value.map(id => 
        axios.delete(`/api/v1/activity/${id}`, {
          headers: {
            Authorization: getToken()
          }
        })
      ))
      
      // 更新本地数据
      activityData.value = activityData.value.filter(activity => !selectedIds.value.includes(activity.id))
      
      ElMessage.success('批量删除成功')
      selectedIds.value = []
      getList()
    } catch (error) {
      console.error('批量删除失败', error)
      ElMessage.error('批量删除失败，请稍后重试')
    }
  }).catch(() => {})
}

// 处理单个活动审核通过（融合新的API请求逻辑）
const handleApproveSingle = (activityId) => {
  ElMessageBox.confirm('确定要审核通过此活动吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'question'
  }).then(async () => {
    try {
      // 发送审核请求，将活动ID作为路径参数，state作为请求体参数
      const newState = 1; // 审核通过，状态设为1
      await axios.post(`/api/v1/admin/users/${activityId}`, {
        state: newState
      }, {
        headers: {
          Authorization: getToken()
        }
      });
      
      // 更新本地数据
      const activityIndex = activityData.value.findIndex(activity => activity.id === activityId)
      if (activityIndex !== -1) {
        activityData.value[activityIndex].state = newState; // 更新为已审核状态
      }
      
      ElMessage.success('审核通过成功');
      getList(); // 刷新列表
    } catch (error) {
      console.error('审核活动失败', error);
      ElMessage.error('审核活动失败，请稍后重试');
    }
  }).catch(() => {});
}

// 处理批量审核（仅支持单个活动审核）
const handleApprove = () => {
  if (selectedIds.value.length !== 1) {
    ElMessage.warning('请选择一个活动进行审核')
    return
  }
  
  handleApproveSingle(selectedIds.value[0])
}

// 处理导出功能（暂未实现）
const handleExport = () => {
  ElMessage.info('导出功能开发中')
}

// 组件挂载时获取数据
onMounted(() => {
  getList()
})
</script>

<style lang="scss" scoped>
.activity-management-container {
  .search-card {
    margin-bottom: 20px;
  }
  
  .stats-card-container {
    margin-bottom: 20px;
    
    .stats-card {
      height: 120px;
      
      .stats-content {
        display: flex;
        flex-direction: column;
        justify-content: center;
        height: 100%;
        
        .stats-title {
          font-size: 14px;
          color: #909399;
          margin-bottom: 5px;
        }
        
        .stats-value {
          font-size: 28px;
          font-weight: bold;
          color: #303133;
          margin-bottom: 10px;
        }
        
        .stats-trend {
          font-size: 12px;
          color: #606266;
          
          &.positive {
            color: #67c23a; /* 绿色 - 正增长 */
          }
          
          &.warning {
            color: #e6a23c; /* 黄色 - 警告 */
          }
          
          &.negative {
            color: #f56c6c; /* 红色 - 负增长 */
          }
        }
      }
    }
  }
  
  .table-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      .right-actions {
        display: flex;
        gap: 10px;
      }
    }
    
    .el-pagination {
      margin-top: 20px;
      justify-content: flex-end;
    }
    
    .text-ellipsis {
      display: block;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      max-width: 180px;
    }
  }
}
</style>