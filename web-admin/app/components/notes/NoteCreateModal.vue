<template>
  <UModal v-model="visible" :ui="{ container: 'items-center', base: 'sm:max-w-3xl w-full' }">
    <UCard>
      <template #header>
        <div class="flex items-start justify-between">
          <div>
            <h3 class="text-xl font-semibold text-gray-900 dark:text-white">
              {{ t("notes.createNote") }}
            </h3>
            <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
              创建一个新的笔记文档
            </p>
          </div>
          <UButton variant="ghost" size="sm" icon="i-heroicons-x-mark" @click="closeModal" />
        </div>
      </template>

      <form class="space-y-6 max-h-[70vh] overflow-y-auto pr-1" @submit.prevent="handleSubmit">
        <UAlert v-if="submitError" color="red" variant="soft" :title="submitError" />

        <UCard>
          <template #header>
            <h4 class="text-lg font-semibold text-gray-900 dark:text-white">基本信息</h4>
          </template>

          <div class="space-y-4">
            <UFormField :label="t('notes.noteTitle')" required>
              <UInput
                v-model="form.title"
                :placeholder="t('notes.noteTitle')"
                size="lg"
                :error="errors.title"
              />
            </UFormField>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <UFormField :label="t('notes.category')">
                <USelect
                  v-model="form.category"
                  :options="categoryOptions"
                  placeholder="选择分类"
                />
              </UFormField>

              <UFormField :label="t('notes.priority')">
                <USelect
                  v-model="form.priority"
                  :options="priorityOptions"
                  placeholder="选择优先级"
                />
              </UFormField>
            </div>

            <UFormField :label="t('notes.tags')">
              <div class="space-y-2">
                <UInput
                  v-model="tagInput"
                  placeholder="输入标签后按回车添加"
                  @keydown.enter.prevent="addTag"
                />
                <div v-if="form.tags.length > 0" class="flex flex-wrap gap-2">
                  <UBadge
                    v-for="(tag, index) in form.tags"
                    :key="index"
                    variant="subtle"
                    color="primary"
                    class="cursor-pointer"
                    @click="removeTag(index)"
                  >
                    {{ tag }}
                    <UIcon name="i-heroicons-x-mark" class="w-3 h-3 ml-1" />
                  </UBadge>
                </div>
              </div>
            </UFormField>
          </div>
        </UCard>

        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h4 class="text-lg font-semibold text-gray-900 dark:text-white">
                {{ t("notes.noteContent") }}
              </h4>
              <div class="flex items-center gap-2">
                <UButton
                  variant="ghost"
                  size="sm"
                  :color="isPreview ? 'gray' : 'primary'"
                  @click="isPreview = false"
                >
                  编辑
                </UButton>
                <UButton
                  variant="ghost"
                  size="sm"
                  :color="isPreview ? 'primary' : 'gray'"
                  @click="isPreview = true"
                >
                  预览
                </UButton>
              </div>
            </div>
          </template>

          <div class="min-h-80">
            <UTextarea
              v-if="!isPreview"
              v-model="form.content"
              :placeholder="t('notes.noteContent')"
              :rows="12"
              resize
              class="font-mono"
              :error="errors.content"
            />
            <div
              v-else
              class="prose prose-gray dark:prose-invert max-w-none min-h-80 p-4 border border-gray-200 dark:border-gray-700 rounded-md overflow-y-auto"
              v-html="renderedContent"
            />
          </div>
        </UCard>

        <UCard>
          <template #header>
            <h4 class="text-lg font-semibold text-gray-900 dark:text-white">发布设置</h4>
          </template>

          <div class="space-y-4">
            <UFormField label="状态">
              <USelect v-model="form.status" :options="statusOptions" placeholder="选择状态" />
            </UFormField>
            <UFormField label="发布时间">
              <UInput v-model="form.publishAt" type="datetime-local" :min="minPublishTime" />
            </UFormField>
          </div>
        </UCard>

        <div class="flex items-center justify-between pt-2">
          <UButton variant="ghost" color="gray" @click="closeModal">
            {{ t("common.cancel") }}
          </UButton>
          <div class="flex items-center gap-3">
            <UButton variant="outline" :loading="saving" @click.prevent="saveDraft">
              保存草稿
            </UButton>
            <UButton type="submit" color="primary" :loading="publishing">
              {{ form.status === "published" ? "发布笔记" : "保存笔记" }}
            </UButton>
          </div>
        </div>
      </form>
    </UCard>
  </UModal>
