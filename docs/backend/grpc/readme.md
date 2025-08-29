# gRPC 集成指南（插件侧） — `grpc.md`

> 目标：在 **当前backend目录结构** 下为插件同时提供
> 1）**PowerX 核心的 gRPC 客户端**（调用 Member/Team 等服务）
> 2）**插件自己的 gRPC Server**（可选，后续供平台/其他插件调用）
> 并与现有 Gin HTTP 一起运行。

---

## 目录规划（在你现有结构上新增）

```
./backend/
├── cmd/
│   └── plugin/
│       └── main.go                  # 入口：同时启动 Gin + gRPC（插件侧）
├── internal/
│   ├── config/
│   │   └── config.go                # 追加 gRPC 配置（见下）
│   ├── grpc/
│   │   ├── client/
│   │   │   └── powerx.go            # ✅ PowerX gRPC 客户端封装（Member/Team）
│   │   └── server/
│   │       └── server.go            # ✅ 插件 gRPC Server（健康+反射，预留注册点）
│   └── ...（原有代码保持不动）
└── etc/
    └── config.yaml                  # 追加 grpc 配置项
```

---

## 1) 配置（`internal/config/config.go`）

在你已有的 `Config` 里补充：

```go
package config

type GRPCUpstream struct {
  Address   string `yaml:"address"`     // PowerX 网关/服务地址，如 "localhost:9001"
  Token     string `yaml:"token"`       // Capability Token（插件安装后下发）
  TenantID  int64  `yaml:"tenant_id"`   // 当前租户
  UseTLS    bool   `yaml:"use_tls"`     // 上线后建议 true
  CACert    string `yaml:"ca_cert"`     // 可选：根证书（UseTLS=true 时）
}

type GRPCServer struct {
  Enable bool   `yaml:"enable"`         // 是否启用插件自己的 gRPC Server
  Addr   string `yaml:"addr"`           // 插件 gRPC 监听，如 ":9101"
  UseTLS bool   `yaml:"use_tls"`
  Cert   string `yaml:"cert"`
  Key    string `yaml:"key"`
}

type Config struct {
  // ... 你已有的字段
  GRPCUpstream GRPCUpstream `yaml:"grpc_upstream"`
  GRPCServer   GRPCServer   `yaml:"grpc_server"`
}
```

`etc/config.yaml` 示例：

```yaml
grpc_upstream:
  address: "localhost:9001"
  token: ""
  tenant_id: 42
  use_tls: false

grpc_server:
  enable: true
  addr: ":9101"
  use_tls: false
```

---

## 2) PowerX gRPC 客户端封装（`internal/grpc/client/powerx.go`）

```go
package client

import (
  "context"
  "crypto/tls"
  "crypto/x509"
  "fmt"
  "io/ioutil"
  "time"

  cfgpkg "github.com/ArtisanCloud/PowerX/Core/Plugins/com.powerx.plugin.scrum/backend/internal/config"

  commonv1 "github.com/ArtisanCloud/PowerX/api/grpc/gen/go/common/v1"
  orgv1 "github.com/ArtisanCloud/PowerX/api/grpc/gen/go/powerx/organization/v1"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials"
  "google.golang.org/grpc/credentials/insecure"
  "google.golang.org/grpc/metadata"
)

type PowerX struct {
  conn    *grpc.ClientConn
  Members orgv1.MemberServiceClient
  Teams   orgv1.TeamServiceClient

  token    string
  tenantID int64
}

// NewPowerX 根据配置拨号 PowerX gRPC
func NewPowerX(ctx context.Context, c cfgpkg.GRPCUpstream) (*PowerX, error) {
  dialOpts := []grpc.DialOption{grpc.WithBlock()}
  if c.UseTLS {
    var creds credentials.TransportCredentials
    if c.CACert != "" {
      pem, err := ioutil.ReadFile(c.CACert)
      if err != nil {
        return nil, fmt.Errorf("read ca cert: %w", err)
      }
      cp := x509.NewCertPool()
      cp.AppendCertsFromPEM(pem)
      creds = credentials.NewTLS(&tls.Config{RootCAs: cp})
    } else {
      creds = credentials.NewTLS(&tls.Config{})
    }
    dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
  } else {
    dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
  }

  ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
  defer cancel()

  conn, err := grpc.DialContext(ctx, c.Address, dialOpts...)
  if err != nil {
    return nil, err
  }

  return &PowerX{
    conn:     conn,
    Members:  orgv1.NewMemberServiceClient(conn),
    Teams:    orgv1.NewTeamServiceClient(conn),
    token:    c.Token,
    tenantID: c.TenantID,
  }, nil
}

func (p *PowerX) Close() error { return p.conn.Close() }

// Outgoing 基础 ctx：附带 auth 头（未来加拦截器时无缝）
func (p *PowerX) Outgoing(ctx context.Context) context.Context {
  if p.token == "" {
    return ctx
  }
  md := metadata.Pairs("authorization", "Bearer "+p.token, "x-powerx-tenant-id", fmt.Sprint(p.tenantID))
  return metadata.NewOutgoingContext(ctx, md)
}

// RC 便捷构造 RequestContext（当前你还未启用服务端拦截器）
func (p *PowerX) RC() *commonv1.RequestContext {
  return &commonv1.RequestContext{
    TenantId:    p.tenantID,
    AccessToken: p.token,
  }
}
```

