package services

import (
	"context"
	"fmt"

	powerxclient "scrum-plugin/internal/grpc/client"
)

// DemoService Demo 服务
type DemoService struct {
	powerxClient *powerxclient.PowerX
}

// NewDemoService 创建 Demo 服务
func NewDemoService(powerxClient *powerxclient.PowerX) *DemoService {
	return &DemoService{
		powerxClient: powerxClient,
	}
}

// HealthCheckRequest 健康检查请求
type HealthCheckRequest struct{}

// HealthCheckResponse 健康检查响应
type HealthCheckResponse struct {
	Status    string `json:"status"`
	Connected bool   `json:"connected"`
	TenantID  int64  `json:"tenant_id"`
	HasToken  bool   `json:"has_token"`
}

// ListMembersRequest 获取成员列表请求
type ListMembersRequest struct {
	Keyword   string `json:"keyword"`
	PageIndex int32  `json:"page_index"`
	PageSize  int32  `json:"page_size"`
}

// ListMembersResponse 获取成员列表响应
type ListMembersResponse struct {
	Members  interface{} `json:"members"`
	GRPCInfo GRPCInfo    `json:"grpc"`
}

// ListTeamsRequest 获取团队列表请求
type ListTeamsRequest struct {
	Keyword   string `json:"keyword"`
	PageIndex int32  `json:"page_index"`
	PageSize  int32  `json:"page_size"`
}

// ListTeamsResponse 获取团队列表响应
type ListTeamsResponse struct {
	Teams    interface{} `json:"teams"`
	GRPCInfo GRPCInfo    `json:"grpc"`
}

// GetMemberRequest 获取单个成员请求
type GetMemberRequest struct {
	ID int64 `json:"id"`
}

// GetMemberResponse 获取单个成员响应
type GetMemberResponse struct {
	Member      interface{} `json:"member"`
	RequestedID string      `json:"requested_id"`
}

// GetTeamRequest 获取单个团队请求
type GetTeamRequest struct {
	ID int64 `json:"id"`
}

// GetTeamResponse 获取单个团队响应
type GetTeamResponse struct {
	Team        interface{} `json:"team"`
	RequestedID string      `json:"requested_id"`
}

// DebugInfoResponse 调试信息响应
type DebugInfoResponse struct {
	GRPCConnection GRPCInfo `json:"grpc_connection"`
	Endpoints      []string `json:"endpoints"`
	Note           string   `json:"note"`
}

// GRPCInfo gRPC 连接信息
type GRPCInfo struct {
	Connected bool  `json:"connected"`
	TenantID  int64 `json:"tenant_id"`
	HasToken  bool  `json:"has_token"`
}

// HealthCheck 健康检查
func (s *DemoService) HealthCheck(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error) {
	if err := s.powerxClient.HealthCheck(ctx); err != nil {
		return nil, fmt.Errorf("PowerX gRPC service unavailable: %w", err)
	}

	return &HealthCheckResponse{
		Status:    "ok",
		Connected: s.powerxClient.IsConnected(),
		TenantID:  s.powerxClient.GetTenantID(),
		HasToken:  s.powerxClient.GetToken() != "",
	}, nil
}

// ListMembers 获取成员列表
func (s *DemoService) ListMembers(ctx context.Context, req *ListMembersRequest) (*ListMembersResponse, error) {
	grpcReq := &powerxclient.ListMembersRequest{
		Ctx: s.powerxClient.RC(),
		Page: &powerxclient.PageRequest{
			PageIndex: req.PageIndex,
			PageSize:  req.PageSize,
		},
		Keyword: req.Keyword,
	}

	resp, err := s.powerxClient.ListMembers(ctx, grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call PowerX gRPC service: %w", err)
	}

	return &ListMembersResponse{
		Members: resp,
		GRPCInfo: GRPCInfo{
			Connected: s.powerxClient.IsConnected(),
			TenantID:  s.powerxClient.GetTenantID(),
			HasToken:  s.powerxClient.GetToken() != "",
		},
	}, nil
}

// ListTeams 获取团队列表
func (s *DemoService) ListTeams(ctx context.Context, req *ListTeamsRequest) (*ListTeamsResponse, error) {
	grpcReq := &powerxclient.ListTeamsRequest{
		Ctx: s.powerxClient.RC(),
		Page: &powerxclient.PageRequest{
			PageIndex: req.PageIndex,
			PageSize:  req.PageSize,
		},
		Keyword: req.Keyword,
	}

	resp, err := s.powerxClient.ListTeams(ctx, grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call PowerX gRPC service: %w", err)
	}

	return &ListTeamsResponse{
		Teams: resp,
		GRPCInfo: GRPCInfo{
			Connected: s.powerxClient.IsConnected(),
			TenantID:  s.powerxClient.GetTenantID(),
			HasToken:  s.powerxClient.GetToken() != "",
		},
	}, nil
}

// GetMember 获取单个成员
func (s *DemoService) GetMember(ctx context.Context, req *GetMemberRequest, requestedID string) (*GetMemberResponse, error) {
	grpcReq := &powerxclient.GetMemberRequest{
		Ctx: s.powerxClient.RC(),
		Id:  req.ID,
	}

	resp, err := s.powerxClient.GetMember(ctx, grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call PowerX gRPC service: %w", err)
	}

	return &GetMemberResponse{
		Member:      resp,
		RequestedID: requestedID,
	}, nil
}

// GetTeam 获取单个团队
func (s *DemoService) GetTeam(ctx context.Context, req *GetTeamRequest, requestedID string) (*GetTeamResponse, error) {
	grpcReq := &powerxclient.GetTeamRequest{
		Ctx: s.powerxClient.RC(),
		Id:  req.ID,
	}

	resp, err := s.powerxClient.GetTeam(ctx, grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call PowerX gRPC service: %w", err)
	}

	return &GetTeamResponse{
		Team:        resp,
		RequestedID: requestedID,
	}, nil
}

// GetDebugInfo 获取调试信息
func (s *DemoService) GetDebugInfo(ctx context.Context) *DebugInfoResponse {
	return &DebugInfoResponse{
		GRPCConnection: GRPCInfo{
			Connected: s.powerxClient.IsConnected(),
			TenantID:  s.powerxClient.GetTenantID(),
			HasToken:  s.powerxClient.GetToken() != "",
		},
		Endpoints: []string{
			"GET /api/v1/demo/grpc/health - 检查 gRPC 连接状态",
			"GET /api/v1/demo/grpc/members - 获取成员列表",
			"GET /api/v1/demo/grpc/members/{id} - 获取单个成员",
			"GET /api/v1/demo/grpc/teams - 获取团队列表",
			"GET /api/v1/demo/grpc/teams/{id} - 获取单个团队",
			"GET /api/v1/demo/grpc/debug - 查看调试信息",
		},
		Note: "当前使用模拟数据，可以通过这些接口测试 gRPC 连接和功能",
	}
}
