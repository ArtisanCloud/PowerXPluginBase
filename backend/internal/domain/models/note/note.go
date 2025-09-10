package note

import "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models"

type Note struct {
	models.BaseModel
	Title    string `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	Content  string `gorm:"type:text;not null;comment:内容" json:"content"`
	Author   string `gorm:"type:varchar(255);not null;comment:作者" json:"author"`
	TeamID   uint64 `gorm:"not null;index;comment:团队ID"           json:"team_id"`
	MemberID uint64 `gorm:"not null;index;comment:成员ID"           json:"member_id"`
}

func (n *Note) TableName() string {
	return models.S(models.TableNote)
}
