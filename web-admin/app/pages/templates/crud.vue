<template>
  <UContainer class="py-10 space-y-8">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">
          {{ $t("templates.crud.title") }}
        </h1>
        <p class="text-gray-600 dark:text-gray-300">
          {{ $t("templates.crud.description") }}
        </p>
      </div>
      <UButton icon="i-heroicons-plus" color="primary" @click="startCreate">
        {{ $t("templates.crud.create") }}
      </UButton>
    </div>

    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-medium">{{ $t("templates.crud.formTitle") }}</span>
          <UButton variant="ghost" size="xs" @click="resetForm">{{ $t("common.reset") }}</UButton>
        </div>
      </template>
      <UForm :state="form" class="grid gap-4" @submit.prevent="handleSubmit">
        <UFormField :label="$t('templates.crud.fields.name')" name="name" required>
          <UInput v-model="form.name" :placeholder="$t('templates.crud.fields.namePlaceholder')" />
        </UFormField>
        <UFormField :label="$t('templates.crud.fields.description')" name="description" required>
          <UTextarea v-model="form.description" :placeholder="$t('templates.crud.fields.descriptionPlaceholder')" />
        </UFormField>
        <UFormField :label="$t('templates.crud.fields.content')" name="content" required>
          <UTextarea v-model="form.content" :placeholder="$t('templates.crud.fields.contentPlaceholder')" :rows="4" />
        </UFormField>
        <div class="flex justify-end gap-2">
          <UButton type="submit" color="primary" :loading="saving">
            {{ editingId ? $t('templates.crud.actions.update') : $t('templates.crud.actions.save') }}
          </UButton>
        </div>
      </UForm>
    </UCard>

    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-medium">{{ $t("templates.crud.listTitle") }}</span>
          <UBadge variant="soft" color="primary">{{ templates.length }}</UBadge>
        </div>
      </template>
      <UTable :columns="columns" :data="templates" :loading="loading">
        <!-- 注意：v3 是 -cell，不是 -data；row.original 才是你的对象 -->
        <template #actions-cell="{ row }">
          <div class="flex gap-2">
            <UButton
              size="xs"
              variant="soft"
              icon="i-heroicons-pencil"
              @click="startEdit(row.original)"
            >
              {{ $t('common.edit') }}
            </UButton>
            <UButton
              size="xs"
              variant="soft"
              color="error"
              icon="i-heroicons-trash"
              @click="confirmDelete(row.original)"
            >
              {{ $t('common.delete') }}
            </UButton>
          </div>
        </template>
      </UTable>

    </UCard>

    <ConfirmDialog
      v-model="deleteDialog"
      :title="$t('templates.crud.deleteTitle')"
      :description="$t('templates.crud.deleteConfirm', { name: selectedTemplate?.name || '' })"
      confirm-color="red"
      @confirm="performDelete"
    />
  </UContainer>
</template>

<script setup lang="ts">
import ConfirmDialog from "~/components/ConfirmDialog.vue"
import { useTemplateApi } from "~/composables/api/useTemplate"
import type { Template } from "~/composables/api/useTemplate"

type TemplateRow = {
  id: number
  name: string
  description: string
  content: string
}

const columns = [
  { accessorKey: 'name', header: 'Name' },
  { accessorKey: 'description', header: 'Description' },
  { accessorKey: 'content', header: 'Content' },
  { id: 'actions', header: '' }
] satisfies any

const {
  baseURL: templateApiBase,
  listTemplates,
  createTemplate: createTemplateApi,
  updateTemplate: updateTemplateApi,
  deleteTemplate: deleteTemplateApi,
} = useTemplateApi()

const templates = ref<Template[]>([])
const loading = ref(false)
const saving = ref(false)
const editingId = ref<number | null>(null)
const deleteDialog = ref(false)
const selectedTemplate = ref<Template | null>(null)

const form = reactive({
  name: "",
  description: "",
  content: "",
})

const makeLogHandlers = (action: string, context: Record<string, any> = {}) => ({
  onRequest({ request, options }: any) {
    console.debug(`[templates/crud] ${action} request`, {
      baseURL: templateApiBase,
      request,
      options,
      context,
    })
  },
  onResponse({ response }: any) {
    console.debug(`[templates/crud] ${action} response`, {
      status: response.status,
      data: response._data,
      headers: typeof response.headers?.get === "function"
        ? {
            "x-request-id": response.headers.get("x-request-id"),
          }
        : undefined,
      context,
    })
  },
  onResponseError({ response }: any) {
    console.error(`[templates/crud] ${action} response error`, {
      status: response?.status,
      data: response?._data,
      context,
    })
  },
})

const fetchTemplates = async () => {
  loading.value = true
  try {
    const query = { page: 1, page_size: 50 }
    console.debug('[templates/crud] fetching templates', {
      baseURL: templateApiBase,
      path: 'templates',
      query,
    })

    const res = await listTemplates(query.page, query.page_size, "", makeLogHandlers("templates:list", { query }))
    // console.log(res?.success , res.data , Array.isArray(res.data.list),res)
    if (res?.success && res.data && Array.isArray(res.data.list)) {
      templates.value = res.data.list
      console.debug('[templates/crud] templates loaded', {
        count: templates.value.length,
      })
    } else {
      templates.value = []
      console.warn('[templates/crud] templates response unexpected', res)
    }
  } catch (error) {
    console.error("[templates/crud] Failed to load templates", error)
    templates.value = []
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  editingId.value = null
  form.name = ""
  form.description = ""
  form.content = ""
}

const startCreate = () => {
  resetForm()
}

const startEdit = (tpl: Template) => {
  editingId.value = tpl.id
  form.name = tpl.name
  form.description = tpl.description
  form.content = tpl.content
}

const handleSubmit = async () => {
  if (!form.name || !form.description || !form.content) {
    return
  }
  saving.value = true
  try {
    const payload = {
      name: form.name,
      description: form.description,
      content: form.content,
    }
    if (editingId.value) {
      const res = await updateTemplateApi(
        editingId.value,
        payload,
        makeLogHandlers("templates:update", { id: editingId.value, payload })
      )
      if (!res?.success) {
        throw new Error(res?.message || "Update template failed")
      }
    } else {
      const res = await createTemplateApi(
        payload,
        makeLogHandlers("templates:create", { payload })
      )
      if (!res?.success) {
        throw new Error(res?.message || "Create template failed")
      }
    }
    await fetchTemplates()
    resetForm()
  } catch (error) {
    console.error("[templates/crud] Failed to save template", error)
  } finally {
    saving.value = false
  }
}

const confirmDelete = (tpl: Template) => {
  selectedTemplate.value = tpl
  deleteDialog.value = true
}

const performDelete = async () => {
  if (!selectedTemplate.value) return
  try {
    const res = await deleteTemplateApi(
      selectedTemplate.value.id,
      makeLogHandlers("templates:delete", { id: selectedTemplate.value.id })
    )
    if (!res?.success) {
      throw new Error(res?.message || "Delete template failed")
    }
    await fetchTemplates()
  } catch (error) {
    console.error("[templates/crud] Failed to delete template", error)
  } finally {
    deleteDialog.value = false
    selectedTemplate.value = null
  }
}

onMounted(() => {
  fetchTemplates()
})
</script>
