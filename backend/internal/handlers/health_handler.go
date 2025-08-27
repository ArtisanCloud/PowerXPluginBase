package handlers
package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/powerx-plugins/scrum/internal/contracts"
	"github.com/powerx-plugins/scrum/internal/db"
)

// HealthHandler 健康检查处理器
type HealthHandler struct{}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck 健康检查端点
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	checks := make(map[string]string)
	
	// 检查数据库连接
	if err := db.Health(); err != nil {
		checks["database"] = "unhealthy: " + err.Error()
		c.JSON(http.StatusServiceUnavailable, contracts.HealthResponse{
			Status:    "unhealthy",
			Service:   "powerx-plugin-scrum",
			Version:   "0.1.0",
			Timestamp: time.Now(),
			Checks:    checks,
		})
		return
	}
	checks["database"] = "healthy"
	
	// 所有检查通过
	c.JSON(http.StatusOK, contracts.HealthResponse{
		Status:    "healthy",
		Service:   "powerx-plugin-scrum",
		Version:   "0.1.0",
		Timestamp: time.Now(),
		Checks:    checks,
	})
}

// Ping 简单的 ping 端点
func (h *HealthHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "pong",
		"timestamp": time.Now(),
		"service":   "powerx-plugin-scrum",
	})
}