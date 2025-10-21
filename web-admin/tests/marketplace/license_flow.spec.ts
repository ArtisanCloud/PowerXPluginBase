import { describe, it, expect, vi, beforeEach, afterEach } from "vitest"
import { flushPromises } from "@vue/test-utils"
import { mountSuspended, mockNuxtImport } from "@nuxt/test-utils/runtime"

mockNuxtImport("useRuntimeConfig", () => () => ({
  public: {
    apiBaseUrl: "http://localhost:8086/api/v1",
  },
}))

mockNuxtImport("useToast", () => () => ({
  add: vi.fn(),
}))

const uiStubs = {
  UContainer: { template: "<div><slot /></div>" },
  UCard: {
    template: `<div class="u-card"><slot name="header" /><div><slot /></div></div>`,
  },
  UButton: {
    props: ["color", "loading"],
    emits: ["click"],
    template: `<button type="button" @click="$emit('click', $event)"><slot /></button>`
  },
  UInput: {
    props: ["modelValue", "type", "min", "max", "step"],
    emits: ["update:modelValue"],
    template: `<input :type="type || 'text'" :value="modelValue" @input="$emit('update:modelValue', $event.target.value)" />`
  },
  UTextarea: {
    props: ["modelValue", "rows"],
    emits: ["update:modelValue"],
    template: `<textarea :rows="rows" :value="modelValue" @input="$emit('update:modelValue', $event.target.value)"></textarea>`
  },
  USelectMenu: {
    props: ["modelValue", "options"],
    emits: ["update:modelValue"],
    template: `<select :value="modelValue" @change="$emit('update:modelValue', $event.target.value)"><option v-for="option in options" :key="option.value" :value="option.value">{{ option.label }}</option></select>`
  },
  USelect: {
    props: ["modelValue", "options"],
    emits: ["update:modelValue"],
    template: `<select :value="modelValue" @change="$emit('update:modelValue', $event.target.value)"><option v-for="option in options" :key="option.value" :value="option.value">{{ option.label }}</option></select>`
  },
  UFormGroup: {
    template: `<label><slot /> <slot name="help" /></label>`
  },
  UTable: {
    props: ["rows"],
    emits: ["select"],
    template: `<table><tbody><tr v-for="row in rows" :key="row.id" @click="$emit('select', row)"><slot name="title-data" :row="row">{{ row.title }}</slot></tr></tbody></table>`
  },
  UBadge: {
    template: `<span class="badge"><slot /></span>`
  },
  UIcon: {
    props: ["name"],
    template: `<span :data-icon="name"><slot /></span>`
  },
  UAlert: {
    template: `<div class="u-alert"><slot /></div>`
  },
}

describe("Marketplace Plans Page", () => {
  const PlansPage = () => import("~/pages/_p/com.powerx.plugins.base/admin/integration/marketplace/plans.vue")
  let fetchSpy: ReturnType<typeof vi.fn>

  beforeEach(() => {
    fetchSpy = vi.fn(async (url: string, options?: any) => {
      if (url.endsWith("/admin/marketplace/listings") && !options?.method) {
        return {
          data: {
            items: [
              {
                id: "listing-1",
                plugin_id: "plugin.demo",
                vendor_id: "vendor-1",
                status: "draft",
                title: "Demo Listing",
                slug: "demo",
                ready_checklist_score: 80,
              },
            ],
          },
        }
      }
      if (url.endsWith("/admin/marketplace/listings/listing-1") && !options?.method) {
        return {
          data: {
            id: "listing-1",
            plugin_id: "plugin.demo",
            vendor_id: "vendor-1",
            status: "draft",
            title: "Demo Listing",
            slug: "demo",
            ready_checklist_score: 80,
            pricing_plans: [
              {
                id: "plan-1",
                listing_id: "listing-1",
                plan_code: "standard",
                plan_type: "subscription",
                currency: "USD",
                amount: 29.99,
                billing_period: "monthly",
                trial_period_days: 14,
                quota_limit: 1000,
                overage_policy: "billable",
                feature_matrix: { seats: 5 },
                is_default: true,
                created_at: new Date().toISOString(),
                updated_at: new Date().toISOString(),
                tiers: [],
              },
            ],
          },
        }
      }
      if (url.endsWith("/admin/marketplace/listings/listing-1") && options?.method === "PATCH") {
        return { message: "ok" }
      }
      return { data: {} }
    })
    ;(globalThis as any).$fetch = fetchSpy
  })

  afterEach(() => {
    vi.restoreAllMocks()
    delete (globalThis as any).$fetch
  })

  it("submits pricing plan payload", async () => {
    const wrapper = await mountSuspended(await PlansPage(), {
      global: { stubs: uiStubs },
    })
    await flushPromises()

    const saveButton = wrapper.findAll("button").find((btn) => btn.text().includes("保存变更"))
    expect(saveButton).toBeTruthy()
    await saveButton!.trigger("click")
    await flushPromises()

    const patchCall = fetchSpy.mock.calls.find(([url, options]) =>
      url.endsWith("/admin/marketplace/listings/listing-1") && options?.method === "PATCH"
    )
    expect(patchCall).toBeTruthy()
    expect(patchCall?.[1].body.pricing_plans[0]).toMatchObject({ plan_code: "standard" })
  })
})

