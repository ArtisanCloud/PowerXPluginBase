import { describe, it, expect, vi, beforeEach, afterEach } from "vitest"
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
    template: `<div class="u-card"><slot name="header" /><div><slot /></div><slot name="footer" /></div>`,
  },
  UButton: {
    emits: ["click"],
    template: `<button type="button" @click="$emit('click', $event)"><slot /></button>`,
  },
  USelectMenu: {
    props: ["modelValue"],
    emits: ["update:modelValue"],
    template: `<select multiple @change="$emit('update:modelValue', Array.from($event.target.selectedOptions).map(o => o.value))"><slot /></select>`,
  },
  USelect: {
    props: ["modelValue", "options"],
    emits: ["update:modelValue"],
    template: `<select @change="$emit('update:modelValue', $event.target.value)"><option v-for="option in options" :key="option.value" :value="option.value">{{ option.label }}</option></select>`,
  },
  UInput: {
    props: ["modelValue", "type"],
    emits: ["update:modelValue"],
    template: `<input :type="type || 'text'" :value="modelValue" @input="$emit('update:modelValue', $event.target.value)" />`,
  },
  UTextarea: {
    props: ["modelValue"],
    emits: ["update:modelValue"],
    template: `<textarea :value="modelValue" @input="$emit('update:modelValue', $event.target.value)"></textarea>`,
  },
  UBadge: {
    template: `<span class="u-badge"><slot /></span>`,
  },
  UIcon: {
    props: ["name"],
    template: `<span :data-icon="name"><slot /></span>`,
  },
  UProgress: {
    props: ["value", "max"],
    template: `<progress :value="value" :max="max"></progress>`,
  },
  UTable: {
    props: ["rows"],
    emits: ["select"],
    template: `<table><tbody><tr v-for="row in rows" :key="row.id" @click="$emit('select', row)"><slot name="title-data" :row="row">{{ row.title }}</slot></tr></tbody></table>`,
  },
  UFormGroup: {
    template: `<label><slot /><slot name="help" /></label>`,
  },
  UToggel: { template: "<div><slot /></div>" }, // fallback typo guard
  UToggle: {
    props: ["modelValue"],
    emits: ["update:modelValue"],
    template: `<input type="checkbox" :checked="modelValue" @change="$emit('update:modelValue', $event.target.checked)" />`,
  },
  ULink: {
    props: ["to"],
    template: `<a :href="to"><slot /></a>`,
  },
  UModal: {
    props: ["modelValue"],
    emits: ["update:modelValue"],
    template: `<div v-if="modelValue" class="u-modal"><slot /></div>`,
  },
  UBadgeCheck: { template: "<span><slot /></span>" },
}

