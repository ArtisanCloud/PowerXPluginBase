import { beforeEach, afterEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

import { useOperationsStore } from '../../app/stores/operations/useOperationsStore'

describe('operations support store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBaseUrl: '/api/v1' } }))
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('fetches playbook data', async () => {
    const payload = {
      channels: [{ channel: 'marketplace_ticket', address: 'https://support.local' }],
      knowledge_base: [{ label: 'FAQ', url: 'https://docs.local/faq' }],
      readiness: [],
    }
    const fetchSpy = vi.fn().mockResolvedValue({ success: true, data: payload })
    vi.stubGlobal('$fetch', fetchSpy)

    const store = useOperationsStore()
    await store.fetchPlaybook()

    expect(fetchSpy).toHaveBeenCalledWith('/api/v1/admin/operations/support/playbook', { credentials: 'include' })
    expect(store.playbook).toEqual(payload)
  })

  it('saves playbook changes', async () => {
    const savePayload = {
      channels: [],
      knowledge_base: [],
    }
    const fetchSpy = vi.fn().mockResolvedValue({ success: true, data: { ...savePayload, readiness: [] } })
    vi.stubGlobal('$fetch', fetchSpy)

    const store = useOperationsStore()
    await store.savePlaybook(savePayload)

    expect(fetchSpy).toHaveBeenCalledWith('/api/v1/admin/operations/support/playbook', {
      method: 'PUT',
      credentials: 'include',
      body: savePayload,
    })
  })
})
