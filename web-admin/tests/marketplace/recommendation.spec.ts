import { describe, it, expect, vi, beforeEach, afterEach } from "vitest"
import { mountSuspended, mockNuxtImport } from "@nuxt/test-utils/runtime"

mockNuxtImport("useRuntimeConfig", () => () => ({
  public: {
    apiBaseUrl: "http://localhost:8086/api/v1",
  },
}))

const toastAdd = vi.fn()
mockNuxtImport("useToast", () => () => ({ add: toastAdd }))

const uiStubs = {
  UContainer: { template: "<div><slot /></div>" },
  UCard: { template: `<div class=\"u-card\"><slot name=header /><div><slot /></div></div>` },
  UProgress: { template: "<div class=progress></div>" },
  UButton: { props: ["loading"], emits: ["click"], template: `<button type=\"button\" @click=\"$emit('click')\"><slot /></button>` },
  UInput: {
    props: ["modelValue", "type"],
    emits: ["update:modelValue"],
    template: `<input :type=\"type || 'text'\" :value=\"modelValue\" @input=\"$emit('update:modelValue', $event.target.value)\" />`,
  },
  UFormGroup: { template: `<label><slot /><slot name=help /></label>` },
  UTable: { props: ["rows"], template: `<table><tbody><tr v-for=\"row in rows\" :key=\"row.id\"><slot name=title-data :row=\"row\">{{ row.title }}</slot></tr></tbody></table>` },
  UBadge: { template: `<span><slot /></span>` },
  UIcon: { props: ["name"], template: `<span :data-icon=\"name\"><slot /></span>` },
}

describe("Marketplace Recommendation Page", () => {
  let fetchSpy: ReturnType<typeof vi.fn>

  beforeEach(() => {
    fetchSpy = vi.fn(async (url: string, options?: any) => {
      if (url.endsWith("/admin/marketplace/recommendation/config")) {
        return {
          config: {
            enabled: true,
            default_weight: 0.5,
            experiment_topic: "topic-a",
            frequency_minutes: 60,
          },
          top_listings: [
            {
              id: "listing-1",
              title: "Listing One",
              status: "published",
              plugin_id: "com.powerx.plugins.example",
              vendor_id: "vendor-1",
              recommended_weight: 0.75,
            },
          ],
        }
      }
      if (url.endsWith("/admin/marketplace/recommendation/sync") && options?.method === "POST") {
        return { updated: 3 }
      }
      if (url.endsWith("/admin/marketplace/recommendation/experiment") && options?.method === "PATCH") {
        return { default_weight: options.body?.default_weight }
      }
      return {}
    })

    mockNuxtImport("useNuxtApp", () => () => ({
      $fetch: fetchSpy,
    }))
    ;(globalThis as any).$fetch = fetchSpy
  })

  afterEach(() => {
    vi.restoreAllMocks()
    toastAdd.mockReset()
  })

  it("loads configuration on mount", async () => {
    const page = await import("~/pages/_p/com.powerx.plugins.base/admin/integration/marketplace/recommendation.vue")
    const wrapper = await mountSuspended(page.default, {
      global: { stubs: uiStubs },
    })

    expect(fetchSpy).toHaveBeenCalledWith(
      "http://localhost:8086/api/v1/admin/marketplace/recommendation/config",
      undefined
    )
    expect(wrapper.html()).toContain("Listing One")
  })

  it("triggers manual sync", async () => {
    const page = await import("~/pages/_p/com.powerx.plugins.base/admin/integration/marketplace/recommendation.vue")
    const wrapper = await mountSuspended(page.default, {
      global: { stubs: uiStubs },
    })
    await wrapper.findAll("button")[1].trigger("click")
    expect(fetchSpy).toHaveBeenCalledWith(
      "http://localhost:8086/api/v1/admin/marketplace/recommendation/sync",
      expect.objectContaining({ method: "POST" })
    )
  })

  it("updates default weight", async () => {
    const page = await import("~/pages/_p/com.powerx.plugins.base/admin/integration/marketplace/recommendation.vue")
    const wrapper = await mountSuspended(page.default, {
      global: { stubs: uiStubs },
    })

    const input = wrapper.find("input")
    await input.setValue("0.9")
    await wrapper.findAll("button")[0].trigger("click")

    expect(fetchSpy).toHaveBeenCalledWith(
      "http://localhost:8086/api/v1/admin/marketplace/recommendation/experiment",
      expect.objectContaining({ method: "PATCH" })
    )
  })
})
