<template>
  <div class="user-management-container">
    <!-- 搜索区域 -->
    <el-card class="search-card">
      <el-form :inline="true" :model="queryParams" class="search-form">
        <el-form-item label="用户名">
          <el-input v-model="queryParams.username" placeholder="请输入用户名" clearable />
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="queryParams.addr" placeholder="请输入地址" clearable />
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

    <!-- 数据统计卡片 -->
    <el-row :gutter="20" class="stats-card-container">
      <el-col :span="8">
        <el-card class="stats-card" shadow="hover">
          <div class="stats-content">
            <div class="stats-title">总用户数</div>
            <div class="stats-value">{{ totalUsers }}</div>
            <div class="stats-trend">
              <i class="el-icon-caret-top" :class="trendClass"></i>
              <span>{{ trendText }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="stats-card" shadow="hover">
          <div class="stats-content">
            <div class="stats-title">活跃用户</div>
            <div class="stats-value">{{ activeUsers }}</div>
            <div class="stats-trend positive">
              <i class="el-icon-caret-top"></i>
              <span>12.5% 较上月</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="stats-card" shadow="hover">
          <div class="stats-content">
            <div class="stats-title">今日新增</div>
            <div class="stats-value">{{ todayNewUsers }}</div>
            <div class="stats-trend positive">
              <i class="el-icon-caret-top"></i>
              <span>5 人</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 表格区域 -->
    <el-card class="table-card">
      <template #header>
        <div class="card-header">
          <span>用户管理列表</span>
          <div class="right-actions">
            <el-button type="primary" @click="handleAdd">
              <el-icon><Plus /></el-icon>
              新增用户
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
        :data="userList"
        @selection-change="handleSelectionChange"
        stripe
        fit
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="用户ID" prop="id" width="80" />
        <el-table-column label="用户名" prop="username" width="120" />
        <el-table-column label="账号" prop="name" width="120" />
        <el-table-column label="性别" width="80">
          <template #default="{ row }">
            <el-tag :type="row.gender === '男' ? 'primary' : 'success'">
              {{ row.gender }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="地址" prop="addr" width="180" />
        <el-table-column label="创建时间" prop="createTime" width="180">
          <template #default="{ row }">
            {{ formatDate(row.createTime) }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : row.status === 'pending' ? 'warning' : 'danger'">
              {{ statusText[row.status] }}
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
// 引入axios和获取token的函数
import axios from 'axios'
import { getToken } from '@/utils/auth'

const router = useRouter()

// 状态文本映射
const statusText = {
  active: '活跃',
  pending: '待审核',
  disabled: '已禁用'
}

// 查询参数
const queryParams = reactive({
  username: '',
  addr: '',
  startTime: '',
  endTime: '',
  pageNum: 1,
  pageSize: 10
})

// 日期范围
const dateRange = ref([])

// 用户列表
const userList = ref([])

// 原始用户数据
const userData = ref([])

// 选中的用户ID
const selectedIds = ref([])

// 总数
const total = ref(0)

// 加载状态
const loading = ref(false)

// 统计数据
const totalUsers = ref(0)
const activeUsers = ref(0)
const todayNewUsers = ref(0)
const trendClass = ref('positive')
const trendText = ref('8.2% 较上月')

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

// 过滤用户数据
const filterUsers = () => {
  let filtered = [...userData.value]
  
  // 按用户名过滤
  if (queryParams.username) {
    filtered = filtered.filter(user => 
      user.username.includes(queryParams.username) || 
      user.name.includes(queryParams.username)
    )
  }
  
  // 按地址过滤
  if (queryParams.addr) {
    filtered = filtered.filter(user => 
      user.addr.includes(queryParams.addr)
    )
  }
  
  // 按时间过滤
  if (queryParams.startTime && queryParams.endTime) {
    const start = new Date(queryParams.startTime)
    const end = new Date(queryParams.endTime)
    end.setDate(end.getDate() + 1) // 包含结束日期
    
    filtered = filtered.filter(user => {
      const createTime = new Date(user.createTime)
      return createTime >= start && createTime < end
    })
  }
  
  // 更新统计数据
  totalUsers.value = filtered.length
  activeUsers.value = filtered.filter(user => user.status === 'active').length
  
  // 分页处理
  const startIndex = (queryParams.pageNum - 1) * queryParams.pageSize
  const endIndex = startIndex + queryParams.pageSize
  
  userList.value = filtered.slice(startIndex, endIndex)
  total.value = filtered.length
}

// 获取用户列表
const getList = async () => {
  loading.value = true
  try {
    // 发起API请求获取用户数据
    const response = await axios.get('/api/v1/admin/users', {
      headers: {
        Authorization: getToken()
      }
    })
    
    userData.value = response.data || []
    console.log('用户数据:', userData.value)
    
    // 更新统计数据
    totalUsers.value = userData.value.length
    activeUsers.value = userData.value.filter(user => user.status === 'active').length
    
    // 应用过滤和分页
    watchDateRange()
    filterUsers()
  } catch (error) {
    console.error('获取用户列表失败', error)
    ElMessage.error('获取用户列表失败，请稍后重试')
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
    addr: '',
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

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 处理查看用户
const handleView = (row) => {
  router.push(`/user/detail/${row.id}`)
}

// 处理新增用户
const handleAdd = () => {
  router.push('/user/add')
}

// 处理编辑用户
const handleEdit = (row) => {
  router.push(`/user/edit/${row.id}`)
}

// 处理删除用户
const handleDelete = (row) => {
  ElMessageBox.confirm(`确定要删除用户 ${row.username} 吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      // 调用API删除用户
      await axios.delete(`/api/v1/admin/users/${row.id}`, {
        headers: {
          Authorization: getToken()
        }
      })
      
      // 从本地数据中移除
      userData.value = userData.value.filter(user => user.id !== row.id)
      
      ElMessage.success('删除成功')
      getList()
    } catch (error) {
      console.error('删除用户失败', error)
      ElMessage.error('删除用户失败，请稍后重试')
    }
  }).catch(() => {})
}

// 处理批量删除
const handleBatchDelete = () => {
  if (selectedIds.value.length === 0) {
    ElMessage.warning('请选择要删除的用户')
    return
  }
  
  ElMessageBox.confirm(`确定要删除选中的${selectedIds.value.length}个用户吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      // 批量删除
      await Promise.all(selectedIds.value.map(id => 
        axios.delete(`/api/v1/admin/users/${id}`, {
          headers: {
            Authorization: getToken()
          }
        })
      ))
      
      // 更新本地数据
      userData.value = userData.value.filter(user => !selectedIds.value.includes(user.id))
      
      ElMessage.success('批量删除成功')
      selectedIds.value = []
      getList()
    } catch (error) {
      console.error('批量删除失败', error)
      ElMessage.error('批量删除失败，请稍后重试')
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
/* 样式保持不变 */
</style>