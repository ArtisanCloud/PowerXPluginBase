<template>
  <div class="p-6 space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
          {{ $t("navigation.activenote") }}
        </h1>
        <p class="text-gray-600 dark:text-gray-400 mt-1">
          管理所有活跃状态的笔记
        </p>
      </div>
      <UButton
        color="primary"
        size="lg"
        to="/notes/create"
        class="flex items-center gap-2"
      >
        <UIcon name="i-heroicons-plus" class="w-4 h-4" />
        {{ $t("notes.createNote") }}
      </UButton>
    </div>

    <!-- 搜索和筛选工具栏 -->
    <UCard>
      <div class="flex flex-col lg:flex-row gap-4 items-center justify-between">
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
            v-model="selectedPriority"
            :options="priorityOptions"
            placeholder="按优先级筛选"
            class="w-full sm:w-48"
          />
        </div>
        <div class="flex gap-2">
          <USelect
            v-model="sortBy"
            :options="sortOptions"
            :placeholder="$t('notes.sortBy')"
            class="w-48"
          />
          <UButton
            variant="outline"
            @click="toggleView"
            :icon="
              viewMode === 'grid'
                ? 'i-heroicons-list-bullet'
                : 'i-heroicons-squares-2x2'
            "
          />
        </div>
      </div>
    </UCard>

    <!-- 笔记列表/网格视图 -->
    <div v-if="filteredNotes.length === 0" class="text-center py-12">
      <UIcon
        name="i-heroicons-document-text"
        class="w-16 h-16 text-gray-400 mx-auto mb-4"
      />
      <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">
        {{ $t("notes.noNotesFound") }}
      </h3>
      <p class="text-gray-500 dark:text-gray-400 mb-6">
        开始创建您的第一个笔记吧
      </p>
      <UButton color="primary" to="/notes/create">
        {{ $t("notes.createNote") }}
      </UButton>
    </div>

    <!-- 网格视图 -->
    <div
      v-else-if="viewMode === 'grid'"
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
    >
      <UCard
        v-for="note in paginatedNotes"
        :key="note.id"
        class="hover:shadow-lg transition-shadow cursor-pointer"
        @click="viewNote(note)"
      >
        <template #header>
          <div class="flex items-start justify-between">
            <div class="flex-1">
              <h3
                class="font-semibold text-gray-900 dark:text-white line-clamp-2"
              >
                {{ note.title }}
              </h3>
              <div class="flex items-center gap-2 mt-2">
                <UBadge
                  :color="getPriorityColor(note.priority)"
                  variant="subtle"
                  size="xs"
                >
                  {{ getPriorityText(note.priority) }}
                </UBadge>
                <UBadge color="green" variant="subtle" size="xs">
                  {{ $t("notes.published") }}
                </UBadge>
              </div>
            </div>
            <UDropdown :items="getDropdownItems(note)">
              <UButton
                variant="ghost"
                size="sm"
                icon="i-heroicons-ellipsis-vertical"
              />
            </UDropdown>
          </div>
        </template>

        <div class="space-y-3">
          <p class="text-sm text-gray-600 dark:text-gray-400 line-clamp-3">
            {{ note.content }}
          </p>

          <div
            class="flex items-center justify-between text-xs text-gray-500 dark:text-gray-400"
          >
            <span class="flex items-center gap-1">
              <UIcon name="i-heroicons-user" class="w-3 h-3" />
              {{ note.author }}
            </span>
            <span>{{ formatDate(note.updatedAt) }}</span>
          </div>

          <div
            v-if="note.tags && note.tags.length > 0"
            class="flex flex-wrap gap-1"
          >
            <UBadge
              v-for="tag in note.tags.slice(0, 3)"
              :key="tag"
              variant="outline"
              size="xs"
              color="gray"
            >
              {{ tag }}
            </UBadge>
            <UBadge
              v-if="note.tags.length > 3"
              variant="outline"
              size="xs"
              color="gray"
            >
              +{{ note.tags.length - 3 }}
            </UBadge>
          </div>
        </div>
      </UCard>
    </div>

    <!-- 列表视图 -->
    <UCard v-else>
      <div class="space-y-4">
        <div
          v-for="note in paginatedNotes"
          :key="note.id"
          class="flex items-center justify-between p-4 border border-gray-200 dark:border-gray-700 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors cursor-pointer"
          @click="viewNote(note)"
        >
          <div class="flex-1">
            <div class="flex items-center gap-3 mb-2">
              <h3 class="font-medium text-gray-900 dark:text-white">
                {{ note.title }}
              </h3>
              <UBadge color="green" variant="subtle" size="xs">
                {{ $t("notes.published") }}
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
                {{ formatDate(note.updatedAt) }}
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

    <!-- 分页 -->
    <div v-if="totalPages > 1" class="flex justify-center">
      <UPagination
        v-model="currentPage"
        :page-count="pageSize"
        :total="filteredNotes.length"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
