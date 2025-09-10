package notes

import (
	"github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"
	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(rg *gin.RouterGroup, deps *app.Deps) {
	h := NewNoteHandler(deps)

	// /notes
	g := rg.Group("/notes")
	{
		g.GET("", h.GetNotes)          // 列表
		g.GET("/:id", h.GetNote)       // 详情
		g.POST("", h.CreateNote)       // 新建
		g.PUT("/:id", h.UpdateNote)    // 更新
		g.DELETE("/:id", h.DeleteNote) // 删除（软删）
	}
}
