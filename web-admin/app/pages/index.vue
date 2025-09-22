<template>
  <div class="p-6 space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
          {{ $t("notes.overview") }}
        </h1>
        <p class="text-gray-600 dark:text-gray-400 mt-1">
          {{ $t("notes.title") }}系统概览
        </p>
      </div>
      <UButton
        color="primary"
        size="lg"
        @click="createNote"
        class="flex items-center gap-2"
      >
        <UIcon name="i-heroicons-plus" class="w-4 h-4" />
        {{ $t("notes.createNote") }}
      </UButton>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <UCard class="bg-gradient-to-r from-blue-500 to-blue-600 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-blue-100 text-sm font-medium">
              {{ $t("notes.totalNotes") }}
            </p>
            <p class="text-2xl font-bold">{{ stats.total }}</p>
          </div>
          <UIcon
            name="i-heroicons-document-text"
            class="w-8 h-8 text-blue-200"
          />
        </div>
      </UCard>

      <UCard class="bg-gradient-to-r from-green-500 to-green-600 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-green-100 text-sm font-medium">
              {{ $t("notes.activeNotes") }}
            </p>
            <p class="text-2xl font-bold">{{ stats.active }}</p>
          </div>
          <UIcon name="i-heroicons-play" class="w-8 h-8 text-green-200" />
        </div>
      </UCard>

      <UCard class="bg-gradient-to-r from-purple-500 to-purple-600 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-purple-100 text-sm font-medium">
              {{ $t("notes.completedNotes") }}
            </p>
            <p class="text-2xl font-bold">{{ stats.completed }}</p>
          </div>
          <UIcon
            name="i-heroicons-check-circle"
            class="w-8 h-8 text-purple-200"
          />
        </div>
      </UCard>

      <UCard class="bg-gradient-to-r from-orange-500 to-orange-600 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-orange-100 text-sm font-medium">草稿笔记</p>
            <p class="text-2xl font-bold">{{ stats.draft }}</p>
          </div>
          <UIcon name="i-heroicons-document" class="w-8 h-8 text-orange-200" />
        </div>
      </UCard>
    </div>

    <!-- 搜索和筛选 -->
    <UCard>
      <div class="flex flex-col sm:flex-row gap-4 items-center justify-between">
        <div class="flex flex-col sm:flex-row gap-4 flex-1">
          <UInput
            v-model="searchQuery"
            :placeholder="$t('notes.searchNotes')"
            icon="i-heroicons-magnifying-glass"
            class="flex-1"
          />
          <USelect
            v-model="selectedCategory"
            :options="categoryOptions"
            :placeholder="$t('notes.filterByCategory')"
            class="w-full sm:w-48"
          />
          <USelect
            v-model="selectedStatus"
            :options="statusOptions"
            :placeholder="$t('notes.filterByStatus')"
            class="w-full sm:w-48"
          />
        </div>
        <USelect
          v-model="sortBy"
          :options="sortOptions"
          :placeholder="$t('notes.sortBy')"
          class="w-full sm:w-48"
        />
      </div>
    </UCard>

    <!-- 最近笔记列表 -->
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ $t("notes.recentNotes") }}
          </h2>
          <UButton variant="ghost" size="sm" to="/notes/active">
            查看全部
            <UIcon name="i-heroicons-arrow-right" class="w-4 h-4 ml-1" />
          </UButton>
        </div>
      </template>

      <div v-if="filteredNotes.length === 0" class="text-center py-8">
        <UIcon
          name="i-heroicons-document-text"
          class="w-12 h-12 text-gray-400 mx-auto mb-4"
        />
        <p class="text-gray-500 dark:text-gray-400">
          {{ $t("notes.noNotesFound") }}
        </p>
      </div>

      <div v-else class="space-y-4">
        <div
          v-for="note in filteredNotes"
          :key="note.id"
          class="flex items-center justify-between p-4 border border-gray-200 dark:border-gray-700 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors cursor-pointer"
          @click="viewNote(note)"
        >
          <div class="flex-1">
            <div class="flex items-center gap-3 mb-2">
              <h3 class="font-medium text-gray-900 dark:text-white">
                {{ note.title }}
              </h3>
              <UBadge
                :color="getStatusColor(note.status)"
                variant="subtle"
                size="xs"
              >
                {{ getStatusText(note.status) }}
              </UBadge>
              <UBadge
                :color="getPriorityColor(note.priority)"
                variant="outline"
                size="xs"
              >
                {{ getPriorityText(note.priority) }}
              </UBadge>
            </div>
            <p
              class="text-sm text-gray-600 dark:text-gray-400 mb-2 line-clamp-2"
            >
              {{ note.content }}
            </p>
            <div
              class="flex items-center gap-4 text-xs text-gray-500 dark:text-gray-400"
            >
              <span class="flex items-center gap-1">
                <UIcon name="i-heroicons-user" class="w-3 h-3" />
                {{ note.author }}
              </span>
              <span class="flex items-center gap-1">
                <UIcon name="i-heroicons-calendar" class="w-3 h-3" />
                {{ formatDate(note.createdAt) }}
              </span>
              <span v-if="note.category" class="flex items-center gap-1">
                <UIcon name="i-heroicons-tag" class="w-3 h-3" />
                {{ note.category }}
              </span>
            </div>
          </div>
          <div class="flex items-center gap-2 ml-4">
            <UButton
              variant="ghost"
              size="sm"
              color="gray"
              @click.stop="editNote(note)"
            >
              <UIcon name="i-heroicons-pencil" class="w-4 h-4" />
            </UButton>
            <UButton
              variant="ghost"
              size="sm"
              color="red"
              @click.stop="deleteNote(note)"
            >
              <UIcon name="i-heroicons-trash" class="w-4 h-4" />
            </UButton>
          </div>
        </div>
      </div>
    </UCard>
    <NoteCreateModal v-model="showCreateModal" @created="handleModalCreated" />
  </div>
