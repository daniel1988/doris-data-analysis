---
name: "naming-convention"
description: "提供后端(Golang)和前端(Vue/TypeScript)的代码命名规范检查和建议。在进行代码审查或需要统一代码命名风格时调用。"
---

# 命名规范

该技能用于检查或指导项目中的代码命名，以确保符合以下团队规范：

## 后端 (Golang)
- **导出变量/函数**：使用 `PascalCase` (大驼峰命名法)。
- **私有变量/函数**：使用 `camelCase` (小驼峰命名法)。
- **数据库字段/JSON Tag**：使用 `snake_case` (下划线命名法)。

## 前端 (TypeScript/Vue)
- **组件文件**：使用 `PascalCase.vue` (例如：`UserList.vue`)。
- **变量与函数**：使用 `camelCase` (小驼峰命名法)。
- **类型与接口**：使用 `PascalCase` (大驼峰命名法)。
- **目录名**：尽量使用小写下划线 (`snake_case`) 或连字符 (`kebab-case`)。

在编写新代码或重构时，请始终遵守上述规范。
