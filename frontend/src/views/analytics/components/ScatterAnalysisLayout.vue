<template>
  <div class="layout">
    <!-- 左侧：指标创建 -->
    <div class="left-pane">
      <!-- 可滚动内容区域 -->
      <div class="left-pane-content">
        <slot name="config"></slot>
      </div>
      
      <!-- 固定在底部的操作按钮区域 -->
      <div class="left-pane-footer">
        <slot name="footer"></slot>
      </div>
    </div>

    <!-- 右侧：时间粒度/时间选择 + 结果展示 -->
    <div class="right-pane">
      <!-- 时间控制区域（固定不滚动） -->
      <div class="right-pane-header">
        <slot name="header"></slot>
      </div>
      
      <!-- 数据结果区域（独立滚动） -->
      <div class="right-pane-content">
        <slot name="results"></slot>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// Generic layout for scatter analysis
</script>

<style scoped lang="scss">
/* 布局容器 - 使用 flex 布局实现固定高度 */
.layout {
  display: flex;
  gap: 32px;
  flex: 1;
  height: calc(100vh - 180px); /* 减去顶部导航、页眉和内边距的估计高度 */
  min-height: 0; /* 重要：允许 flex item 正确缩小 */
  overflow: hidden;
}

/* 左侧面板 - 使用 flex 布局分为内容区 and 底部按钮区 */
.left-pane {
  width: 520px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

/* 左侧可滚动内容区域 */
.left-pane-content {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding-right: 8px;
  min-height: 0;
}

.left-pane-content > * {
  flex-shrink: 0;
}

.left-pane-content::-webkit-scrollbar {
  width: 6px;
}

.left-pane-content::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.left-pane-content::-webkit-scrollbar-thumb {
  background: var(--el-color-primary-light-5);
  border-radius: 3px;
  transition: background 0.2s ease;
}

.left-pane-content::-webkit-scrollbar-thumb:hover {
  background: var(--el-color-primary-light-3);
}

/* 左侧固定在底部的操作按钮区域 */
.left-pane-footer {
  flex-shrink: 0;
  background: var(--el-bg-color);
  border-radius: 0 0 12px 12px;
  box-shadow: 0 -4px 12px rgba(0, 0, 0, 0.05);
  z-index: 10;
  overflow: hidden;
}

/* 右侧面板 - 分为固定头部和滚动内容 */
.right-pane {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

/* 右侧固定头部（时间控制区域） */
.right-pane-header {
  flex-shrink: 0;
  margin-bottom: 12px;
}

/* 右侧可滚动内容区域（数据结果） */
.right-pane-content {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  min-height: 0;
  padding-right: 4px;
}

.right-pane-content::-webkit-scrollbar {
  width: 6px;
}

.right-pane-content::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.right-pane-content::-webkit-scrollbar-thumb {
  background: var(--el-color-success-light-5);
  border-radius: 3px;
  transition: background 0.2s ease;
}

.right-pane-content::-webkit-scrollbar-thumb:hover {
  background: var(--el-color-success-light-3);
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .left-pane {
    width: 480px;
  }
}

@media (max-width: 992px) {
  .layout {
    flex-direction: column;
    overflow-y: auto;
    height: auto;
  }
  
  .left-pane {
    width: 100%;
    height: auto;
    max-height: none;
    overflow: visible;
  }
  
  .left-pane-content {
    max-height: none;
    overflow-y: visible;
  }
  
  .right-pane {
    width: 100%;
    height: auto;
    min-height: 400px;
    overflow: visible;
  }
  
  .right-pane-header {
    margin-bottom: 12px;
  }
  
  .right-pane-content {
    overflow-y: visible;
  }
}
</style>
