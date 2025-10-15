package runtime_ops

import "time"

// MCPSession models the MCP connection for a plugin instance.
type MCPSession struct {
	ID                string
	RuntimeAssignment string
	TenantID          string
	State             string
	JWTID             string
	CapabilitiesHash  string
	MissedHeartbeats  int
	LastPingAt        time.Time
	ClosedAt          *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