</template>

<script setup lang="ts">
import NoteCreateModal from "~/components/notes/NoteCreateModal.vue";

const { t } = useI18n();

// 页面元数据
definePageMeta({
  title: "笔记概览",
  layout: "default",
  alias: ["/dashboard"],
});

// 响应式数据
const showCreateModal = ref(false);
const searchQuery = ref("");
const selectedCategory = ref("");
const selectedStatus = ref("");
const sortBy = ref("date");

// 统计数据
const stats = ref({
  total: 156,
  active: 42,
  completed: 89,
  draft: 25,
});

// 模拟笔记数据
const notes = ref([
  {
    id: 1,
    title: "项目需求分析文档",
    content:
      "详细分析了新项目的功能需求和技术要求，包括用户故事、验收标准等内容...",
    status: "published",
    priority: "high",
    author: "张三",
    category: "项目管理",
    createdAt: "2024-01-15T10:30:00Z",
    updatedAt: "2024-01-16T14:20:00Z",
    tags: ["需求", "分析", "项目"],
  },
  {
    id: 2,
    title: "技术架构设计方案",
    content:
      "基于微服务架构的系统设计方案，包含服务拆分、数据库设计、API设计等...",
    status: "draft",
    priority: "high",
    author: "李四",
    category: "技术文档",
    createdAt: "2024-01-14T09:15:00Z",
    updatedAt: "2024-01-15T16:45:00Z",
    tags: ["架构", "设计", "微服务"],
  },
  {
    id: 3,
    title: "用户体验优化建议",
    content: "基于用户反馈和数据分析，提出的界面和交互优化建议...",
    status: "published",
    priority: "medium",
    author: "王五",
    category: "产品设计",
    createdAt: "2024-01-13T14:20:00Z",
    updatedAt: "2024-01-14T11:30:00Z",
    tags: ["UX", "优化", "用户体验"],
  },
  {
    id: 4,
    title: "数据库性能优化记录",
    content: "记录了数据库查询优化的过程和结果，包括索引优化、查询重写等...",
    status: "archived",
    priority: "low",
    author: "赵六",
    category: "技术文档",
    createdAt: "2024-01-12T16:45:00Z",
    updatedAt: "2024-01-13T09:20:00Z",
    tags: ["数据库", "性能", "优化"],
  },
  {
    id: 5,
    title: "团队会议纪要",
    content: "记录了本周团队会议的讨论内容、决策和行动项...",
    status: "published",
    priority: "medium",
    author: "孙七",
    category: "会议记录",
    createdAt: "2024-01-11T10:00:00Z",
    updatedAt: "2024-01-11T10:30:00Z",
    tags: ["会议", "纪要", "团队"],
  },
]);

