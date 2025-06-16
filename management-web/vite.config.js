import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  server: {
    //设置开发服务器监听的端口
    port: 5173,
    //
    proxy: {
      // 当请求路径以 `/api` 开头时触发代理
      '/api': {
        //配置代理规则
        target: 'http://localhost:8081',
        changeOrigin: true,
        rewrite: (path) => path
      }
    }
  }
})
