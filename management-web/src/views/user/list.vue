<template>
  <div class="container mx-auto px-4 py-8 max-w-7xl">
    <!-- 页面标题 -->
    <header class="mb-10 text-center">
      <h1 class="text-[clamp(2rem,5vw,3rem)] font-bold text-neutral-800 mb-3 bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
        用户管理系统
      </h1>
      <p class="text-neutral-600 max-w-2xl mx-auto text-lg">
        查看和管理系统用户信息，包括用户名、性别、地址和创建时间等数据
      </p>
    </header>

    <!-- 数据概览卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-8 mb-12">
      <div class="bg-white rounded-2xl shadow-xl p-8 transition-all duration-300 hover:shadow-2xl hover:-translate-y-1">
        <div class="flex items-start justify-between">
          <div>
            <p class="text-neutral-500 text-sm mb-2">总用户数</p>
            <h3 class="text-4xl font-bold text-neutral-800 mt-1">{{ users.length }}</h3>
          </div>
          <div class="w-14 h-14 rounded-full bg-primary/10 flex items-center justify-center shadow-lg">
            <i class="fa fa-users text-primary text-2xl"></i>
          </div>
        </div>
        <div class="mt-6 pt-6 border-t border-neutral-100">
          <p class="text-neutral-500 text-sm flex items-center">
            <i class="fa fa-clock-o mr-2 text-primary"></i> 最近更新: {{ lastUpdateTime }}
          </p>
        </div>
      </div>
      
      <div class="bg-white rounded-2xl shadow-xl p-8 transition-all duration-300 hover:shadow-2xl hover:-translate-y-1">
        <div class="flex items-start justify-between">
          <div>
            <p class="text-neutral-500 text-sm mb-2">男性用户</p>
            <h3 class="text-4xl font-bold text-neutral-800 mt-1">{{ maleCount }}</h3>
          </div>
          <div class="w-14 h-14 rounded-full bg-blue-100 flex items-center justify-center shadow-lg">
            <i class="fa fa-male text-blue-500 text-2xl"></i>
          </div>
        </div>
        <div class="mt-6 pt-6 border-t border-neutral-100">
          <div class="w-full bg-neutral-200 rounded-full h-3 mt-2">
            <div class="bg-blue-500 h-3 rounded-full" :style="{ width: malePercentage + '%' }" 
                 class="transition-all duration-1000 ease-out"></div>
          </div>
          <p class="text-neutral-500 text-sm mt-2 flex items-center">
            <span class="text-blue-500 font-medium mr-1">{{ malePercentage }}%</span> 占比
          </p>
        </div>
      </div>
      
      <div class="bg-white rounded-2xl shadow-xl p-8 transition-all duration-300 hover:shadow-2xl hover:-translate-y-1">
        <div class="flex items-start justify-between">
          <div>
            <p class="text-neutral-500 text-sm mb-2">女性用户</p>
            <h3 class="text-4xl font-bold text-neutral-800 mt-1">{{ femaleCount }}</h3>
          </div>
          <div class="w-14 h-14 rounded-full bg-pink-100 flex items-center justify-center shadow-lg">
            <i class="fa fa-female text-pink-500 text-2xl"></i>
          </div>
        </div>
        <div class="mt-6 pt-6 border-t border-neutral-100">
          <div class="w-full bg-neutral-200 rounded-full h-3 mt-2">
            <div class="bg-pink-500 h-3 rounded-full" :style="{ width: femalePercentage + '%' }" 
                 class="transition-all duration-1000 ease-out"></div>
          </div>
          <p class="text-neutral-500 text-sm mt-2 flex items-center">
            <span class="text-pink-500 font-medium mr-1">{{ femalePercentage }}%</span> 占比
          </p>
        </div>
      </div>
    </div>

    <!-- 搜索和筛选区域 -->
    <div class="bg-white rounded-2xl shadow-xl p-8 mb-12 relative overflow-hidden">
      <div class="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-primary to-secondary"></div>
      <div class="flex flex-col md:flex-row gap-6">
        <div class="relative flex-1">
          <i class="fa fa-search absolute left-4 top-1/2 transform -translate-y-1/2 text-neutral-400"></i>
          <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="搜索用户名、姓名或地址..." 
            class="w-full pl-12 pr-6 py-4 rounded-xl border border-neutral-200 focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary transition-all duration-300 shadow-sm"
          >
        </div>
        <div class="flex flex-col sm:flex-row gap-6">
          <select 
            v-model="filterByGender" 
            class="w-full sm:w-60 px-6 py-4 rounded-xl border border-neutral-200 focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary transition-all duration-300 appearance-none bg-white shadow-sm"
          >
            <option value="">全部性别</option>
            <option value="男">男</option>
            <option value="女">女</option>
          </select>
          <button 
            @click="refreshData" 
            class="w-full sm:w-auto px-8 py-4 bg-primary text-white rounded-xl hover:bg-primary/90 transition-all duration-300 flex items-center justify-center shadow-md hover:shadow-lg"
          >
            <i class="fa fa-refresh mr-3"></i> 刷新数据
          </button>
        </div>
      </div>
    </div>

    <!-- 用户列表 -->
    <div class="bg-white rounded-2xl shadow-xl overflow-hidden mb-12">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="bg-neutral-50 text-left">
              <th class="px-8 py-5 text-sm font-semibold text-neutral-600 border-b">ID</th>
              <th class="px-8 py-5 text-sm font-semibold text-neutral-600 border-b">用户名</th>
              <th class="px-8 py-5 text-sm font-semibold text-neutral-600 border-b">姓名</th>
              <th class="px-8 py-5 text-sm font-semibold text-neutral-600 border-b">性别</th>
              <th class="px-8 py-5 text-sm font-semibold text-neutral-600 border-b">地址</th>
              <th class="px-8 py-5 text-sm font-semibold text-neutral-600 border-b">创建时间</th>
              <th class="px-8 py-5 text-sm font-semibold text-neutral-600 border-b text-right">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr 
              v-for="user in filteredUsers" 
              :key="user.id" 
              class="border-t border-neutral-100 hover:bg-neutral-50 transition-all duration-200 hover:shadow-inner"
            >
              <td class="px-8 py-5 text-sm text-neutral-700 font-medium">{{ user.id }}</td>
              <td class="px-8 py-5">
                <div class="flex items-center">
                  <div class="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center text-primary font-medium shadow-sm">
                    {{ user.username.charAt(0) }}
                  </div>
                  <div class="ml-4">
                    <p class="text-sm font-medium text-neutral-800">{{ user.username }}</p>
                    <p class="text-xs text-neutral-500">用户账号</p>
                  </div>
                </div>
              </td>
              <td class="px-8 py-5 text-sm text-neutral-700">{{ user.name }}</td>
              <td class="px-8 py-5">
                <span 
                  class="px-3 py-1.5 text-xs rounded-full"
                  :class="user.gender === '男' ? 'bg-blue-100 text-blue-800' : 'bg-pink-100 text-pink-800'"
                >
                  {{ user.gender }}
                </span>
              </td>
              <td class="px-8 py-5 text-sm text-neutral-700">{{ user.addr }}</td>
              <td class="px-8 py-5 text-sm text-neutral-700">{{ formatDateTime(user.createTime) }}</td>
              <td class="px-8 py-5 text-right">
                <button 
                  @click="viewUserDetails(user)" 
                  class="text-primary hover:text-primary/80 transition-colors duration-200"
                >
                  <i class="fa fa-eye text-lg"></i>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <!-- 空状态 -->
      <div v-if="filteredUsers.length === 0 && !isLoading" class="py-20 text-center">
        <div class="w-20 h-20 mx-auto mb-6 bg-neutral-100 rounded-full flex items-center justify-center shadow-md">
          <i class="fa fa-search text-neutral-400 text-2xl"></i>
        </div>
        <h3 class="text-xl font-semibold text-neutral-800 mb-3">未找到匹配的用户</h3>
        <p class="text-neutral-500 max-w-md mx-auto text-base">尝试调整搜索条件或清除筛选器以查看更多结果</p>
        <button 
          @click="filterByGender = ''; searchQuery = ''" 
          class="mt-6 px-6 py-3 border border-neutral-300 rounded-xl text-neutral-700 hover:bg-neutral-50 transition-all duration-300 shadow-sm"
        >
          清除筛选
        </button>
      </div>
      
      <!-- 加载状态 -->
      <div v-if="isLoading" class="py-20 text-center">
        <div class="inline-block animate-spin rounded-full h-10 w-10 border-b-2 border-primary"></div>
        <p class="mt-4 text-neutral-500 text-base">正在加载用户数据...</p>
      </div>
    </div>

    <!-- 分页控制 -->
    <div class="flex flex-col sm:flex-row justify-between items-center bg-white rounded-2xl shadow-xl p-8">
      <div class="text-sm text-neutral-500 mb-4 sm:mb-0">
        显示 {{ (currentPage - 1) * pageSize + 1 }} 到 {{ Math.min(currentPage * pageSize, users.length) }} 条，共 {{ users.length }} 条
      </div>
      <div class="flex items-center space-x-2">
        <button 
          @click="prevPage" 
          class="w-12 h-12 rounded-xl border border-neutral-200 flex items-center justify-center text-neutral-500 hover:bg-neutral-50 transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed shadow-sm"
          :disabled="currentPage === 1"
        >
          <i class="fa fa-chevron-left"></i>
        </button>
        <div class="flex">
          <button 
            v-for="page in pageNumbers" 
            :key="page" 
            @click="currentPage = page" 
            class="w-12 h-12 rounded-xl flex items-center justify-center mx-1 text-sm"
            :class="{
              'bg-primary text-white shadow-md': page === currentPage,
              'text-neutral-700 hover:bg-neutral-50 transition-colors': page !== currentPage
            }"
          >
            {{ page }}
          </button>
        </div>
        <button 
          @click="nextPage" 
          class="w-12 h-12 rounded-xl border border-neutral-200 flex items-center justify-center text-neutral-500 hover:bg-neutral-50 transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed shadow-sm"
          :disabled="currentPage === totalPages"
        >
          <i class="fa fa-chevron-right"></i>
        </button>
      </div>
    </div>

    <!-- 用户详情模态框 -->
    <div 
      v-if="selectedUser" 
      class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center opacity-0 pointer-events-none transition-opacity duration-300"
      :class="{ 'opacity-100 pointer-events-auto': selectedUser }"
    >
      <div class="bg-white rounded-2xl shadow-2xl max-w-md w-full max-h-[90vh] overflow-y-auto mx-4 transform transition-all duration-500 scale-95"
           :class="{ 'scale-100': selectedUser }">
        <div class="p-8 border-b border-neutral-100">
          <div class="flex justify-between items-center">
            <h3 class="text-xl font-semibold text-neutral-800">用户详情</h3>
            <button @click="selectedUser = null" class="text-neutral-400 hover:text-neutral-600 transition-colors duration-200">
              <i class="fa fa-times text-xl"></i>
            </button>
          </div>
        </div>
        <div class="p-8">
          <div class="flex items-center mb-6">
            <div class="w-20 h-20 rounded-full bg-primary/10 flex items-center justify-center text-primary text-2xl font-medium shadow-lg">
              {{ selectedUser.username.charAt(0) }}
            </div>
            <div class="ml-5">
              <h4 class="font-semibold text-lg text-neutral-800">{{ selectedUser.username }}</h4>
              <p class="text-neutral-500 text-sm">{{ selectedUser.name }}</p>
            </div>
          </div>
          
          <div class="space-y-5">
            <div class="grid grid-cols-3 gap-3 items-center">
              <span class="text-sm text-neutral-500">用户ID:</span>
              <span class="col-span-2 text-sm text-neutral-700 font-medium">{{ selectedUser.id }}</span>
            </div>
            <div class="grid grid-cols-3 gap-3 items-center">
              <span class="text-sm text-neutral-500">性别:</span>
              <span class="col-span-2 text-sm text-neutral-700">
                <span 
                  class="px-3 py-1.5 text-xs rounded-full"
                  :class="selectedUser.gender === '男' ? 'bg-blue-100 text-blue-800' : 'bg-pink-100 text-pink-800'"
                >
                  {{ selectedUser.gender }}
                </span>
              </span>
            </div>
            <div class="grid grid-cols-3 gap-3 items-center">
              <span class="text-sm text-neutral-500">地址:</span>
              <span class="col-span-2 text-sm text-neutral-700">{{ selectedUser.addr }}</span>
            </div>
            <div class="grid grid-cols-3 gap-3 items-center">
              <span class="text-sm text-neutral-500">创建时间:</span>
              <span class="col-span-2 text-sm text-neutral-700">{{ formatDateTime(selectedUser.createTime) }}</span>
            </div>
          </div>
        </div>
        <div class="p-8 border-t border-neutral-100 flex justify-end">
          <button 
            @click="selectedUser = null" 
            class="px-8 py-3 bg-primary text-white rounded-xl hover:bg-primary/90 transition-all duration-300 shadow-md"
          >
            关闭
          </button>
        </div>
      </div>
    </div>

    <!-- 页脚 -->
    <footer class="mt-16 text-center text-neutral-500 text-sm py-6 border-t border-neutral-200">
      <p>© 2025 用户管理系统 | 数据更新于: {{ lastUpdateTime }}</p>
    </footer>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'UserManagementSystem',
  data() {
    return {
      users: [],
      searchQuery: '',
      filterByGender: '',
      currentPage: 1,
      pageSize: 5,
      selectedUser: null,
      isLoading: true,
      error: null
    }
  },
  computed: {
    totalPages() {
      return Math.ceil(this.users.length / this.pageSize);
    },
    
    paginatedUsers() {
      const startIndex = (this.currentPage - 1) * this.pageSize;
      const endIndex = startIndex + this.pageSize;
      return this.users.slice(startIndex, endIndex);
    },
    
    filteredUsers() {
      let filtered = this.paginatedUsers;
      
      if (this.filterByGender) {
        filtered = filtered.filter(user => user.gender === this.filterByGender);
      }
      
      if (this.searchQuery) {
        const query = this.searchQuery.toLowerCase();
        filtered = filtered.filter(user => 
          user.username.toLowerCase().includes(query) || 
          user.name.toLowerCase().includes(query) ||
          user.addr.toLowerCase().includes(query)
        );
      }
      
      return filtered;
    },
    
    pageNumbers() {
      const maxVisiblePages = 5;
      let startPage = Math.max(1, this.currentPage - Math.floor(maxVisiblePages / 2));
      let endPage = startPage + maxVisiblePages - 1;
      
      if (endPage > this.totalPages) {
        endPage = this.totalPages;
        startPage = Math.max(1, endPage - maxVisiblePages + 1);
      }
      
      return Array.from({ length: endPage - startPage + 1 }, (_, i) => startPage + i);
    },
    
    maleCount() {
      return this.users.filter(user => user.gender === '男').length;
    },
    
    femaleCount() {
      return this.users.filter(user => user.gender === '女').length;
    },
    
    malePercentage() {
      if (this.users.length === 0) return 0;
      return Math.round((this.maleCount / this.users.length) * 100);
    },
    
    femalePercentage() {
      if (this.users.length === 0) return 0;
      return Math.round((this.femaleCount / this.users.length) * 100);
    },
    
    lastUpdateTime() {
      const now = new Date();
      return now.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      });
    }
  },
  methods: {
    formatDateTime(dateTimeString) {
      if (!dateTimeString) return '';
      
      const date = new Date(dateTimeString);
      
      if (isNaN(date.getTime())) {
        const parts = dateTimeString.split('T');
        if (parts.length === 2) {
          const datePart = parts[0].split('-');
          const timePart = parts[1].split('.')[0].split(':');
          
          if (datePart.length === 3 && timePart.length >= 2) {
            return `${datePart[0]}-${datePart[1]}-${datePart[2]} ${timePart[0]}:${timePart[1]}`;
          }
        }
        return dateTimeString;
      }
      
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      });
    },
    
    async fetchUsers() {
      this.isLoading = true;
      this.error = null;
      
      try {
        await new Promise(resolve => setTimeout(resolve, 800));
        
        const response = [
          {
            "id": 1,
            "username": "张治国",
            "password": "8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918",
            "name": "user003",
            "gender": "男",
            "addr": "大同",
            "headImg": "",
            "createTime": "2025-06-13T10:24:59.395547+08:00",
            "lat": 0,
            "lon": 0
          },
          {
            "id": 2,
            "username": "test1",
            "password": "1b4f0e9851971998e732078544c96b36c3d01cedf7caa332359d6f1d83567014",
            "name": "user001",
            "gender": "男",
            "addr": "太原",
            "headImg": "",
            "createTime": "2025-06-13T10:52:59.8995384+08:00",
            "lat": 0,
            "lon": 0
          },
          {
            "id": 3,
            "username": "test2",
            "password": "60303ae22b998861bce3b28f33eec1be758a213c86c93c076dbe9f558c11c752",
            "name": "user002",
            "gender": "女",
            "addr": "晋中",
            "headImg": "",
            "createTime": "2025-06-13T10:53:04.7467386+08:00",
            "lat": 0,
            "lon": 0
          },
          {
            "id": 4,
            "username": "test5",
            "password": "fd61a03af4f77d870fc21e05e7e80678095c92d808cfb3b5c279ee04c74aca13",
            "name": "user005",
            "gender": "男",
            "addr": "吕梁",
            "headImg": "",
            "createTime": "2025-06-13T10:53:08.8553382+08:00",
            "lat": 0,
            "lon": 0
          },
          {
            "id": 5,
            "username": "test4",
            "password": "a4e624d686e03ed2767c0abd85c14426b0b1157d2ce81d27bb4fe4f6f01d688a",
            "name": "user004",
            "gender": "女",
            "addr": "忻州",
            "headImg": "",
            "createTime": "2025-06-13T10:53:12.3484393+08:00",
            "lat": 0,
            "lon": 0
          }
        ];
        
        const formattedUsers = response.map(user => ({
          ...user,
          createTime: this.formatDateTime(user.createTime)
        }));
        
        this.users = formattedUsers;
      } catch (error) {
        console.error('获取用户数据失败:', error);
        this.error = '获取用户数据失败，请重试';
      } finally {
        this.isLoading = false;
      }
    },
    
    viewUserDetails(user) {
      this.selectedUser = { ...user };
      // 触发模态框动画
      this.$nextTick(() => {
        const modal = this.$el.querySelector('.fixed');
        if (modal) {
          modal.classList.add('opacity-100', 'pointer-events-auto');
        }
      });
    },
    
    prevPage() {
      if (this.currentPage > 1) {
        this.currentPage--;
      }
    },
    
    nextPage() {
      if (this.currentPage < this.totalPages) {
        this.currentPage++;
      }
    },
    
    refreshData() {
      this.currentPage = 1;
      this.searchQuery = '';
      this.filterByGender = '';
      this.fetchUsers();
    }
  },
  mounted() {
    this.fetchUsers();
  }
}
</script>

<style scoped>
@layer utilities {
  .content-auto {
    content-visibility: auto;
  }
  .shadow-card {
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.08);
  }
  .transition-custom {
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }
  .bg-gradient-blue {
    background: linear-gradient(135deg, #165DFF 0%, #36CFC9 100%);
  }
}
</style>