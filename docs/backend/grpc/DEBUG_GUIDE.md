# PowerX gRPC 接口调试指南

## 概述

已经成功对接了 PowerX 的 Member Service 和 Team Service！目前使用模拟数据响应，当你启动真实的 PowerX gRPC 服务后，可以直接切换到真实调用。

## 🚀 启动应用

```bash
cd backend
go run ./cmd/plugin
```

应用将在以下端口启动：

- HTTP 服务器: `http://localhost:8086`
- gRPC 服务器: `localhost:9101`

## 📋 调试接口清单

### 1. 健康检查接口

```bash
curl http://localhost:8086/api/v1/demo/grpc/health
```

**响应示例：**

```json
{
  "status": "ok",
  "connected": true,
  "tenant_id": 1,
  "has_token": false
}
```

### 2. 获取成员列表

```bash
curl http://localhost:8086/api/v1/demo/grpc/members
```

**带查询参数：**

```bash
curl "http://localhost:8086/api/v1/demo/grpc/members?keyword=alice"
```

**响应示例：**

```json
{
  "success": true,
  "data": {
    "members": [
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
    "total_count": 2,
    "page_index": 0,
    "page_size": 20
  },
  "grpc": {
    "connected": true,
    "tenant_id": 1
  }
}
```

### 3. 获取单个成员

```bash
curl http://localhost:8086/api/v1/demo/grpc/members/123
```

**响应示例：**

```json
{
  "success": true,
  "data": {
    "member": {
      "id": 1,
      "user_id": 101,
      "name": "Alice Johnson",
      "email": "alice@example.com",
      "position": "Software Engineer",
      "department": "Engineering",
      "status": "active"
    }
  },
  "requested_id": "123"
}
```

### 4. 获取团队列表

```bash
curl http://localhost:8086/api/v1/demo/grpc/teams
```

**带查询参数：**

```bash
curl "http://localhost:8086/api/v1/demo/grpc/teams?keyword=dev"
```

**响应示例：**

```json
{
  "success": true,
  "data": {
    "teams": [
      {
        "id": 1,
        "name": "Development Team",
        "description": "Core development team",
        "leader_id": 101,
        "status": "active",
        "member_count": 5,
        "created_at": "2023-01-01T10:00:00Z",
        "updated_at": "2024-01-15T10:00:00Z"
      }
    ],
    "total_count": 2,
    "page_index": 0,
    "page_size": 20
  }
}
```

### 5. 获取单个团队

```bash
curl http://localhost:8086/api/v1/demo/grpc/teams/456
```

### 6. 调试信息接口

```bash
curl http://localhost:8086/api/v1/demo/grpc/debug
```

**响应示例：**

```json
{
  "grpc_connection": {
    "connected": true,
    "tenant_id": 1,
    "has_token": false
  },
  "endpoints": [
    "GET /api/v1/demo/grpc/health - 检查 gRPC 连接状态",
    "GET /api/v1/demo/grpc/members - 获取成员列表",
    "GET /api/v1/demo/grpc/members/{id} - 获取单个成员",
    "GET /api/v1/demo/grpc/teams - 获取团队列表",
    "GET /api/v1/demo/grpc/teams/{id} - 获取单个团队",
    "GET /api/v1/demo/grpc/debug - 查看调试信息"
  ],
  "note": "当前使用模拟数据，可以通过这些接口测试 gRPC 连接和功能"
}
```

## 🔧 配置说明

### 环境变量配置

```bash
# PowerX gRPC 服务配置
export PX_GRPC_UPSTREAM_ADDRESS="localhost:9001"
export PX_GRPC_UPSTREAM_TOKEN="your_capability_token"
export PX_GRPC_UPSTREAM_TENANT_ID="1"
export PX_GRPC_UPSTREAM_USE_TLS="false"

# 插件配置
export PX_BIND_ADDR=":8086"
export PX_LOG_LEVEL="debug"
```

### YAML 配置文件 (`backend/etc/config.yaml`)

```yaml
grpc_upstream:
  address: "localhost:9001"
  token: ""
  tenant_id: 1
  use_tls: false

server:
  bind_addr: ":8086"
  log_level: "debug"
  dev_mode: true
```

