<template>
  <div class="order-list-container">
    <!-- 搜索区域 -->
    <el-card class="search-card">
      <el-form :inline="true" :model="queryParams" class="search-form">
        <el-form-item label="用户名">
          <el-input v-model="queryParams.username" placeholder="请输入用户名" clearable />
        </el-form-item>
        <el-form-item label="活动名称">
          <el-input v-model="queryParams.activityName" placeholder="请输入活动名称" clearable />
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

    <!-- 表格区域 -->
    <el-card class="table-card">
      <template #header>
        <div class="card-header">
          <span>参与记录列表</span>
          <div class="right">
            <el-button type="primary" @click="handleAdd">
              <el-icon><Plus /></el-icon>
              新增记录
            </el-button>
            <el-button type="danger" :disabled="!selectedIds.length" @click="handleBatchDelete">
              <el-icon><Delete /></el-icon>
              批量删除
            </el-button>
            <el-button type="success" @click="handleExport">
              <el-icon><Download /></el-icon>
              导出
            </el-button>
          </div>
        </div>
      </template>
      
      <!-- 表格 -->
      <el-table
        v-loading="loading"
        :data="orderList"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="记录ID" prop="id" width="80" />
        <el-table-column label="用户" width="120">
          <template #default="{ row }">
            <div class="user-info">
              <span>{{ row.user?.username || '未知用户' }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="活动信息">
          <template #default="{ row }">
            <div class="activity-info">
              <div>{{ row.product?.name || '未知活动' }}</div>
              <div class="activity-category">{{ row.product?.category || '未知分类' }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="参与人数" prop="count" width="100" />
        <el-table-column label="总人数" width="100">
          <template #default="{ row }">
            {{ row.totalPrice || 0 }}
          </template>
        </el-table-column>
        <el-table-column label="参与时间" prop="createTime" width="180" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">
              <el-icon><View /></el-icon>
              查看
            </el-button>
            <el-button type="primary" link @click="handleEdit(row)">
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-button type="danger" link @click="handleDelete(row)">
              <el-icon><Delete /></el-icon>
              删除
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
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getOrderList, deleteOrder } from '@/api/order'

const router = useRouter()

// 查询参数
const queryParams = reactive({
  username: '',
  activityName: '',
  startTime: '',
  endTime: '',
  pageNum: 1,
  pageSize: 10
})

// 日期范围
const dateRange = ref([])

// 参与记录列表
const orderList = ref([])

// 选中的记录ID
const selectedIds = ref([])

// 总数
const total = ref(0)

// 加载状态
const loading = ref(false)

// 监听日期范围变化
const watchDateRange = () => {
  if (dateRange.value && dateRange.value.length === 2) {
    queryParams.startTime = dateRange.value[0]
    queryParams.endTime = dateRange.value[1]
  } else {
    queryParams.startTime = ''
    queryParams.endTime = ''
  }
}

// 获取参与记录列表
const getList = async () => {
  loading.value = true
  watchDateRange()
  try {
    const res = await getOrderList({
      ...queryParams,
      pageNum: queryParams.pageNum - 1 // 后端从0开始计数
    })
    orderList.value = res.data.records || []
    total.value = res.data.total || 0
  } catch (error) {
    console.error('获取参与记录列表失败', error)
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
  Object.assign(queryParams, {
    username: '',
    activityName: '',
    startTime: '',
    endTime: '',
    pageNum: 1,
    pageSize: 10
  })
  dateRange.value = []
  getList()
}

// 处理选择变化
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

// 处理查看
const handleView = (row) => {
  router.push(`/order/detail/${row.id}`)
}

// 处理新增
const handleAdd = () => {
  ElMessage.info('功能开发中')
}

// 处理编辑
const handleEdit = (row) => {
  ElMessage.info('功能开发中')
}

// 处理删除
const handleDelete = (row) => {
  ElMessageBox.confirm(`确定要删除该参与记录吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteOrder(row.id)
      ElMessage.success('删除成功')
      getList()
    } catch (error) {
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

// 处理批量删除
const handleBatchDelete = () => {
  if (selectedIds.value.length === 0) {
    ElMessage.warning('请选择要删除的记录')
    return
  }
  
  ElMessageBox.confirm(`确定要删除选中的${selectedIds.value.length}条记录吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      // 批量删除
      await Promise.all(selectedIds.value.map(id => deleteOrder(id)))
      ElMessage.success('批量删除成功')
      getList()
    } catch (error) {
      ElMessage.error('批量删除失败')
    }
  }).catch(() => {})
}

// 处理导出
const handleExport = () => {
  ElMessage.info('导出功能开发中')
}

onMounted(() => {
  getList()
})
</script>

<style lang="scss" scoped>
.order-list-container {
  .search-card {
    margin-bottom: 20px;
  }
  
  .table-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      .right {
        display: flex;
        gap: 10px;
      }
    }
    
    .user-info {
      display: flex;
      flex-direction: column;
    }
    
    .activity-info {
      display: flex;
      flex-direction: column;
      
      .activity-category {
        font-size: 12px;
        color: #909399;
        margin-top: 5px;
      }
    }
    
    .el-pagination {
      margin-top: 20px;
      justify-content: flex-end;
    }
  }
}
</style> 