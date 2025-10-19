package integration

import (
	"strings"

	authx "github.com/ArtisanCloud/PowerXPlugin/internal/middleware"
)

// RBACEntries 返回 Integration 管理端的权限映射。
func RBACEntries(prefix string) map[string]authx.Permission {
	base := strings.TrimRight(prefix, "/") + "/admin/integration"
	return map[string]authx.Permission{
		"GET:" + base + "/approvals":              {Resource: "integration.approvals", Action: "read"},
		"POST:" + base + "/approvals/:id/approve": {Resource: "integration.approvals", Action: "manage"},
		"POST:" + base + "/approvals/:id/reject":  {Resource: "integration.approvals", Action: "manage"},
		"GET:" + base + "/grant-matrix":           {Resource: "integration.grant_matrix", Action: "read"},
	}
}
