<template>
  <div class="p-6 max-w-4xl mx-auto">
    <!-- 页面标题 -->
    <div class="mb-6">
      <div class="flex items-center gap-2 mb-2">
        <UButton
          variant="ghost"
          size="sm"
          to="/notes/active"
          icon="i-heroicons-arrow-left"
        />
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
          {{ $t("notes.createNote") }}
        </h1>
      </div>
      <p class="text-gray-600 dark:text-gray-400">创建一个新的笔记文档</p>
    </div>

    <!-- 创建表单 -->
    <form @submit.prevent="handleSubmit" class="space-y-6">
      <!-- 基本信息 -->
      <UCard>
        <template #header>
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
            基本信息
          </h2>
        </template>

        <div class="space-y-4">
          <!-- 标题 -->
          <UFormGroup :label="$t('notes.noteTitle')" required>
            <UInput
              v-model="form.title"
              :placeholder="$t('notes.noteTitle')"
              size="lg"
              :error="errors.title"
            />
          </UFormGroup>

          <!-- 分类和优先级 -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <UFormGroup :label="$t('notes.category')">
              <USelect
                v-model="form.category"
                :options="categoryOptions"
                placeholder="选择分类"
              />
            </UFormGroup>

            <UFormGroup :label="$t('notes.priority')">
              <USelect
                v-model="form.priority"
                :options="priorityOptions"
                placeholder="选择优先级"
              />
            </UFormGroup>
          </div>

          <!-- 标签 -->
          <UFormGroup :label="$t('notes.tags')">
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
          </UFormGroup>
        </div>
      </UCard>

      <!-- 内容编辑 -->
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ $t("notes.noteContent") }}
            </h2>
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

        <div class="min-h-96">
          <!-- 编辑模式 -->
          <UTextarea
            v-if="!isPreview"
            v-model="form.content"
            :placeholder="$t('notes.noteContent')"
            :rows="20"
            resize
            class="font-mono"
            :error="errors.content"
          />

          <!-- 预览模式 -->
          <div
            v-else
            class="prose prose-gray dark:prose-invert max-w-none min-h-96 p-4 border border-gray-200 dark:border-gray-700 rounded-md"
            v-html="renderedContent"
          />
        </div>
      </UCard>

      <!-- 发布设置 -->
      <UCard>
        <template #header>
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
            发布设置
          </h2>
        </template>

        <div class="space-y-4">
          <UFormGroup label="状态">
            <USelect
              v-model="form.status"
              :options="statusOptions"
              placeholder="选择状态"
            />
          </UFormGroup>

          <UFormGroup label="发布时间">
            <UInput
              v-model="form.publishAt"
              type="datetime-local"
              :min="new Date().toISOString().slice(0, 16)"
            />
          </UFormGroup>
        </div>
      </UCard>

      <!-- 操作按钮 -->
      <div
        class="flex items-center justify-between pt-6 border-t border-gray-200 dark:border-gray-700"
      >
        <UButton variant="ghost" color="gray" to="/notes/active">
          {{ $t("common.cancel") }}
        </UButton>

        <div class="flex items-center gap-3">
          <UButton variant="outline" @click="saveDraft" :loading="saving">
            保存草稿
          </UButton>
          <UButton type="submit" color="primary" :loading="publishing">
            {{ form.status === "published" ? "发布笔记" : "保存笔记" }}
          </UButton>
        </div>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
const { t } = useI18n();

// 页面元数据
definePageMeta({
  title: "创建笔记",
  layout: "default",
});

// 响应式数据
const isPreview = ref(false);
const tagInput = ref("");
const saving = ref(false);
const publishing = ref(false);

// 表单数据
const form = reactive({
  title: "",
  content: "",
  category: "",
  priority: "medium",
  tags: [] as string[],
  status: "draft",
  publishAt: "",
});

// 错误信息
const errors = reactive({
  title: "",
  content: "",
});

// 选项数据
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

// 计算属性
const renderedContent = computed(() => {
  if (!form.content) return "<p class='text-gray-500'>暂无内容</p>";

  // 简单的 Markdown 渲染（实际项目中应使用专业的 Markdown 解析器）
  return form.content
    .replace(/\n/g, "<br>")
    .replace(/\*\*(.*?)\*\*/g, "<strong>$1</strong>")
    .replace(/\*(.*?)\*/g, "<em>$1</em>")
    .replace(
      /`(.*?)`/g,
      '<code class="bg-gray-100 dark:bg-gray-800 px-1 rounded">$1</code>'
    );
});

// 方法
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

const saveDraft = async () => {
  if (!validateForm()) return;

  saving.value = true;
  try {
    // 模拟保存草稿
    await new Promise((resolve) => setTimeout(resolve, 1000));

    // 这里应该调用 API 保存草稿
    console.log("保存草稿:", { ...form, status: "draft" });

    // 显示成功消息
    // toast.success("草稿保存成功");

    // 跳转到笔记列表
    navigateTo("/notes/active");
  } catch (error) {
    console.error("保存失败:", error);
    // toast.error("保存失败，请重试");
  } finally {
    saving.value = false;
  }
};

const handleSubmit = async () => {
  if (!validateForm()) return;

  publishing.value = true;
  try {
    // 模拟发布笔记
    await new Promise((resolve) => setTimeout(resolve, 1500));

    // 这里应该调用 API 创建笔记
    console.log("创建笔记:", form);

    // 显示成功消息
    // toast.success("笔记创建成功");

    // 跳转到笔记列表
    navigateTo("/notes/active");
  } catch (error) {
    console.error("创建失败:", error);
    // toast.error("创建失败，请重试");
  } finally {
    publishing.value = false;
  }
};

// 初始化发布时间
onMounted(() => {
  const now = new Date();
  form.publishAt = now.toISOString().slice(0, 16);
});
</script>

<style scoped>
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
