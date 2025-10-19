package integration

import (
	"net/http"

	"github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"
	"github.com/gin-gonic/gin"
)

// Handler 暂存 integration Admin 后台的 HTTP 端点。
type Handler struct {
	deps *app.Deps
}

// NewHandler 构造 Handler。
func NewHandler(deps *app.Deps) *Handler {
	return &Handler{deps: deps}
}

// ListApprovals 列出待审批项。
func (h *Handler) ListApprovals(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "pending",
		"message": "integration approval listing not implemented",
	})
}

// Approve 通过审批。
func (h *Handler) Approve(c *gin.Context) {
	c.JSON(http.StatusAccepted, gin.H{
		"status":  "pending",
		"message": "integration approval approve endpoint not implemented",
	})
}

// Reject 拒绝审批。
func (h *Handler) Reject(c *gin.Context) {
	c.JSON(http.StatusAccepted, gin.H{
		"status":  "pending",
		"message": "integration approval reject endpoint not implemented",
	})
}

// ListGrantMatrix 返回当前策略。
func (h *Handler) ListGrantMatrix(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "pending",
		"message": "integration admin grant-matrix view not implemented",
	})
}
