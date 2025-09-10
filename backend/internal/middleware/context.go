package middleware

// internal/middleware/context.go

import "github.com/gin-gonic/gin"

type TenantContext struct {
	TenantID      int64    `json:"tenant_id"`
	UserID        int64    `json:"user_id"`
	Roles         []string `json:"roles"`
	Permissions   []string `json:"permissions"`
	PolicyVersion string   `json:"policy_version"`
}

const (
	ctxKeyTenant = "tenant_ctx"
	ctxKeyToken  = "raw_bearer_token"
)

func SetTenantContext(c *gin.Context, tc TenantContext) { c.Set(ctxKeyTenant, tc) }
func GetTenantContext(c *gin.Context) (TenantContext, bool) {
	v, ok := c.Get(ctxKeyTenant)
	if !ok || v == nil {
		return TenantContext{}, false
	}
	tc, ok := v.(TenantContext)
	return tc, ok
}
func SetRawBearerToken(c *gin.Context, token string) {
	if token != "" {
		c.Set(ctxKeyToken, token)
	}
}
func GetRawBearerToken(c *gin.Context) (string, bool) {
	v, ok := c.Get(ctxKeyToken)
	if !ok || v == nil {
		return "", false
	}
	s, ok := v.(string)
	return s, ok && s != ""
}