</template>

<script setup lang="ts">
const { t } = useI18n();

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits<{
  (event: "update:modelValue", value: boolean): void;
  (event: "created"): void;
}>();

const visible = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit("update:modelValue", value),
});

const isPreview = ref(false);
const tagInput = ref("");
const saving = ref(false);
const publishing = ref(false);
const submitError = ref<string | null>(null);

const form = reactive({
  title: "",
  content: "",
  category: "",
  priority: "medium" as "low" | "medium" | "high",
  tags: [] as string[],
  status: "draft" as "draft" | "published",
  publishAt: "",
});

const errors = reactive({
  title: "",
  content: "",
});

const categoryOptions = [
  { label: "项目管理", value: "项目管理" },
  { label: "技术文档", value: "技术文档" },
  { label: "产品设计", value: "产品设计" },
  { label: "会议记录", value: "会议记录" },
  { label: "学习笔记", value: "学习笔记" },
  { label: "其他", value: "其他" },
];

const priorityOptions = [
  { label: t("notes.high"), value: "high" },
  { label: t("notes.medium"), value: "medium" },
  { label: t("notes.low"), value: "low" },
];

const statusOptions = [
  { label: t("notes.draft"), value: "draft" },
  { label: t("notes.published"), value: "published" },
];

const renderedContent = computed(() => {
  if (!form.content) {
    return "<p class='text-gray-500'>暂无内容</p>";
  }

  return form.content
    .replace(/\n/g, "<br>")
    .replace(/\*\*(.*?)\*\*/g, "<strong>$1</strong>")
    .replace(/\*(.*?)\*/g, "<em>$1</em>")
    .replace(
      /`(.*?)`/g,
      '<code class="bg-gray-100 dark:bg-gray-800 px-1 rounded">$1</code>'
    );
});

const minPublishTime = computed(() => new Date().toISOString().slice(0, 16));

const resetForm = () => {
  form.title = "";
  form.content = "";
  form.category = "";
  form.priority = "medium";
  form.tags = [];
  form.status = "draft";
  form.publishAt = new Date().toISOString().slice(0, 16);
  tagInput.value = "";
  errors.title = "";
  errors.content = "";
  submitError.value = null;
  isPreview.value = false;
};

watch(
  () => visible.value,
  (isOpen) => {
    if (isOpen) {
      resetForm();
    }
  }
);

const closeModal = () => {
  visible.value = false;
};

const addTag = () => {
  const tag = tagInput.value.trim();
  if (tag && !form.tags.includes(tag)) {
    form.tags.push(tag);
    tagInput.value = "";
  }
};

const removeTag = (index: number) => {
  form.tags.splice(index, 1);
};

const validateForm = () => {
  errors.title = "";
  errors.content = "";
  submitError.value = null;

  if (!form.title.trim()) {
    errors.title = "请输入笔记标题";
    return false;
  }

  if (!form.content.trim()) {
    errors.content = "请输入笔记内容";
    return false;
  }

  return true;
};

const simulateNetwork = (delay = 1000) =>
  new Promise((resolve) => setTimeout(resolve, delay));

const saveDraft = async () => {
  if (!validateForm()) return;

  saving.value = true;
  try {
    await simulateNetwork();
    emit("created");
    closeModal();
  } catch (error) {
    console.error("保存草稿失败", error);
    submitError.value = "保存草稿失败，请稍后重试";
  } finally {
    saving.value = false;
  }
};

const handleSubmit = async () => {
  if (!validateForm()) return;

  publishing.value = true;
  try {
    await simulateNetwork(1200);
    emit("created");
    closeModal();
  } catch (error) {
    console.error("创建笔记失败", error);
    submitError.value = "创建笔记失败，请稍后重试";
  } finally {
    publishing.value = false;
  }
};
</script>

<style scoped>
@reference "@/assets/css/main.css";

.prose {
  @apply text-gray-900 dark:text-gray-100;
}

.prose code {
  @apply text-sm;
}

.prose strong {
  @apply font-semibold text-gray-900 dark:text-white;
}

.prose em {
  @apply italic;
}
</style>
