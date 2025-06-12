<template>
  <div class="dashboard-container">
    <div class="welcome-section">
      <h2>欢迎使用兴趣小队管理系统</h2>
      <p>今天是 {{ currentDate }}，{{ greeting }}</p>
    </div>
    
    <!-- 数据概览卡片 -->
    <el-row :gutter="20">
      <el-col :span="6" v-for="(item, index) in statsCards" :key="index">
        <el-card shadow="hover" class="stats-card">
          <div class="stats-icon" :style="{ backgroundColor: item.color }">
            <el-icon>
              <component :is="item.icon"></component>
            </el-icon>
          </div>
          <div class="stats-info">
            <div class="stats-value">{{ item.value }}</div>
            <div class="stats-title">{{ item.title }}</div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 图表区域 -->
    <el-row :gutter="20" class="chart-row">
      <el-col :span="16">
        <el-card shadow="hover" class="chart-card">
          <template #header>
            <div class="card-header">
              <span>活动趋势</span>
              <el-radio-group v-model="activityChartType" size="small">
                <el-radio-button label="week">本周</el-radio-button>
                <el-radio-button label="month">本月</el-radio-button>
                <el-radio-button label="year">全年</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div class="chart-container" ref="activityChartRef"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover" class="chart-card">
          <template #header>
            <div class="card-header">
              <span>活动分类分布</span>
            </div>
          </template>
          <div class="chart-container" ref="categoryChartRef"></div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 最新活动 -->
    <el-card shadow="hover" class="recent-card">
      <template #header>
        <div class="card-header">
          <span>最新活动</span>
          <el-button type="primary" link @click="viewMore">查看更多</el-button>
        </div>
      </template>
      <el-table :data="recentActivities" style="width: 100%" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="活动名称" />
        <el-table-column prop="category" label="分类" width="120" />
        <el-table-column prop="startTime" label="开始时间" width="180" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button link type="primary" size="small" @click="viewDetail(scope.row)">
              查看
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'
import axios from 'axios'

const router = useRouter()
const activityChartRef = ref(null)
const categoryChartRef = ref(null)
const activityChart = ref(null)
const categoryChart = ref(null)
const activityChartType = ref('week')

// 数据加载状态
const loading = ref(false)

// 获取当前日期和问候语
const currentDate = computed(() => {
  const now = new Date()
  return now.toLocaleDateString('zh-CN', { 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric', 
    weekday: 'long' 
  })
})

const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 6) return '凌晨好'
  if (hour < 9) return '早上好'
  if (hour < 12) return '上午好'
  if (hour < 14) return '中午好'
  if (hour < 17) return '下午好'
  if (hour < 19) return '傍晚好'
  return '晚上好'
})

// 数据统计卡片
const statsCards = ref([
  { title: '总活动数', value: 0, icon: 'Calendar', color: '#409EFF' },
  { title: '用户数量', value: 0, icon: 'User', color: '#67C23A' },
  { title: '小队数量', value: 0, icon: 'UserFilled', color: '#E6A23C' },
  { title: '本月新增活动', value: 0, icon: 'Plus', color: '#F56C6C' }
])

// 最近活动列表
const recentActivities = ref([])

