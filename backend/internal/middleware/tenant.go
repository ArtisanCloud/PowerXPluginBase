package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/powerx-plugins/scrum/internal/config"
	"github.com/powerx-plugins/scrum/internal/logger"
)

// TenantContext 租户上下文
type TenantContext struct {
	TenantID    int64     `json:"tenant_id"`
	UserID      int64     `json:"user_id,omitempty"`
	Roles       []string  `json:"roles,omitempty"`
	Permissions []string  `json:"permissions,omitempty"`
	IssuedAt    time.Time `json:"iat,omitempty"`
	ExpiresAt   time.Time `json:"exp,omitempty"`
}

// PowerXContext PowerX 上下文头部信息
type PowerXContext struct {
	TenantID    int64    `json:"tenant_id"`
	UserID      int64    `json:"user_id,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	Timestamp   int64    `json:"timestamp"`
}

const (
	// HTTP 头部常量
	HeaderPowerXContext = "X-PowerX-CTX"
	HeaderPowerXJWT     = "X-PowerX-CTX-JWT"

	// 上下文键
	ContextKeyTenant = "tenant"
	ContextKeyUser   = "user"
)

// TenantMiddleware 创建租户中间件
func TenantMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开发模式旁路
		if cfg.DevMode {
			handleDevMode(c)
			return
		}

		var tenantCtx *TenantContext
		var err error

		// 根据配置选择验证方式
		if cfg.IsJWTMode() {
			tenantCtx, err = validateJWT(c, cfg)
		} else if cfg.IsHMACMode() {
			tenantCtx, err = validateHMAC(c, cfg)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "No authentication method configured",
			})
			c.Abort()
			return
		}

		if err != nil {
			logger.AuthMiddleware().WithError(err).Warn("Authentication failed")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "Invalid authentication",
				"detail": err.Error(),
			})
			c.Abort()
			return
		}

		if tenantCtx == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing tenant context",
			})
			c.Abort()
			return
		}

		// 将租户上下文存储到 Gin 上下文
		c.Set(ContextKeyTenant, tenantCtx)

		// 记录请求日志
		logger.AuthMiddleware().WithFields(logger.Fields{
			"tenant_id": tenantCtx.TenantID,
			"user_id":   tenantCtx.UserID,
			"path":      c.Request.URL.Path,
			"method":    c.Request.Method,
		}).Info("Authenticated request")

		c.Next()
	}
}

// handleDevMode 处理开发模式
func handleDevMode(c *gin.Context) {
	// 开发模式下，可以通过查询参数或头部注入租户 ID
	var tenantID int64 = 1 // 默认租户 ID

	if tid := c.Query("tenant_id"); tid != "" {
		if parsed, err := strconv.ParseInt(tid, 10, 64); err == nil {
			tenantID = parsed
		}
	}

	if tid := c.GetHeader("X-Dev-Tenant-ID"); tid != "" {
		if parsed, err := strconv.ParseInt(tid, 10, 64); err == nil {
			tenantID = parsed
		}
	}

	tenantCtx := &TenantContext{
		TenantID:    tenantID,
		UserID:      1, // 默认用户 ID
		Roles:       []string{"admin"},
		Permissions: []string{"*"},
	}

	c.Set(ContextKeyTenant, tenantCtx)

	logger.AuthMiddleware().WithField("tenant_id", tenantID).Debug("Dev mode: using mock tenant context")
	c.Next()
}

// validateJWT 验证 JWT token
func validateJWT(c *gin.Context, cfg *config.Config) (*TenantContext, error) {
	tokenString := c.GetHeader(HeaderPowerXJWT)
	if tokenString == "" {
		return nil, fmt.Errorf("missing JWT token")
	}

	// 解析 JWT（这里简化实现，实际应该从 JWKS 获取公钥）
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// TODO: 从 JWKS URL 获取公钥
		// 这里需要实现 JWKS 客户端
		return nil, fmt.Errorf("JWKS validation not implemented yet")
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tenantCtx := &TenantContext{}

		if tid, ok := claims["tenant_id"].(float64); ok {
			tenantCtx.TenantID = int64(tid)
		}

		if uid, ok := claims["user_id"].(float64); ok {
			tenantCtx.UserID = int64(uid)
		}

		if roles, ok := claims["roles"].([]interface{}); ok {
			for _, role := range roles {
				if roleStr, ok := role.(string); ok {
					tenantCtx.Roles = append(tenantCtx.Roles, roleStr)
				}
			}
		}

		return tenantCtx, nil
	}

	return nil, fmt.Errorf("invalid JWT claims")
}

// validateHMAC 验证 HMAC 签名
func validateHMAC(c *gin.Context, cfg *config.Config) (*TenantContext, error) {
	contextHeader := c.GetHeader(HeaderPowerXContext)
	if contextHeader == "" {
		return nil, fmt.Errorf("missing context header")
	}

	// 解码 Base64
	contextData, err := base64.StdEncoding.DecodeString(contextHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode context: %w", err)
	}

	// 解析 JSON
	var pxCtx PowerXContext
	if err := json.Unmarshal(contextData, &pxCtx); err != nil {
		return nil, fmt.Errorf("failed to parse context: %w", err)
	}

	// 验证时间戳（防重放）
	now := time.Now().Unix()
	if abs(now-pxCtx.Timestamp) > int64(cfg.Context.TTL.Seconds()) {
		return nil, fmt.Errorf("context expired")
	}

	// 生成预期的签名

	// 从头部获取实际签名（简化实现，实际可能在单独头部）
	// TODO: 实现完整的 HMAC 签名验证逻辑

	return &TenantContext{
		TenantID:    pxCtx.TenantID,
		UserID:      pxCtx.UserID,
		Roles:       pxCtx.Roles,
		Permissions: pxCtx.Permissions,
	}, nil
}

// generateHMACSignature 生成 HMAC 签名
func generateHMACSignature(data []byte, secret string) string {
	secretBytes, _ := base64.StdEncoding.DecodeString(secret)
	h := hmac.New(sha256.New, secretBytes)
	h.Write(data)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// verifyHMACSignature 验证 HMAC 签名
func verifyHMACSignature(data []byte, signature, secret string) bool {
	expected := generateHMACSignature(data, secret)
	return subtle.ConstantTimeCompare([]byte(signature), []byte(expected)) == 1
}

// abs 返回整数的绝对值
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// GetTenantContext 从 Gin 上下文获取租户上下文
func GetTenantContext(c *gin.Context) (*TenantContext, bool) {
	if ctx, exists := c.Get(ContextKeyTenant); exists {
		if tenantCtx, ok := ctx.(*TenantContext); ok {
			return tenantCtx, true
		}
	}
	return nil, false
}

// RequireTenantContext 要求租户上下文的中间件
func RequireTenantContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := GetTenantContext(c); !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Tenant context required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetTenantID 从上下文获取租户 ID
func GetTenantID(c *gin.Context) (int64, error) {
	tenantCtx, exists := GetTenantContext(c)
	if !exists {
		return 0, fmt.Errorf("tenant context not found")
	}
	return tenantCtx.TenantID, nil
}

// GetUserID 从上下文获取用户 ID
func GetUserID(c *gin.Context) (int64, error) {
	tenantCtx, exists := GetTenantContext(c)
	if !exists {
		return 0, fmt.Errorf("tenant context not found")
	}
	return tenantCtx.UserID, nil
}
