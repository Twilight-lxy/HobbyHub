<template>
  <div class="user-container">
    <!-- 搜索区域 -->
    <el-card class="search-card">
      <el-form :inline="true" :model="queryParams" class="search-form">
        <el-form-item label="用户名">
          <el-input v-model="queryParams.username" placeholder="请输入用户名" clearable />
        </el-form-item>
        <el-form-item label="姓名">
          <el-input v-model="queryParams.name" placeholder="请输入姓名" clearable />
        </el-form-item>
        <el-form-item label="性别">
          <el-select v-model="queryParams.gender" placeholder="请选择性别" clearable>
            <el-option label="男" :value="1" />
            <el-option label="女" :value="0" />
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
          <span>用户列表</span>
          <div class="right">
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
              导出
            </el-button>
            <el-button type="warning" @click="handleImport">
              <el-icon><Upload /></el-icon>
              导入
            </el-button>
          </div>
        </div>
      </template>
      
      <!-- 表格 -->
      <el-table
        v-loading="loading"
        :data="userList"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="用户ID" prop="id" width="80" />
        <el-table-column label="用户信息" width="200">
          <template #default="{ row }">
            <div class="user-info">
              <el-avatar :size="32" :src="row.avatar || defaultAvatar" />
              <div class="user-detail">
                <div class="username">{{ row.username }}</div>
                <div class="nickname">{{ row.name }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="性别" prop="gender" width="80">
          <template #default="{ row }">
            <el-tag :type="row.gender === '男' ? 'primary' : 'success'">
              {{ row.gender }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="地址" prop="addr" width="120" />
        <el-table-column label="创建时间" prop="createTime" width="180" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-button type="primary" link @click="handleResetPwd(row)">
              <el-icon><Key /></el-icon>
              重置密码
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

    <!-- 重置密码对话框 -->
    <el-dialog
      title="重置密码"
      v-model="resetPwdDialogVisible"
      width="500px"
    >
      <el-form :model="pwdForm" :rules="pwdRules" ref="pwdFormRef" label-width="100px">
        <el-form-item label="新密码" prop="password">
          <el-input v-model="pwdForm.password" type="password" placeholder="请输入新密码" show-password />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="pwdForm.confirmPassword" type="password" placeholder="请再次输入新密码" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resetPwdDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitResetPwd">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'
import { getToken } from '@/utils/auth'

// 默认头像
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

// 查询参数
const queryParams = reactive({
  username: '',
  name: '',
  gender: '',
  startTime: '',
  endTime: '',
  pageNum: 1,
  pageSize: 10
})

// 日期范围
const dateRange = ref([])

// 用户列表
const userList = ref([])

// 选中的用户ID
const selectedIds = ref([])

// 总数
const total = ref(0)

// 加载状态
const loading = ref(false)

// 对话框可见性
const resetPwdDialogVisible = ref(false)

// 密码表单
const pwdForm = reactive({
  userId: '',
  password: '',
  confirmPassword: ''
})

// 当前操作的用户
const currentUser = reactive({
  id: '',
  username: ''
})

// 密码表单验证规则
const pwdRules = {
  password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== pwdForm.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

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

// 获取用户列表
const getList = async () => {
  loading.value = true
  watchDateRange()
  try {
    // 修正API请求格式
    const response = await axios.get('/api/v1/admin/users', {
      headers: {
        'Authorization': `Bearer ${getToken()}`
      },
      params: {
        username: queryParams.username,
        name: queryParams.name,
        gender: queryParams.gender,
        startTime: queryParams.startTime,
        endTime: queryParams.endTime,
        pageNum: queryParams.pageNum,
        pageSize: queryParams.pageSize
      }
    })
    
    if (response.data.code === 200) {
      // 转换API数据以匹配表格结构
      userList.value = response.data.data.map(user => ({
        ...user,
        // 格式化日期
        createTime: formatDate(user.createTime)
      }))
      total.value = response.data.total || userList.value.length
    } else {
      ElMessage.error(response.data.msg || '获取用户列表失败')
    }
  } catch (error) {
    console.error('获取用户列表失败', error)
    ElMessage.error('获取用户列表失败，请检查网络连接')
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
    name: '',
    gender: '',
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

// 处理编辑
const handleEdit = (row) => {
  ElMessage.info(`编辑用户: ${row.username}`)
}

// 处理删除
const handleDelete = (row) => {
  ElMessageBox.confirm(`确定要删除用户"${row.username}"吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    // 模拟删除操作
    const index = userList.value.findIndex(item => item.id === row.id)
    if (index !== -1) {
      userList.value.splice(index, 1)
      total.value--
      ElMessage.success('删除成功')
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
  }).then(() => {
    // 模拟批量删除操作
    userList.value = userList.value.filter(item => !selectedIds.value.includes(item.id))
    total.value = userList.value.length
    selectedIds.value = []
    ElMessage.success('批量删除成功')
  }).catch(() => {})
}

// 处理重置密码
const handleResetPwd = (row) => {
  pwdForm.userId = row.id
  pwdForm.password = ''
  pwdForm.confirmPassword = ''
  currentUser.username = row.username
  resetPwdDialogVisible.value = true
}

// 提交重置密码
const submitResetPwd = () => {
  pwdFormRef.value.validate((valid) => {
    if (valid) {
      // 模拟API请求
      ElMessage.success(`用户"${currentUser.username}"的密码已重置`)
      resetPwdDialogVisible.value = false
    }
  })
}

// 处理导出
const handleExport = () => {
  ElMessage.success('用户数据导出成功')
}

// 处理导入
const handleImport = () => {
  ElMessage.info('导入功能开发中...')
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString()
}

onMounted(() => {
  getList()
})
</script>

<style lang="scss" scoped>
.user-container {
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
      align-items: center;
      
      .user-detail {
        margin-left: 10px;
        
        .username {
          font-weight: bold;
        }
        
        .nickname {
          font-size: 12px;
          color: #909399;
        }
      }
    }
    
    .el-pagination {
      margin-top: 20px;
      justify-content: flex-end;
    }
  }
}
</style>