<template>
  <div class="team-list-container">
    <div class="team-header">
      <h2>小队列表</h2>
    </div>
    
    <!-- 搜索和筛选 -->
    <el-card class="filter-container" shadow="hover">
      <el-form :model="queryParams" ref="queryForm" :inline="true">
        <el-form-item label="小队名称" prop="name">
          <el-input
            v-model="queryParams.name"
            placeholder="请输入小队名称"
            clearable
            @keyup.enter="handleQuery"
          />
        </el-form-item>
        <el-form-item label="关联活动" prop="activityId">
          <el-select
            v-model="queryParams.activityId"
            placeholder="请选择关联活动"
            clearable
          >
            <el-option
              v-for="item in activityOptions"
              :key="item.id"
              :label="item.title"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="小队状态" prop="status">
          <el-select
            v-model="queryParams.status"
            placeholder="请选择小队状态"
            clearable
          >
            <el-option label="招募中" :value="0" />
            <el-option label="已组满" :value="1" />
            <el-option label="已解散" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="创建时间">
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
    
    <!-- 小队表格 -->
    <el-card class="table-container" shadow="hover">
      <el-table
        v-loading="loading"
        :data="teamList"
        border
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="小队名称" min-width="180" show-overflow-tooltip>
          <template #default="scope">
            <el-link type="primary" @click="handleView(scope.row)">{{ scope.row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="activityTitle" label="关联活动" width="180" show-overflow-tooltip />
        <el-table-column label="成员数量" width="120" align="center">
          <template #default="scope">
            {{ scope.row.memberCount }}/{{ scope.row.maxCount }}
          </template>
        </el-table-column>
        <el-table-column prop="leaderName" label="队长" width="120" />
        <el-table-column prop="createTime" label="创建时间" width="180" />
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
            >解散</el-button>
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
const teamList = ref([])
const total = ref(0)
const dateRange = ref([])

// 查询参数
const queryParams = reactive({
  pageNum: 1,
  pageSize: 10,
  name: '',
  activityId: undefined,
  status: undefined,
  startTime: undefined,
  endTime: undefined
})

// 活动选项
const activityOptions = ref([])

// 获取活动选项
const fetchActivityOptions = async () => {
  try {
    const res = await request.get('/api/product/list')
    if (res && res.code === 200) {
      activityOptions.value = (res.data || []).map(item => ({
        id: item.id,
        title: item.name
      }))
    }
  } catch (error) {
    console.error('获取活动列表失败', error)
  }
}

// 获取状态类型和文本
const getStatusType = (status) => {
  const types = ['success', 'warning', 'info']
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = ['招募中', '已组满', '已解散']
  return texts[status] || '未知'
}

// 查询小队列表
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
    // 调用小队列表API - 使用正确的API路径
    const res = await request.get('/api/product_cart/list')
    
    if (res && res.code === 200) {
      const teams = res.data || []
      
      // 处理数据格式，将product_cart数据转换为team格式
      const formattedTeams = teams.map(item => {
        return {
          id: item.id,
          name: item.product?.name || '未命名小队',
          activityTitle: item.product?.name || '未知活动',
          activityId: item.productId,
          memberCount: item.count || 0,
          maxCount: item.product?.price ? item.product.price * 10 : 0,
          leaderName: item.user?.nickname || '未知',
          createTime: item.createTime,
          status: item.status || 0
        }
      })
      
      // 过滤数据
      let filteredList = formattedTeams
      
      // 按名称筛选
      if (queryParams.name) {
        filteredList = filteredList.filter(item => 
          item.name && item.name.toLowerCase().includes(queryParams.name.toLowerCase())
        )
      }
      
      // 按活动筛选
      if (queryParams.activityId) {
        filteredList = filteredList.filter(item => item.activityId === queryParams.activityId)
      }
      
      // 按状态筛选
      if (queryParams.status !== undefined) {
        filteredList = filteredList.filter(item => item.status === queryParams.status)
      }
      
      // 按时间筛选
      if (queryParams.startTime && queryParams.endTime) {
        filteredList = filteredList.filter(item => {
          const createTime = new Date(item.createTime).getTime()
          const startTime = new Date(queryParams.startTime).getTime()
          const endTime = new Date(queryParams.endTime).getTime()
          return createTime >= startTime && createTime <= endTime
        })
      }
      
      // 分页
      total.value = filteredList.length
      const start = (queryParams.pageNum - 1) * queryParams.pageSize
      const end = start + queryParams.pageSize
      
      teamList.value = filteredList.slice(start, end)
    } else {
      ElMessage.error(res.msg || '获取小队列表失败')
    }
  } catch (error) {
    console.error('获取小队列表失败', error)
    ElMessage.error('获取小队列表失败')
  } finally {
    loading.value = false
  }
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
    name: '',
    activityId: undefined,
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

// 创建小队
const handleCreate = () => {
  router.push('/team/create')
}

// 查看小队
const handleView = (row) => {
  router.push(`/team/detail/${row.id}`)
}

// 编辑小队
const handleEdit = (row) => {
  router.push(`/team/edit/${row.id}`)
}

// 解散小队
const handleDelete = (row) => {
  ElMessageBox.confirm(`确定要解散"${row.name}"小队吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      // 调用实际的API - 使用正确的API路径
      const res = await request.delete(`/api/product_cart/remove/ids?ids=${row.id}`)
      
      if (res && res.code === 200) {
        ElMessage.success('小队已解散')
        getList() // 重新加载列表
      } else {
        ElMessage.error(res.msg || '解散失败')
      }
    } catch (error) {
      console.error('解散失败', error)
      ElMessage.error('解散失败')
    }
  }).catch(() => {})
}

// 初始化
onMounted(() => {
  fetchActivityOptions()
  getList()
})
</script>

<style lang="scss" scoped>
.team-list-container {
  .team-header {
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