import { describe, it, expect, vi, beforeEach, afterEach } from "vitest"
import { flushPromises } from "@vue/test-utils"
import { mountSuspended, mockNuxtImport } from "@nuxt/test-utils/runtime"
import { createTestingPinia } from "@pinia/testing"

mockNuxtImport("useRuntimeConfig", () => () => ({
  public: {
    apiBaseUrl: "http://localhost:8086/api/v1",
  },
}))

const toastAdd = vi.fn()
mockNuxtImport("useToast", () => () => ({
  add: toastAdd,
}))

mockNuxtImport("useRoute", () => () => ({
  query: {},
}))

const uiStubs = {
  UContainer: { template: `<div class="u-container"><slot /></div>` },
  UCard: {
    template: `<div class="u-card"><slot name="header" /><div><slot /></div></div>`,
  },
  UButton: {
    props: ["loading"],
    emits: ["click"],
    template: `<button type="button" @click="$emit('click')"><slot /></button>`
  },
  UInput: {
    props: ["modelValue", "label", "type", "placeholder"],
    emits: ["update:modelValue"],
    template: `<label class="u-input"><span class="sr-only">{{ label }}</span><input :type="type || 'text'" :placeholder="placeholder" :value="modelValue" @input="$emit('update:modelValue', $event.target.value)" /></label>`
  },
  USelectMenu: {
    props: ["modelValue", "options"],
    emits: ["update:modelValue"],
    template: `<select :value="modelValue" @change="$emit('update:modelValue', $event.target.value)"><option v-for="option in options" :key="option.value" :value="option.value">{{ option.label || option.value }}</option></select>`
  },
  UAlert: {
    template: `<div class="u-alert"><slot name="title" /><slot /></div>`
  },
  UBadge: {
    template: `<span class="u-badge"><slot /></span>`
  },
  UTable: {
    props: ["rows", "columns"],
    template: `<table><thead><tr><th v-for="column in columns" :key="column.key">{{ column.label }}</th></tr></thead><tbody><tr v-for="row in rows" :key="JSON.stringify(row)"><td v-for="column in columns" :key="column.key">{{ row[column.key] }}</td></tr></tbody></table>`
  },
  UIcon: {
    props: ["name"],
    template: `<span :data-icon="name"><slot /></span>`
  },
  USkeleton: {
    template: `<div class="u-skeleton"></div>`
  },
}

describe("Marketplace Analytics Dashboard", () => {
  const DashboardPage = () => import("~/pages/_p/com.powerx.plugins.base/admin/integration/marketplace/dashboard.vue")
  let fetchSpy: ReturnType<typeof vi.fn>

  beforeEach(() => {
    fetchSpy = vi.fn(async (url: string) => {
      if (url.includes("/usage/tenants/")) {
        return {
          data: {
            series: [
              {
                timestamp: new Date().toISOString(),
                metric: "calls",
                value: 250,
                revenue: 270,
                quota_remaining: 50,
                window: "day",
                currency: "USD",
              },
            ],
            alerts: [
              {
                code: "usage_spike",
                severity: "warning",
                message: "Usage spike detected",
              },
            ],
          },
        }
      }
      if (url.includes("/revenue-share/reports")) {
        return {
          data: [
            {
              id: "report-1",
              tenant_id: "tenant-1",
              vendor_id: "vendor-1",
              period_start: "2025-01-01",
              period_end: "2025-01-31",
              gross_amount: 270,
              vendor_share: 216,
              platform_share: 40.5,
              fees: 13.5,
              currency: "USD",
              status: "ready",
              generated_at: new Date().toISOString(),
            },
          ],
        }
      }
      return { data: {} }
    })
    ;(globalThis as any).$fetch = fetchSpy
  })

  afterEach(() => {
    vi.restoreAllMocks()
    toastAdd.mockReset()
    delete (globalThis as any).$fetch
  })

  it("loads metrics and displays alerts", async () => {
    const wrapper = await mountSuspended(await DashboardPage(), {
      global: {
        stubs: uiStubs,
        plugins: [createTestingPinia()],
      },
    })

    const tenantInput = wrapper.find("input[placeholder='tenant-123']")
    const licenseInput = wrapper.find("input[placeholder='license-abc']")

    await tenantInput.setValue("tenant-1")
    await licenseInput.setValue("license-1")

    const loadButton = wrapper.find("button")
    await loadButton.trigger("click")
    await flushPromises()

    expect(fetchSpy).toHaveBeenCalledWith(
      expect.stringContaining("/usage/tenants/tenant-1/licenses/license-1/metrics"),
      expect.any(Object)
    )
    expect(fetchSpy).toHaveBeenCalledWith(expect.stringContaining("/revenue-share/reports"), expect.any(Object))

    expect(wrapper.text()).toContain("Usage spike detected")
    expect(wrapper.text()).toContain("Vendor")
    expect(wrapper.text()).toContain("270.00 USD")
  })
})
