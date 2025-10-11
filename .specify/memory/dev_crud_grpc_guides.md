# CRUD over gRPC — Plugin Guide

> 面向插件侧 gRPC 传输层的人读说明。  
> 目标：统一 proto 命名/拦截器/metadata/错误映射/生成路径；与 HTTP 共用 Service（PG-SVC-001）。

## 1. proto 命名空间与目录

- 建议命名空间：`api/grpc/<vendor>/<plugin>/v1`  
- **权威源**：以该目录下的 `.proto` 为单一事实来源（SoT）  
- **Go 生成输出**：可采用 `api/grpc/gen/go` 或插件内 `backend/api/grpc/gen/go`（保持 `paths = source_relative`）

> 若采用 buf：`buf.yaml`/`buf.gen.yaml` 放在 `api/grpc/contracts/` 或插件侧 `api/grpc/`，二选一，保持一致即可。

## 2. Service / Message 设计

- **Service（TemplateService）**：
  - `CreateTemplate(CreateTemplateRequest) returns (TemplateResponse)`
  - `UpdateTemplate(UpdateTemplateRequest) returns (TemplateResponse)`
  - `DeleteTemplate(DeleteTemplateRequest) returns (google.protobuf.Empty)`
  - `GetTemplate(GetTemplateRequest) returns (TemplateResponse)`
  - `ListTemplates(ListTemplateRequest) returns (ListTemplateResponse)`
- **消息命名**尽量与 HTTP DTO 对齐（字段一致，语义一致）。

## 3. Metadata 与多租户

- 入站拦截器从 `metadata` 读取并验证：`authorization`（JWT/HMAC）、`x-tenant-id`、`x-request-id` 等  
- 通过拦截器把 `tenant_id` 注入 DB 会话（RLS）：`SET LOCAL app.tenant_id = ?`  
- 与 HTTP 一致（PG-CTX-001）

## 4. 拦截器链（服务端）

- `auth`（验签/上下文注入）  
- `tenant`（RLS/会话变量）  
- `logging`（结构化日志，含 request_id/tenant_id）  
- `recovery`（panic 捕获，统一错误）

> 所有业务实现均**复用同一 Service**（PG-SVC-001），禁止把业务写在 gRPC handler 中。

## 5. 错误映射

- 使用 `google.rpc.Status` 或自定义错误码与 HTTP 对齐：  
  - 权限：`PERMISSION_DENIED` ↔ 403  
  - 未认证：`UNAUTHENTICATED` ↔ 401  
  - 资源不存在：`NOT_FOUND` ↔ 404  
  - 参数错误：`INVALID_ARGUMENT` ↔ 400  
  - 冲突：`ALREADY_EXISTS` ↔ 409  
  - 服务错误：`INTERNAL` ↔ 500

## 6. 生成与装配

- **生成命令**：`make proto-gen`（或 buf 等价）  
- **服务装配**：统一在 `internal/grpc/server/*.go` 或集中入口注册 `Register*ServiceServer(...)`  
- **依赖注入**：通过构造函数把 Service 注入到 gRPC 层

## 7. 测试策略

- **Contract Test**：proto/schema/生成产物路径  
- **Server Test**：拦截器链验证（auth/tenant/logging/recovery）  
- **E2E**：与 HTTP 的行为保持一致（单一 Service 源），尤其在多租户读写上

## 8. 合规清单（Checklist）

- [ ] proto 命名空间与输出路径固定，生成可重现  
- [ ] metadata 注入/校验（JWT/HMAC + tenant_id + request_id）  
- [ ] 拦截器链完整  
- [ ] 与 HTTP 复用同一 Service  
- [ ] 错误码映射一致  
- [ ] 有合同/拦截器/行为一致性测试

（相关 Gates：PG-CTX-001 / PG-SVC-001）
