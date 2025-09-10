package iam

import (
	"context"

	"github.com/ArtisanCloud/PowerXPlugin/internal/grpc/client"
)

// ===== 请求 / 响应模型（与 team.go 风格一致） =====

type ListMembersRequest struct {
	Ctx      *client.RequestContext `json:"ctx"`
	Page     *client.PageRequest    `json:"page"`
	Keyword  string                 `json:"keyword,omitempty"`
	TeamIds  []int64                `json:"team_ids,omitempty"`
	Statuses []string               `json:"statuses,omitempty"`
}

type Member struct {
	Id          int64   `json:"id"`
	Username    string  `json:"username"`
	DisplayName string  `json:"display_name"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	Status      string  `json:"status"`
	TeamIds     []int64 `json:"team_ids,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type ListMembersResponse struct {
	Members    []*Member `json:"members"`
	TotalCount int64     `json:"total_count"`
	PageIndex  int32     `json:"page_index"`
	PageSize   int32     `json:"page_size"`
}

type GetMemberRequest struct {
	Ctx *client.RequestContext `json:"ctx"`
	Id  int64                  `json:"id"`
}

type GetMemberResponse struct {
	Member *Member `json:"member"`
}

// ===== 客户端（与 TeamServiceClient 同款薄包装） =====

type MemberServiceClient struct {
	pxClient *client.PowerXServiceClient
}

func NewMemberServiceClient(pxClient *client.PowerXServiceClient) *MemberServiceClient {
	return &MemberServiceClient{pxClient: pxClient}
}

func (p *MemberServiceClient) ListMembers(ctx context.Context, req *ListMembersRequest) (*ListMembersResponse, error) {
	resp := &ListMembersResponse{}
	err := p.pxClient.InvokeGRPC(ctx, "powerx.organization.v1.MemberService", "ListMembers", req, resp)
	return resp, err
}

func (p *MemberServiceClient) GetMember(ctx context.Context, req *GetMemberRequest) (*GetMemberResponse, error) {
	resp := &GetMemberResponse{}
	err := p.pxClient.InvokeGRPC(ctx, "powerx.organization.v1.MemberService", "GetMember", req, resp)
	return resp, err
}
