---
name: "frontend-design"
description: "提供前端 Vue3 + TypeScript 项目的设计与开发规范。在进行前端组件开发、重构或代码审查时调用。"
---

# frontend-design (前端设计与开发规范)

该技能用于指导前端项目（Vue 3 + TypeScript）的架构设计、组件开发和编码规范。当涉及前端开发或代码审查时，请务必遵循以下原则。

## 1. 技术栈与基础规范
- **核心框架**：统一使用 **Vue 3** + **TypeScript** + **Composition API** (`<script setup lang="ts">`)。
- **状态管理**：全局状态使用 **Pinia**，复杂的组件业务逻辑应当抽取为 **Composables** (如 `useXxx.ts`)。
- **UI 组件库**：
  - 统一使用 **Element Plus** 组件库。
  - 图标统一使用 `@element-plus/icons-vue`。

## 2. 组件结构与代码组织
- **文件长度**：单个前端 Vue/TS 文件原则上不应超过 **500 行**。超过 300 行的 UI 模板应考虑拆分子组件。
- **代码顺序**：始终保持以下顺序：
  1. `<template>`
  2. `<script setup lang="ts">`
  3. `<style scoped lang="scss">`
- **样式规范**：
  - 必须使用 `scoped` 属性隔离组件样式。
  - 推荐使用 SCSS 预处理器，并充分利用 CSS 变量。
  - 页面级组件应优先使用 Flex/Grid 布局实现自适应，避免硬编码固定像素值。

## 3. 组件通信
- **Props**：必须使用 TypeScript 接口进行严格类型定义。
  ```ts
  interface Props {
    modelValue: string;
    items?: any[];
  }
  const props = defineProps<Props>()
  ```
- **Emits**：必须显式定义所有触发的事件类型。
  ```ts
  const emit = defineEmits<{
    (e: 'update:modelValue', val: string): void
    (e: 'change', item: any): void
  }>()
  ```

## 4. 分析模块专属规范
在处理特定的业务模块（如数据分析、图表等）时：
- **Context 模式**：所有分析模块（如 Event, Scatter 等）必须使用 `Provide/Inject` 模式共享状态，遵循 `AnalysisContext` 接口（包含 `state` 和 `actions`）。
- **指标规范化**：处理分析指标字段时，**必须**使用 `payloadNormalizer.ts` 中的 `normalizeMetricField` 函数，将前端的占位符（如 `__TOTAL_TIMES__`）自动转换为后端实际所需的 `field` 和 `formula` 配置。
- **核心常量**：必须使用预定义的枚举（如 `TableNames`, `FormulaTypes`, `VizModes` 等）代替硬编码字符串。

## 5. API 交互与错误处理
- 前后端交互必须统一使用 **JSON** 格式，字段命名统一为 **`snake_case`** (下划线)。
- 前端在发起请求前应进行基础的字段必填与格式校验。
- 接口报错或业务异常时，统一使用 Element Plus 的 `ElMessage` 组件进行提示反馈。