describe("MarketplaceChecklistRunner", () => {
  const ChecklistRunner = () => import("~/components/marketplace/ChecklistRunner.vue")
  let fetchSpy: ReturnType<typeof vi.fn>

  beforeEach(() => {
    fetchSpy = vi.fn(async (_url: string, options: any) => {
      const op = options?.body?.operationName
      switch (op) {
        case "ChecklistRuns":
          return {
            data: {
              checklistRuns: [
                {
                  id: "run-1",
                  listing_id: "listing-1",
                  run_number: 1,
                  trigger_source: "vendor",
                  status: "passed",
                  summary: "Initial pass",
                  started_at: "2025-10-20T10:00:00Z",
                  completed_at: "2025-10-20T10:01:00Z",
                  ci_pipeline_id: null,
                  items: [
                    {
                      id: "item-1",
                      code: "ASSET_COVER",
                      description: "Cover asset present",
                      result: "passed",
                      evidence_uri: null,
                      notes: null,
                      auto_fix_link: null,
                    },
                  ],
                },
              ],
            },
          }
        case "LatestChecklistRun":
          return {
            data: {
              latestChecklistRun: {
                id: "run-1",
                listing_id: "listing-1",
                run_number: 1,
                trigger_source: "vendor",
                status: "passed",
                summary: "Initial pass",
                started_at: "2025-10-20T10:00:00Z",
                completed_at: "2025-10-20T10:01:00Z",
                ci_pipeline_id: null,
                items: [],
              },
            },
          }
        case "TriggerChecklist":
          return {
            data: {
              triggerChecklist: {
                id: "run-2",
                listing_id: "listing-1",
                run_number: 2,
                trigger_source: "ci",
                status: "pending",
                summary: "CI pipeline trigger",
                started_at: "2025-10-21T08:00:00Z",
                completed_at: null,
                ci_pipeline_id: "ci-123",
                items: [],
              },
            },
          }
        default:
          return { data: {} }
      }
    })
    ;(globalThis as any).$fetch = fetchSpy
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it("loads checklist runs for the given listing and triggers new run", async () => {
    const component = await ChecklistRunner()
    const wrapper = await mountSuspended(component, {
      props: { listingId: "listing-1" },
      global: { stubs: uiStubs },
    })

    expect(fetchSpy).toHaveBeenCalledWith(
      "http://localhost:8086/api/v1/admin/marketplace/checklist/graphql",
      expect.objectContaining({
        method: "POST",
        body: expect.objectContaining({ operationName: "ChecklistRuns" }),
      })
    )

    const button = wrapper.find("button")
    await button.trigger("click")

    expect(fetchSpy).toHaveBeenCalledWith(
      "http://localhost:8086/api/v1/admin/marketplace/checklist/graphql",
      expect.objectContaining({
        body: expect.objectContaining({ operationName: "TriggerChecklist" }),
      })
    )
  })
})

describe("Marketplace Listings Page", () => {
  const ListingsPage = () =>
    import("~/pages/_p/com.powerx.plugins.base/admin/integration/marketplace/listings.vue")
  let fetchSpy: ReturnType<typeof vi.fn>

  beforeEach(() => {
    fetchSpy = vi.fn(async (url: string, options?: any) => {
      if (url.endsWith("/admin/marketplace/listings") && !options?.method) {
        return {
          data: {
            items: [
              {
                id: "listing-1",
                plugin_id: "com.powerx.plugins.base",
                vendor_id: "vendor-1",
                status: "draft",
                title: "Base Plugin",
                slug: "base-plugin",
                summary: "A base plugin",
                description: "",
                locale: "en",
                categories: ["devops"],
                tags: ["base"],
                ready_checklist_score: 80,
                recommended_weight: 0,
                created_at: "2025-10-20T10:00:00Z",
                updated_at: "2025-10-20T12:00:00Z",
                assets: [],
                pricing_plans: [],
              },
            ],
            total: 1,
          },
        }
      }

      if (url.endsWith("/admin/marketplace/listings") && options?.method === "POST") {
        return {
          data: {
            id: "listing-2",
            plugin_id: options.body.plugin_id,
            vendor_id: options.body.vendor_id,
            status: "draft",
            title: options.body.title,
            slug: options.body.slug,
            locale: options.body.locale,
            categories: options.body.categories,
            tags: options.body.tags,
            ready_checklist_score: 0,
            recommended_weight: 0,
            created_at: "2025-10-21T08:00:00Z",
            updated_at: "2025-10-21T08:00:00Z",
            assets: options.body.assets,
            pricing_plans: [],
          },
        }
      }

      if (url.includes("/admin/marketplace/checklist/graphql")) {
        return { data: { checklistRuns: [] } }
      }

      return { data: {} }
    })

    mockNuxtImport("useNuxtApp", () => () => ({
      $fetch: fetchSpy,
    }))

    ;(globalThis as any).$fetch = fetchSpy
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it("submits new listing payload with transformed arrays", async () => {
    const component = await ListingsPage()
    const wrapper = await mountSuspended(component, {
      global: { stubs: uiStubs },
    })

    const vm: any = wrapper.vm
    vm.createForm.vendor_id = "vendor-2"
    vm.createForm.title = "New Marketplace Listing"
    vm.createForm.slug = "new-marketplace-listing"
    vm.createForm.categories = "analytics,automation"
    vm.createForm.tags = "ai,marketplace"
    vm.createForm.assets = [
      {
        key: "asset-1",
        asset_type: "logo",
        storage_uri: "https://cdn.example/logo.png",
        is_primary: true,
        locale: "en",
        weight: 0,
      },
    ]

    await vm.submitCreate()

    const createCall = fetchSpy.mock.calls.find(([, options]) => options?.method === "POST")
    expect(createCall).toBeTruthy()
    const [, options] = createCall!
    expect(options.body.categories).toEqual(["analytics", "automation"])
    expect(options.body.tags).toEqual(["ai", "marketplace"])
    expect(options.body.assets).toHaveLength(1)
    expect(options.body.assets[0]).toMatchObject({
      asset_type: "logo",
      storage_uri: "https://cdn.example/logo.png",
      is_primary: true,
    })
  })
})
