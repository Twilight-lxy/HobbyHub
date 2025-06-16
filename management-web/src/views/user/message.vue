<template>
  <div class="bg-gray-50 min-h-screen">
    <!-- 主要内容区 -->
    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div class="bg-white shadow-sm rounded-lg overflow-hidden">
        <!-- 评论表格 -->
        <div class="overflow-x-auto">
          <div class="p-4 sm:p-6">
            <table class="min-w-full divide-y divide-gray-200 rounded-lg overflow-hidden">
              <thead class="bg-gray-50">
                <tr>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider bg-gray-50 border-b border-gray-200">
                    评论内容
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider bg-gray-50 border-b border-gray-200">
                    用户
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider bg-gray-50 border-b border-gray-200">
                    小组
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider bg-gray-50 border-b border-gray-200">
                    活动
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider bg-gray-50 border-b border-gray-200">
                    时间
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider bg-gray-50 border-b border-gray-200">
                    状态
                  </th>
                  <th scope="col" class="px-6 py-3 text-right text-xs font-semibold text-gray-500 uppercase tracking-wider bg-gray-50 border-b border-gray-200">
                    操作
                  </th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr v-for="comment in filteredComments" :key="comment.id" class="hover:bg-gray-50 transition-colors duration-200 hover:shadow-sm">
                  <td class="px-6 py-4 whitespace-nowrap border-b border-gray-100">
                    <div class="text-sm text-gray-900 max-w-xs truncate">{{ comment.content }}</div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap border-b border-gray-100">
                    <div class="flex items-center">
                      <img class="h-10 w-10 rounded-full object-cover border-2 border-gray-100" :src="comment.user.avatar" :alt="comment.user.name">
                      <div class="ml-3">
                        <div class="text-sm font-medium text-gray-900">{{ comment.user.name }}</div>
                        <div class="text-xs text-gray-500">{{ comment.user.email }}</div>
                      </div>
                    </div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap border-b border-gray-100">
                    <div class="text-sm text-gray-900">{{ comment.group.name }}</div>
                    <div class="text-xs text-gray-500">ID: {{ comment.group.id }}</div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap border-b border-gray-100">
                    <div class="text-sm text-gray-900">{{ comment.activity.title }}</div>
                    <div class="text-xs text-gray-500">ID: {{ comment.activity.id }}</div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 border-b border-gray-100">
                    {{ formatDate(comment.createdAt) }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap border-b border-gray-100">
                    <span v-if="comment.status === 'approved'" class="px-2.5 py-0.5 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                      已审核
                    </span>
                    <span v-else-if="comment.status === 'pending'" class="px-2.5 py-0.5 inline-flex text-xs leading-5 font-semibold rounded-full bg-yellow-100 text-yellow-800">
                      待审核
                    </span>
                    <span v-else-if="comment.status === 'rejected'" class="px-2.5 py-0.5 inline-flex text-xs leading-5 font-semibold rounded-full bg-red-100 text-red-800">
                      已拒绝
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium border-b border-gray-100">
                    <button @click="viewComment(comment)" class="text-primary hover:text-primary-dark mr-3 transition-colors flex items-center">
                      <i class="fa fa-eye mr-1.5"></i>查看
                    </button>
                    <button v-if="comment.status === 'pending'" @click="approveComment(comment.id)" class="text-green-600 hover:text-green-900 mr-3 transition-colors flex items-center">
                      <i class="fa fa-check mr-1.5"></i>通过
                    </button>
                    <button v-if="comment.status === 'pending'" @click="rejectComment(comment.id)" class="text-red-600 hover:text-red-900 transition-colors flex items-center">
                      <i class="fa fa-times mr-1.5"></i>拒绝
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- 分页 -->
        <div class="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6">
          <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
            <div>
              <p class="text-sm text-gray-700">
                显示第 <span class="font-medium">{{ (currentPage - 1) * pageSize + 1 }}</span> 到 <span class="font-medium">{{ Math.min(currentPage * pageSize, totalComments) }}</span> 条，共 <span class="font-medium">{{ totalComments }}</span> 条评论
              </p>
            </div>
            <div>
              <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
                <button @click="prevPage" :disabled="currentPage === 1" class="relative inline-flex items-center px-3 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors">
                  <i class="fa fa-chevron-left text-xs"></i>
                </button>
                <button v-for="page in totalPages" :key="page" @click="goToPage(page)" :class="page === currentPage ? 'z-10 bg-primary text-white border-primary' : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'" class="relative inline-flex items-center px-4 py-2 border text-sm font-medium transition-colors">
                  {{ page }}
                </button>
                <button @click="nextPage" :disabled="currentPage === totalPages" class="relative inline-flex items-center px-3 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors">
                  <i class="fa fa-chevron-right text-xs"></i>
                </button>
              </nav>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- 评论详情模态框 -->
    <div v-if="showCommentModal" class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true"></div>
        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <div class="sm:flex sm:items-start">
              <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left w-full">
                <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                  评论详情
                </h3>
                <div class="mt-2 border-t border-gray-100 pt-4">
                  <div class="flex items-center mb-4">
                    <img class="h-10 w-10 rounded-full object-cover" :src="selectedComment.user.avatar" :alt="selectedComment.user.name">
                    <div class="ml-3">
                      <div class="text-sm font-medium text-gray-900">{{ selectedComment.user.name }}</div>
                      <div class="text-xs text-gray-500">{{ selectedComment.user.email }}</div>
                    </div>
                    <div class="ml-auto text-xs text-gray-500">
                      {{ formatDate(selectedComment.createdAt) }}
                    </div>
                  </div>
                  <div class="mb-4">
                    <h4 class="text-sm font-medium text-gray-500 mb-1">评论内容：</h4>
                    <p class="text-sm text-gray-900 whitespace-pre-wrap">{{ selectedComment.content }}</p>
                  </div>
                  <div class="mb-4">
                    <h4 class="text-sm font-medium text-gray-500 mb-1">所属小组：</h4>
                    <p class="text-sm text-gray-900">{{ selectedComment.group.name }}</p>
                  </div>
                  <div class="mb-4">
                    <h4 class="text-sm font-medium text-gray-500 mb-1">关联活动：</h4>
                    <p class="text-sm text-gray-900">{{ selectedComment.activity.title }}</p>
                  </div>
                  <div>
                    <h4 class="text-sm font-medium text-gray-500 mb-1">评论状态：</h4>
                    <span v-if="selectedComment.status === 'approved'" class="px-3 py-1 inline-flex text-sm leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                      <i class="fa fa-check-circle mr-1"></i>已审核
                    </span>
                    <span v-else-if="selectedComment.status === 'pending'" class="px-3 py-1 inline-flex text-sm leading-5 font-semibold rounded-full bg-yellow-100 text-yellow-800">
                      <i class="fa fa-clock-o mr-1"></i>待审核
                    </span>
                    <span v-else-if="selectedComment.status === 'rejected'" class="px-3 py-1 inline-flex text-sm leading-5 font-semibold rounded-full bg-red-100 text-red-800">
                      <i class="fa fa-times-circle mr-1"></i>已拒绝
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button type="button" @click="closeCommentModal" class="w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary sm:ml-3 sm:w-auto sm:text-sm transition-colors">
              关闭
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'CommentManagement',
  data() {
    return {
      comments: [],
      searchQuery: '',
      filterStatus: 'all',
      filterGroup: 'all',
      currentPage: 1,
      pageSize: 10,
      showCommentModal: false,
      selectedComment: null,
      // 小组列表
      groups: [
        { id: 1, name: '摄影爱好者' },
        { id: 2, name: '读书会' },
        { id: 3, name: '登山队' },
        { id: 4, name: '羽毛球俱乐部' },
        { id: 5, name: '美食分享会' }
      ],
      // 模拟数据
      mockComments: [
        {
          id: 1,
          content: "这次活动组织得非常好，拍摄地点很漂亮，期待下次活动！",
          user: {
            id: 101,
            name: "张三",
            email: "zhangsan@example.com",
            avatar: "https://picsum.photos/id/1005/40/40"
          },
          group: {
            id: 1,
            name: "摄影爱好者"
          },
          activity: {
            id: 201,
            title: "春季郊外摄影活动"
          },
          status: "approved",
          createdAt: "2025-06-14T10:30:00Z"
        },
        {
          id: 2,
          content: "我觉得可以多安排一些室内拍摄的活动，下雨天也能参加。",
          user: {
            id: 102,
            name: "李四",
            email: "lisi@example.com",
            avatar: "https://picsum.photos/id/1012/40/40"
          },
          group: {
            id: 1,
            name: "摄影爱好者"
          },
          activity: {
            id: 201,
            title: "春季郊外摄影活动"
          },
          status: "pending",
          createdAt: "2025-06-14T14:45:00Z"
        },
        {
          id: 3,
          content: "这本书非常值得一读，作者的观点很独特，大家有时间可以看看。",
          user: {
            id: 103,
            name: "王五",
            email: "wangwu@example.com",
            avatar: "https://picsum.photos/id/1025/40/40"
          },
          group: {
            id: 2,
            name: "读书会"
          },
          activity: {
            id: 202,
            title: "《人类简史》读书分享会"
          },
          status: "approved",
          createdAt: "2025-06-14T16:20:00Z"
        },
        {
          id: 4,
          content: "我不太同意作者关于历史发展的观点，感觉有些片面。",
          user: {
            id: 104,
            name: "赵六",
            email: "zhaoliu@example.com",
            avatar: "https://picsum.photos/id/1027/40/40"
          },
          group: {
            id: 2,
            name: "读书会"
          },
          activity: {
            id: 202,
            title: "《人类简史》读书分享会"
          },
          status: "rejected",
          createdAt: "2025-06-14T18:10:00Z"
        },
        {
          id: 5,
          content: "上周的登山活动非常累但也很有趣，山顶的风景真不错！",
          user: {
            id: 105,
            name: "孙七",
            email: "sunqi@example.com",
            avatar: "https://picsum.photos/id/1062/40/40"
          },
          group: {
            id: 3,
            name: "登山队"
          },
          activity: {
            id: 203,
            title: "白云山登山活动"
          },
          status: "pending",
          createdAt: "2025-06-15T09:15:00Z"
        },
        {
          id: 6,
          content: "建议下次活动可以安排一些难度更高的路线，适合有经验的成员。",
          user: {
            id: 106,
            name: "周八",
            email: "zhouba@example.com",
            avatar: "https://picsum.photos/id/1074/40/40"
          },
          group: {
            id: 3,
            name: "登山队"
          },
          activity: {
            id: 203,
            title: "白云山登山活动"
          },
          status: "approved",
          createdAt: "2025-06-15T11:30:00Z"
        },
        {
          id: 7,
          content: "这次羽毛球比赛组织得很好，期待下次再参加！",
          user: {
            id: 107,
            name: "吴九",
            email: "wujiu@example.com",
            avatar: "https://picsum.photos/id/1074/40/40"
          },
          group: {
            id: 4,
            name: "羽毛球俱乐部"
          },
          activity: {
            id: 204,
            title: "春季羽毛球友谊赛"
          },
          status: "approved",
          createdAt: "2025-06-15T13:45:00Z"
        },
        {
          id: 8,
          content: "场地有点小，打球时有些拥挤，下次可以考虑换个更大的场地。",
          user: {
            id: 108,
            name: "郑十",
            email: "zhengshi@example.com",
            avatar: "https://picsum.photos/id/1083/40/40"
          },
          group: {
            id: 4,
            name: "羽毛球俱乐部"
          },
          activity: {
            id: 204,
            title: "春季羽毛球友谊赛"
          },
          status: "pending",
          createdAt: "2025-06-15T15:20:00Z"
        },
        {
          id: 9,
          content: "上次的美食分享太棒了，大家做的菜都很美味，希望下次还能参加！",
          user: {
            id: 109,
            name: "钱十一",
            email: "qian11@example.com",
            avatar: "https://picsum.photos/id/1084/40/40"
          },
          group: {
            id: 5,
            name: "美食分享会"
          },
          activity: {
            id: 205,
            title: "家常菜制作交流活动"
          },
          status: "pending",
          createdAt: "2025-06-15T16:45:00Z"
        },
        {
          id: 10,
          content: "建议下次可以尝试一些异国料理的制作，比如意大利面或者寿司。",
          user: {
            id: 110,
            name: "孙十二",
            email: "sun12@example.com",
            avatar: "https://picsum.photos/id/1080/40/40"
          },
          group: {
            id: 5,
            name: "美食分享会"
          },
          activity: {
            id: 205,
            title: "家常菜制作交流活动"
          },
          status: "approved",
          createdAt: "2025-06-15T18:10:00Z"
        }
      ]
    }
  },
  computed: {
    // 计算总评论数
    totalComments() {
      return this.filteredComments.length;
    },
    // 计算总页数
    totalPages() {
      return Math.ceil(this.totalComments / this.pageSize);
    },
    // 过滤和搜索后的评论
    filteredComments() {
      let filtered = this.comments;
      
      // 按小组过滤
      if (this.filterGroup !== 'all') {
        filtered = filtered.filter(comment => comment.group.id === parseInt(this.filterGroup));
      }
      
      // 按状态过滤
      if (this.filterStatus !== 'all') {
        filtered = filtered.filter(comment => comment.status === this.filterStatus);
      }
      
      // 按搜索词过滤
      if (this.searchQuery) {
        const query = this.searchQuery.toLowerCase();
        filtered = filtered.filter(comment => 
          comment.content.toLowerCase().includes(query) || 
          comment.user.name.toLowerCase().includes(query) || 
          comment.group.name.toLowerCase().includes(query) ||
          comment.activity.title.toLowerCase().includes(query)
        );
      }
      
      return filtered;
    },
    // 当前页的评论
    paginatedComments() {
      const startIndex = (this.currentPage - 1) * this.pageSize;
      const endIndex = startIndex + this.pageSize;
      return this.filteredComments.slice(startIndex, endIndex);
    }
  },
  created() {
    // 初始化加载评论数据
    this.loadComments();
  },
  methods: {
    // 加载评论数据
    loadComments() {
      // 模拟API请求延迟
      setTimeout(() => {
        this.comments = [...this.mockComments];
      }, 500);
    },
    // 刷新评论
    refreshComments() {
      this.currentPage = 1;
      this.loadComments();
    },
    // 格式化日期
    formatDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      });
    },
    // 查看评论详情
    viewComment(comment) {
      this.selectedComment = { ...comment };
      this.showCommentModal = true;
    },
    // 关闭评论详情模态框
    closeCommentModal() {
      this.showCommentModal = false;
      this.selectedComment = null;
    },
    // 审核通过评论
    approveComment(id) {
      const commentIndex = this.comments.findIndex(comment => comment.id === id);
      if (commentIndex !== -1) {
        this.comments[commentIndex].status = 'approved';
        // 这里可以添加实际的API请求代码
        console.log(`评论 ${id} 已审核通过`);
      }
    },
    // 拒绝评论
    rejectComment(id) {
      const commentIndex = this.comments.findIndex(comment => comment.id === id);
      if (commentIndex !== -1) {
        this.comments[commentIndex].status = 'rejected';
        // 这里可以添加实际的API请求代码
        console.log(`评论 ${id} 已拒绝`);
      }
    },
    // 上一页
    prevPage() {
      if (this.currentPage > 1) {
        this.currentPage--;
      }
    },
    // 下一页
    nextPage() {
      if (this.currentPage < this.totalPages) {
        this.currentPage++;
      }
    },
    // 跳转到指定页
    goToPage(page) {
      if (page >= 1 && page <= this.totalPages) {
        this.currentPage = page;
      }
    }
  }
}
</script>

<style scoped>
/* 这里可以添加自定义样式 */
.primary {
  color: #165DFF;
}
.primary-dark {
  color: #0E42D2;
}
</style>