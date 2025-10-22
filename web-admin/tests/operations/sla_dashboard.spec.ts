import { beforeEach, afterEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

import { useSlaStore } from '../../app/stores/operations/useSlaStore'
import type { SlaProfile } from '../../app/types/operations'

describe('operations SLA store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBaseUrl: '/api/v1' } }))
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('fetches profiles', async () => {
    const store = useSlaStore()
    const profiles: SlaProfile[] = [
      {
        id: '1',
        pluginId: 'plugin',
        planType: 'real_time',
        uptimeTarget: 99.9,
        uptimeActual: 99.8,
        responseTargetMs: 500,
        responseActualMs: 450,
        successTargetPct: 99.5,
        successActualPct: 99.1,
        supportFrtTargetHours: 4,
        supportFrtActualHours: 3.5,
        slaScore: 90,
        computedAt: new Date().toISOString(),
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        incentiveAppliedAt: null,
        penaltyAppliedAt: null,
        notes: '',
      },
    ]
    const fetchSpy = vi.fn().mockResolvedValue(profiles)
    vi.stubGlobal('$fetch', fetchSpy)

    await store.fetchProfiles()

    expect(fetchSpy).toHaveBeenCalledWith('/api/v1/admin/operations/sla/profiles', { credentials: 'include' })
    expect(store.profiles).toEqual(profiles)
  })

  it('upserts targets', async () => {
    const store = useSlaStore()
    const fetchSpy = vi.fn()
    fetchSpy
      .mockResolvedValueOnce({})
      .mockResolvedValueOnce([])
    vi.stubGlobal('$fetch', fetchSpy)

    await store.upsertProfile({
      planType: 'real_time',
      targets: {
        uptimeTarget: 99.9,
        responseTargetMs: 600,
        successTargetPct: 99.5,
        supportFrtTargetHours: 4,
      },
    })

    expect(fetchSpy).toHaveBeenNthCalledWith(1, '/api/v1/admin/operations/sla/profiles', {
      method: 'POST',
      credentials: 'include',
      body: {
        planType: 'real_time',
        targets: {
          uptimeTarget: 99.9,
          responseTargetMs: 600,
          successTargetPct: 99.5,
          supportFrtTargetHours: 4,
        },
      },
    })
  })

  it('updates actual metrics', async () => {
    const store = useSlaStore()
    const fetchSpy = vi.fn()
    fetchSpy
      .mockResolvedValueOnce({})
      .mockResolvedValueOnce([])
    vi.stubGlobal('$fetch', fetchSpy)

    await store.updateActuals({
      planType: 'real_time',
      actuals: {
        uptimeActual: 99.9,
        responseActualMs: 500,
        successActualPct: 99.4,
        supportFrtActualHours: 3.5,
      },
    })

    expect(fetchSpy).toHaveBeenNthCalledWith(1, '/api/v1/admin/operations/sla/profiles/actuals', {
      method: 'PATCH',
      credentials: 'include',
      body: {
        planType: 'real_time',
        actuals: {
          uptimeActual: 99.9,
          responseActualMs: 500,
          successActualPct: 99.4,
          supportFrtActualHours: 3.5,
        },
      },
    })
    expect(fetchSpy).toHaveBeenNthCalledWith(2, '/api/v1/admin/operations/sla/profiles', { credentials: 'include' })
  })
})