describe("Marketplace Purchase Page", () => {
  const PurchasePage = () => import("~/pages/_p/com.powerx.plugins.base/admin/integration/marketplace/purchase.vue")
  let fetchSpy: ReturnType<typeof vi.fn>

  beforeEach(() => {
    fetchSpy = vi.fn(async (url: string, options?: any) => {
      if (url.endsWith("/admin/marketplace/listings") && !options?.method) {
        return {
          data: {
            items: [
              {
                id: "listing-1",
                plugin_id: "plugin.demo",
                vendor_id: "vendor-1",
                status: "published",
                title: "Demo Listing",
                slug: "demo",
                ready_checklist_score: 95,
              },
            ],
          },
        }
      }
      if (url.endsWith("/admin/marketplace/listings/listing-1") && !options?.method) {
        return {
          data: {
            id: "listing-1",
            plugin_id: "plugin.demo",
            vendor_id: "vendor-1",
            status: "published",
            title: "Demo Listing",
            slug: "demo",
            ready_checklist_score: 95,
            pricing_plans: [
              {
                id: "plan-1",
                listing_id: "listing-1",
                plan_code: "enterprise",
                plan_type: "subscription",
                currency: "USD",
                amount: 59.99,
                billing_period: "monthly",
                trial_period_days: 7,
                quota_limit: 2000,
                overage_policy: "throttle",
                feature_matrix: { seats: 10 },
                is_default: true,
                created_at: new Date().toISOString(),
                updated_at: new Date().toISOString(),
                tiers: [],
              },
            ],
          },
        }
      }
      if (url.endsWith("/marketplace/licenses") && options?.method === "POST") {
        return {
          data: {
            id: "license-1",
            listing_id: "listing-1",
            plan_id: "plan-1",
            status: "active",
            expires_at: new Date().toISOString(),
            token: "token-demo",
            renewal_token: "renew-1",
            settlement_currency: "USD",
          },
        }
      }
      if (url.includes("/marketplace/licenses/license-1") && !options?.method) {
        return {
          data: {
            id: "license-1",
            listing_id: "listing-1",
            plan_id: "plan-1",
            status: "active",
            expires_at: new Date().toISOString(),
            token: "token-demo",
            renewal_token: "renew-1",
            settlement_currency: "USD",
          },
        }
      }
      return { data: {} }
    })
    ;(globalThis as any).$fetch = fetchSpy
  })

  afterEach(() => {
    vi.restoreAllMocks()
    delete (globalThis as any).$fetch
  })

  it("sends purchase payload", async () => {
    const wrapper = await mountSuspended(await PurchasePage(), {
      global: { stubs: uiStubs },
    })
    await flushPromises()

    const purchaseButton = wrapper.findAll("button").find((btn) => btn.text().includes("立即购买"))
    expect(purchaseButton).toBeTruthy()
    await purchaseButton!.trigger("click")
    await flushPromises()

    const postCall = fetchSpy.mock.calls.find(([url, options]) =>
      url.endsWith("/marketplace/licenses") && options?.method === "POST"
    )
    expect(postCall).toBeTruthy()
    expect(postCall?.[1].body).toMatchObject({
      listing_id: "listing-1",
      plan_id: "plan-1",
      payment_intent_id: "pi_mock_123",
    })
  })
})
