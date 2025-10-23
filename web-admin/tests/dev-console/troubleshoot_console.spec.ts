import { beforeEach, afterEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

import { useDevConsoleTroubleshootStore } from '../../app/stores/dev-console/troubleshoot'

describe('dev console troubleshoot store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.useFakeTimers()
    vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBaseUrl: '/api/v1' } }))
  })

  afterEach(() => {
    vi.useRealTimers()
    vi.unstubAllGlobals()
    vi.clearAllMocks()
  })

  it('retries a failed job run and prepends the new run', async () => {
    const retryResponse = {
      success: true,
      data: {
        id: 'run-2',
        retry_of: 'run-1',
        status: 'pending',
        job_type: 'webhook_replay',
        trigger_source: 'manual',
        created_at: '2025-01-01T00:00:00Z',
      },
    }
    const fetchSpy = vi.fn().mockResolvedValue(retryResponse)
    vi.stubGlobal('$fetch', fetchSpy)

    const store = useDevConsoleTroubleshootStore()
    store.runs = [
      {
        id: 'run-1',
        status: 'failed',
        job_type: 'webhook_replay',
        trigger_source: 'manual',
        created_at: '2025-01-01T00:00:00Z',
      },
    ] as any

    await store.retryRun('run-1')

    expect(fetchSpy).toHaveBeenCalledWith('/api/v1/admin/dev-console/jobs/runs/run-1/retry', {
      method: 'POST',
      credentials: 'include',
      body: {},
    })
    expect(store.runs[0]?.id).toBe('run-2')
    expect(store.lastRetry?.id).toBe('run-2')
  })

  it('auto refreshes troubleshooting summary based on interval', async () => {
    const summaryResponse = {
      success: true,
      data: {
        refreshed_at: '2025-01-01T00:00:00Z',
        refresh_interval_seconds: 30,
        health: [],
        quota: [],
        webhook_delivery: { success_rate: 1, retry_rate: 0, dlq_rate: 0, recent_failures: [] },
        guidance: [],
      },
    }
    const fetchSpy = vi.fn().mockResolvedValue(summaryResponse)
    vi.stubGlobal('$fetch', fetchSpy)

    const store = useDevConsoleTroubleshootStore()
    await store.fetchSummary()

    expect(fetchSpy).toHaveBeenCalledTimes(1)

    await vi.runOnlyPendingTimersAsync()

    expect(fetchSpy).toHaveBeenCalledTimes(2)
    store.stopAutoRefresh()
  })
})
