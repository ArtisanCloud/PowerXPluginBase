# gRPC Client 使用说明

## 概述

已成功实现了基于文档要求的 gRPC client，包括：

1. **PowerX gRPC 客户端封装** - 用于调用 PowerX 核心服务
2. **插件 gRPC 服务器** - 插件自己的 gRPC 服务（可选）
3. **配置支持** - 完整的 YAML 配置和环境变量支持
4. **并发启动** - HTTP 和 gRPC 服务器并行运行
5. **演示路由** - 用于测试 gRPC 连接和功能

## 架构

```
./backend/
├── internal/
│   ├── grpc/
│   │   ├── client/
│   │   │   └── powerx.go          # PowerX gRPC 客户端封装
│   │   └── server/
│   │       └── server.go          # 插件 gRPC 服务器
│   ├── config/config.go           # 添加了 gRPC 配置支持
│   └── api/http/demo/routes.go    # gRPC 演示路由
└── cmd/plugin/main.go             # 集成 gRPC 客户端和服务器启动
```

## 配置

### YAML 配置文件 (`etc/config.yaml`)

```yaml
# gRPC 上游配置（连接 PowerX 核心服务）
grpc_upstream:
  address: "localhost:9001"    # PowerX gRPC 服务地址
  token: ""                    # Capability Token（插件安装后获得）
  tenant_id: 1                 # 租户 ID
  use_tls: false               # 开发环境关闭 TLS
  ca_cert: ""                  # TLS 根证书路径（生产环境使用）

# gRPC 服务器配置（插件提供的 gRPC 服务）
grpc_server:
  enable: true                 # 是否启用插件 gRPC 服务器
  addr: ":9101"               # 插件 gRPC 监听地址
  use_tls: false               # 开发环境关闭 TLS
  cert: ""                    # TLS 证书路径（生产环境使用）
  key: ""                     # TLS 私钥路径（生产环境使用）
```

### 环境变量

```bash
# PowerX gRPC 上游配置
export POWERX_GRPC_UPSTREAM_ADDRESS="localhost:9001"
export POWERX_GRPC_UPSTREAM_TOKEN="your_capability_token"
export POWERX_GRPC_UPSTREAM_TENANT_ID="1"
export POWERX_GRPC_UPSTREAM_USE_TLS="false"

# 插件 gRPC 服务器配置
export POWERX_GRPC_SERVER_ENABLE="true"
export POWERX_GRPC_SERVER_ADDR=":9101"
export POWERX_GRPC_SERVER_USE_TLS="false"
```

## 功能特性

### PowerX gRPC 客户端

- ✅ **连接管理** - 自动连接、重连、健康检查
- ✅ **认证支持** - Bearer Token 和租户 ID 自动附加
- ✅ **TLS 支持** - 生产环境 TLS/mTLS 支持
- ✅ **上下文管理** - 正确的元数据传递
- ✅ **错误处理** - 连接失败和重试机制

### 插件 gRPC 服务器

- ✅ **健康检查** - 内置健康检查服务
- ✅ **反射服务** - 开发调试支持
- ✅ **优雅关闭** - 支持信号中断的优雅关闭
- ✅ **TLS 支持** - 可配置的 TLS 支持
- ✅ **并发运行** - 与 HTTP 服务器并行运行

## API 端点

### gRPC 演示端点

- `GET /api/v1/demo/grpc/health` - 检查 PowerX gRPC 连接状态
- `GET /api/v1/demo/grpc/members` - 演示调用 PowerX 成员服务（模拟）
- `GET /api/v1/demo/grpc/teams` - 演示调用 PowerX 团队服务（模拟）

### 示例响应

```bash
# 检查 gRPC 连接健康状态
curl http://localhost:8086/api/v1/demo/grpc/health

# 响应示例
{
  "status": "ok",
  "connected": true,
  "tenant_id": 1,
  "has_token": false
}
```

## 开发使用

### 1. 启动应用

```bash
cd backend
go run ./cmd/plugin
```

应用将并行启动：
- HTTP 服务器：`localhost:8086`
- gRPC 服务器：`localhost:9101`

### 2. 测试 gRPC 连接

```bash
# 检查 PowerX gRPC 连接状态
curl http://localhost:8086/api/v1/demo/grpc/health

# 测试演示端点
curl http://localhost:8086/api/v1/demo/grpc/members
```

### 3. 使用 grpcurl 测试插件 gRPC 服务器

```bash
# 列出可用服务
grpcurl -plaintext localhost:9101 list

# 健康检查
grpcurl -plaintext localhost:9101 grpc.health.v1.Health/Check
```

## 集成 PowerX

### 当前状态

- ✅ gRPC 客户端框架已完成
- ✅ 配置系统已集成
- ✅ 并发启动已实现
- ⏳ **需要 PowerX proto 文件** 来实现真实的服务调用

### 下一步

1. **获取 PowerX proto 文件**
   ```bash
   # 示例：从 Core 仓库获取生成的 Go 代码
   go get github.com/ArtisanCloud/PowerX/api/grpc/gen/go/...
   ```

2. **实现真实的服务调用**
   ```go
   // 在 powerx.go 中取消注释并实现
   func (p *PowerX) ListMembers(ctx context.Context, req *orgv1.ListMembersRequest) (*orgv1.ListMembersResponse, error) {
       return p.Members.ListMembers(p.Outgoing(ctx), req)
   }
   ```

3. **添加插件自己的 proto 服务**
   ```bash
   # 定义插件 proto 文件
   # 生成 Go 代码
   # 在 server.go 中注册服务
   ```

## 生产环境配置

### TLS 配置

```yaml
grpc_upstream:
  use_tls: true
  ca_cert: "/path/to/ca.crt"

grpc_server:
  use_tls: true
  cert: "/path/to/server.crt"
  key: "/path/to/server.key"
```

### 安全考虑

- 生产环境务必启用 TLS
- 使用有效的 Capability Token
- 配置合适的超时和重试策略
- 监控 gRPC 连接状态

## 故障排除

### 常见问题

1. **连接失败**
   - 检查 PowerX 服务是否运行
   - 验证地址和端口配置
   - 查看日志中的连接错误

2. **认证失败**
   - 确认 Capability Token 有效
   - 检查租户 ID 是否正确
   - 验证权限配置

3. **TLS 错误**
   - 检查证书路径和有效性
   - 验证 CA 证书配置
   - 确认 TLS 版本兼容性

### 日志查看

```bash
# 查看应用日志
cd backend
go run ./cmd/plugin | grep grpc

# 调试模式运行
POWERX_LOG_LEVEL=debug go run ./cmd/plugin
```

## 完成状态

✅ **所有文档要求已实现**：

1. ✅ 配置文件支持 gRPC 客户端和服务器配置
2. ✅ 创建了完整的 gRPC 目录结构
3. ✅ 添加了所有必需的 gRPC 依赖
4. ✅ 实现了 PowerX gRPC 客户端封装
5. ✅ 实现了插件 gRPC 服务器
6. ✅ 更新了 main.go 集成 gRPC 启动
7. ✅ 更新了配置示例文件
8. ✅ 提供了演示路由和使用示例

代码已完全可运行，等待 PowerX proto 文件的集成以实现完整功能。