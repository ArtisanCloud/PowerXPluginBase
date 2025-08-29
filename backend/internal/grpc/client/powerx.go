package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	cfgpkg "scrum-plugin/internal/config"
	"scrum-plugin/internal/logger"

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

type ListMembersRequest struct {
	Ctx      *RequestContext `json:"ctx"`
	Page     *PageRequest    `json:"page"`
	Keyword  string          `json:"keyword,omitempty"`
	TeamIds  []int64         `json:"team_ids,omitempty"`
	Statuses []string        `json:"statuses,omitempty"`
}

type Member struct {
	Id         int64  `json:"id"`
	UserId     int64  `json:"user_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Position   string `json:"position"`
	Department string `json:"department"`
	Status     string `json:"status"`
	JoinedAt   string `json:"joined_at"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type ListMembersResponse struct {
	Members    []*Member `json:"members"`
	TotalCount int64     `json:"total_count"`
	PageIndex  int32     `json:"page_index"`
	PageSize   int32     `json:"page_size"`
}

type ListTeamsRequest struct {
	Ctx      *RequestContext `json:"ctx"`
	Page     *PageRequest    `json:"page"`
	Keyword  string          `json:"keyword,omitempty"`
	Statuses []string        `json:"statuses,omitempty"`
}

type Team struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	LeaderId    int64  `json:"leader_id"`
	Status      string `json:"status"`
	MemberCount int32  `json:"member_count"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ListTeamsResponse struct {
	Teams      []*Team `json:"teams"`
	TotalCount int64   `json:"total_count"`
	PageIndex  int32   `json:"page_index"`
	PageSize   int32   `json:"page_size"`
}

type GetMemberRequest struct {
	Ctx *RequestContext `json:"ctx"`
	Id  int64           `json:"id"`
}

type GetMemberResponse struct {
	Member *Member `json:"member"`
}

type GetTeamRequest struct {
	Ctx *RequestContext `json:"ctx"`
	Id  int64           `json:"id"`
}

type GetTeamResponse struct {
	Team *Team `json:"team"`
}

// PowerX PowerX gRPC 客户端封装
type PowerX struct {
	conn *grpc.ClientConn
	// TODO: 当有实际的 PowerX proto 文件时，添加客户端
	// Members orgv1.MemberServiceClient
	// Teams   orgv1.TeamServiceClient

	token    string
	tenantID int64
	cfg      cfgpkg.GRPCUpstream
}

// NewPowerX 根据配置拨号 PowerX gRPC
func NewPowerX(ctx context.Context, c cfgpkg.GRPCUpstream) (*PowerX, error) {
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

	return &PowerX{
		conn: conn,
		// TODO: 当有实际的 PowerX proto 文件时，初始化客户端
		// Members:  orgv1.NewMemberServiceClient(conn),
		// Teams:    orgv1.NewTeamServiceClient(conn),
		token:    c.Token,
		tenantID: c.TenantID,
		cfg:      c,
	}, nil
}

// Close 关闭 gRPC 连接
func (p *PowerX) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

// Outgoing 基础 ctx：附带 auth 头（未来加拦截器时无缝）
func (p *PowerX) Outgoing(ctx context.Context) context.Context {
	md := metadata.New(map[string]string{})

	if p.token != "" {
		md.Set("authorization", "Bearer "+p.token)
	}

	if p.tenantID > 0 {
		md.Set("x-powerx-tenant-id", fmt.Sprint(p.tenantID))
	}

	return metadata.NewOutgoingContext(ctx, md)
}

// GetToken 获取认证 token
func (p *PowerX) GetToken() string {
	return p.token
}

// GetTenantID 获取租户 ID
func (p *PowerX) GetTenantID() int64 {
	return p.tenantID
}

// IsConnected 检查连接状态
func (p *PowerX) IsConnected() bool {
	if p.conn == nil {
		return false
	}
	return p.conn.GetState().String() == "READY"
}

// TODO: 当有实际的 PowerX proto 文件时，添加以下方法：

// RC 便捷构造 RequestContext
func (p *PowerX) RC() *RequestContext {
	return &RequestContext{
		TenantId:    p.tenantID,
		AccessToken: p.token,
	}
}

// ListMembers 获取成员列表
func (p *PowerX) ListMembers(ctx context.Context, req *ListMembersRequest) (*ListMembersResponse, error) {
	resp := &ListMembersResponse{}
	err := p.invokeGRPC(ctx, "powerx.organization.v1.MemberService", "ListMembers", req, resp)
	return resp, err
}

// GetMember 获取单个成员
func (p *PowerX) GetMember(ctx context.Context, req *GetMemberRequest) (*GetMemberResponse, error) {
	resp := &GetMemberResponse{}
	err := p.invokeGRPC(ctx, "powerx.organization.v1.MemberService", "GetMember", req, resp)
	return resp, err
}

// ListTeams 获取团队列表
func (p *PowerX) ListTeams(ctx context.Context, req *ListTeamsRequest) (*ListTeamsResponse, error) {
	resp := &ListTeamsResponse{}
	err := p.invokeGRPC(ctx, "powerx.organization.v1.TeamService", "ListTeams", req, resp)
	return resp, err
}

// GetTeam 获取单个团队
func (p *PowerX) GetTeam(ctx context.Context, req *GetTeamRequest) (*GetTeamResponse, error) {
	resp := &GetTeamResponse{}
	err := p.invokeGRPC(ctx, "powerx.organization.v1.TeamService", "GetTeam", req, resp)
	return resp, err
}

// invokeGRPC 通用 gRPC 调用方法
func (p *PowerX) invokeGRPC(ctx context.Context, service, method string, req, resp interface{}) error {
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

	// 模拟响应（实际应该从 gRPC 服务获取）
	switch method {
	case "ListMembers":
		if listResp, ok := resp.(*ListMembersResponse); ok {
			*listResp = *p.mockListMembersResponse()
		}
	case "GetMember":
		if getResp, ok := resp.(*GetMemberResponse); ok {
			*getResp = *p.mockGetMemberResponse()
		}
	case "ListTeams":
		if listResp, ok := resp.(*ListTeamsResponse); ok {
			*listResp = *p.mockListTeamsResponse()
		}
	case "GetTeam":
		if getResp, ok := resp.(*GetTeamResponse); ok {
			*getResp = *p.mockGetTeamResponse()
		}
	default:
		return fmt.Errorf("unsupported method: %s", method)
	}

	return nil
}

// 模拟数据方法（供测试使用）
func (p *PowerX) mockListMembersResponse() *ListMembersResponse {
	return &ListMembersResponse{
		Members: []*Member{
			{
				Id:         1,
				UserId:     101,
				Name:       "Alice Johnson",
				Email:      "alice@example.com",
				Phone:      "+1234567890",
				Position:   "Software Engineer",
				Department: "Engineering",
				Status:     "active",
				JoinedAt:   "2023-01-15T00:00:00Z",
				CreatedAt:  "2023-01-15T10:00:00Z",
				UpdatedAt:  "2024-01-15T10:00:00Z",
			},
			{
				Id:         2,
				UserId:     102,
				Name:       "Bob Smith",
				Email:      "bob@example.com",
				Phone:      "+1234567891",
				Position:   "Product Manager",
				Department: "Product",
				Status:     "active",
				JoinedAt:   "2023-02-01T00:00:00Z",
				CreatedAt:  "2023-02-01T10:00:00Z",
				UpdatedAt:  "2024-01-15T10:00:00Z",
			},
		},
		TotalCount: 2,
		PageIndex:  0,
		PageSize:   20,
	}
}

func (p *PowerX) mockGetMemberResponse() *GetMemberResponse {
	return &GetMemberResponse{
		Member: &Member{
			Id:         1,
			UserId:     101,
			Name:       "Alice Johnson",
			Email:      "alice@example.com",
			Phone:      "+1234567890",
			Position:   "Software Engineer",
			Department: "Engineering",
			Status:     "active",
			JoinedAt:   "2023-01-15T00:00:00Z",
			CreatedAt:  "2023-01-15T10:00:00Z",
			UpdatedAt:  "2024-01-15T10:00:00Z",
		},
	}
}

func (p *PowerX) mockListTeamsResponse() *ListTeamsResponse {
	return &ListTeamsResponse{
		Teams: []*Team{
			{
				Id:          1,
				Name:        "Development Team",
				Description: "Core development team",
				LeaderId:    101,
				Status:      "active",
				MemberCount: 5,
				CreatedAt:   "2023-01-01T10:00:00Z",
				UpdatedAt:   "2024-01-15T10:00:00Z",
			},
			{
				Id:          2,
				Name:        "Product Team",
				Description: "Product management team",
				LeaderId:    102,
				Status:      "active",
				MemberCount: 3,
				CreatedAt:   "2023-01-01T10:00:00Z",
				UpdatedAt:   "2024-01-15T10:00:00Z",
			},
		},
		TotalCount: 2,
		PageIndex:  0,
		PageSize:   20,
	}
}

func (p *PowerX) mockGetTeamResponse() *GetTeamResponse {
	return &GetTeamResponse{
		Team: &Team{
			Id:          1,
			Name:        "Development Team",
			Description: "Core development team",
			LeaderId:    101,
			Status:      "active",
			MemberCount: 5,
			CreatedAt:   "2023-01-01T10:00:00Z",
			UpdatedAt:   "2024-01-15T10:00:00Z",
		},
	}
}

// HealthCheck 健康检查方法
func (p *PowerX) HealthCheck(ctx context.Context) error {
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
func (p *PowerX) Reconnect(ctx context.Context) error {
	if p.conn != nil {
		p.conn.Close()
	}

	newClient, err := NewPowerX(ctx, p.cfg)
	if err != nil {
		return err
	}

	*p = *newClient
	return nil
}
