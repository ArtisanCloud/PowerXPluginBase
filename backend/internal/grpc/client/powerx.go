package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	cfgpkg "github.com/ArtisanCloud/PowerXPlugin/internal/config"
	"github.com/ArtisanCloud/PowerXPlugin/internal/logger"
	"io/ioutil"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// 通用请求和响应结构（与 PowerX proto 兼容）
type RequestContext struct {
	TenantId    int64  `json:"tenant_id"`
	AccessToken string `json:"access_token"`
}

type PageRequest struct {
	PageIndex int32 `json:"page_index"`
	PageSize  int32 `json:"page_size"`
}

// PowerX PowerX gRPC 客户端封装
type PowerXServiceClient struct {
    conn *grpc.ClientConn
	// TODO: 当有实际的 PowerX proto 文件时，添加客户端
	// Members orgv1.MemberServiceClient
	// Teams   orgv1.TeamServiceClient

    token    string
    tenantID int64
    cfg      *cfgpkg.GRPCUpstream
    tm       *TokenManager
}

// NewPowerX 根据配置拨号 PowerX gRPC
func NewPowerXServiceClient(ctx context.Context, c *cfgpkg.GRPCUpstream) (*PowerXServiceClient, error) {
	if c.Address == "" {
		return nil, fmt.Errorf("grpc upstream address is required")
	}

	dialOpts := []grpc.DialOption{grpc.WithBlock()}

	if c.UseTLS {
		var creds credentials.TransportCredentials
		if c.CACert != "" {
			pem, err := ioutil.ReadFile(c.CACert)
			if err != nil {
				return nil, fmt.Errorf("read ca cert: %w", err)
			}
			cp := x509.NewCertPool()
			if !cp.AppendCertsFromPEM(pem) {
				return nil, fmt.Errorf("failed to append ca cert")
			}
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

	logger.WithField("address", c.Address).Info("Connecting to PowerX gRPC service")

    conn, err := grpc.DialContext(ctx, c.Address, dialOpts...)
    if err != nil {
        return nil, fmt.Errorf("failed to dial PowerX gRPC: %w", err)
    }

	logger.Info("Successfully connected to PowerX gRPC service")

    p := &PowerXServiceClient{
        conn: conn,
		// TODO: 当有实际的 PowerX proto 文件时，初始化客户端
		// Members:  orgv1.NewMemberServiceClient(conn),
		// Teams:    orgv1.NewTeamServiceClient(conn),
        token:    c.Token,
        tenantID: c.TenantID,
        cfg:      c,
    }

    // 初始化 STS TokenManager（若提供了 client_id/secret）
    if c.STSClientID != "" && c.STSClientSecret != "" {
        p.tm = NewTokenManager(c.STSClientID, c.STSClientSecret, c.STSAudience, c.STSScope, c.STSTTL, func(ctx context.Context, req *STSExchangeRequest) (*STSExchangeResponse, error) {
            // 使用通用调用与约定的服务/方法名
            var resp STSExchangeResponse
            // 直接调用，不附带旧的 Authorization 头
            // 这里复用 InvokeGRPC 的日志和错误语义；未来替换为真实 proto 客户端
            if err := p.InvokeGRPC(ctx, "powerx.auth.sts.v1.STSService", "Exchange", req, &resp); err != nil {
                return nil, err
            }
            return &resp, nil
        })
    }

    return p, nil
}

// Close 关闭 gRPC 连接
func (p *PowerXServiceClient) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

func (p *PowerXServiceClient) RC() *RequestContext {
	return &RequestContext{
		TenantId:    p.tenantID,
		AccessToken: p.token,
	}
}

// Outgoing 基础 ctx：附带 auth 头（未来加拦截器时无缝）
func (p *PowerXServiceClient) Outgoing(ctx context.Context) context.Context {
    md := metadata.New(map[string]string{})

    // 优先使用 STS token，其次使用静态 Token
    bearer := p.token
    if p.tm != nil {
        if tok, err := p.tm.GetToken(ctx); err == nil && tok != "" {
            bearer = tok
        }
    }
    if bearer != "" {
        md.Set("authorization", "Bearer "+bearer)
    }

	if p.tenantID > 0 {
		md.Set("x-powerx-tenant-id", fmt.Sprint(p.tenantID))
	}

	return metadata.NewOutgoingContext(ctx, md)
}

// GetToken 获取认证 token
func (p *PowerXServiceClient) GetToken() string {
    if p.token != "" {
        return p.token
    }
    if p.tm != nil && p.tm.HasValid() {
        // 为避免阻塞，这里不触发刷新，只返回空串代表未知
        // 上层可通过 HealthCheck/实际调用来驱动刷新
        return "sts"
    }
    return ""
}

// HasToken 是否配置/具备可用的访问凭据（静态或临时）
func (p *PowerXServiceClient) HasToken() bool {
    if p.token != "" {
        return true
    }
    if p.tm != nil && p.tm.HasValid() {
        return true
    }
    return false
}

// GetTenantID 获取租户 ID
func (p *PowerXServiceClient) GetTenantID() int64 {
	return p.tenantID
}

// IsConnected 检查连接状态
func (p *PowerXServiceClient) IsConnected() bool {
	if p.conn == nil {
		return false
	}
	return p.conn.GetState().String() == "READY"
}

// invokeGRPC 通用 gRPC 调用方法
func (p *PowerXServiceClient) InvokeGRPC(ctx context.Context, service, method string, req, resp interface{}) error {
	if p.conn == nil {
		return fmt.Errorf("gRPC connection is not established")
	}

	// 添加认证头部
	ctx = p.Outgoing(ctx)

	// 将请求序列化为 JSON（简化版本，实际应使用 protobuf）
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// 记录调用日志
	logger.WithField("service", service).WithField("method", method).Info("Calling PowerX gRPC service")

	// 这里应该使用实际的 gRPC 调用
	// 目前返回模拟数据供测试
	_ = reqBytes

	return nil
}

// HealthCheck 健康检查方法
func (p *PowerXServiceClient) HealthCheck(ctx context.Context) error {
	if p.conn == nil {
		return fmt.Errorf("grpc connection is nil")
	}

	state := p.conn.GetState()
	if state.String() != "READY" && state.String() != "IDLE" {
		return fmt.Errorf("grpc connection state is %s", state.String())
	}

	return nil
}

// Reconnect 重新连接
func (p *PowerXServiceClient) Reconnect(ctx context.Context) error {
	if p.conn != nil {
		p.conn.Close()
	}

	newClient, err := NewPowerXServiceClient(ctx, p.cfg)
	if err != nil {
		return err
	}

	*p = *newClient
	return nil
}