## 🧪 测试工具

### 使用 curl 测试

```bash
# 快速测试所有接口
curl http://localhost:8086/api/v1/demo/grpc/health
curl http://localhost:8086/api/v1/demo/grpc/members
curl http://localhost:8086/api/v1/demo/grpc/teams
curl http://localhost:8086/api/v1/demo/grpc/debug
```

### 使用 HTTPie 测试

```bash
http GET localhost:8086/api/v1/demo/grpc/members keyword==alice
http GET localhost:8086/api/v1/demo/grpc/teams keyword==dev
```

### 使用 Postman 测试

导入以下 collection：

```json
{
  "info": {
    "name": "PowerX Note Plugin - gRPC Debug"
  },
  "item": [
    {
      "name": "Health Check",
      "request": {
        "method": "GET",
        "url": "http://localhost:8086/api/v1/demo/grpc/health"
      }
    },
    {
      "name": "List Members",
      "request": {
        "method": "GET",
        "url": "http://localhost:8086/api/v1/demo/grpc/members"
      }
    },
    {
      "name": "Get Member",
      "request": {
        "method": "GET",
        "url": "http://localhost:8086/api/v1/demo/grpc/members/123"
      }
    },
    {
      "name": "List Teams",
      "request": {
        "method": "GET",
        "url": "http://localhost:8086/api/v1/demo/grpc/teams"
      }
    }
  ]
}
```

## 🔍 日志调试

### 查看 gRPC 调用日志

```bash
# 启用调试日志
PX_LOG_LEVEL=debug go run ./cmd/plugin

# 日志示例
INFO[2024-01-15T10:00:00Z] Calling PowerX gRPC service method=ListMembers service="powerx.organization.v1.MemberService"
INFO[2024-01-15T10:00:00Z] Starting HTTP server... addr=":8086"
INFO[2024-01-15T10:00:00Z] Starting gRPC server addr=":9101"
```

### 常见错误排查

**1. 连接失败**

```json
{
  "error": "PowerX gRPC service unavailable",
  "details": "connection is not established"
}
```

- 检查 PowerX 服务是否运行在 `localhost:9001`
- 检查网络连接和防火墙设置

**2. 认证失败**

```json
{
  "error": "Failed to call PowerX gRPC service",
  "details": "authentication failed"
}
```

- 检查 `PX_GRPC_UPSTREAM_TOKEN` 是否设置正确
- 检查 `PX_GRPC_UPSTREAM_TENANT_ID` 是否正确

## 🚀 切换到真实 PowerX 服务

当你的 PowerX gRPC 服务运行后，只需要修改 `internal/grpc/client/powerx.go` 中的 `invokeGRPC` 方法：

```go
// 替换模拟调用为真实 gRPC 调用
func (p *PowerX) invokeGRPC(ctx context.Context, service, method string, req, resp interface{}) error {
    // 使用真实的 gRPC 客户端调用
    switch service {
    case "powerx.organization.v1.MemberService":
        // 实际调用 MemberService
    case "powerx.organization.v1.TeamService":
        // 实际调用 TeamService
    }
    // ...
}
```

## 📊 性能监控

### gRPC 性能指标

- 连接状态：通过 `/api/v1/demo/grpc/health` 监控
- 调用延迟：在日志中查看调用耗时
- 错误率：通过返回的错误信息统计

### 建议监控指标

1. gRPC 连接存活状态
2. 接口调用成功率
3. 平均响应时间
4. 错误类型分布

---

## 总结

所有接口都已经准备好，你可以：

1. **启动应用**：`go run ./cmd/plugin`
2. **测试连接**：`curl http://localhost:8086/api/v1/demo/grpc/health`
3. **测试功能**：使用上面的 curl 命令测试各个接口
4. **查看日志**：观察 gRPC 调用日志
5. **配置真实服务**：修改配置连接到真实的 PowerX gRPC 服务

目前使用模拟数据，当你启动真实的 PowerX 服务后，可以通过修改 `invokeGRPC` 方法来切换到真实调用！
