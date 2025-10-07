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
        <template #actions-data="{ row }">
          <div class="flex gap-2">
            <UButton size="xs" variant="soft" icon="i-heroicons-pencil" @click="startEdit(row)">
              {{ $t("common.edit") }}
            </UButton>
            <UButton size="xs" variant="soft" color="red" icon="i-heroicons-trash" @click="confirmDelete(row)">
              {{ $t("common.delete") }}
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

const columns = [
  { key: "name", id: "name", label: "Name" },
  { key: "description", id: "description", label: "Description" },
  { key: "content", id: "content", label: "Content" },
  { key: "actions", id: "actions", label: "" },
]

const { public: { apiBaseUrl } } = useRuntimeConfig()

const templates = ref<any[]>([])
const loading = ref(false)
const saving = ref(false)
const editingId = ref<number | null>(null)
const deleteDialog = ref(false)
const selectedTemplate = ref<any | null>(null)

const form = reactive({
  name: "",
  description: "",
  content: "",
})

const fetchTemplates = async () => {
  loading.value = true
  try {
    const res = await $fetch(`${apiBaseUrl}/templates`, {
      query: { page: 1, page_size: 50 },
    })
    if (res && res.data && Array.isArray(res.data.list)) {
      templates.value = res.data.list
    } else {
      templates.value = []
    }
  } catch (error) {
    console.error("Failed to load templates", error)
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

const startEdit = (tpl: any) => {
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
    if (editingId.value) {
      await $fetch(`${apiBaseUrl}/templates/${editingId.value}`, {
        method: "PUT",
        body: form,
      })
    } else {
      await $fetch(`${apiBaseUrl}/templates`, {
        method: "POST",
        body: form,
      })
    }
    await fetchTemplates()
    resetForm()
  } catch (error) {
    console.error("Failed to save template", error)
  } finally {
    saving.value = false
  }
}

const confirmDelete = (tpl: any) => {
  selectedTemplate.value = tpl
  deleteDialog.value = true
}

const performDelete = async () => {
  if (!selectedTemplate.value) return
  try {
    await $fetch(`${apiBaseUrl}/templates/${selectedTemplate.value.id}`, {
      method: "DELETE",
    })
    await fetchTemplates()
  } catch (error) {
    console.error("Failed to delete template", error)
  } finally {
    deleteDialog.value = false
    selectedTemplate.value = null
  }
}

onMounted(() => {
  fetchTemplates()
})
</script>
