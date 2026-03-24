<template>
  <el-config-provider :locale="zhCn">
    <component :is="layout">
      <router-view />
    </component>
  </el-config-provider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { ElConfigProvider } from 'element-plus'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
import SidebarLayout from '@/layout/SidebarLayout.vue'

const route = useRoute()

const layout = computed(() => {
  // 登录页等 meta.public 为 true 的不使用布局
  if (route.meta.public) {
    return 'div'
  }
  return SidebarLayout
})
</script>

<style>
body {
  margin: 0;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  height: 100vh;
}

/* 统一卡片间距 */
.el-card {
  margin-bottom: 20px;
}
</style>
