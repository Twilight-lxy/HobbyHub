<template>
  <div class="role-container">
    <div class="role-header">
      <h2>角色管理</h2>
      <el-button type="primary" @click="handleAdd">新增角色</el-button>
    </div>
    
    <!-- 角色表格 -->
    <el-card class="table-container" shadow="hover">
      <el-table
        v-loading="loading"
        :data="roleList"
        border
        style="width: 100%"
      >
        <el-table-column prop="id" label="角色ID" width="80" />
        <el-table-column prop="name" label="角色名称" width="150" />
        <el-table-column prop="code" label="角色编码" width="150" />
        <el-table-column prop="description" label="角色描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="createTime" label="创建时间" width="180" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
              {{ scope.row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="scope">
            <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button 
              size="small" 
              type="primary" 
              @click="handlePermission(scope.row)"
            >权限设置</el-button>
            <el-button 
              size="small" 
              type="danger" 
              @click="handleDelete(scope.row)"
              :disabled="scope.row.code === 'admin'"
            >删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 角色表单对话框 -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="500px"
    >
      <el-form 
        ref="formRef" 
        :model="form" 
        :rules="rules" 
        label-width="80px"
      >
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="角色编码" prop="code">
          <el-input v-model="form.code" placeholder="请输入角色编码" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="角色描述" prop="description">
          <el-input 
            v-model="form.description" 
            type="textarea" 
            placeholder="请输入角色描述" 
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm">确定</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 权限设置对话框 -->
    <el-dialog
      title="权限设置"
      v-model="permissionDialogVisible"
      width="600px"
    >
      <div v-if="currentRole.id">
        <div class="permission-header">
          <span>角色：{{ currentRole.name }}</span>
        </div>
        <el-tree
          ref="permissionTreeRef"
          :data="permissionTree"
          show-checkbox
          node-key="id"
          :default-checked-keys="checkedPermissions"
          :props="{ label: 'name', children: 'children' }"
        />
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="permissionDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="savePermissions">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'

const loading = ref(false)
const roleList = ref([])
const dialogVisible = ref(false)
const permissionDialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)
const permissionTreeRef = ref(null)
const currentRole = ref({})
const checkedPermissions = ref([])

// 表单数据
const form = reactive({
  id: null,
  name: '',
  code: '',
  description: '',
  status: 1
})

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入角色名称', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入角色编码', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '角色编码只能包含字母、数字和下划线', trigger: 'blur' }
  ]
}

// 对话框标题
const dialogTitle = computed(() => {
  return isEdit.value ? '编辑角色' : '新增角色'
})

// 权限树数据
const permissionTree = ref([
  {
    id: 1,
    name: '系统管理',
    children: [
      { id: 11, name: '用户管理' },
      { id: 12, name: '角色管理' },
      { id: 13, name: '菜单管理' },
      { id: 14, name: '系统日志' }
    ]
  },
  {
    id: 2,
    name: '活动管理',
    children: [
      { id: 21, name: '活动列表' },
      { id: 22, name: '活动分类' },
      { id: 23, name: '活动审核' }
    ]
  },
  {
    id: 3,
    name: '小队管理',
    children: [
      { id: 31, name: '小队列表' },
      { id: 32, name: '小队审核' }
    ]
  }
])

// 获取角色列表
const getRoleList = async () => {
  loading.value = true
  try {
    // 这里替换为实际的API调用
    // const res = await axios.get('/api/system/role/list')
    
    // 模拟数据
    const res = {
      data: [
        {
          id: 1,
          name: '超级管理员',
          code: 'admin',
          description: '系统超级管理员，拥有所有权限',
          createTime: '2023-01-01 00:00:00',
          status: 1
        },
        {
          id: 2,
          name: '普通管理员',
          code: 'manager',
          description: '普通管理员，拥有部分管理权限',
          createTime: '2023-01-02 00:00:00',
          status: 1
        },
        {
          id: 3,
          name: '活动管理员',
          code: 'activity_manager',
          description: '活动管理员，负责活动相关管理',
          createTime: '2023-01-03 00:00:00',
          status: 1
        },
        {
          id: 4,
          name: '小队管理员',
          code: 'team_manager',
          description: '小队管理员，负责小队相关管理',
          createTime: '2023-01-04 00:00:00',
          status: 1
        },
        {
          id: 5,
          name: '访客',
          code: 'visitor',
          description: '访客角色，仅有查看权限',
          createTime: '2023-01-05 00:00:00',
          status: 0
        }
      ]
    }
    
    roleList.value = res.data
  } catch (error) {
    console.error('获取角色列表失败', error)
    ElMessage.error('获取角色列表失败')
  } finally {
    loading.value = false
  }
}