> **使用要点**：目前你还没写拦截器，所以**请求消息里务必带 `RequestContext`**；同时我们也把 `authorization` 头加上，未来加拦截器时可以直接从头里取，无需改调用方。

---

## 3) 插件 gRPC Server（`internal/grpc/server/server.go`）

> 先提供健康检查与反射；**业务服务注册点已留好**（以后你定义好插件自己的 .proto 再注册）。

```go
package server

import (
  "context"
  "crypto/tls"
  "fmt"
  "net"

  cfgpkg "github.com/ArtisanCloud/PowerX/Core/Plugins/com.powerx.plugin.scrum/backend/internal/config"
  "github.com/ArtisanCloud/PowerX/internal/logger"

  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials"
  "google.golang.org/grpc/credentials/insecure"
  "google.golang.org/grpc/health"
  healthpb "google.golang.org/grpc/health/grpc_health_v1"
  "google.golang.org/grpc/reflection"
)

type Server struct {
  *grpc.Server
  lis net.Listener
}

func New(ctx context.Context, c cfgpkg.GRPCServer /* 传入你的 deps 也行 */) (*Server, error) {
  if !c.Enable {
    return nil, nil
  }
  lis, err := net.Listen("tcp", c.Addr)
  if err != nil {
    return nil, err
  }

  var opts []grpc.ServerOption
  if c.UseTLS {
    creds, err := credentials.NewServerTLSFromFile(c.Cert, c.Key)
    if err != nil {
      return nil, fmt.Errorf("tls: %w", err)
    }
    opts = append(opts, grpc.Creds(creds))
  } else {
    // 明确声明：开发期不加 TLS
    _ = insecure.NewCredentials()
  }

  s := grpc.NewServer(opts...)

  // 注册健康检查与反射
  healthpb.RegisterHealthServer(s, health.NewServer())
  reflection.Register(s)

  // TODO: 在这里注册你的插件 gRPC 服务
  // pluginv1.RegisterScrumPluginServiceServer(s, NewScrumServer(deps))

  return &Server{Server: s, lis: lis}, nil
}

func (s *Server) Serve(ctx context.Context) error {
  logger.From(ctx).Infof("[gRPC][plugin] listening on %s", s.lis.Addr().String())
  return s.Server.Serve(s.lis)
}
```

---

## 4) 在 `cmd/plugin/main.go` 里启动（Gin + gRPC 并行）

在你现有的入口中，加入客户端初始化与（可选）Server 启动：

