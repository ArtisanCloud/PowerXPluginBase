export interface ConsentTokenResponse {
  id: string
  tenantId: string
  token: string
  scope: string[]
  status: string
  expiresAt?: string
  issuedAt: string
  issuedBy: string
  revokedAt?: string
  revokedReason?: string
}

export interface ConsentTokenListResponse {
  data: ConsentTokenResponse[]
}

export interface LifecycleEventResponse {
  id: string
  tenantId: string
  eventType: string
  assetKey: string
  status: string
  occurredAt: string
  recordedBy: string
  payload?: any
}

export interface LifecycleEventListResponse {
  data: LifecycleEventResponse[]
}

export interface AuditReport {
  id: string
  baseline_id: string
  initiated_by: string
  status: string
  findings?: Record<string, any>
  artifact_path?: string
  sarif_path?: string
  report_hash?: string
  checklist_version: string
  created_at: string
}

export interface AuditReportListResponse {
  data: AuditReport[]
}
