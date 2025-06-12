<template>
  <div class="log-container">
    <div class="log-header">
      <h2>系统日志</h2>
    </div>
    
    <!-- 搜索和筛选 -->
    <el-card class="filter-container" shadow="hover">
      <el-form :model="queryParams" ref="queryForm" :inline="true">
        <el-form-item label="操作人" prop="username">
          <el-input
            v-model="queryParams.username"
            placeholder="请输入操作人"
            clearable
            @keyup.enter="handleQuery"
          />
        </el-form-item>
        <el-form-item label="操作类型" prop="operationType">
          <el-select
            v-model="queryParams.operationType"
            placeholder="请选择操作类型"
            clearable
          >
            <el-option label="登录" value="LOGIN" />
            <el-option label="登出" value="LOGOUT" />
            <el-option label="新增" value="INSERT" />
            <el-option label="修改" value="UPDATE" />
            <el-option label="删除" value="DELETE" />
            <el-option label="导出" value="EXPORT" />
            <el-option label="导入" value="IMPORT" />
            <el-option label="其他" value="OTHER" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select
            v-model="queryParams.status"
            placeholder="请选择状态"
            clearable
          >
            <el-option label="成功" :value="0" />
            <el-option label="失败" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item label="操作时间">
          <el-date-picker
            v-model="dateRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
          <el-button type="danger" @click="handleClear">清空日志</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <!-- 日志表格 -->
    <el-card class="table-container" shadow="hover">
      <el-table
        v-loading="loading"
        :data="logList"
        border
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column prop="id" label="日志编号" width="80" />
        <el-table-column prop="title" label="操作模块" width="150" />
        <el-table-column prop="operationType" label="操作类型" width="100">
          <template #default="scope">
            <el-tag :type="getOperationTypeTag(scope.row.operationType)">
              {{ getOperationTypeText(scope.row.operationType) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="method" label="请求方法" width="180" show-overflow-tooltip />
        <el-table-column prop="requestUrl" label="请求URL" min-width="200" show-overflow-tooltip />
        <el-table-column prop="username" label="操作人" width="120" />
        <el-table-column prop="ip" label="IP地址" width="130" />
        <el-table-column prop="operationTime" label="操作时间" width="180" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="scope">
            <el-tag :type="scope.row.status === 0 ? 'success' : 'danger'">
              {{ scope.row.status === 0 ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="scope">
            <el-button size="small" @click="handleView(scope.row)">详情</el-button>
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
    
    <!-- 日志详情对话框 -->
    <el-dialog
      title="日志详情"
      v-model="dialogVisible"
      width="700px"
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item label="日志编号">{{ currentLog.id }}</el-descriptions-item>
        <el-descriptions-item label="操作模块">{{ currentLog.title }}</el-descriptions-item>
        <el-descriptions-item label="操作类型">
          <el-tag :type="getOperationTypeTag(currentLog.operationType)">
            {{ getOperationTypeText(currentLog.operationType) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="操作状态">
          <el-tag :type="currentLog.status === 0 ? 'success' : 'danger'">
            {{ currentLog.status === 0 ? '成功' : '失败' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="操作人">{{ currentLog.username }}</el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ currentLog.ip }}</el-descriptions-item>
        <el-descriptions-item label="请求方法">{{ currentLog.method }}</el-descriptions-item>
        <el-descriptions-item label="操作时间">{{ currentLog.operationTime }}</el-descriptions-item>
        <el-descriptions-item label="请求URL" :span="2">{{ currentLog.requestUrl }}</el-descriptions-item>
        <el-descriptions-item label="请求参数" :span="2">
          <div class="code-block">{{ formatJson(currentLog.requestParams) }}</div>
        </el-descriptions-item>
        <el-descriptions-item v-if="currentLog.status === 1" label="异常信息" :span="2">
          <div class="error-block">{{ currentLog.errorMsg }}</div>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'

const loading = ref(false)
const logList = ref([])
const total = ref(0)
const dateRange = ref([])
const dialogVisible = ref(false)
const currentLog = ref({})

// 查询参数
const queryParams = reactive({
  pageNum: 1,
  pageSize: 10,
  username: '',
  operationType: '',
  status: undefined,
  startTime: undefined,
  endTime: undefined
})

// 获取操作类型标签
const getOperationTypeTag = (type) => {
  const map = {
    'LOGIN': 'success',
    'LOGOUT': 'info',
    'INSERT': 'primary',
    'UPDATE': 'warning',
    'DELETE': 'danger',
    'EXPORT': '',
    'IMPORT': '',
    'OTHER': 'info'
  }
  return map[type] || 'info'
}

// 获取操作类型文本
const getOperationTypeText = (type) => {
  const map = {
    'LOGIN': '登录',
    'LOGOUT': '登出',
    'INSERT': '新增',
    'UPDATE': '修改',
    'DELETE': '删除',
    'EXPORT': '导出',
    'IMPORT': '导入',
    'OTHER': '其他'
  }
  return map[type] || '未知'
}

// 格式化JSON
const formatJson = (jsonString) => {
  if (!jsonString) return ''
  try {
    const obj = typeof jsonString === 'string' ? JSON.parse(jsonString) : jsonString
    return JSON.stringify(obj, null, 2)
  } catch (e) {
    return jsonString
  }
}

// 查询日志列表
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
    // 这里替换为实际的API调用
    // const res = await axios.get('/api/system/log/list', {
    //   params: queryParams
    // })
    
    // 模拟数据
    const res = {
      data: {
        total: 100,
        list: Array.from({ length: 10 }, (_, index) => {
          const operationTypes = ['LOGIN', 'LOGOUT', 'INSERT', 'UPDATE', 'DELETE', 'EXPORT', 'IMPORT', 'OTHER']
          const methods = ['GET', 'POST', 'PUT', 'DELETE']
          const modules = ['用户管理', '角色管理', '活动管理', '小队管理', '系统设置']
          const users = ['admin', 'manager', 'user1', 'user2']
          
          return {
            id: 1000 + index,
            title: modules[Math.floor(Math.random() * modules.length)],
            operationType: operationTypes[Math.floor(Math.random() * operationTypes.length)],
            method: methods[Math.floor(Math.random() * methods.length)],
            requestUrl: `/api/${['user', 'role', 'activity', 'team', 'system'][Math.floor(Math.random() * 5)]}/${['list', 'save', 'update', 'delete'][Math.floor(Math.random() * 4)]}`,
            username: users[Math.floor(Math.random() * users.length)],
            ip: `192.168.1.${Math.floor(Math.random() * 255)}`,
            operationTime: new Date(Date.now() - Math.floor(Math.random() * 7 * 24 * 60 * 60 * 1000)).toLocaleString(),
            status: Math.random() > 0.9 ? 1 : 0,
            requestParams: JSON.stringify({ id: Math.floor(Math.random() * 100), name: '测试数据' }),
            errorMsg: '系统内部错误，请联系管理员'
          }
        })
      }
    }
    
    logList.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    console.error('获取日志列表失败', error)
    ElMessage.error('获取日志列表失败')
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
    operationType: '',
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

// 查看日志详情
const handleView = (row) => {
  currentLog.value = row
  dialogVisible.value = true
}

// 清空日志
const handleClear = () => {
  ElMessageBox.confirm('确定要清空所有系统日志吗？此操作不可恢复！', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      // 这里替换为实际的API调用
      // await axios.delete('/api/system/log/clear')
      
      // 模拟清空
      logList.value = []
      total.value = 0
      
      ElMessage.success('日志清空成功')
    } catch (error) {
      console.error('日志清空失败', error)
      ElMessage.error('日志清空失败')
    }
  }).catch(() => {})
}

// 初始化
onMounted(() => {
  getList()
})
</script>

<style lang="scss" scoped>
.log-container {
  .log-header {
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
  
  .code-block {
    background-color: #f5f7fa;
    border: 1px solid #e4e7ed;
    border-radius: 4px;
    padding: 10px;
    font-family: monospace;
    white-space: pre-wrap;
    max-height: 200px;
    overflow-y: auto;
  }
  
  .error-block {
    background-color: #fef0f0;
    border: 1px solid #fde2e2;
    border-radius: 4px;
    padding: 10px;
    color: #f56c6c;
    max-height: 200px;
    overflow-y: auto;
  }
}
</style> 