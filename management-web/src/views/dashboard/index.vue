<template>
  <div class="dashboard-container">
    <div class="welcome-section">
      <h2>欢迎来到HobbyHub管理系统</h2>
      <p>今天是 {{ currentDate }}，{{ greeting }}</p>
    </div>
    
    <!-- 数据概览卡片 -->
    <el-row :gutter="24">
      <el-col :span="6" v-for="(item, index) in statsCards" :key="index">
        <el-card 
          shadow="never" 
          class="stats-card"
          :class="{ 'stats-card-hover': isCardHovering[index] }"
          @mouseenter="handleCardHover(index, true)"
          @mouseleave="handleCardHover(index, false)"
        >
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
    <el-row :gutter="24" class="chart-row">
      <el-col :span="16">
        <el-card 
          shadow="never" 
          class="chart-card"
          :class="{ 'chart-card-hover': isChartHovering[0] }"
          @mouseenter="handleChartHover(0, true)"
          @mouseleave="handleChartHover(0, false)"
        >
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
        <el-card 
          shadow="never" 
          class="chart-card"
          :class="{ 'chart-card-hover': isChartHovering[1] }"
          @mouseenter="handleChartHover(1, true)"
          @mouseleave="handleChartHover(1, false)"
        >
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
    <el-card 
      shadow="never" 
      class="recent-card"
      :class="{ 'recent-card-hover': isRecentHovering }"
      @mouseenter="handleRecentHover(true)"
      @mouseleave="handleRecentHover(false)"
    >
      <template #header>
        <div class="card-header">
          <span>最新活动</span>
          <el-button type="primary" link @click="viewMore">查看更多</el-button>
        </div>
      </template>
      <el-table 
        :data="recentActivities" 
        style="width: 100%" 
        stripe 
        class="activity-table"
      >
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
import { ref, onMounted, computed, watch, onUnmounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'
import axios from 'axios'

const router = useRouter()
const activityChartRef = ref(null)
const categoryChartRef = ref(null)
const activityChart = ref(null)
const categoryChart = ref(null)
const activityChartType = ref('week')

// 卡片悬停状态
const isCardHovering = ref([false, false, false, false])
const isChartHovering = ref([false, false])
const isRecentHovering = ref(false)

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
  { title: '总活动数', value: 0, icon: 'Calendar', color: '#722ED1' },
  { title: '用户数量', value: 0, icon: 'User', color: '#36CFC9' },
  { title: '小队数量', value: 0, icon: 'UserFilled', color: '#F7BA1E' },
  { title: '本月新增活动', value: 0, icon: 'Plus', color: '#F56C6C' }
])

// 最近活动列表
const recentActivities = ref([])

// 获取活动列表数据
const fetchActivityData = async () => {
  loading.value = true
  try {
    // 模拟数据请求
    await new Promise(resolve => setTimeout(resolve, 800))
    
    // 模拟数据
    const activities = [
      { 
        id: 1, 
        name: '户外徒步探险活动', 
        category: ['户外活动'],
        startTime: '2023-06-15 09:00:00',
        endTime: '2023-06-15 17:00:00',
        createTime: '2023-06-10'
      },
      { 
        id: 2, 
        name: '绘画艺术交流分享会', 
        category: ['文化交流', '艺术'],
        startTime: '2023-06-16 14:00:00',
        endTime: '2023-06-16 16:30:00',
        createTime: '2023-06-11'
      },
      { 
        id: 3, 
        name: '篮球友谊比赛', 
        category: ['体育竞技'],
        startTime: '2023-06-17 16:00:00',
        endTime: '2023-06-17 18:30:00',
        createTime: '2023-06-12'
      },
      { 
        id: 4, 
        name: '编程入门培训课程', 
        category: ['技能培训', '科技'],
        startTime: '2023-06-18 10:00:00',
        endTime: '2023-06-18 12:00:00',
        createTime: '2023-06-13'
      },
      { 
        id: 5, 
        name: '读书分享会', 
        category: ['文化交流'],
        startTime: '2023-06-19 19:00:00',
        endTime: '2023-06-19 21:00:00',
        createTime: '2023-06-14'
      }
    ]
    
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
    recentActivities.value = activities.map(item => {
      return {
        id: item.id,
        title: item.name,
        category: item.category ? item.category[0] : '未分类',
        startTime: formatDate(item.startTime),
        status: getActivityStatus(item.startTime, item.endTime)
      }
    })
    
    // 更新图表数据
    updateChartData(activities)
  } catch (error) {
    console.error('获取活动数据失败:', error)
  } finally {
    loading.value = false
  }
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
    minute: '2-digit'
  })
}

// 获取用户数据
const fetchUserData = async () => {
  try {
    // 模拟数据
    await new Promise(resolve => setTimeout(resolve, 500))
    statsCards.value[1].value = 892
  } catch (error) {
    console.error('获取用户数据失败:', error)
  }
}

