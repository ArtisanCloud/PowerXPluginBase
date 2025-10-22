export type IncidentSeverity = 'sev0' | 'sev1' | 'sev2' | 'sev3' | 'sev4'
export type IncidentStatus = 'detected' | 'acknowledged' | 'mitigated' | 'monitoring' | 'resolved' | 'closed'
export type IncidentChannel = 'support_hub' | 'status_page' | 'security_email' | 'hotline' | 'webhook'

export interface ChecklistSummary {
  support_ready: boolean
  incident_ready: boolean
  sla_ready: boolean
  blocking_items: string[]
}

export interface IncidentRecord {
  id: string
  severity: IncidentSeverity
  status: IncidentStatus
  detection_source: string
  summary: string
  impact?: Record<string, any>
  mitigation?: string
  root_cause?: string
  labels?: Record<string, boolean>
  confidentiality?: string
  detected_at: string
  acknowledged_at?: string | null
  mitigated_at?: string | null
  resolved_at?: string | null
  closed_at?: string | null
  next_update_at?: string | null
}

export interface IncidentTimelineEntry {
  id: string
  incident_id: string
  entry_type: string
  message: string
  stakeholder_channel?: IncidentChannel | ''
  author_role?: string
  posted_at: string
  metadata?: Record<string, any>
}

export interface IncidentChecklistItem {
  id: string
  incident_id: string
  item_key: string
  description: string
  status: string
  completed_at?: string | null
}

export interface IncidentResponse {
  incident: IncidentRecord
  timeline: IncidentTimelineEntry[]
  checklist: IncidentChecklistItem[]
  checklist_status: ChecklistSummary
}

export interface IncidentDraftPayload {
  severity: IncidentSeverity
  detection_source: string
  summary: string
  tenant_id?: string | null
  labels?: Record<string, boolean>
  mitigation?: string
  confidentiality?: string
  impact?: Record<string, any>
  next_update_at?: string | null
}

export interface IncidentUpdatePayload {
  status?: IncidentStatus
  mitigation?: string
  root_cause?: string
  next_update_at?: string | null
  confidentiality?: string
  labels?: Record<string, boolean>
}

export interface TimelineCreatePayload {
  entry_type: string
  message: string
  stakeholder_channel?: IncidentChannel | ''
  author_role?: string
  metadata?: Record<string, any>
}
