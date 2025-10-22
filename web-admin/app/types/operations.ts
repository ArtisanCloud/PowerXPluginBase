export type SlaPlanType = 'real_time' | 'transactional' | 'utility'

export interface SlaProfile {
  id: string
  pluginId: string
  planType: SlaPlanType
  uptimeTarget: number
  uptimeActual: number
  responseTargetMs: number
  responseActualMs: number
  successTargetPct: number
  successActualPct: number
  supportFrtTargetHours: number
  supportFrtActualHours: number
  slaScore: number
  incentiveAppliedAt?: string | null
  penaltyAppliedAt?: string | null
  notes?: string
  computedAt: string
  createdAt: string
  updatedAt: string
}

export interface SlaProfileUpdatePayload {
  planType: SlaPlanType
  targets: {
    uptimeTarget: number
    responseTargetMs: number
    successTargetPct: number
    supportFrtTargetHours: number
  }
}

export interface SlaActualsPayload {
  planType: SlaPlanType
  actuals: {
    uptimeActual: number
    responseActualMs: number
    successActualPct: number
    supportFrtActualHours: number
  }
}

export interface PublicSlaRecord {
  pluginId: string
  planType: SlaPlanType
  uptime: number
  responseMs: number
  successRate: number
  supportFrtHours: number
  slaScore: number
  lastUpdated: string
}
