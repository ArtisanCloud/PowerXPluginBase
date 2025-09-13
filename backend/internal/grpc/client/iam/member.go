package iam

import (
    "context"
    "time"

    commonv1 "github.com/ArtisanCloud/PowerX/api/grpc/gen/go/common/v1"
    iamv1 "github.com/ArtisanCloud/PowerX/api/grpc/gen/go/powerx/iam/v1"

    "github.com/ArtisanCloud/PowerXPlugin/internal/grpc/client"
    utils "github.com/ArtisanCloud/PowerXPlugin/internal/shared/utils"
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
    // 转为 proto 请求
    var pageIndex, pageSize int32
    if req.Page != nil { pageIndex = req.Page.PageIndex; pageSize = req.Page.PageSize }
    pr := &iamv1.ListMembersRequest{
        Ctx:     &commonv1.RequestContext{TenantId: req.Ctx.TenantId, AccessToken: req.Ctx.AccessToken},
        Keyword: req.Keyword,
        Page:    &commonv1.PageRequest{PageSize: pageSize, Offset: pageIndex * pageSize},
    }
    // 仅取首个 teamId（proto 只支持单个）
    if len(req.TeamIds) > 0 && req.TeamIds[0] > 0 {
        pr.TeamId = uint64(req.TeamIds[0])
    }
    // 仅取首个 status（proto 单值枚举）
    if len(req.Statuses) > 0 {
        pr.Status = parseMemberStatus(req.Statuses[0])
    }

    cli, err := p.pxClient.MembersClient(ctx)
    if err != nil { return nil, err }
    // 携带认证与租户信息
    ctx = p.pxClient.Outgoing(ctx)
    presp, err := cli.ListMembers(ctx, pr)
    if err != nil { return nil, err }

    // 组装响应
    var members []*Member
    if presp != nil && presp.Data != nil {
        for _, m := range presp.Data.GetItems() {
            members = append(members, convertMember(m))
        }
    }
    total := int64(0)
    if presp != nil && presp.Data != nil && presp.Data.Page != nil { total = presp.Data.Page.GetTotal() }
    return &ListMembersResponse{Members: members, TotalCount: total, PageIndex: pageIndex, PageSize: pageSize}, nil
}

func (p *MemberServiceClient) GetMember(ctx context.Context, req *GetMemberRequest) (*GetMemberResponse, error) {
    pr := &iamv1.GetMemberRequest{Ctx: &commonv1.RequestContext{TenantId: req.Ctx.TenantId, AccessToken: req.Ctx.AccessToken}}
    if req.Id > 0 {
        pr.Selector = &iamv1.GetMemberRequest_Id{Id: uint64(req.Id)}
    }
    cli, err := p.pxClient.MembersClient(ctx)
    if err != nil { return nil, err }
    ctx = p.pxClient.Outgoing(ctx)
    presp, err := cli.GetMember(ctx, pr)
    if err != nil { return nil, err }
    if presp == nil || presp.Data == nil || presp.Data.Member == nil { return &GetMemberResponse{Member: nil}, nil }
    return &GetMemberResponse{Member: convertMember(presp.Data.Member)}, nil
}

// —— 辅助 —— //
func pickMap(m map[string]any, key string) map[string]any {
	if v, ok := m[key]; ok {
		if mv, ok2 := v.(map[string]any); ok2 {
			return mv
		}
	}
	return nil
}

// 已将散落的转换工具迁移到 internal/shared/utils/jsonx.go

// —— 强类型 -> 本地简化类型转换 ——
func convertMember(src *iamv1.Member) *Member {
    if src == nil { return nil }
    var id int64
    if src.Ref != nil { id = utils.ToInt64(any(src.Ref.GetId())) }
    // 时间格式化为 ISO8601
    var createdAt, updatedAt string
    if src.CreatedAt != nil {
        t := time.Unix(src.CreatedAt.Seconds, int64(src.CreatedAt.Nanos))
        createdAt = t.UTC().Format(time.RFC3339)
    }
    if src.UpdatedAt != nil {
        t := time.Unix(src.UpdatedAt.Seconds, int64(src.UpdatedAt.Nanos))
        updatedAt = t.UTC().Format(time.RFC3339)
    }
    // team ids
    var teamIDs []int64
    for _, v := range src.TeamIds { teamIDs = append(teamIDs, int64(v)) }
    // status -> 字符串
    status := src.Status.String()
    return &Member{
        Id:          id,
        Username:    src.GetUsername(),
        DisplayName: src.GetDisplayName(),
        Email:       src.GetEmail(),
        Phone:       src.GetPhone(),
        Status:      status,
        TeamIds:     teamIDs,
        CreatedAt:   createdAt,
        UpdatedAt:   updatedAt,
    }
}

func parseMemberStatus(s string) iamv1.MemberStatus {
    switch s {
    case "ACTIVE", "MEMBER_STATUS_ACTIVE":
        return iamv1.MemberStatus_MEMBER_STATUS_ACTIVE
    case "INACTIVE", "MEMBER_STATUS_INACTIVE":
        return iamv1.MemberStatus_MEMBER_STATUS_INACTIVE
    case "SUSPENDED", "MEMBER_STATUS_SUSPENDED":
        return iamv1.MemberStatus_MEMBER_STATUS_SUSPENDED
    default:
        return iamv1.MemberStatus_MEMBER_STATUS_UNSPECIFIED
    }
}
