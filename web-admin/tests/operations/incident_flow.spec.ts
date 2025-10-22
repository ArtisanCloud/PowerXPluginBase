import { beforeEach, afterEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

import { useIncidentsStore } from '../../app/stores/operations/useIncidentsStore'
import type { IncidentResponse } from '../../app/types/operations'

describe('operations incidents store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBaseUrl: '/api/v1' } }))
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('fetches incident list', async () => {
    const store = useIncidentsStore()
    const fetchSpy = vi.fn().mockResolvedValue({ success: true, data: [{ id: 'inc-1', severity: 'sev1', status: 'detected', detection_source: 'monitoring', summary: 'API down', detected_at: new Date().toISOString() }] })
    vi.stubGlobal('$fetch', fetchSpy)

    await store.fetchIncidents()

    expect(fetchSpy).toHaveBeenCalledWith('/api/v1/admin/operations/incidents', { credentials: 'include', query: {} })
    expect(store.incidents).toHaveLength(1)
  })

  it('creates an incident and refreshes list', async () => {
    const store = useIncidentsStore()
    const response: IncidentResponse = {
      incident: {
        id: 'inc-1',
        severity: 'sev1',
        status: 'detected',
        detection_source: 'monitoring',
        summary: 'New incident',
        detected_at: new Date().toISOString(),
        labels: {},
      },
      timeline: [],
      checklist: [],
      checklist_status: {
        support_ready: false,
        incident_ready: true,
        sla_ready: false,
        blocking_items: [],
      },
    }
    const fetchSpy = vi.fn()
    fetchSpy
      .mockResolvedValueOnce({ success: true, data: response })
      .mockResolvedValueOnce({ success: true, data: [] })
    vi.stubGlobal('$fetch', fetchSpy)

    await store.createIncident({ severity: 'sev1', detection_source: 'monitoring', summary: 'New incident' })

    expect(fetchSpy).toHaveBeenNthCalledWith(1, '/api/v1/admin/operations/incidents', {
      method: 'POST',
      credentials: 'include',
      body: { severity: 'sev1', detection_source: 'monitoring', summary: 'New incident' },
    })
    expect(fetchSpy).toHaveBeenNthCalledWith(2, '/api/v1/admin/operations/incidents', { credentials: 'include', query: {} })
    expect(store.selected?.incident.id).toBe('inc-1')
  })

  it('appends timeline entry and reloads incident', async () => {
    const store = useIncidentsStore()
    vi.stubGlobal('$fetch', vi.fn().mockResolvedValue({ success: true, data: [] }))

    await store.appendTimeline('inc-1', { entry_type: 'announcement', message: 'ack', stakeholder_channel: 'support_hub' })

    expect($fetch).toHaveBeenNthCalledWith(1, '/api/v1/admin/operations/incidents/inc-1/timeline', {
      method: 'POST',
      credentials: 'include',
      body: { entry_type: 'announcement', message: 'ack', stakeholder_channel: 'support_hub' },
    })
    expect($fetch).toHaveBeenNthCalledWith(2, '/api/v1/admin/operations/incidents/inc-1', { credentials: 'include' })
  })
})
