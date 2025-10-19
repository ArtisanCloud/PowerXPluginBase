package integration

import (
	"github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 挂载 Integration 管理端路由。
func RegisterRoutes(rg *gin.RouterGroup, deps *app.Deps) {
	if rg == nil {
		return
	}

	handler := NewHandler(deps)
	group := rg.Group("/integration")
	{
		group.GET("/approvals", handler.ListApprovals)
		group.POST("/approvals/:id/approve", handler.Approve)
		group.POST("/approvals/:id/reject", handler.Reject)

		group.GET("/grant-matrix", handler.ListGrantMatrix)
		group.GET("/webhooks", handler.ListWebhooks)
		group.POST("/webhooks", handler.CreateWebhook)
		group.PUT("/webhooks/:id", handler.UpdateWebhook)
		group.DELETE("/webhooks/:id", handler.DeleteWebhook)
		group.GET("/webhooks/:id/attempts", handler.ListWebhookAttempts)
		group.POST("/webhooks/attempts/:attemptId/replay", handler.ReplayAttempt)
	}
}
