package middleware

import (
	authx "github.com/ArtisanCloud/PowerXPlugin/internal/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

// JWTAuth: 从 HTTP 头解析 Bearer 或 Signed-Context，注入 TenantContext 与原始 Bearer
func JWTAuth(cfg authx.JWTAuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		tc, bearer, ok := authx.ParseFromHeaders(c.GetHeader, cfg)
		if ok {
			authx.SetTenantContext(c, tc)
			authx.SetRawBearerToken(c, bearer)
			c.Next()
			return
		}
		if cfg.Optional {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
}
