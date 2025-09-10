<template>
  <div class="space-y-6">
    <!-- 页面头部 -->
    <div class="border-b border-gray-200 dark:border-gray-700 pb-4">
      <div class="flex justify-between items-center">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
            {{ $t("navigation.notesList") }}
          </h1>
          <p class="text-gray-600 dark:text-gray-400 mt-1">
            管理和查看所有笔记内容
          </p>
        </div>
        <UButton color="primary" @click="createNote">
          <UIcon name="i-heroicons-plus" class="w-4 h-4 mr-2" />
          {{ $t("notes.createNote") }}
        </UButton>
      </div>
    </div>

    <!-- API 基础信息 -->
    <UCard v-if="baseURL">
      <div class="text-sm text-gray-600 dark:text-gray-400">
        <strong>API 基础路径:</strong> {{ baseURL }}
      </div>
    </UCard>

    <!-- 加载状态 -->
    <div v-if="loading" class="flex justify-center py-12">
      <UIcon
        name="i-heroicons-arrow-path"
        class="w-8 h-8 animate-spin text-blue-500"
      />
    </div>

    <!-- 错误状态 -->
    <UAlert v-if="error" color="red" variant="soft" :title="error" />

    <template v-else>
      <!-- 笔记统计 -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
        <UCard>
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <UIcon
                name="i-heroicons-document-text"
                class="w-8 h-8 text-blue-500"
              />
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600 dark:text-gray-400">
                {{ $t("notes.totalNotes") }}
              </p>
              <p class="text-2xl font-bold text-gray-900 dark:text-white">
                {{ noteStats.total }}
              </p>
            </div>
          </div>
        </UCard>

        <UCard>
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <UIcon name="i-heroicons-play" class="w-8 h-8 text-green-500" />
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600 dark:text-gray-400">
                {{ $t("notes.activeNotes") }}
              </p>
              <p class="text-2xl font-bold text-gray-900 dark:text-white">
                {{ noteStats.active }}
              </p>
            </div>
          </div>
        </UCard>

        <UCard>
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <UIcon
                name="i-heroicons-document-check"
                class="w-8 h-8 text-purple-500"
              />
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600 dark:text-gray-400">
                {{ $t("notes.published") }}
              </p>
              <p class="text-2xl font-bold text-gray-900 dark:text-white">
                {{ noteStats.published }}
              </p>
            </div>
          </div>
        </UCard>

        <UCard>
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <UIcon
                name="i-heroicons-archive-box"
                class="w-8 h-8 text-orange-500"
              />
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600 dark:text-gray-400">
                {{ $t("notes.archived") }}
              </p>
              <p class="text-2xl font-bold text-gray-900 dark:text-white">
                {{ noteStats.archived }}
              </p>
            </div>
          </div>
        </UCard>
      </div>

      <!-- 搜索和筛选 -->
      <UCard>
        <div class="flex flex-col sm:flex-row gap-4">
          <UInput
            v-model="searchQuery"
            :placeholder="$t('notes.searchNotes')"
            icon="i-heroicons-magnifying-glass"
            class="flex-1"
          />
          <USelectMenu
            v-model="selectedCategory"
            :options="categoryOptions"
            :placeholder="$t('notes.filterByCategory')"
            class="w-full sm:w-48"
          />
          <USelectMenu
            v-model="selectedStatus"
            :options="statusOptions"
            :placeholder="$t('notes.filterByStatus')"
            class="w-full sm:w-48"
          />
          <USelectMenu
            v-model="selectedPriority"
            :options="priorityOptions"
            placeholder="按优先级筛选"
            class="w-full sm:w-48"
          />
        </div>
      </UCard>

      <!-- 笔记列表 -->
      <UCard>
        <template #header>
          <div class="flex justify-between items-center">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
              笔记列表 ({{ notes.length }})
            </h2>
            <div class="flex items-center space-x-2">
              <UButton variant="ghost" size="sm" @click="toggleView">
                <UIcon
                  :name="
                    viewMode === 'grid'
                      ? 'i-heroicons-list-bullet'
                      : 'i-heroicons-squares-2x2'
                  "
                  class="w-4 h-4"
                />
              </UButton>
              <USelectMenu v-model="sortBy" :options="sortOptions" size="sm" />
            </div>
          </div>
        </template>

        <!-- 网格视图 -->
        <div
          v-if="viewMode === 'grid'"
          class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
        >
          <div
            v-for="note in notes"
            :key="note.id"
            class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:shadow-md transition-shadow cursor-pointer"
            @click="viewNote(note)"
          >
            <div class="flex justify-between items-start mb-3">
              <h3
                class="font-medium text-gray-900 dark:text-white line-clamp-2"
              >
                {{ note.title }}
              </h3>
              <UDropdown :items="getNoteActions(note)">
                <UButton variant="ghost" size="xs" @click.stop>
                  <UIcon name="i-heroicons-ellipsis-vertical" class="w-4 h-4" />
                </UButton>
              </UDropdown>
            </div>

            <p
              class="text-sm text-gray-600 dark:text-gray-400 mb-3 line-clamp-3"
            >
              {{ note.content }}
            </p>

            <div
              class="flex items-center justify-between text-xs text-gray-500"
            >
              <div class="flex items-center space-x-2">
                <UBadge
                  :color="getStatusColor(note.status)"
                  variant="soft"
                  size="xs"
                >
                  {{ getStatusText(note.status) }}
                </UBadge>
                <UBadge
                  :color="getPriorityColor(note.priority)"
                  variant="soft"
                  size="xs"
                >
                  {{ getPriorityText(note.priority) }}
                </UBadge>
              </div>
              <span>{{ formatDate(note.updated_at) }}</span>
            </div>

            <div
              class="flex items-center justify-between mt-3 pt-3 border-t border-gray-100 dark:border-gray-800"
            >
              <div class="flex items-center space-x-2">
                <UAvatar :alt="note.author_name" size="xs" />
                <span class="text-xs text-gray-600 dark:text-gray-400">{{
                  note.author_name
                }}</span>
              </div>
              <div class="flex items-center space-x-1">
                <UIcon name="i-heroicons-tag" class="w-3 h-3 text-gray-400" />
                <span class="text-xs text-gray-500">{{ note.category }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 列表视图 -->
        <div v-else class="overflow-x-auto">
          <table
            class="min-w-full divide-y divide-gray-200 dark:divide-gray-700"
          >
            <thead class="bg-gray-50 dark:bg-gray-800">
              <tr>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                >
                  {{ $t("notes.noteTitle") }}
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                >
                  {{ $t("notes.author") }}
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                >
                  {{ $t("notes.category") }}
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                >
                  {{ $t("notes.noteStatus") }}
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                >
                  {{ $t("notes.priority") }}
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                >
                  {{ $t("notes.updatedAt") }}
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                >
                  操作
                </th>
              </tr>
            </thead>
            <tbody
              class="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-700"
            >
              <tr
                v-for="note in notes"
                :key="note.id"
                class="hover:bg-gray-50 dark:hover:bg-gray-800 cursor-pointer"
                @click="viewNote(note)"
              >
                <td class="px-6 py-4">
                  <div>
                    <div
                      class="text-sm font-medium text-gray-900 dark:text-white line-clamp-1"
                    >
                      {{ note.title }}
                    </div>
                    <div class="text-sm text-gray-500 line-clamp-2">
                      {{ note.content }}
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <UAvatar :alt="note.author_name" size="xs" class="mr-2" />
                    <span class="text-sm text-gray-900 dark:text-white">{{
                      note.author_name
                    }}</span>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <UBadge variant="soft" size="sm">{{ note.category }}</UBadge>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <UBadge
                    :color="getStatusColor(note.status)"
                    variant="soft"
                    size="sm"
                  >
                    {{ getStatusText(note.status) }}
                  </UBadge>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <UBadge
                    :color="getPriorityColor(note.priority)"
                    variant="soft"
                    size="sm"
                  >
                    {{ getPriorityText(note.priority) }}
                  </UBadge>
                </td>
                <td
                  class="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-400"
                >
                  {{ formatDate(note.updated_at) }}
                </td>
                <td
                  class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium"
                >
                  <UDropdown :items="getNoteActions(note)">
                    <UButton variant="ghost" size="xs" @click.stop>
                      操作
                      <UIcon
                        name="i-heroicons-chevron-down"
                        class="w-4 h-4 ml-1"
                      />
                    </UButton>
                  </UDropdown>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- 空状态 -->
        <div v-if="notes.length === 0 && !loading" class="text-center py-12">
          <UIcon
            name="i-heroicons-document-text"
            class="w-12 h-12 text-gray-400 mx-auto mb-4"
          />
          <p class="text-gray-500 dark:text-gray-400">
            {{ $t("notes.noNotesFound") }}
          </p>
          <UButton color="primary" class="mt-4" @click="createNote">
            {{ $t("notes.createNote") }}
          </UButton>
        </div>
      </UCard>
    </template>
  </div>
</template>

<script setup lang="ts">
import { useNoteApi, type Note, type Page } from "~/composables/api";

// 页面元数据
definePageMeta({
  title: "notesList",
});

// API 实例
const {
  listNotes,
  createNote: createNoteApi,
  deleteNote: deleteNoteApi,
  archiveNote,
  publishNote,
  duplicateNote,
  baseURL,
} = useNoteApi();

// 响应式数据
const searchQuery = ref("");
const selectedCategory = ref("all");
const selectedStatus = ref("all");
const selectedPriority = ref("all");
const sortBy = ref("updatedAt");
const viewMode = ref("grid");
const loading = ref(false);
const error = ref<string | null>(null);

// 笔记数据
const notes = ref<Note[]>([]);
const noteStats = ref({
  total: 0,
  active: 0,
  published: 0,
  archived: 0,
});

// 分页数据
const currentPage = ref(1);
const pageSize = ref(20);
const totalNotes = ref(0);

// 选项数据
const categoryOptions = [
  { label: "全部分类", value: "all" },
  { label: "产品设计", value: "产品设计" },
  { label: "技术文档", value: "技术文档" },
  { label: "用户体验", value: "用户体验" },
  { label: "项目管理", value: "项目管理" },
];

const statusOptions = [
  { label: "全部状态", value: "all" },
  { label: "草稿", value: "draft" },
  { label: "已发布", value: "published" },
  { label: "已归档", value: "archived" },
];

const priorityOptions = [
  { label: "全部优先级", value: "all" },
  { label: "高优先级", value: "high" },
  { label: "中优先级", value: "medium" },
  { label: "低优先级", value: "low" },
];

const sortOptions = [
  { label: "按更新时间", value: "updatedAt" },
  { label: "按创建时间", value: "createdAt" },
  { label: "按标题", value: "title" },
  { label: "按优先级", value: "priority" },
];

// 获取笔记数据
const fetchNotes = async () => {
  try {
    loading.value = true;
    error.value = null;

    const filters = {
      search: searchQuery.value || undefined,
      category:
        selectedCategory.value !== "all" ? selectedCategory.value : undefined,
      status: selectedStatus.value !== "all" ? selectedStatus.value : undefined,
      priority:
        selectedPriority.value !== "all" ? selectedPriority.value : undefined,
    };

    const response = await listNotes(
      currentPage.value,
      pageSize.value,
      filters
    );
    notes.value = response.items;
    totalNotes.value = response.total;

    // 更新统计数据
    updateStats();
  } catch (err) {
    error.value = err instanceof Error ? err.message : "获取笔记失败";
    console.error("获取笔记失败:", err);

    // 使用模拟数据作为后备
    notes.value = getMockNotes();
    updateStats();
  } finally {
    loading.value = false;
  }
};

// 模拟数据
const getMockNotes = (): Note[] => [
  {
    id: "1",
    title: "产品需求分析文档",
    content:
      "本文档详细分析了新产品功能的需求，包括用户故事、功能点和技术实现方案...",
    category: "产品设计",
    status: "published",
    priority: "high",
    created_at: "2024-02-10T00:00:00Z",
    updated_at: "2024-02-15T00:00:00Z",
    author_name: "张三",
    tags: ["需求", "产品", "分析"],
  },
  {
    id: "2",
    title: "技术架构设计方案",
    content:
      "系统架构设计文档，包含前端、后端、数据库等各个层面的技术选型和设计思路...",
    category: "技术文档",
    status: "draft",
    priority: "medium",
    created_at: "2024-02-12T00:00:00Z",
    updated_at: "2024-02-14T00:00:00Z",
    author_name: "李四",
    tags: ["架构", "技术", "设计"],
  },
  {
    id: "3",
    title: "用户体验优化建议",
    content: "基于用户反馈和数据分析，提出的用户体验优化建议和改进方案...",
    category: "用户体验",
    status: "published",
    priority: "low",
    created_at: "2024-02-08T00:00:00Z",
    updated_at: "2024-02-13T00:00:00Z",
    author_name: "王五",
    tags: ["UX", "优化", "用户"],
  },
];

// 更新统计数据
const updateStats = () => {
  const total = notes.value.length;
  const published = notes.value.filter(
    (note) => note.status === "published"
  ).length;
  const draft = notes.value.filter((note) => note.status === "draft").length;
  const archived = notes.value.filter(
    (note) => note.status === "archived"
  ).length;

  noteStats.value = {
    total: totalNotes.value || total,
    active: draft,
    published,
    archived,
  };
};

// 方法
const getStatusColor = (status?: string) => {
  const colorMap = {
    draft: "yellow",
    published: "green",
    archived: "gray",
  };
  return colorMap[status as keyof typeof colorMap] || "gray";
};

const getStatusText = (status?: string) => {
  const textMap = {
    draft: "草稿",
    published: "已发布",
    archived: "已归档",
  };
  return textMap[status as keyof typeof textMap] || status;
};

const getPriorityColor = (priority?: string) => {
  const colorMap = {
    high: "red",
    medium: "yellow",
    low: "green",
  };
  return colorMap[priority as keyof typeof colorMap] || "gray";
};

const getPriorityText = (priority?: string) => {
  const textMap = {
    high: "高",
    medium: "中",
    low: "低",
  };
  return textMap[priority as keyof typeof textMap] || priority;
};

const getNoteActions = (note: Note) => [
  [
    {
      label: "查看",
      icon: "i-heroicons-eye",
      click: () => viewNote(note),
    },
    {
      label: "编辑",
      icon: "i-heroicons-pencil",
      click: () => editNote(note),
    },
  ],
  [
    {
      label: "复制",
      icon: "i-heroicons-document-duplicate",
      click: () => handleDuplicateNote(note),
    },
    {
      label: "分享",
      icon: "i-heroicons-share",
      click: () => shareNote(note),
    },
  ],
  [
    {
      label: note.status === "archived" ? "取消归档" : "归档",
      icon: "i-heroicons-archive-box",
      click: () => handleToggleArchive(note),
    },
    {
      label: "删除",
      icon: "i-heroicons-trash",
      click: () => handleDeleteNote(note),
    },
  ],
];

const createNote = () => {
  navigateTo("/notes/create");
};

const viewNote = (note: Note) => {
  navigateTo(`/notes/${note.id}`);
};

const editNote = (note: Note) => {
  navigateTo(`/notes/${note.id}/edit`);
};

const handleDuplicateNote = async (note: Note) => {
  try {
    await duplicateNote(note.id);
    await fetchNotes();
    // 显示成功提示
  } catch (err) {
    console.error("复制笔记失败:", err);
    // 显示错误提示
  }
};

const shareNote = (note: Note) => {
  // 实现分享功能
  console.log("分享笔记", note);
};

const handleToggleArchive = async (note: Note) => {
  try {
    if (note.status === "archived") {
      // 取消归档 - 发布笔记
      await publishNote(note.id);
    } else {
      // 归档笔记
      await archiveNote(note.id);
    }
    await fetchNotes();
    // 显示成功提示
  } catch (err) {
    console.error("切换归档状态失败:", err);
    // 显示错误提示
  }
};

const handleDeleteNote = async (note: Note) => {
  // 确认删除
  const confirmed = confirm(`确定要删除笔记"${note.title}"吗？`);
  if (!confirmed) return;

  try {
    await deleteNoteApi(note.id);
    await fetchNotes();
    // 显示成功提示
  } catch (err) {
    console.error("删除笔记失败:", err);
    // 显示错误提示
  }
};

const toggleView = () => {
  viewMode.value = viewMode.value === "grid" ? "list" : "grid";
};

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString();
};

// 监听筛选条件变化
watch(
  [searchQuery, selectedCategory, selectedStatus, selectedPriority],
  () => {
    currentPage.value = 1;
    fetchNotes();
  },
  { debounce: 300 }
);

// 监听排序变化
watch(sortBy, () => {
  fetchNotes();
});

// 初始化数据
onMounted(() => {
  fetchNotes();
});
</script>

<style scoped>
.line-clamp-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

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