// 获取小队数据
const fetchTeamData = async () => {
  try {
    // 模拟数据
    await new Promise(resolve => setTimeout(resolve, 500))
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
        },
        backgroundColor: 'rgba(255, 255, 255, 0.9)',
        borderColor: '#ebedf0',
        borderWidth: 1,
        textStyle: {
          color: '#303133'
        },
        padding: 10
      },
      legend: {
        data: ['新增活动', '活动参与人数'],
        textStyle: {
          color: '#606266'
        }
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true,
        borderColor: '#ebeef5',
        borderWidth: 1
      },
      xAxis: {
        type: 'category',
        data: chartData.xAxis,
        axisLine: {
          lineStyle: {
            color: '#ebeef5'
          }
        },
        axisTick: {
          show: false
        },
        axisLabel: {
          color: '#606266',
          margin: 10
        }
      },
      yAxis: [
        {
          type: 'value',
          name: '活动数',
          position: 'left',
          axisLine: {
            lineStyle: {
              color: '#ebeef5'
            }
          },
          axisTick: {
            show: false
          },
          splitLine: {
            lineStyle: {
              color: '#ebeef5'
            }
          },
          axisLabel: {
            formatter: '{value}',
            color: '#606266'
          }
        },
        {
          type: 'value',
          name: '人数',
          position: 'right',
          axisLine: {
            lineStyle: {
              color: '#ebeef5'
            }
          },
          axisTick: {
            show: false
          },
          splitLine: {
            lineStyle: {
              color: '#ebeef5'
            }
          },
          axisLabel: {
            formatter: '{value}',
            color: '#606266'
          }
        }
      ],
      series: [
        {
          name: '新增活动',
          type: 'bar',
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: '#722ED1' },
              { offset: 1, color: '#5221B1' }
            ]),
            borderRadius: 4
          },
          data: chartData.activities
        },
        {
          name: '活动参与人数',
          type: 'line',
          yAxisIndex: 1,
          itemStyle: {
            color: '#36CFC9'
          },
          lineStyle: {
            color: '#36CFC9'
          },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(54, 207, 201, 0.3)' },
              { offset: 1, color: 'rgba(54, 207, 201, 0.05)' }
            ])
          },
          data: chartData.participants
        }
      ]
    }
    
    activityChart.value.setOption(option)
  }
}

// 根据类型获取活动图表数据
const getActivityChartData = (type) => {
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
        formatter: '{a} <br/>{b}: {c} ({d}%)',
        backgroundColor: 'rgba(255, 255, 255, 0.9)',
        borderColor: '#ebedf0',
        borderWidth: 1,
        textStyle: {
          color: '#303133'
        },
        padding: 10
      },
      legend: {
        orient: 'vertical',
        left: 10,
        data: ['户外活动', '文化交流', '体育竞技', '技能培训', '其他'],
        textStyle: {
          color: '#606266'
        }
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
            position: 'center',
            color: '#303133',
            fontStyle: 'bold'
          },
          emphasis: {
            label: {
              show: true,
              fontSize: '16',
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

// 卡片悬停处理
const handleCardHover = (index, isHovering) => {
  isCardHovering.value[index] = isHovering
}

const handleChartHover = (index, isHovering) => {
  isChartHovering.value[index] = isHovering
}

const handleRecentHover = (isHovering) => {
  isRecentHovering.value = isHovering
}

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
:root {
  --primary-gradient: linear-gradient(135deg, #722ED1 0%, #5221B1 100%);
  --secondary-gradient: linear-gradient(135deg, #36CFC9 0%, #1C949A 100%);
  --accent-gradient: linear-gradient(135deg, #F7BA1E 0%, #E69E0E 100%);
  --danger-gradient: linear-gradient(135deg, #F56C6C 0%, #E94F4F 100%);
  --card-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  --card-shadow-hover: 0 10px 30px rgba(0, 0, 0, 0.12);
}

.dashboard-container {
  width: 100%;
  height: 100%;
  padding: 30px;
  background-color: #f5f7fa;
  display: flex;
  flex-direction: column;
  
  .welcome-section {
    margin-bottom: 30px;
    animation: fadeIn 0.8s ease-out;
    
    h2 {
      font-size: 28px;
      margin-bottom: 10px;
      color: #303133;
      font-weight: 700;
    }
    
    p {
      font-size: 18px;
      color: #606266;
    }
  }
  
  .stats-card {
    display: flex;
    align-items: center;
    padding: 24px;
    border-radius: 16px;
    margin-bottom: 24px;
    background-color: white;
    transition: all 0.3s ease;
    box-shadow: var(--card-shadow);
    
    .stats-icon {
      width: 64px;
      height: 64px;
      border-radius: 12px;
      display: flex;
      justify-content: center;
      align-items: center;
      margin-right: 20px;
      box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
      
      .el-icon {
        font-size: 32px;
        color: white;
      }
    }
    
    .stats-info {
      flex: 1;
      
      .stats-value {
        font-size: 28px;
        font-weight: 700;
        margin-bottom: 6px;
        color: #303133;
      }
      
      .stats-title {
        font-size: 14px;
        color: #606266;
      }
    }
  }
  
  .stats-card-hover {
    transform: translateY(-5px);
    box-shadow: var(--card-shadow-hover);
  }
  
  .chart-row {
    margin-bottom: 30px;
  }
  
  .chart-card {
    background-color: white;
    border-radius: 16px;
    margin-bottom: 24px;
    transition: all 0.3s ease;
    box-shadow: var(--card-shadow);
    
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 0 0 16px 0;
      
      span {
        font-size: 18px;
        font-weight: 600;
        color: #303133;
      }
    }
    
    .chart-container {
      height: 350px;
    }
  }
  
  .chart-card-hover {
    transform: translateY(-5px);
    box-shadow: var(--card-shadow-hover);
  }
  
  .recent-card {
    background-color: white;
    border-radius: 16px;
    transition: all 0.3s ease;
    box-shadow: var(--card-shadow);
    
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 0 0 16px 0;
      
      span {
        font-size: 18px;
        font-weight: 600;
        color: #303133;
      }
    }
  }
  
  .recent-card-hover {
    transform: translateY(-5px);
    box-shadow: var(--card-shadow-hover);
  }
  
  .activity-table {
    th {
      background-color: #f5f7fa !important;
      color: #606266;
      font-weight: 500;
    }
    
    td {
      padding: 16px 0 !important;
      color: #606266;
      font-size: 14px;
    }
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>