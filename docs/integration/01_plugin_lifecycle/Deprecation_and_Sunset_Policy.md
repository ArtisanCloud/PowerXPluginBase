# 插件生命周期终止与退役策略（01_plugin_lifecycle/Deprecation_and_Sunset_Policy.md）

> 本文档定义 PowerX 插件从“弃用（Deprecation）”到“退役（Sunset）”的管理规范，  
> 适用于所有基于 **PowerXPluginBase** 的插件（含自研与第三方）。

---

## 🧭 一、目标

确保插件在停止维护、下架、替换或版本合并时：

- 不影响宿主系统与租户业务；
- 有可追溯的计划、通知与回滚路径；
- 能在 Marketplace、宿主 PowerX、以及插件本体三方保持状态一致。

---

## 🧩 二、关键阶段定义

| 阶段 | 描述 | 插件状态 |
|------|------|----------|
| **Active（活跃）** | 正常维护与发布新版本 | 可安装、可更新 |
| **Deprecated（弃用）** | 宣布不再推荐使用，将停止更新 | 可继续运行，不再新增功能 |
| **Sunset（退役）** | 彻底下架，不再提供支持或安装包 | 不可安装，不可更新，仅保留历史记录 |
| **Archived（归档）** | 代码仓库与文档冻结，仅供审计 | 不再接受 PR 或 Issue |

---

## 🧱 三、Deprecation 触发条件

- 插件被替换为新版或合并到主干功能；
- 插件依赖的 PowerX 核心接口或 SDK 版本不再兼容；
- 插件存在安全漏洞、合规风险；
- 插件长期无维护或测试失败；
- Marketplace 侧策略性下架（商业/品牌原因）。

---

## ⚙️ 四、弃用（Deprecation）流程

| 阶段 | 动作 | 责任方 |
|------|------|--------|
| 1️⃣ 内部评估 | 核心组或 Vendor 提出弃用申请（Deprecation Proposal） | 插件维护者 |
| 2️⃣ 通知 | 在 Marketplace 与 PowerX Admin 插件详情页标记 “Deprecated” | Marketplace 团队 |
| 3️⃣ 技术迁移 | 提供迁移建议（替代插件 / API 映射 / 数据导出方案） | 插件维护者 |
| 4️⃣ 发布公告 | 公告应包含弃用日期（Deprecation Date）与退役日期（Sunset Date） | Marketplace 团队 |
| 5️⃣ 限制更新 | 禁止新版本发布，仅允许安全修复 | Marketplace 自动化 |
| 6️⃣ 持续监测 | 确保租户业务无异常，收集反馈 | PowerX 核心组 |

### 通知模板示例

```yaml
notice:
  type: deprecation
  title: "插件弃用公告"
  plugin_id: com.powerx.plugin.crm
  deprecated_at: "2025-12-01"
  sunset_at: "2026-03-01"
  replacement: com.powerx.plugin.crm.v2
  migration_guide: https://docs.powerx.cloud/migrate/crm-v2
````

---

## ☀️ 五、退役（Sunset）流程

| 阶段          | 动作                            | 说明        |
| ----------- | ----------------------------- | --------- |
| 1️⃣ 冻结      | 插件在 Marketplace 中下架；宿主无法安装新副本 | 已安装租户仍可运行 |
| 2️⃣ 停止更新    | 插件仓库标记 `sunset=true`，停止 CI/CD | 禁止提交新版本   |
| 3️⃣ 终止服务    | 后端 API 停止响应（可返回 410 Gone）     | 需提前30天公告  |
| 4️⃣ 数据迁移/导出 | 提供数据导出工具或中转插件                 | 防止租户数据丢失  |
| 5️⃣ 归档      | 插件转入只读状态（Archive）             | 代码与文档冻结   |

---

## 🧰 六、插件自身的 Deprecation 信号机制

插件可在自身 manifest 中声明 Deprecation 元信息：

```yaml
lifecycle:
  status: deprecated
  deprecated_at: "2025-12-01"
  sunset_at: "2026-03-01"
  replacement: com.powerx.plugin.crm.v2
  message: "本插件已弃用，请尽快迁移至 CRM Plugin v2"
```

宿主 PowerX 会：

- 在插件详情中显示弃用提示；
- 向租户管理员发送通知；
- 自动推荐替代插件（若提供 replacement 字段）。

---

## 🧩 七、迁移与替代建议

| 场景         | 建议处理                     |
| ---------- | ------------------------ |
| 业务模块合并     | 将旧插件数据导出并导入新插件 Schema    |
| 接口替换       | 提供兼容层（Adapter）以维持短期兼容    |
| SDK 变更     | 更新 PowerX gRPC SDK 调用封装  |
| 前端替换       | 保留相同路由结构与入口，减少 UI 变动     |
| Agent 工具迁移 | 通过 MCP 注册新 Agent 工具并废弃旧的 |

> 建议提供 `migration.md` 文档放于 `docs/migration/` 目录。

---

## 🧩 八、Marketplace 状态与策略

Marketplace 应同步维护每个版本的状态：

| 字段                  | 示例                                      | 说明     |
| ------------------- | --------------------------------------- | ------ |
| `status`            | `active` / `deprecated` / `sunset`      | 生命周期状态 |
| `deprecated_at`     | `2025-12-01`                            | 弃用时间   |
| `sunset_at`         | `2026-03-01`                            | 退役时间   |
| `replacement`       | `com.powerx.plugin.crm.v2`              | 推荐替代   |
| `support_level`     | `maintenance`                           | 支持级别   |
| `security_advisory` | `https://market/p/advisory/AC-2025-001` | 安全公告链接 |

---

## 🔒 九、安全与合规要求

- 插件退役后 **必须撤销所有 ToolGrant 与 Token 权限**；
- 插件如涉及数据存储，应在退役时清理租户缓存与日志；
- 若插件由外部 Vendor 提供，应签署 **Sunset & Data Retention 协议**；
- Marketplace 应保存退役插件的 `.pxp` 包供合规审计；
- 所有退役操作须在 PowerX 宿主日志中留存记录。

---

## 🧠 十、文档与公告模板（示例）

路径：`docs/lifecycle/NOTICE_DEPRECATION.md`（通过 `make sync-lifecycle-docs` 同步至本目录）

```markdown
# 插件弃用公告：CRM Plugin

本插件（com.powerx.plugin.crm）将于 **2025-12-01** 起标记为弃用，计划于 **2026-03-01** 完全退役。

推荐替代插件：**com.powerx.plugin.crm.v2**

迁移文档：  
👉 [CRM v1 → v2 迁移指南](https://docs.powerx.cloud/migrate/crm-v2)

PowerX 团队与 ArtisanCloud 将继续提供技术支持直至退役完成。
```

---

## 📚 延伸阅读

- [Manifest_and_Metadata.md](./Manifest_and_Metadata.md)
- [Versioning_and_Publishing.md](./Versioning_and_Publishing.md)
- [04_security_and_compliance/Vulnerability_Response.md](../04_security_and_compliance/Vulnerability_Response.md)

---

> **文档版本：** v1.0.0
> **适用范围：** PowerX ≥ 0.9.0
> **维护团队：** PluginBase 核心组 & Marketplace 团队
> **最后更新：** 2025-10