```go
package main

import (
  "context"
  "errors"
  "log"
  "net/http"
  "time"

  "golang.org/x/sync/errgroup"

  cfgpkg "github.com/ArtisanCloud/PowerX/Core/Plugins/com.powerx.plugin.scrum/backend/internal/config"
  powerxgrpc "github.com/ArtisanCloud/PowerX/Core/Plugins/com.powerx.plugin.scrum/backend/internal/grpc/client"
  plugingrpc "github.com/ArtisanCloud/PowerX/Core/Plugins/com.powerx.plugin.scrum/backend/internal/grpc/server"

  "github.com/gin-gonic/gin"
)

func main() {
  ctx := context.Background()

  // 1) 读取配置（按你项目方式）
  cfg := cfgpkg.Load() // 假设你有这个方法；或自行读取 etc/config.yaml

  // 2) 初始化 PowerX gRPC 客户端
  pxc, err := powerxgrpc.NewPowerX(ctx, cfg.GRPCUpstream)
  if err != nil { log.Fatalf("grpc dial powerx: %v", err) }
  defer pxc.Close()

  // 3) HTTP（现有 Gin）
  r := gin.New()
  httpSrv := &http.Server{Addr: ":8080", Handler: r} // 你的端口

  // 示例：HTTP 路由里调用 PowerX gRPC
  r.GET("/members", func(c *gin.Context) {
    resp, err := pxc.Members.ListMembers(
      pxc.Outgoing(c.Request.Context()),
      &orgv1.ListMembersRequest{
        Ctx:  pxc.RC(),
        Page: &commonv1.PageRequest{PageSize: 20},
      },
    )
    if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
    c.JSON(200, resp)
  })

  // 4) （可选）插件 gRPC Server
  gs, err := plugingrpc.New(ctx, cfg.GRPCServer)
  if err != nil { log.Fatalf("grpc server build: %v", err) }

  // 5) 并发启动
  var g errgroup.Group
  g.Go(func() error { return httpSrv.ListenAndServe() })
  if gs != nil {
    g.Go(func() error { return gs.Serve(ctx) })
  }

  // 6) 等待退出（省略优雅关闭代码，你可按需补充）
  if err := g.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
    log.Fatalf("server exit: %v", err)
  }
}
```

---

## 5) 依赖（`go.mod`）

确保引入：

```bash
go get google.golang.org/grpc \
       google.golang.org/grpc/credentials \
       google.golang.org/grpc/health \
       google.golang.org/grpc/reflection
```

> **生成的 pb 依赖**：本指南默认你能从 **PowerX Core 仓库** 引用到生成物包：
> `github.com/ArtisanCloud/PowerX/api/grpc/gen/go/...`
> 如果插件仓库无法直接 import，可选两种方案：
>
> * **方案 A（推荐）**：在 Core 仓库给生成物打 tag/release，让插件以模块依赖方式引入；
> * **方案 B**：在插件仓库下新建 `third_party/powerx-contracts`（或 git submodule）放 `contracts/`，执行 `buf generate` 生成本地 `gen/go` 后从本地包导入。

---

## 6) 本地联调 & 验证

* 启动 **Core（:9001）** 与 **插件（HTTP :8080 / gRPC :9101）**
* `grpcurl` 验证 Core：

  ```bash
  grpcurl -plaintext localhost:9001 list
  grpcurl -plaintext -d '{"ctx":{"tenantId":42},"page":{"pageSize":10}}' \
    localhost:9001 powerx.organization.v1.MemberService.ListMembers
  ```
* 访问插件 HTTP：

  ```
  curl "http://localhost:8080/members"
  ```

---

## 7) 常见问题

* **`no matching versions for query "latest"`**
  说明插件的 `go.mod` 找不到 Core 的生成物包版本。请按“依赖”章节处理（发布 tag 或使用本地生成）。
* **包名/目录不匹配（buf lint 报错）**
  `common.v1` → `common/v1/`；`powerx.organization.v1` → `powerx/organization/v1/`。
* **请求没携带租户**
  在未启用拦截器前，**务必在请求消息里填 `RequestContext.TenantId`**，否则服务端会拒绝。

---

## 8) 下一步建议

* 给客户端加**重试/超时**与**封装 Helper**（如：`client.ListMembers(ctx, keyword, pageSize)`）
* 插件 gRPC Server 定义自己的 `.proto`（如 `plugin.scrum.v1`），在 `server.go` 注册
* 上线环境启用 **TLS/mTLS**；Core 侧建议走统一网关（Envoy / gRPC-Gateway）

---
