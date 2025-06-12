<template>
  <div class="app-wrapper">
    <!-- 侧边栏 -->
    <div class="sidebar-container" :class="{ 'is-collapsed': isCollapse }">
      <div class="logo-container">
        <img src="@/assets/logo.svg" alt="Logo" class="logo" />
        <h1 class="title" v-show="!isCollapse">兴趣小队</h1>
      </div>
      
      <el-scrollbar>
        <el-menu
          :default-active="activeMenu"
          :collapse="isCollapse"
          :unique-opened="true"
          background-color="#304156"
          text-color="#bfcbd9"
          active-text-color="#409EFF"
          router
        >
          <sidebar-item 
            v-for="route in routes" 
            :key="route.path" 
            :item="route" 
            :base-path="route.path" 
          />
        </el-menu>
      </el-scrollbar>
    </div>
    
    <!-- 主区域 -->
    <div class="main-container">
      <!-- 顶部导航栏 -->
      <div class="navbar">
        <div class="left-menu">
          <el-icon class="hamburger" @click="toggleSidebar">
            <component :is="isCollapse ? 'Expand' : 'Fold'" />
          </el-icon>
          <breadcrumb />
        </div>
        
        <div class="right-menu">
          <el-dropdown trigger="click">
            <div class="avatar-container">
              <el-avatar :size="30" :src="userAvatar" />
              <span class="username">{{ userInfo.nickname || userInfo.username }}</span>
              <el-icon><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="goToProfile">
                  <el-icon><User /></el-icon>
                  个人中心
                </el-dropdown-item>
                <el-dropdown-item divided @click="logout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
      
      <!-- 内容区域 -->
      <div class="app-main">
        <router-view v-slot="{ Component }">
          <transition name="fade-transform" mode="out-in">
            <keep-alive :include="cachedViews">
              <component :is="Component" />
            </keep-alive>
          </transition>
        </router-view>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import SidebarItem from './components/SidebarItem.vue'
import Breadcrumb from './components/Breadcrumb.vue'

// 路由信息
const router = useRouter()
const route = useRoute()

// 侧边栏折叠状态
const isCollapse = ref(false)

// 用户信息
const userInfo = ref({})
const userAvatar = computed(() => {
  return userInfo.value.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
})

// 路由菜单
const routes = [
  {
    path: '/dashboard',
    component: () => import('@/views/dashboard/index.vue'),
    meta: { title: '首页', icon: 'HomeFilled' }
  },
  {
    path: '/activity',
    meta: { title: '活动管理', icon: 'Calendar' },
    children: [
      {
        path: '/activity/list',
        component: () => import('@/views/activity/list.vue'),
        meta: { title: '活动列表' }
      },
      {
        path: '/activity/category',
        component: () => import('@/views/activity/category.vue'),
        meta: { title: '活动分类' }
      }
    ]
  },
  {
    path: '/user',
    meta: { title: '用户管理', icon: 'User' },
    children: [
      {
        path: '/user/list',
        component: () => import('@/views/user/list.vue'),
        meta: { title: '用户列表' }
      }
    ]
  }
]

// 缓存的视图
const cachedViews = ref([])

// 当前激活的菜单
const activeMenu = computed(() => {
  return route.path
})

// 切换侧边栏折叠状态
const toggleSidebar = () => {
  isCollapse.value = !isCollapse.value
}

// 前往个人中心
const goToProfile = () => {
  router.push('/profile')
}

// 退出登录
const logout = () => {
  ElMessageBox.confirm('确定要退出登录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    localStorage.removeItem('token')
    localStorage.removeItem('adminInfo')
    router.push('/login')
    ElMessage.success('退出登录成功')
  }).catch(() => {})
}

// 获取用户信息
onMounted(() => {
  const adminInfo = localStorage.getItem('adminInfo')
  if (adminInfo) {
    userInfo.value = JSON.parse(adminInfo)
  }
})
</script>

<style lang="scss" scoped>
.app-wrapper {
  position: relative;
  height: 100vh;
  width: 100%;
  display: flex;
  margin: 0;
  padding: 0;
  overflow: hidden;
  
  .sidebar-container {
    width: 210px;
    height: 100%;
    background-color: #304156;
    transition: width 0.3s;
    overflow: hidden;
    flex-shrink: 0;
    
    &.is-collapsed {
      width: 64px;
    }
    
    .logo-container {
      height: 60px;
      display: flex;
      align-items: center;
      padding: 0 15px;
      background-color: #2b3649;
      
      .logo {
        width: 32px;
        height: 32px;
      }
      
      .title {
        margin-left: 10px;
        color: #fff;
        font-size: 18px;
        font-weight: bold;
        white-space: nowrap;
      }
    }
  }
  
  .main-container {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    
    .navbar {
      height: 60px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 0 20px;
      box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
      background-color: #fff;
      
      .left-menu {
        display: flex;
        align-items: center;
        
        .hamburger {
          font-size: 20px;
          cursor: pointer;
          margin-right: 15px;
          
          &:hover {
            color: #409EFF;
          }
        }
      }
      
      .right-menu {
        .avatar-container {
          display: flex;
          align-items: center;
          cursor: pointer;
          
          .username {
            margin: 0 5px;
            font-size: 14px;
          }
        }
      }
    }
    
    .app-main {
      flex: 1;
      padding: 0;
      overflow-y: auto;
      background-color: #fff;
    }
  }
}

// 路由过渡动画
.fade-transform-enter-active,
.fade-transform-leave-active {
  transition: all 0.3s;
}

.fade-transform-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.fade-transform-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
</style> 