<template>
  <div v-if="!item.hidden">
    <!-- 没有子菜单的情况 -->
    <template v-if="!hasChildren(item)">
      <el-menu-item :index="resolvePath(item.path)" :class="{'submenu-title-noDropdown': !isNest}">
        <el-icon v-if="item.meta && item.meta.icon">
          <component :is="item.meta.icon" />
        </el-icon>
        <template #title>
          <span>{{ item.meta.title }}</span>
        </template>
      </el-menu-item>
    </template>
    
    <!-- 有子菜单的情况 -->
    <el-sub-menu v-else :index="resolvePath(item.path)" popper-append-to-body>
      <template #title>
        <el-icon v-if="item.meta && item.meta.icon">
          <component :is="item.meta.icon" />
        </el-icon>
        <span>{{ item.meta.title }}</span>
      </template>
      
      <!-- 递归渲染子菜单 -->
      <sidebar-item
        v-for="child in item.children"
        :key="child.path"
        :item="child"
        :base-path="resolvePath(child.path)"
        :is-nest="true"
      />
    </el-sub-menu>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import path from 'path-browserify'

const props = defineProps({
  item: {
    type: Object,
    required: true
  },
  basePath: {
    type: String,
    default: ''
  },
  isNest: {
    type: Boolean,
    default: false
  }
})

// 判断是否有子菜单
const hasChildren = (item) => {
  return item.children && item.children.length > 0
}

// 解析路径
const resolvePath = (routePath) => {
  if (/^(https?:|mailto:|tel:)/.test(routePath)) {
    return routePath
  }
  return path.resolve(props.basePath, routePath)
}
</script>

<style lang="scss" scoped>
.el-menu-item, .el-sub-menu {
  .el-icon {
    margin-right: 10px;
  }
}
</style> 