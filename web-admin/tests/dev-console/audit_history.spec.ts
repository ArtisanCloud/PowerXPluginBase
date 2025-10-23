import { beforeEach, afterEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

import { useDevConsoleAuditStore } from '../../app/stores/dev-console/audit'

describe('dev console audit store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBaseUrl: '/api/v1' } }))
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('fetches audit events', async () => {
    const data = {
      events: [
        {
          id: 'evt-1',
          action: 'config.section.update',
          permission_code: 'operations.plugin.admin',
          resource_type: 'config.section',
          occurred_at: '2025-10-22T10:00:00Z',
          actor: { id: 'user:1', name: 'Jane Doe' },
        },
      ],
      next_cursor: 'cursor-2',
    }
    const fetchSpy = vi.fn().mockResolvedValue({ success: true, data })
    vi.stubGlobal('$fetch', fetchSpy)
    vi.stubGlobal('atob', (str: string) => Buffer.from(str, 'base64').toString('binary'))

    const store = useDevConsoleAuditStore()
    await store.fetchEvents()

    expect(fetchSpy).toHaveBeenCalledWith('/api/v1/admin/dev-console/audit/events', { credentials: 'include' })
    expect(store.events).toEqual(data.events)
    expect(store.nextCursor).toBe('cursor-2')
  })

  it('exports in selected format', async () => {
    const payload = {
      success: true,
      data: {
        filename: 'audit.json',
        content_type: 'application/json',
        content_base64: Buffer.from('test', 'utf8').toString('base64'),
      },
    }
    const fetchSpy = vi.fn().mockResolvedValue(payload)
    vi.stubGlobal('$fetch', fetchSpy)
    vi.stubGlobal('atob', (str: string) => Buffer.from(str, 'base64').toString('binary'))

    const appendSpy = vi.fn()
    const removeSpy = vi.fn()
    const clickMock = vi.fn()
    const fakeAnchor = {
      href: '',
      download: '',
      click: clickMock,
    } as unknown as HTMLAnchorElement
    vi.stubGlobal('document', {
      createElement: vi.fn().mockReturnValue(fakeAnchor),
      body: {
        appendChild: appendSpy,
        removeChild: removeSpy,
      },
    })
    const originalURL = global.URL
    vi.stubGlobal('URL', { createObjectURL: vi.fn().mockReturnValue('blob://test'), revokeObjectURL: vi.fn() })

    const store = useDevConsoleAuditStore()
    await store.exportEvents('json')

    expect(fetchSpy).toHaveBeenCalled()
    expect(appendSpy).toHaveBeenCalledWith(fakeAnchor)
    expect(clickMock).toHaveBeenCalled()
    expect(removeSpy).toHaveBeenCalledWith(fakeAnchor)

    global.URL = originalURL
  })
})
