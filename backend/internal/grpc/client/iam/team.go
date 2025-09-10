package iam

import (
	"context"
	"github.com/ArtisanCloud/PowerXPlugin/internal/grpc/client"
)

type ListTeamsRequest struct {
	Ctx      *client.RequestContext `json:"ctx"`
	Page     *client.PageRequest    `json:"page"`
	Keyword  string                 `json:"keyword,omitempty"`
	Statuses []string               `json:"statuses,omitempty"`
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

type GetTeamRequest struct {
	Ctx *client.RequestContext `json:"ctx"`
	Id  int64                  `json:"id"`
}

type GetTeamResponse struct {
	Team *Team `json:"team"`
}

type TeamServiceClient struct {
	pxClient *client.PowerXServiceClient
}

func NewTeamServiceClient(pxClient *client.PowerXServiceClient) *TeamServiceClient {
	return &TeamServiceClient{pxClient: pxClient}
}

// ListTeams 获取团队列表
func (p *TeamServiceClient) ListTeams(ctx context.Context, req *ListTeamsRequest) (*ListTeamsResponse, error) {
	resp := &ListTeamsResponse{}
	err := p.pxClient.InvokeGRPC(ctx, "powerx.organization.v1.TeamService", "ListTeams", req, resp)
	return resp, err
}

// GetTeam 获取单个团队
func (p *TeamServiceClient) GetTeam(ctx context.Context, req *GetTeamRequest) (*GetTeamResponse, error) {
	resp := &GetTeamResponse{}
	err := p.pxClient.InvokeGRPC(ctx, "powerx.organization.v1.TeamService", "GetTeam", req, resp)
	return resp, err
}
