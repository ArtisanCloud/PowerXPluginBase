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
