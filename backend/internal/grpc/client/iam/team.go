package iam

import (
    "context"
    "time"

    commonv1 "github.com/ArtisanCloud/PowerX/api/grpc/gen/go/common/v1"
    iamv1 "github.com/ArtisanCloud/PowerX/api/grpc/gen/go/powerx/iam/v1"

    "github.com/ArtisanCloud/PowerXPlugin/internal/grpc/client"
    utils "github.com/ArtisanCloud/PowerXPlugin/internal/shared/utils"
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
    var pageIndex, pageSize int32
    if req.Page != nil { pageIndex = req.Page.PageIndex; pageSize = req.Page.PageSize }
    pr := &iamv1.ListTeamsRequest{
        Ctx:     &commonv1.RequestContext{TenantId: req.Ctx.TenantId, AccessToken: req.Ctx.AccessToken},
        Keyword: req.Keyword,
        Page:    &commonv1.PageRequest{PageSize: pageSize, Offset: pageIndex * pageSize},
    }
    cli, err := p.pxClient.TeamsClient(ctx)
    if err != nil { return nil, err }
    ctx = p.pxClient.Outgoing(ctx)
    presp, err := cli.ListTeams(ctx, pr)
    if err != nil { return nil, err }
    // 映射
    out := &ListTeamsResponse{PageIndex: pageIndex, PageSize: pageSize}
    if presp != nil && presp.Data != nil {
        for _, t := range presp.Data.GetItems() {
            out.Teams = append(out.Teams, convertTeam(t))
        }
        if presp.Data.Page != nil { out.TotalCount = presp.Data.Page.GetTotal() }
    }
    return out, nil
}

// GetTeam 获取单个团队
func (p *TeamServiceClient) GetTeam(ctx context.Context, req *GetTeamRequest) (*GetTeamResponse, error) {
    pr := &iamv1.GetTeamRequest{Ctx: &commonv1.RequestContext{TenantId: req.Ctx.TenantId, AccessToken: req.Ctx.AccessToken}}
    if req.Id > 0 { pr.Selector = &iamv1.GetTeamRequest_Id{Id: uint64(req.Id)} }
    cli, err := p.pxClient.TeamsClient(ctx)
    if err != nil { return nil, err }
    ctx = p.pxClient.Outgoing(ctx)
    presp, err := cli.GetTeam(ctx, pr)
    if err != nil { return nil, err }
    if presp == nil || presp.Data == nil || presp.Data.Team == nil { return &GetTeamResponse{Team: nil}, nil }
    return &GetTeamResponse{Team: convertTeam(presp.Data.Team)}, nil
}

// —— 强类型 -> 本地简化类型转换 ——
func convertTeam(src *iamv1.Team) *Team {
    if src == nil { return nil }
    var id, leaderID int64
    if src.Ref != nil { id = utils.ToInt64(any(src.Ref.GetId())) }
    // 注意：上游 proto 可能已移除 Leader 字段，若需 leaderId 请在新字段中补取
    // 这里保持兼容，默认置 0。
    var createdAt, updatedAt string
    if src.CreatedAt != nil {
        t := time.Unix(src.CreatedAt.Seconds, int64(src.CreatedAt.Nanos))
        createdAt = t.UTC().Format(time.RFC3339)
    }
    if src.UpdatedAt != nil {
        t := time.Unix(src.UpdatedAt.Seconds, int64(src.UpdatedAt.Nanos))
        updatedAt = t.UTC().Format(time.RFC3339)
    }
    return &Team{
        Id:          id,
        Name:        src.GetName(),
        Description: src.GetDescription(),
        LeaderId:    leaderID,
        // 部分版本 Team 可能不再包含 Status 字段，这里保持空串兼容
        Status:      "",
        MemberCount: int32(src.GetMemberCount()),
        CreatedAt:   createdAt,
        UpdatedAt:   updatedAt,
    }
}