// 选项数据
const categoryOptions = [
  { label: "全部分类", value: "" },
  { label: "项目管理", value: "项目管理" },
  { label: "技术文档", value: "技术文档" },
  { label: "产品设计", value: "产品设计" },
  { label: "会议记录", value: "会议记录" },
];

const statusOptions = [
  { label: "全部状态", value: "" },
  { label: t("notes.draft"), value: "draft" },
  { label: t("notes.published"), value: "published" },
  { label: t("notes.archived"), value: "archived" },
];

const sortOptions = [
  { label: t("notes.sortByDate"), value: "date" },
  { label: t("notes.sortByTitle"), value: "title" },
  { label: t("notes.sortByPriority"), value: "priority" },
];

// 计算属性
const filteredNotes = computed(() => {
  let filtered = notes.value;

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    filtered = filtered.filter(
      (note) =>
        note.title.toLowerCase().includes(query) ||
        note.content.toLowerCase().includes(query) ||
        note.author.toLowerCase().includes(query)
    );
  }

  // 分类过滤
  if (selectedCategory.value) {
    filtered = filtered.filter(
      (note) => note.category === selectedCategory.value
    );
  }

  // 状态过滤
  if (selectedStatus.value) {
    filtered = filtered.filter((note) => note.status === selectedStatus.value);
  }

  // 排序
  filtered.sort((a, b) => {
    switch (sortBy.value) {
      case "title":
        return a.title.localeCompare(b.title);
      case "priority":
        const priorityOrder = { high: 3, medium: 2, low: 1 };
        return priorityOrder[b.priority] - priorityOrder[a.priority];
      case "date":
      default:
        return (
          new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
        );
    }
  });

  return filtered.slice(0, 10); // 只显示前10条
});

// 方法
const createNote = () => {
  showCreateModal.value = true;
};

const handleModalCreated = () => {
  showCreateModal.value = false;
};

const viewNote = (note: any) => {
  // 跳转到笔记详情页面
  navigateTo(`/notes/${note.id}`);
};

const editNote = (note: any) => {
  // 跳转到编辑笔记页面
  navigateTo(`/notes/${note.id}/edit`);
};

const deleteNote = async (note: any) => {
  // 删除笔记逻辑
  const confirmed = confirm(t("message.confirmDelete"));
  if (confirmed) {
    // 这里应该调用API删除笔记
    console.log("删除笔记:", note.id);
    // 从列表中移除
    const index = notes.value.findIndex((n) => n.id === note.id);
    if (index > -1) {
      notes.value.splice(index, 1);
    }
  }
};

const getStatusColor = (status: string) => {
  switch (status) {
    case "published":
      return "green";
    case "draft":
      return "yellow";
    case "archived":
      return "gray";
    default:
      return "blue";
  }
};

const getStatusText = (status: string) => {
  switch (status) {
    case "published":
      return t("notes.published");
    case "draft":
      return t("notes.draft");
    case "archived":
      return t("notes.archived");
    default:
      return status;
  }
};

const getPriorityColor = (priority: string) => {
  switch (priority) {
    case "high":
      return "red";
    case "medium":
      return "yellow";
    case "low":
      return "green";
    default:
      return "gray";
  }
};

const getPriorityText = (priority: string) => {
  switch (priority) {
    case "high":
      return t("notes.high");
    case "medium":
      return t("notes.medium");
    case "low":
      return t("notes.low");
    default:
      return priority;
  }
};

const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  return date.toLocaleDateString("zh-CN", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
};
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
