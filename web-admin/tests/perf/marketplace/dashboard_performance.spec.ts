import { describe, it, expect, vi, beforeEach, afterEach } from "vitest"
import metricsPlugin from "~/plugins/metrics.client"

describe("Marketplace dashboard performance metrics", () => {
  const hooks: Record<string, Array<() => void>> = {}
  const nuxtApp: any = {
    hook(event: string, handler: () => void) {
      hooks[event] ||= []
      hooks[event].push(handler)
    },
    $router: {
      currentRoute: {
        value: { fullPath: "/admin/integration/marketplace/dashboard" },
      },
    },
  }

  let nowSpy: ReturnType<typeof vi.fn>

  beforeEach(() => {
    ;(globalThis as any).window = {}
    nowSpy = vi.fn()
    nowSpy.mockReturnValueOnce(100).mockReturnValueOnce(260)
    ;(globalThis as any).performance = {
      now: nowSpy,
    }
  })

  afterEach(() => {
    delete (globalThis as any).window
    delete (globalThis as any).performance
    vi.restoreAllMocks()
    for (const key of Object.keys(hooks)) {
      hooks[key] = []
    }
  })

  it("records first paint duration when dashboard route finishes", () => {
    metricsPlugin(nuxtApp as any)
    expect(hooks["page:finish"]).toBeTruthy()
    hooks["page:finish"].forEach((handler) => handler())

    const metrics = (window as any).__pxMetrics
    expect(metrics.events.length).toBeGreaterThan(0)
    expect(metrics.events[0]).toMatchObject({ name: "dashboard_first_paint" })
    expect(nowSpy).toHaveBeenCalled()
  })
})
