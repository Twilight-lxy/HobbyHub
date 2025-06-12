<template>
  <div class="user-list-container">
    <div class="user-header">
      <h2>用户列表</h2>
    </div>
    
    <!-- 搜索和筛选 -->
    <el-card class="filter-container" shadow="hover">
      <el-form :model="queryParams" ref="queryForm" :inline="true">
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="queryParams.username"
            placeholder="请输入用户名"
            clearable
            @keyup.enter="handleQuery"
          />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input
            v-model="queryParams.phone"
            placeholder="请输入手机号"
            clearable
            @keyup.enter="handleQuery"
          />
        </el-form-item>
        <el-form-item label="用户状态" prop="status">
          <el-select
            v-model="queryParams.status"
            placeholder="请选择用户状态"
            clearable
          >
            <el-option label="正常" :value="0" />
            <el-option label="禁用" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item label="注册时间">
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
    
    <!-- 用户表格 -->
    <el-card class="table-container" shadow="hover">
      <el-table
        v-loading="loading"
        :data="userList"
        border
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="用户信息" min-width="200">
          <template #default="scope">
            <div class="user-info">
              <el-avatar :size="40" :src="scope.row.avatar" />
              <div class="user-detail">
                <div class="username">{{ scope.row.nickname || scope.row.username }}</div>
                <div class="user-id">ID: {{ scope.row.id }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="phone" label="手机号" width="120" />
        <el-table-column prop="email" label="邮箱" width="180" show-overflow-tooltip />
        <el-table-column prop="gender" label="性别" width="80">
          <template #default="scope">
            {{ scope.row.gender === 1 ? '男' : scope.row.gender === 2 ? '女' : '未知' }}
          </template>
        </el-table-column>
        <el-table-column prop="registerTime" label="注册时间" width="180" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="scope">
            <el-button size="small" @click="handleView(scope.row)">查看</el-button>
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
const userList = ref([])
const total = ref(0)
const dateRange = ref([])

// 查询参数
const queryParams = reactive({
  pageNum: 1,
  pageSize: 10,
  username: '',
  phone: '',
  status: undefined,
  startTime: undefined,
  endTime: undefined
})

// 查询用户列表
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
    // 获取用户列表
    const res = await request.get('/api/user/list', {
      params: queryParams
    })
    
    if (res && res.code === 200) {
      const userData = res.data || []
      
      // 过滤数据
      let filteredList = userData
      
      // 按用户名筛选
      if (queryParams.username) {
        filteredList = filteredList.filter(item => 
          (item.username && item.username.toLowerCase().includes(queryParams.username.toLowerCase())) ||
          (item.nickname && item.nickname.toLowerCase().includes(queryParams.username.toLowerCase()))
        )
      }
      
      // 按手机号筛选
      if (queryParams.phone) {
        filteredList = filteredList.filter(item => 
          item.phone && item.phone.includes(queryParams.phone)
        )
      }
      
      // 按状态筛选
      if (queryParams.status !== undefined) {
        filteredList = filteredList.filter(item => item.status === queryParams.status)
      }
      
      // 按时间筛选
      if (queryParams.startTime && queryParams.endTime) {
        filteredList = filteredList.filter(item => {
          const registerTime = new Date(item.registerTime).getTime()
          const startTime = new Date(queryParams.startTime).getTime()
          const endTime = new Date(queryParams.endTime).getTime()
          return registerTime >= startTime && registerTime <= endTime
        })
      }
      
      // 分页
      total.value = filteredList.length
      const start = (queryParams.pageNum - 1) * queryParams.pageSize
      const end = start + queryParams.pageSize
      
      userList.value = filteredList.slice(start, end)
    } else {
      ElMessage.error(res.msg || '获取用户列表失败')
    }
  } catch (error) {
    console.error('获取用户列表失败', error)
    ElMessage.error('获取用户列表失败')
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
    username: '',
    phone: '',
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

// 查看用户
const handleView = (row) => {
  router.push(`/user/detail/${row.id}`)
}

// 修改用户状态
const handleStatusChange = (row) => {
  const statusText = row.status === 0 ? '禁用' : '启用'
  ElMessageBox.confirm(`确定要${statusText}用户"${row.nickname || row.username}"吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      // 调用实际的API
      const res = await request.put(`/api/user/status`, {
        id: row.id,
        status: row.status === 0 ? 1 : 0
      })
      
      if (res && res.code === 200) {
        row.status = row.status === 0 ? 1 : 0
        ElMessage.success(`${statusText}成功`)
      } else {
        ElMessage.error(res.msg || `${statusText}失败`)
      }
    } catch (error) {
      console.error(`${statusText}失败`, error)
      ElMessage.error(`${statusText}失败`)
    }
  }).catch(() => {})
}

// 重置密码
const handleReset = (row) => {
  ElMessageBox.confirm(`确定要重置用户"${row.nickname || row.username}"的密码吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      // 调用实际的API
      const res = await request.put(`/api/user/reset-password/${row.id}`)
      
      if (res && res.code === 200) {
        ElMessage.success('密码重置成功，新密码已发送至用户手机')
      } else {
        ElMessage.error(res.msg || '密码重置失败')
      }
    } catch (error) {
      console.error('密码重置失败', error)
      ElMessage.error('密码重置失败')
    }
  }).catch(() => {})
}

// 初始化
onMounted(() => {
  getList()
})
</script>

<style lang="scss" scoped>
.user-list-container {
  .user-header {
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
    .user-info {
      display: flex;
      align-items: center;
      
      .user-detail {
        margin-left: 10px;
        
        .username {
          font-size: 14px;
          font-weight: 500;
        }
        
        .user-id {
          font-size: 12px;
          color: #909399;
        }
      }
    }
    
    .pagination-container {
      margin-top: 20px;
      display: flex;
      justify-content: flex-end;
    }
  }
}
</style> 