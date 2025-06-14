<template>
  <div class="user-container">
    <!-- 搜索区域 -->
    <el-card class="search-card">
      <el-form :inline="true" :model="queryParams" class="search-form">
        <el-form-item label="用户名">
          <el-input v-model="queryParams.username" placeholder="请输入用户名" clearable />
        </el-form-item>
        <el-form-item label="手机号码">
          <el-input v-model="queryParams.phone" placeholder="请输入手机号码" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="请选择状态" clearable>
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
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
              <el-avatar :size="32" :src="row.avatar" />
              <div class="user-detail">
                <div class="username">{{ row.username }}</div>
                <div class="nickname">{{ row.nickname }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="手机号码" prop="phone" width="120" />
        <el-table-column label="邮箱" prop="email" width="180" />
        <el-table-column label="角色" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.roleId === 1" type="danger">超级管理员</el-tag>
            <el-tag v-else-if="row.roleId === 2" type="warning">管理员</el-tag>
            <el-tag v-else-if="row.roleId === 3" type="success">普通用户</el-tag>
            <el-tag v-else type="info">访客</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
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
            <el-button type="primary" link @click="handleAssignRole(row)">
              <el-icon><User /></el-icon>
              分配角色
            </el-button>
            <el-button
              type="danger"
              link
              @click="handleDelete(row)"
              :disabled="row.username === 'admin'"
            >
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

    <!-- 用户表单对话框 -->
    <el-dialog
      :title="userForm.id ? '编辑用户' : '新增用户'"
      v-model="dialogVisible"
      width="600px"
    >
      <el-form :model="userForm" :rules="rules" ref="userFormRef" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="userForm.username" placeholder="请输入用户名" :disabled="userForm.id" />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="userForm.nickname" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item v-if="!userForm.id" label="密码" prop="password">
          <el-input v-model="userForm.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="手机号码" prop="phone">
          <el-input v-model="userForm.phone" placeholder="请输入手机号码" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="角色" prop="roleId">
          <el-select v-model="userForm.roleId" placeholder="请选择角色">
            <el-option v-for="role in roleOptions" :key="role.id" :label="role.name" :value="role.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="userForm.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>

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

    <!-- 分配角色对话框 -->
    <el-dialog
      title="分配角色"
      v-model="assignRoleDialogVisible"
      width="500px"
    >
      <el-form label-width="80px">
        <el-form-item label="用户名">
          <el-input v-model="currentUser.username" disabled />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="currentUser.roleId" placeholder="请选择角色">
            <el-option v-for="role in roleOptions" :key="role.id" :label="role.name" :value="role.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="assignRoleDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAssignRole">确定</el-button>
      </template>
    </el-dialog>

    <!-- 导入对话框 -->
    <el-dialog
      title="导入用户"
      v-model="importDialogVisible"
      width="400px"
    >
      <el-upload
        class="upload-demo"
        action="#"
        :http-request="handleFileUpload"
        :show-file-list="false"
        accept=".xlsx, .xls"
      >
        <el-button type="primary">选择文件</el-button>
        <template #tip>
          <div class="el-upload__tip">
            请上传Excel文件，仅支持.xlsx或.xls格式
          </div>
        </template>
      </el-upload>
      <template #footer>
        <el-button @click="importDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="downloadTemplate">下载模板</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { setToken, setAdmin,getToken } from '@/utils/auth'
// 查询参数
const queryParams = reactive({
  username: '',
  phone: '',
  status: '',
  startTime: '',
  endTime: '',
  pageNum: 1,
  pageSize: 10
})

// 日期范围
const dateRange = ref([])

// 用户列表（模拟数据）
const userList = ref([
  {
    id: 1,
    username: getToken(),
    nickname: getToken(),
    avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
    phone: '13800138000',
    email: 'admin@example.com',
    roleId: 1,
    status: 1,
    createTime: '2023-01-01 00:00:00'
  },
  {
    id: 2,
    username: 'manager',
    nickname: '运营主管',
    avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
    phone: '13800138001',
    email: 'manager@example.com',
    roleId: 2,
    status: 1,
    createTime: '2023-01-02 00:00:00'
  },
  {
    id: 3,
    username: 'user001',
    nickname: '张三',
    avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
    phone: '13800138002',
    email: 'user001@example.com',
    roleId: 3,
    status: 1,
    createTime: '2023-01-03 00:00:00'
  },
  {
    id: 4,
    username: 'user002',
    nickname: '李四',
    avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
    phone: '13800138003',
    email: 'user002@example.com',
    roleId: 3,
    status: 1,
    createTime: '2023-01-04 00:00:00'
  },
  {
    id: 5,
    username: 'visitor',
    nickname: '访客',
    avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
    phone: '13800138004',
    email: 'visitor@example.com',
    roleId: 4,
    status: 0,
    createTime: '2023-01-05 00:00:00'
  }
])

// 角色选项（模拟数据）
const roleOptions = [
  { id: 1, name: '超级管理员' },
  { id: 2, name: '管理员' },
  { id: 3, name: '普通用户' },
  { id: 4, name: '访客' }
]

// 选中的用户ID
const selectedIds = ref([])

// 总数
const total = ref(5)

// 加载状态
const loading = ref(false)

// 对话框可见性
const dialogVisible = ref(false)
const resetPwdDialogVisible = ref(false)
const assignRoleDialogVisible = ref(false)
const importDialogVisible = ref(false)

// 用户表单
const userForm = reactive({
  id: '',
  username: '',
  nickname: '',
  password: '',
  phone: '',
  email: '',
  roleId: '',
  status: 1
})

// 密码表单
const pwdForm = reactive({
  userId: '',
  password: '',
  confirmPassword: ''
})

// 当前操作的用户
const currentUser = reactive({
  id: '',
  username: '',
  roleId: ''
})

// 表单引用
const userFormRef = ref(null)
const pwdFormRef = ref(null)

// 表单验证规则
const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6个字符', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  roleId: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
}

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
    // 模拟API请求
    setTimeout(() => {
      loading.value = false
    }, 500)
  } catch (error) {
    console.error('获取用户列表失败', error)
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
    phone: '',
    status: '',
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

// 处理新增
const handleAdd = () => {
  resetForm()
  dialogVisible.value = true
}

// 处理编辑
const handleEdit = (row) => {
  resetForm()
  Object.assign(userForm, { ...row })
  dialogVisible.value = true
}

// 处理删除
const handleDelete = (row) => {
  if (row.username === 'admin') {
    ElMessage.warning('超级管理员账号不能删除')
    return
  }
  
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
  
  // 检查是否包含admin账号
  const hasAdmin = userList.value.some(item => 
    selectedIds.value.includes(item.id) && item.username === 'admin'
  )
  
  if (hasAdmin) {
    ElMessage.warning('选中的用户中包含超级管理员账号，无法删除')
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

// 处理状态变更
const handleStatusChange = (row) => {
  if (row.username === 'admin' && row.status === 0) {
    ElMessage.warning('超级管理员账号不能禁用')
    row.status = 1
    return
  }
  
  const statusText = row.status === 1 ? '启用' : '禁用'
  ElMessage.success(`已${statusText}用户"${row.username}"`)
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

// 处理分配角色
const handleAssignRole = (row) => {
  Object.assign(currentUser, {
    id: row.id,
    username: row.username,
    roleId: row.roleId
  })
  assignRoleDialogVisible.value = true
}

// 提交分配角色
const submitAssignRole = () => {
  // 模拟API请求
  const index = userList.value.findIndex(item => item.id === currentUser.id)
  if (index !== -1) {
    userList.value[index].roleId = currentUser.roleId
    ElMessage.success('角色分配成功')
  }
  assignRoleDialogVisible.value = false
}

// 提交表单
const submitForm = () => {
  userFormRef.value.validate((valid) => {
    if (valid) {
      if (userForm.id) {
        // 编辑用户
        const index = userList.value.findIndex(item => item.id === userForm.id)
        if (index !== -1) {
          userList.value[index] = { ...userList.value[index], ...userForm }
          ElMessage.success('编辑成功')
        }
      } else {
        // 新增用户
        const newUser = {
          ...userForm,
          id: userList.value.length + 1,
          avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
          createTime: new Date().toLocaleString()
        }
        userList.value.push(newUser)
        total.value++
        ElMessage.success('新增成功')
      }
      dialogVisible.value = false
    }
  })
}

// 重置表单
const resetForm = () => {
  Object.assign(userForm, {
    id: '',
    username: '',
    nickname: '',
    password: '',
    phone: '',
    email: '',
    roleId: '',
    status: 1
  })
  if (userFormRef.value) {
    userFormRef.value.resetFields()
  }
}

// 处理导出
const handleExport = () => {
  ElMessage.success('用户数据导出成功')
}

// 处理导入
const handleImport = () => {
  importDialogVisible.value = true
}

// 处理文件上传
const handleFileUpload = (options) => {
  const file = options.file
  // 模拟文件上传
  setTimeout(() => {
    ElMessage.success(`成功导入${Math.floor(Math.random() * 10) + 1}条用户数据`)
    importDialogVisible.value = false
  }, 1000)
}

// 下载模板
const downloadTemplate = () => {
  ElMessage.success('模板下载成功')
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