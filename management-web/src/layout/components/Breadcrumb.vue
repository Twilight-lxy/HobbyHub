<template>
  <el-breadcrumb separator="/">
    <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
    <el-breadcrumb-item 
      v-for="(item, index) in breadcrumbs" 
      :key="item.path"
      :to="index < breadcrumbs.length - 1 ? { path: item.path } : null"
    >
      {{ item.title }}
    </el-breadcrumb-item>
  </el-breadcrumb>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const breadcrumbs = ref([])

// 获取面包屑数据
const getBreadcrumbs = () => {
  const matched = route.matched.filter(item => item.meta && item.meta.title)
  
  // 如果没有匹配的路由，直接返回
  if (!matched.length) {
    return []
  }
  
  return matched.map(item => {
    return {
      path: item.path,
      title: item.meta.title
    }
  })
}

// 监听路由变化，更新面包屑
watch(
  () => route.path,
  () => {
    breadcrumbs.value = getBreadcrumbs()
  },
  { immediate: true }
)
</script>

<style lang="scss" scoped>
.el-breadcrumb {
  display: inline-block;
  line-height: 60px;
  font-size: 14px;
}
</style> 