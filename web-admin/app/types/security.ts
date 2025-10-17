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