// 新增角色
const handleAdd = () => {
  isEdit.value = false
  resetForm()
  dialogVisible.value = true
}

// 编辑角色
const handleEdit = (row) => {
  isEdit.value = true
  resetForm()
  Object.assign(form, row)
  dialogVisible.value = true
}

// 删除角色
const handleDelete = (row) => {
  if (row.code === 'admin') {
    ElMessage.warning('超级管理员角色不能删除')
    return
  }
  
  ElMessageBox.confirm(`确定要删除角色"${row.name}"吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      // 这里替换为实际的API调用
      // await axios.delete(`/api/system/role/${row.id}`)
      
      // 模拟删除
      const index = roleList.value.findIndex(item => item.id === row.id)
      if (index !== -1) {
        roleList.value.splice(index, 1)
      }
      
      ElMessage.success('删除成功')
    } catch (error) {
      console.error('删除失败', error)
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

// 权限设置
const handlePermission = (row) => {
  currentRole.value = row
  
  // 模拟获取角色权限
  if (row.id === 1) { // 超级管理员
    checkedPermissions.value = [11, 12, 13, 14, 21, 22, 23, 31, 32]
  } else if (row.id === 2) { // 普通管理员
    checkedPermissions.value = [11, 14, 21, 22, 31]
  } else if (row.id === 3) { // 活动管理员
    checkedPermissions.value = [21, 22, 23]
  } else if (row.id === 4) { // 小队管理员
    checkedPermissions.value = [31, 32]
  } else {
    checkedPermissions.value = []
  }
  
  permissionDialogVisible.value = true
}

// 保存权限设置
const savePermissions = async () => {
  try {
    const checkedKeys = permissionTreeRef.value.getCheckedKeys()
    const halfCheckedKeys = permissionTreeRef.value.getHalfCheckedKeys()
    const allCheckedKeys = [...checkedKeys, ...halfCheckedKeys]
    
    // 这里替换为实际的API调用
    // await axios.post(`/api/system/role/permission`, {
    //   roleId: currentRole.value.id,
    //   permissionIds: allCheckedKeys
    // })
    
    ElMessage.success('权限设置成功')
    permissionDialogVisible.value = false
  } catch (error) {
    console.error('权限设置失败', error)
    ElMessage.error('权限设置失败')
  }
}

// 提交表单
const submitForm = () => {
  formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        // 这里替换为实际的API调用
        // if (isEdit.value) {
        //   await axios.put(`/api/system/role/${form.id}`, form)
        // } else {
        //   await axios.post('/api/system/role', form)
        // }
        
        // 模拟提交
        if (isEdit.value) {
          const index = roleList.value.findIndex(item => item.id === form.id)
          if (index !== -1) {
            roleList.value[index] = { ...form }
          }
        } else {
          const newId = Math.max(...roleList.value.map(item => item.id)) + 1
          roleList.value.push({
            ...form,
            id: newId,
            createTime: new Date().toLocaleString()
          })
        }
        
        ElMessage.success(isEdit.value ? '编辑成功' : '新增成功')
        dialogVisible.value = false
        getRoleList()
      } catch (error) {
        console.error(isEdit.value ? '编辑失败' : '新增失败', error)
        ElMessage.error(isEdit.value ? '编辑失败' : '新增失败')
      }
    }
  })
}

// 重置表单
const resetForm = () => {
  form.id = null
  form.name = ''
  form.code = ''
  form.description = ''
  form.status = 1
  
  if (formRef.value) {
    formRef.value.resetFields()
  }
}

// 初始化
onMounted(() => {
  getRoleList()
})
</script>

<style lang="scss" scoped>
.role-container {
  .role-header {
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
  
  .table-container {
    margin-bottom: 20px;
  }
  
  .permission-header {
    margin-bottom: 20px;
    font-size: 16px;
    font-weight: 500;
  }
}
</style>