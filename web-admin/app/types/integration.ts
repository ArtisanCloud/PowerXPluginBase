export interface IntegrationEnvelopePreview {
  messageId: string
  toolScope: string
  tenantId: string
  issuedAt: string
}

export interface IntegrationApprovalStatus {
  id: string
  target: string
  status: 'PENDING' | 'APPROVED' | 'REJECTED'
  submittedBy: string
  submittedAt: string
}

export interface IntegrationWebhookSubscription {
  id: string
  tenant_id: string
  event_type: string
  target_url: string
  retry_policy?: number[]
  status: string
  metadata?: Record<string, any>
  created_at: string
  updated_at: string
}

export interface IntegrationWebhookAttempt {
  id: string
  subscription_id: string
  status: string
  retry_count: number
  last_error?: string
  next_delivery_at?: string | null
  payload_snapshot?: Record<string, any> | null
  created_at: string
  updated_at: string
}

export interface IntegrationSecret {
  id: string
  tenant_id: string
  integration_type: string
  current_secret_ref?: string
  pending_secret_ref?: string
  rotation_interval_days: number
  last_rotated_at?: string | null
  next_rotation_due_at?: string | null
  status: string
  metadata?: Record<string, any>
  created_at: string
  updated_at: string
}

export interface IntegrationSecretAuditEntry {
  action: string
  actor: string
  timestamp: string
  details?: Record<string, any>
}
