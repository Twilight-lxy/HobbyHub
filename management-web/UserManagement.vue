9++<template>
  <div class="user-management">
    <el-card>
      <div class="header">
        <h2>用户管理</h2>
        <el-button type="primary" @click="showAddDialog">添加用户</el-button>
      </div>
      
      <el-table :data="userList" border style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="phone" label="手机号" />
        <el-table-column prop="status" label="状态">
          <template #default="{row}">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180">
          <template #default="{row}">
            <el-button size="small" @click="editUser(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="deleteUser(row.id)">删除</el-button>
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
    
    <!-- 添加/编辑用户对话框 -->
    <user-dialog 
      v-model="dialogVisible"
      :form-data="currentUser"
      @submit="handleSubmit" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getUserList, deleteUser } from '@/api/user'
import UserDialog from './components/UserDialog.vue'

const userList = ref([])
const pagination = ref({
  current: 1,
  size: 10,
  total: 0
})
const dialogVisible = ref(false)
const currentUser = ref({})

const fetchUserList = async () => {
  const params = {
    page: pagination.value.current,
    size: pagination.value.size
  }
  const res = await getUserList(params)
  userList.value = res.data.list
  pagination.value.total = res.data.total
}

// 初始化加载数据
onMounted(() => {
  fetchUserList()
})

// 分页处理
const handleSizeChange = (val) => {
  pagination.value.size = val
  fetchUserList()
}

const handleCurrentChange = (val) => {
  pagination.value.current = val
  fetchUserList()
}

// 用户操作
const showAddDialog = () => {
  currentUser.value = {}
  dialogVisible.value = true
}

const editUser = (user) => {
  currentUser.value = { ...user }
  dialogVisible.value = true
}

const handleSubmit = async (formData) => {
  // 调用API保存用户
  await saveUser(formData)
  dialogVisible.value = false
  fetchUserList()
}

const deleteUser = async (id) => {
  await ElMessageBox.confirm('确定删除该用户吗？', '提示', {
    type: 'warning'
  })
  await deleteUser(id)
  ElMessage.success('删除成功')
  fetchUserList()
}
</script>

<style lang="scss" scoped>
.user-management {
  padding: 20px;
  
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }
}
</style>