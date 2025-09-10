# IAM API 调试指南 - 新架构

## 概述

按照你的要求重新组织了代码架构：

- `internal/api/http/<domain>/api.go` - API 路由注册
- `internal/transport/http/<domain>/<handler>.go` - HTTP 传输层处理器
- `internal/services/<service>.go` - 业务逻辑服务
- Handler 中集成了 internal service 和 PowerX gRPC client

## 🏗️ 新架构结构

```
backend/
├── internal/
│   ├── api/http/admin/iam/
│   │   └── api.go                          # IAM API 路由注册
│   ├── transport/http/admin/iam/
│   │   └── member_handler.go               # 成员处理器（集成内部服务+gRPC）
│   ├── services/
│   │   └── member_service.go               # 成员业务服务
│   └── grpc/client/
│       └── powerx.go                       # PowerX gRPC 客户端
```

## 🚀 启动应用

```bash
cd backend
go run ./cmd/plugin
```

应用启动后可以访问：

- HTTP 服务器: `http://localhost:8086`
- gRPC 服务器: `localhost:9101`

## 📋 新 IAM API 接口

### 1. IAM 模块状态

```bash
curl http://localhost:8086/api/v1/admin/iam/status
```

**响应示例：**

```json
{
  "module": "iam",
  "version": "1.0.0",
  "services": ["members", "teams"],
  "endpoints": [
    "GET /iam/members - 获取成员列表",
    "GET /iam/members/:id - 获取单个成员",
    "GET /iam/members/search - 搜索成员",
    "GET /iam/members/connection/check - 检查连接",
    "GET /iam/status - 模块状态"
  ],
  "note": "集成 PowerX gRPC 服务进行成员和团队管理"
}
```

### 2. 获取成员列表（通过 PowerX gRPC）

```bash
curl http://localhost:8086/api/v1/admin/iam/members
```

**带分页参数：**

```bash
curl "http://localhost:8086/api/v1/admin/iam/members?page_size=10&page_index=0"
```

**带搜索关键词：**

```bash
curl "http://localhost:8086/api/v1/admin/iam/members?keyword=alice&page_size=20"
```

**响应示例：**

```json
{
  "success": true,
  "data": {
    "data": [
      {
        "id": 1,
        "user_id": 101,
        "name": "Alice Johnson",
        "email": "alice@example.com",
        "phone": "+1234567890",
        "position": "Software Engineer",
        "department": "Engineering",
        "status": "active",
        "joined_at": "2023-01-15T00:00:00Z",
        "created_at": "2023-01-15T10:00:00Z",
        "updated_at": "2024-01-15T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 2,
      "total_pages": 1
    }
  },
  "message": "Successfully retrieved members list",
  "timestamp": "2024-01-15T10:30:00Z",
  "request_id": "req_123456789"
}
```

### 3. 获取单个成员

```bash
curl http://localhost:8086/api/v1/admin/iam/members/123
```

**响应示例：**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 101,
    "name": "Alice Johnson",
    "email": "alice@example.com",
    "position": "Software Engineer",
    "department": "Engineering",
    "status": "active"
  },
  "message": "Successfully retrieved member details",
  "timestamp": "2024-01-15T10:30:00Z",
  "request_id": "req_123456789"
}
```

### 4. 高级成员搜索

```bash
curl "http://localhost:8086/api/v1/admin/iam/members/search?keyword=alice&department=engineering&status=active"
```

**带团队 ID 搜索：**

```bash
curl "http://localhost:8086/api/v1/admin/iam/members/search?keyword=dev&team_ids=1"
```

**响应示例：**

```json
{
  "success": true,
  "data": {
    "members": [
      {
        "id": 1,
        "name": "Alice Johnson",
        "department": "Engineering",
        "status": "active"
      }
    ],
    "total_count": 1,
    "page_size": 50
  },
  "query": {
    "keyword": "alice",
    "department": "engineering",
    "status": "active",
    "team_ids": [1]
  },
  "meta": {
    "source": "powerx_grpc_search"
  }
}
```

### 5. 检查 PowerX 连接状态

```bash
curl http://localhost:8086/api/v1/admin/iam/members/connection/check
```

**响应示例：**

```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "service": "PowerX MemberService",
    "version": "1.0.0",
    "timestamp": "2024-01-15T10:30:00Z",
    "checks": {
      "grpc_connected": "ok",
      "tenant_id": "configured",
      "has_token": "not_configured"
    }
  },
  "message": "PowerX MemberService connection is healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "request_id": "req_123456789"
}
```

## 🔄 对比：新架构 vs 演示路由

### 演示路由（临时测试）

```bash
# 基础演示路由
curl http://localhost:8086/api/v1/demo/grpc/members
curl http://localhost:8086/api/v1/demo/grpc/health
```

### 新架构 API（生产就绪）

```bash
# 标准 IAM API
curl http://localhost:8086/api/v1/admin/iam/members
curl http://localhost:8086/api/v1/admin/iam/members/connection/check
```

## 📊 架构优势

### 1. 清晰的职责分离

- **API Layer** (`internal/api/http/admin/iam/api.go`): 路由注册和组织
- **Transport Layer** (`internal/transport/http/admin/iam/member_handler.go`): HTTP 传输层处理
- **Service Layer** (`internal/services/member_service.go`): 业务逻辑处理
- **gRPC Client** (`internal/grpc/client/powerx.go`): 外部服务调用

### 2. 双重数据源集成

- **PowerX gRPC**: 获取核心成员数据
- **Internal Service**: 处理插件特有的业务逻辑
- **统一接口**: Handler 层统一协调两个数据源

### 3. 可扩展性

```go
// Handler 中可以轻松组合多个数据源
func (h *MemberHandler) GetMemberWithStats(c *gin.Context) {
    // 1. 从 PowerX 获取成员基础信息
    member, err := h.powerxClient.GetMember(ctx, req)

    // 2. 从内部服务获取统计信息
    stats, err := h.memberService.GetMemberStats(ctx, memberID)

    // 3. 组合返回
    c.JSON(200, gin.H{
        "member": member,
        "stats":  stats,
    })
}
```

## 🧪 测试场景

### 场景 1: 基础成员查询

```bash
# 获取所有成员
curl http://localhost:8086/api/v1/admin/iam/members

