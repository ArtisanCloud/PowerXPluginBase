# Web Admin 前端

这个目录将包含插件的前端管理界面。

## 技术栈建议

- **框架**: Vue 3 + Nuxt 3 或 React + Next.js
- **UI 组件库**: Element Plus, Ant Design, 或 Tailwind UI
- **状态管理**: Pinia (Vue) 或 Zustand (React)
- **HTTP 客户端**: Axios
- **构建工具**: Vite

## 目录结构

```
web-admin/
├── components/          # 可复用组件
│   ├── common/         # 通用组件
│   ├── task/           # 任务相关组件
│   └── sprint/         # Sprint 相关组件
├── pages/              # 页面组件
│   ├── dashboard.vue   # 仪表板
│   ├── tasks/          # 任务管理页面
│   └── sprints/        # Sprint 管理页面
├── layouts/            # 布局组件
├── composables/        # 组合式函数 (Vue)
├── stores/            # 状态管理
├── utils/             # 工具函数
├── assets/            # 静态资源
└── public/            # 公共文件
```

## 开发指南

1. **安装依赖**
   ```bash
   npm install
   # 或
   yarn install
   ```

2. **开发模式**
   ```bash
   npm run dev
   # 或
   yarn dev
   ```

3. **构建生产版本**
   ```bash
   npm run build
   # 或
   yarn build
   ```

4. **与后端 API 集成**
   - 基础 URL: `/_p/com.powerx.plugins.scrum/api/v1`
   - 需要在请求头中包含 PowerX 提供的认证信息

## 主要功能页面

### 1. 仪表板 (Dashboard)
- Sprint 概览
- 任务统计
- 燃尽图
- 团队速度图表

### 2. 任务管理 (Tasks)
- 任务列表
- 任务看板
- 任务详情
- 创建/编辑任务

### 3. Sprint 管理 (Sprints)
- Sprint 列表
- 创建/编辑 Sprint
- Sprint 计划
- Sprint 回顾

### 4. 报告 (Reports)
- 速度报告
- 燃尽图
- 累积流图
- 团队效能分析

## 集成说明

本前端应用将通过 PowerX 的插件机制集成：

1. **路由集成**: 通过 `/_p/:id/admin/*` 访问
2. **认证集成**: 使用 PowerX 提供的用户上下文
3. **权限控制**: 基于 PowerX 的 RBAC 系统
4. **主题集成**: 遵循 PowerX 的设计规范

## TODO

- [ ] 选择前端技术栈
- [ ] 初始化前端项目
- [ ] 实现基础布局
- [ ] 集成 PowerX 认证
- [ ] 实现各功能模块
- [ ] 编写测试
- [ ] 优化性能
- [ ] 文档完善