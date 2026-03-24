<template>
  <div class="funnel-layout">
    <!-- 左侧配置区 -->
    <div class="left-pane">
      <div class="left-pane-content">
        <slot name="config"></slot>
      </div>
      <div class="left-pane-footer">
        <slot name="footer"></slot>
      </div>
    </div>

    <!-- 右侧结果区 -->
    <div class="right-pane">
      <div class="right-pane-header">
        <slot name="header"></slot>
      </div>
      <div class="right-pane-content">
        <slot name="results"></slot>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.funnel-layout {
  display: flex;
  gap: 32px;
  flex: 1;
  height: 100%; /* 使用 100% 高度填充父容器 */
  min-height: 0;
  overflow: hidden;
}

.left-pane {
  width: 520px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

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

.left-pane-footer {
  flex-shrink: 0;
  background: var(--el-bg-color);
  border-radius: 0 0 12px 12px;
  box-shadow: 0 -4px 12px rgba(0, 0, 0, 0.05);
  z-index: 10;
  overflow: hidden;
}

.right-pane {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.right-pane-header {
  flex-shrink: 0;
  margin-bottom: 16px;
}

.right-pane-content {
  flex: 1;
  overflow-y: auto;
  min-height: 0;
  padding-right: 4px;
}
</style>