// 获取活动列表数据
const fetchActivityData = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/product/list')
    if (res.data && res.data.code === 200) {
      const activities = res.data.data || []
      
      // 更新统计卡片数据
      statsCards.value[0].value = activities.length
      
      // 获取本月新增活动数量
      const currentMonth = new Date().getMonth()
      const thisMonthActivities = activities.filter(item => {
        const activityMonth = new Date(item.createTime).getMonth()
        return activityMonth === currentMonth
      })
      statsCards.value[3].value = thisMonthActivities.length
      
      // 更新最近活动列表
      recentActivities.value = activities.slice(0, 5).map(item => {
        return {
          id: item.id,
          title: item.name,
          category: item.category ? item.category[0] : '未分类',
          startTime: item.startTime || new Date().toLocaleString(),
          status: getActivityStatus(item.startTime, item.endTime)
        }
      })
      
      // 更新图表数据
      updateChartData(activities)
    }
  } catch (error) {
    console.error('获取活动数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取用户数据
const fetchUserData = async () => {
  try {
    // 这里应该调用用户列表API，但由于没有找到合适的API，暂时使用模拟数据
    // const res = await axios.get('/api/user/list')
    // if (res.data && res.data.code === 200) {
    //   statsCards.value[1].value = res.data.data.length
    // }
    
    // 模拟数据
    statsCards.value[1].value = 892
  } catch (error) {
    console.error('获取用户数据失败:', error)
  }
}

// 获取小队数据
const fetchTeamData = async () => {
  try {
    // 这里应该调用小队列表API，但由于没有找到合适的API，暂时使用模拟数据
    // const res = await axios.get('/api/team/list')
    // if (res.data && res.data.code === 200) {
    //   statsCards.value[2].value = res.data.data.length
    // }
    
    // 模拟数据
    statsCards.value[2].value = 32
  } catch (error) {
    console.error('获取小队数据失败:', error)
  }
}

// 根据开始和结束时间判断活动状态
const getActivityStatus = (startTime, endTime) => {
  const now = new Date().getTime()
  const start = startTime ? new Date(startTime).getTime() : now - 86400000 // 默认昨天开始
  const end = endTime ? new Date(endTime).getTime() : now + 86400000 // 默认明天结束
  
  if (now < start) return 0 // 未开始
  if (now > end) return 2 // 已结束
  return 1 // 进行中
}

// 状态处理
const getStatusType = (status) => {
  const types = ['info', 'success', 'warning', 'danger']
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = ['未开始', '进行中', '已结束', '已取消']
  return texts[status] || '未知'
}

// 更新图表数据
const updateChartData = (activities) => {
  // 活动分类分布数据
  const categoryMap = {}
  activities.forEach(item => {
    if (item.category && item.category.length > 0) {
      item.category.forEach(cat => {
        if (categoryMap[cat]) {
          categoryMap[cat]++
        } else {
          categoryMap[cat] = 1
        }
      })
    }
  })
  
  const categoryData = Object.keys(categoryMap).map(key => {
    return { name: key, value: categoryMap[key] }
  })
  
  // 更新饼图数据
  if (categoryChart.value) {
    const option = categoryChart.value.getOption()
    option.series[0].data = categoryData.length > 0 ? categoryData : [
      { value: 20, name: '户外活动' },
      { value: 15, name: '文化交流' },
      { value: 12, name: '体育竞技' },
      { value: 8, name: '技能培训' },
      { value: 5, name: '其他' }
    ]
    option.legend.data = categoryData.map(item => item.name)
    categoryChart.value.setOption(option)
  }
  
  // 更新活动趋势图数据
  initActivityChart()
}

// 查看更多
const viewMore = () => {
  router.push('/activity/list')
}

// 查看详情
const viewDetail = (row) => {
  router.push(`/activity/detail/${row.id}`)
}

// 初始化活动趋势图表
const initActivityChart = () => {
  if (activityChartRef.value) {
    activityChart.value = echarts.init(activityChartRef.value)
    
    // 根据类型获取不同的数据
    const chartData = getActivityChartData(activityChartType.value)
    
    const option = {
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        }
      },
      legend: {
        data: ['新增活动', '活动参与人数']
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: chartData.xAxis
      },
      yAxis: [
        {
          type: 'value',
          name: '活动数',
          position: 'left',
          axisLabel: {
            formatter: '{value}'
          }
        },
        {
          type: 'value',
          name: '人数',
          position: 'right',
          axisLabel: {
            formatter: '{value}'
          }
        }
      ],
      series: [
        {
          name: '新增活动',
          type: 'bar',
          data: chartData.activities
        },
        {
          name: '活动参与人数',
          type: 'line',
          yAxisIndex: 1,
          data: chartData.participants
        }
      ]
    }
    
    activityChart.value.setOption(option)
  }
}