# 分页查询
curl "http://localhost:8086/api/v1/admin/iam/members?page_size=5&page_index=0"
```

### 场景 2: 搜索和过滤

```bash
# 关键词搜索
curl "http://localhost:8086/api/v1/admin/iam/members/search?keyword=alice"

# 部门过滤
curl "http://localhost:8086/api/v1/admin/iam/members/search?department=engineering"

# 复合条件
curl "http://localhost:8086/api/v1/admin/iam/members/search?keyword=dev&status=active&team_ids=1"
```

### 场景 3: 单个成员详情

```bash
# 获取特定成员
curl http://localhost:8086/api/v1/admin/iam/members/123

# 不存在的成员（测试错误处理）
curl http://localhost:8086/api/v1/admin/iam/members/999999
```

### 场景 4: 连接状态监控

```bash
# 检查 gRPC 连接
curl http://localhost:8086/api/v1/admin/iam/members/connection/check

# 查看 IAM 模块状态
curl http://localhost:8086/api/v1/admin/iam/status
```

## 🔧 配置示例

### 环境变量

```bash
# PowerX gRPC 配置
export PX_GRPC_UPSTREAM_ADDRESS="localhost:9001"
export PX_GRPC_UPSTREAM_TOKEN="your_capability_token"
export PX_GRPC_UPSTREAM_TENANT_ID="1"

# 应用配置
export PX_BIND_ADDR=":8086"
export PX_LOG_LEVEL="debug"
export PX_DB_DSN="postgres://user:pass@localhost:5432/powerx_plugin_note?sslmode=disable"
```

### YAML 配置 (`backend/etc/config.yaml`)

```yaml
server:
  bind_addr: ":8086"
  log_level: "debug"
  dev_mode: true

grpc_upstream:
  address: "localhost:9001"
  token: ""
  tenant_id: 1
  use_tls: false

database:
  dsn: "postgres://user:pass@localhost:5432/powerx_plugin_note?sslmode=disable"
  schema: "note"
```

## 📈 后续扩展

### 1. 添加团队管理

```bash
# 预留的团队路由
curl http://localhost:8086/api/v1/admin/iam/teams
```

### 2. 添加更多业务逻辑

- 成员任务统计
- 成员活动记录
- 权限管理
- 审计日志

### 3. 性能优化

- 缓存策略
- 批量查询
- 异步处理

---

## 总结

新架构提供了：

✅ **清晰的分层结构** - API/Transport/Service 三层分离  
✅ **双数据源集成** - PowerX gRPC + 内部服务无缝结合  
✅ **标准化接口** - 生产级别的 API 设计  
✅ **易于扩展** - 新增功能只需添加对应的 handler 和路由  
✅ **完整测试支持** - 提供丰富的测试接口和场景  
✅ **统一响应格式** - 使用 contracts/api.go 定义的标准响应格式

## 📋 标准 API 响应格式优势

### 1. 统一的响应结构

所有接口都使用 `contracts.APIResponse` 结构，包含：

- `success`: 操作是否成功
- `data`: 实际数据内容
- `error`: 错误信息（包含错误码、消息、详情）
- `message`: 操作消息
- `timestamp`: 响应时间戳
- `request_id`: 请求追踪 ID

### 2. 标准化错误处理

- 使用预定义的错误码 (`contracts.ErrCode*`)
- 结构化的错误信息
- 便于客户端统一处理

### 3. 完善的分页支持

- 使用 `contracts.ListResponse` 和 `contracts.PaginationResponse`
- 1 基础的页码系统
- 包含总数、总页数等完整信息

### 4. 健康检查标准化

- 使用 `contracts.HealthResponse` 格式
- 包含服务状态、版本、检查项等详细信息

你现在可以通过这些 API 接口来测试 gRPC 调用，同时享受标准化、规范化的响应格式！
