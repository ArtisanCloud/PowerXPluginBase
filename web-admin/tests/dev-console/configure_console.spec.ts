import { beforeEach, afterEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

import { useDevConsoleConfigStore } from '~/stores/dev-console/config'

describe('dev console config store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBaseUrl: '/api/v1' } }))
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('loads configuration sections', async () => {
    const payload = {
      sections: [
        {
          key: 'admin_console.retention',
          title: 'Retention',
          fields: [],
          current_values: {},
        },
      ],
    }
    const fetchSpy = vi.fn().mockResolvedValue({ success: true, data: payload })
    vi.stubGlobal('$fetch', fetchSpy)

    const store = useDevConsoleConfigStore()
    await store.fetchSections()

    expect(fetchSpy).toHaveBeenCalledWith('/api/v1/admin/dev-console/config/sections', { credentials: 'include' })
    expect(store.sections).toEqual(payload.sections)
  })

  it('updates a section', async () => {
    const updateResponse = {
      success: true,
      data: {
        key: 'admin_console.retention',
        title: 'Retention',
        fields: [],
        current_values: { audit_retention_days: 200 },
      },
    }
    const fetchSpy = vi.fn().mockResolvedValue(updateResponse)
    vi.stubGlobal('$fetch', fetchSpy)

    const store = useDevConsoleConfigStore()
    await store.updateSection('admin_console.retention', { values: { audit_retention_days: 200 } })

    expect(fetchSpy).toHaveBeenCalledWith('/api/v1/admin/dev-console/config/sections/admin_console.retention', {
      method: 'PUT',
      credentials: 'include',
      body: { values: { audit_retention_days: 200 } },
    })
    expect(store.sectionByKey('admin_console.retention')?.current_values).toEqual({ audit_retention_days: 200 })
  })
})