// 根据类型获取活动图表数据
const getActivityChartData = (type) => {
  // 这里应该从API获取数据，这里使用模拟数据
  if (type === 'week') {
    return {
      xAxis: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
      activities: [2, 4, 3, 5, 6, 8, 4],
      participants: [20, 42, 35, 50, 65, 80, 45]
    }
  } else if (type === 'month') {
    return {
      xAxis: ['第1周', '第2周', '第3周', '第4周'],
      activities: [10, 15, 12, 18],
      participants: [120, 150, 130, 180]
    }
  } else {
    return {
      xAxis: ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月'],
      activities: [8, 6, 10, 12, 15, 20, 18, 22, 16, 14, 10, 12],
      participants: [80, 60, 100, 120, 150, 200, 180, 220, 160, 140, 100, 120]
    }
  }
}

// 初始化分类分布图表
const initCategoryChart = () => {
  if (categoryChartRef.value) {
    categoryChart.value = echarts.init(categoryChartRef.value)
    
    const option = {
      tooltip: {
        trigger: 'item',
        formatter: '{a} <br/>{b}: {c} ({d}%)'
      },
      legend: {
        orient: 'vertical',
        left: 10,
        data: ['户外活动', '文化交流', '体育竞技', '技能培训', '其他']
      },
      series: [
        {
          name: '活动分类',
          type: 'pie',
          radius: ['50%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: {
            show: false,
            position: 'center'
          },
          emphasis: {
            label: {
              show: true,
              fontSize: 16,
              fontWeight: 'bold'
            }
          },
          labelLine: {
            show: false
          },
          data: [
            { value: 20, name: '户外活动' },
            { value: 15, name: '文化交流' },
            { value: 12, name: '体育竞技' },
            { value: 8, name: '技能培训' },
            { value: 5, name: '其他' }
          ]
        }
      ]
    }
    
    categoryChart.value.setOption(option)
  }
}

// 监听图表类型变化
watch(activityChartType, () => {
  initActivityChart()
})

// 监听窗口大小变化，重绘图表
const handleResize = () => {
  activityChart.value && activityChart.value.resize()
  categoryChart.value && categoryChart.value.resize()
}

onMounted(() => {
  // 获取数据
  fetchActivityData()
  fetchUserData()
  fetchTeamData()
  
  // 初始化图表
  initActivityChart()
  initCategoryChart()
  
  // 监听窗口大小变化
  window.addEventListener('resize', handleResize)
})

// 组件卸载时移除监听
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  activityChart.value && activityChart.value.dispose()
  categoryChart.value && categoryChart.value.dispose()
})
</script>

<style lang="scss" scoped>
.dashboard-container {
  width: 100%;
  height: 100%;
  padding: 20px;
  
  .welcome-section {
    margin-bottom: 20px;
    
    h2 {
      font-size: 24px;
      margin-bottom: 10px;
    }
    
    p {
      font-size: 16px;
      color: #606266;
    }
  }
  
  .stats-card {
    display: flex;
    align-items: center;
    padding: 20px;
    margin-bottom: 20px;
    
    .stats-icon {
      width: 60px;
      height: 60px;
      border-radius: 10px;
      display: flex;
      justify-content: center;
      align-items: center;
      margin-right: 15px;
      
      .el-icon {
        font-size: 30px;
        color: white;
      }
    }
    
    .stats-info {
      flex: 1;
      
      .stats-value {
        font-size: 24px;
        font-weight: bold;
        margin-bottom: 5px;
      }
      
      .stats-title {
        font-size: 14px;
        color: #606266;
      }
    }
  }
  
  .chart-row {
    margin-bottom: 20px;
  }
  
  .chart-card {
    margin-bottom: 20px;
    
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      span {
        font-size: 16px;
        font-weight: bold;
      }
    }
    
    .chart-container {
      height: 350px;
    }
  }
  
  .recent-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      span {
        font-size: 16px;
        font-weight: bold;
      }
    }
  }
}
</style> 