const { t } = useI18n();

// 页面元数据
definePageMeta({
  title: "活跃笔记",
  layout: "default",
});

// 响应式数据
const searchQuery = ref("");
const selectedCategory = ref("");
const selectedPriority = ref("");
const sortBy = ref("date");
const viewMode = ref("grid");
const currentPage = ref(1);
const pageSize = 12;

// 模拟活跃笔记数据
const notes = ref([
  {
    id: 1,
    title: "项目需求分析文档",
    content:
      "详细分析了新项目的功能需求和技术要求，包括用户故事、验收标准等内容。这是一个重要的项目文档，需要团队所有成员仔细阅读和理解。",
    status: "published",
    priority: "high",
    author: "张三",
    category: "项目管理",
    createdAt: "2024-01-15T10:30:00Z",
    updatedAt: "2024-01-16T14:20:00Z",
    tags: ["需求", "分析", "项目", "文档"],
  },
  {
    id: 2,
    title: "技术架构设计方案",
    content:
      "基于微服务架构的系统设计方案，包含服务拆分、数据库设计、API设计等关键技术决策。",
    status: "published",
    priority: "high",
    author: "李四",
    category: "技术文档",
    createdAt: "2024-01-14T09:15:00Z",
    updatedAt: "2024-01-15T16:45:00Z",
    tags: ["架构", "设计", "微服务", "API"],
  },
  {
    id: 3,
    title: "用户体验优化建议",
    content:
      "基于用户反馈和数据分析，提出的界面和交互优化建议，旨在提升用户满意度。",
    status: "published",
    priority: "medium",
    author: "王五",
    category: "产品设计",
    createdAt: "2024-01-13T14:20:00Z",
    updatedAt: "2024-01-14T11:30:00Z",
    tags: ["UX", "优化", "用户体验", "界面"],
  },
  // 更多笔记数据...
]);

// 选项数据
const categoryOptions = [
  { label: "全部分类", value: "" },
  { label: "项目管理", value: "项目管理" },
  { label: "技术文档", value: "技术文档" },
  { label: "产品设计", value: "产品设计" },
  { label: "会议记录", value: "会议记录" },
];

const priorityOptions = [
  { label: "全部优先级", value: "" },
  { label: t("notes.high"), value: "high" },
  { label: t("notes.medium"), value: "medium" },
  { label: t("notes.low"), value: "low" },
];

const sortOptions = [
  { label: t("notes.sortByDate"), value: "date" },
  { label: t("notes.sortByTitle"), value: "title" },
  { label: t("notes.sortByPriority"), value: "priority" },
];

// 计算属性
const filteredNotes = computed(() => {
  let filtered = notes.value.filter((note) => note.status === "published");

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

  // 优先级过滤
  if (selectedPriority.value) {
    filtered = filtered.filter(
      (note) => note.priority === selectedPriority.value
    );
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
          new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime()
        );
    }
  });

  return filtered;
});

const totalPages = computed(() =>
  Math.ceil(filteredNotes.value.length / pageSize)
);

const paginatedNotes = computed(() => {
  const start = (currentPage.value - 1) * pageSize;
  const end = start + pageSize;
  return filteredNotes.value.slice(start, end);
});

// 方法
const toggleView = () => {
  viewMode.value = viewMode.value === "grid" ? "list" : "grid";
};

const viewNote = (note: any) => {
  navigateTo(`/notes/${note.id}`);
};

const editNote = (note: any) => {
  navigateTo(`/notes/${note.id}/edit`);
};

const deleteNote = async (note: any) => {
  const confirmed = confirm(t("message.confirmDelete"));
  if (confirmed) {
    const index = notes.value.findIndex((n) => n.id === note.id);
    if (index > -1) {
      notes.value.splice(index, 1);
    }
  }
};

const getDropdownItems = (note: any) => [
  [
    {
      label: "查看",
      icon: "i-heroicons-eye",
      click: () => viewNote(note),
    },
  ],
  [
    {
      label: "编辑",
      icon: "i-heroicons-pencil",
      click: () => editNote(note),
    },
    {
      label: "归档",
      icon: "i-heroicons-archive-box",
      click: () => archiveNote(note),
    },
  ],
  [
    {
      label: "删除",
      icon: "i-heroicons-trash",
      click: () => deleteNote(note),
    },
  ],
];

const archiveNote = (note: any) => {
  // 归档笔记逻辑
  note.status = "archived";
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

// 监听页面变化，重置到第一页
watch([searchQuery, selectedCategory, selectedPriority, sortBy], () => {
  currentPage.value = 1;
});
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
