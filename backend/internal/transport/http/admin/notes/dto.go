package notes

// Note 领域 HTTP DTO（仅含请求体/查询参数）

type NoteListRequest struct {
    Page     int    `form:"page,default=1"`
    PageSize int    `form:"page_size,default=20"`
    Q        string `form:"q"`
    TeamID   uint64 `form:"team_id,default=0"`
    MemberID uint64 `form:"member_id,default=0"`
}

type CreateNoteRequest struct {
    Title    string `json:"title"    binding:"required"`
    Content  string `json:"content"  binding:"required"`
    Author   string `json:"author"   binding:"required"`
    TeamID   uint64 `json:"team_id"  binding:"required"`
    MemberID uint64 `json:"member_id" binding:"required"`
}

type UpdateNoteRequest struct {
    Title    string `json:"title"    binding:"required"`
    Content  string `json:"content"  binding:"required"`
    Author   string `json:"author"   binding:"required"`
    TeamID   uint64 `json:"team_id"  binding:"required"`
    MemberID uint64 `json:"member_id" binding:"required"`
}

