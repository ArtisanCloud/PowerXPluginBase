<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">
          {{ $t("navigation.backlog") }}
        </h1>
        <p class="mt-2 text-sm text-gray-600 dark:text-gray-400">
          {{ $t("backlog.description") }}
        </p>
      </div>
      <div class="flex gap-2">
        <UButton
          color="primary"
          icon="i-heroicons-plus"
          class="ml-2"
          @click.prevent="createUserStory"
        >
          {{ $t("backlog.createUserStory") }}
        </UButton>
        <UButton
          color="purple"
          variant="outline"
          icon="i-heroicons-bookmark"
          @click.prevent="createEpic"
        >
          {{ $t("backlog.createEpic") }}
        </UButton>
      </div>
    </div>

    <!-- 过滤和排序 -->
    <UCard>
      <div class="flex flex-col sm:flex-row gap-4">
        <div class="flex-1">
          <UInput
            v-model="searchQuery"
            :placeholder="$t('common.search')"
            icon="i-heroicons-magnifying-glass"
          />
        </div>
        <div class="flex gap-2">
          <USelect
            v-model="selectedPriority"
            :options="priorityOptions"
            :placeholder="$t('task.priority')"
            size="sm"
          />
          <USelect
            v-model="selectedStatus"
            :options="statusOptions"
            :placeholder="$t('task.status')"
            size="sm"
          />
        </div>
      </div>
    </UCard>

    <!-- Epic 分组 -->
    <div v-for="epic in epics" :key="epic.id" class="space-y-4">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <UIcon
                name="i-heroicons-bookmark"
                class="h-5 w-5 text-purple-500"
              />
              <div>
                <h3 class="text-lg font-medium text-gray-900 dark:text-white">
                  {{ epic.title }}
                </h3>
                <p class="text-sm text-gray-500 dark:text-gray-400">
                  {{ epic.description }}
                </p>
              </div>
            </div>
            <UBadge
              :color="epic.status === 'active' ? 'green' : 'gray'"
              variant="soft"
            >
              {{ $t(`status.${epic.status}`) }}
            </UBadge>
          </div>
        </template>

        <!-- 用户故事列表 -->
        <div class="space-y-3">
          <div
            v-for="story in epic.userStories"
            :key="story.id"
            class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:bg-gray-50 dark:hover:bg-gray-800 cursor-pointer transition-colors"
            @click="editUserStory(story)"
          >
            <div class="flex items-start justify-between">
              <div class="flex-1">
                <div class="flex items-center space-x-2 mb-2">
                  <UBadge size="xs" variant="outline">
                    {{ story.id }}
                  </UBadge>
                  <UBadge
                    :color="getPriorityColor(story.priority)"
                    size="xs"
                    variant="soft"
                  >
                    {{ $t(`priority.${story.priority}`) }}
                  </UBadge>
                </div>
                <h4
                  class="text-sm font-medium text-gray-900 dark:text-white mb-1"
                >
                  {{ story.title }}
                </h4>
                <p class="text-sm text-gray-600 dark:text-gray-400 mb-2">
                  {{ story.description }}
                </p>
                <div
                  class="flex items-center space-x-4 text-xs text-gray-500 dark:text-gray-400"
                >
                  <span
                    >{{ $t("task.storyPoints") }}: {{ story.storyPoints }}</span
                  >
                  <span
                    >{{ $t("task.assignee") }}:
                    {{ story.assignee || $t("common.unassigned") }}</span
                  >
                </div>
              </div>
              <div class="flex items-center space-x-2">
                <UButton
                  icon="i-heroicons-pencil"
                  size="xs"
                  variant="ghost"
                  @click.stop="editUserStory(story)"
                />
                <UButton
                  icon="i-heroicons-trash"
                  size="xs"
                  variant="ghost"
                  color="red"
                  @click.stop="deleteUserStory(story)"
                />
              </div>
            </div>
          </div>
        </div>
      </UCard>
    </div>

    <!-- 空状态 -->
    <UCard v-if="epics.length === 0">
      <div class="text-center py-12">
        <UIcon
          name="i-heroicons-document-text"
          class="mx-auto h-12 w-12 text-gray-400"
        />
        <h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white">
          {{ $t("backlog.noItems") }}
        </h3>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
          {{ $t("backlog.noItemsDescription") }}
        </p>
        <div class="mt-6">
          <UButton color="primary" @click.prevent="createUserStory">
            {{ $t("backlog.createFirstStory") }}
          </UButton>
        </div>
      </div>
    </UCard>

    <!-- 创建/编辑用户故事模态框 -->
    <UModal
      v-model:open="showStoryModal"
      :ui="{ wrapper: 'flex flex-col' }"
      :title="
        editingStory?.id
          ? $t('backlog.editUserStory')
          : $t('backlog.createUserStory')
      "
    >
      <!-- 弹窗主体内容 -->
      <template #body>
        <UForm
          id="story-form"
          :schema="storySchema"
          :state="storyForm"
          class="mt-4"
          @submit="saveUserStory"
        >
          <!-- Epic 选择 -->
          <UFormField label="Epic" name="epicId" required>
            <USelect
              v-model="storyForm.epicId"
              :options="epicSelectOptions"
              placeholder="选择所属 Epic"
            />
          </UFormField>

          <!-- 标题 -->
          <UFormField :label="$t('task.title')" name="title" required>
            <UInput
              v-model="storyForm.title"
              :placeholder="$t('backlog.storyTitlePlaceholder')"
            />
          </UFormField>

          <!-- 描述 -->
          <UFormField :label="$t('task.description')" name="description">
            <UTextarea
              v-model="storyForm.description"
              :rows="4"
              :placeholder="$t('backlog.storyDescriptionPlaceholder')"
            />
          </UFormField>

          <!-- 优先级 -->
          <UFormField :label="$t('task.priority')" name="priority" required>
            <USelect
              v-model="storyForm.priority"
              :options="prioritySelectOptions"
            />
          </UFormField>

          <!-- 故事点数 -->
          <UFormField :label="$t('task.storyPoints')" name="storyPoints">
            <USelect
              v-model="storyForm.storyPoints"
              :options="storyPointsOptions"
              placeholder="选择故事点数"
            />
          </UFormField>

          <!-- 经办人 -->
          <UFormField :label="$t('task.assignee')" name="assigneeId">
            <USelect
              v-model="storyForm.assigneeId"
              :options="assigneeOptions"
              :placeholder="$t('common.unassigned')"
            />
          </UFormField>

          <!-- 标签 -->
          <UFormField :label="$t('task.labels')" name="labels">
            <UInput
              v-model="labelsInput"
              :placeholder="$t('backlog.labelsPlaceholder')"
              @keypress.enter="addLabel"
            />
            <div
              v-if="storyForm.labels.length"
              class="mt-2 flex flex-wrap gap-1"
            >
              <UBadge
                v-for="(label, index) in storyForm.labels"
                :key="index"
                variant="soft"
                color="primary"
                size="sm"
              >
                {{ label }}
                <UButton
                  icon="i-heroicons-x-mark"
                  size="2xs"
                  color="primary"
                  variant="ghost"
                  class="ml-1"
                  @click="removeLabel(index)"
                />
              </UBadge>
            </div>
          </UFormField>

          <!-- 验收标准 -->
          <UFormField
            :label="$t('backlog.acceptanceCriteria')"
            name="acceptanceCriteria"
          >
            <UTextarea
              v-model="storyForm.acceptanceCriteria"
              :rows="3"
              :placeholder="$t('backlog.acceptanceCriteriaPlaceholder')"
            />
          </UFormField>

          <!-- 业务价值 -->
          <UFormField :label="$t('task.businessValue')" name="businessValue">
            <UInput
              v-model.number="storyForm.businessValue"
              type="number"
              :placeholder="$t('backlog.businessValuePlaceholder')"
            />
          </UFormField>

          <!-- 估算工时 -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <UFormField :label="$t('task.originalHours')" name="originalHours">
              <UInput
                v-model.number="storyForm.originalHours"
                type="number"
                :placeholder="$t('task.originalHoursPlaceholder')"
              />
            </UFormField>
            <UFormField
              :label="$t('task.remainingHours')"
              name="remainingHours"
            >
              <UInput
                v-model.number="storyForm.remainingHours"
                type="number"
                :placeholder="$t('task.remainingHoursPlaceholder')"
              />
            </UFormField>
          </div>

          <!-- 时间管理（使用UPopover+UCalendar组合实现日期选择器） -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <UFormField :label="$t('task.startDate')" name="startDate">
              <UPopover>
                <UButton
                  variant="outline"
                  :icon="storyForm.startDate ? undefined : 'i-heroicons-calendar'"
                  class="w-full justify-start"
                >
                  {{ storyForm.startDate ? new Date(storyForm.startDate).toLocaleDateString() : $t('task.selectDate') }}
                </UButton>
                <template #panel>
                  <UCalendar v-model="storyForm.startDate" />
                </template>
              </UPopover>
            </UFormField>
            <UFormField :label="$t('task.status')" name="status" required>
              <USelect
                v-model="storyForm.status"
                :options="statusSelectOptions"
              />
            </UFormField>
            <UFormField :label="$t('task.dueDate')" name="dueDate">
              <UPopover>
                <UButton
                  variant="outline"
                  :icon="storyForm.dueDate ? undefined : 'i-heroicons-calendar'"
                  class="w-full justify-start"
                >
                  {{ storyForm.dueDate ? new Date(storyForm.dueDate).toLocaleDateString() : $t('task.selectDate') }}
                </UButton>
                <template #panel>
                  <UCalendar v-model="storyForm.dueDate" />
                </template>
              </UPopover>
            </UFormField>
          </div>
        </UForm>
      </template>

      <!-- 底部按钮（用 slot 提供的 close() 或表单原生提交） -->
      <template #footer="{ close }">
        <UButton variant="outline" @click="close()">
          {{ $t("common.cancel") }}
        </UButton>
        <!-- 让外部按钮提交 UForm：用 form 属性指向 form 的 id -->
        <UButton type="submit" form="story-form" :loading="saving">
          {{ editingStory?.id ? $t("common.update") : $t("common.create") }}
        </UButton>
      </template>
    </UModal>
    <!-- 创建 Epic 模态框 -->
    <UModal
      v-model:open="showEpicModal"
      :ui="{ wrapper: 'flex flex-col' }"
      :title="$t('backlog.createEpic')"
    >
      <template #body>
        <UForm
          id="epic-form"
          :schema="epicSchema"
          :state="epicForm"
          class="mt-4"
          @submit="saveEpic"
        >
          <div class="space-y-4">
            <!-- 标题 -->
            <UFormField :label="$t('task.title')" name="title" required>
              <UInput
                v-model="epicForm.title"
                :placeholder="$t('backlog.epicTitlePlaceholder')"
              />
            </UFormField>

            <!-- 描述 -->
            <UFormField :label="$t('task.description')" name="description">
              <UTextarea
                v-model="epicForm.description"
                :placeholder="$t('backlog.epicDescriptionPlaceholder')"
                :rows="4"
              />
            </UFormField>

            <!-- 优先级 -->
            <UFormField :label="$t('task.priority')" name="priority" required>
              <USelect
                v-model="epicForm.priority"
                :options="prioritySelectOptions"
              />
            </UFormField>

            <!-- Epic颜色 -->
            <UFormField :label="$t('task.epicColor')" name="epicColor">
              <UColorPicker v-model="epicForm.epicColor" />
            </UFormField>
          </div>
        </UForm>
      </template>
      <template #footer="{ close }">
        <UButton variant="outline" @click="close()">
          {{ $t("common.cancel") }}
        </UButton>
        <UButton type="submit" form="epic-form" :loading="saving">
          {{ $t("common.create") }}
        </UButton>
      </template>
    </UModal>
  </div>
</template>

<script setup>
import { z } from "zod";
import { watch } from "vue";

// 国际化
const { t } = useI18n();

// 通知和确认框
const showToast = useToast();
const $confirm = useConfirm();

// 页面元数据
useSeoMeta({
  title: () => `${t("common.appName")} - ${t("navigation.backlog")}`,
  description: () => t("backlog.description"),
});

// 搜索和过滤
const searchQuery = ref("");
const selectedPriority = ref("");
const selectedStatus = ref("");

// 优先级选项 - 过滤用
const priorityOptions = computed(() => [
  { label: t("common.all"), value: "" },
  { label: t("priority.critical"), value: "critical" },
  { label: t("priority.high"), value: "high" },
  { label: t("priority.medium"), value: "medium" },
  { label: t("priority.low"), value: "low" },
]);

// 优先级选项 - 表单用
const prioritySelectOptions = computed(() => [
  { label: t("priority.critical"), value: "critical" },
  { label: t("priority.urgent"), value: "urgent" },
  { label: t("priority.high"), value: "high" },
  { label: t("priority.medium"), value: "medium" },
  { label: t("priority.low"), value: "low" },
  { label: t("priority.lowest"), value: "lowest" },
]);

// 状态选项 - 表单用
const statusSelectOptions = computed(() => [
  { label: t("status.todo"), value: "todo" },
  { label: t("status.inProgress"), value: "in_progress" },
  { label: t("status.review"), value: "review" },
  { label: t("status.testing"), value: "testing" },
  { label: t("status.blocked"), value: "blocked" },
  { label: t("status.done"), value: "done" },
]);

// 状态选项
const statusOptions = computed(() => [
  { label: t("common.all"), value: "" },
  { label: t("status.todo"), value: "todo" },
  { label: t("status.inProgress"), value: "in_progress" },
  { label: t("status.review"), value: "review" },
  { label: t("status.testing"), value: "testing" },
  { label: t("status.blocked"), value: "blocked" },
  { label: t("status.done"), value: "done" },
]);

// Epic 选择选项
const epicSelectOptions = computed(() =>
  epics.value.map((epic) => ({
    label: epic.title,
    value: epic.id,
  }))
);

// 经办人选项
const assigneeOptions = computed(() => [
  { label: t("common.unassigned"), value: null },
  ...teamMembers.value,
]);

// 故事点数选项
const storyPointsOptions = [
  { label: "1", value: 1 },
  { label: "2", value: 2 },
  { label: "3", value: 3 },
  { label: "5", value: 5 },
  { label: "8", value: 8 },
  { label: "13", value: 13 },
  { label: "21", value: 21 },
];

// 表单状态
const showStoryModal = ref(false);
const showEpicModal = ref(false);
const editingStory = ref(null);
const saving = ref(false);
const labelsInput = ref("");

// 监听模态框状态变化
watch(showStoryModal, (newVal) => {
  console.log("showStoryModal changed:", newVal);
});

// 用户故事表单
const storyForm = ref({
  epicId: "",
  title: "",
  description: "",
  priority: "medium",
  storyPoints: null,
  assigneeId: null,
  labels: [],
  acceptanceCriteria: "",
  projectId: 1, // 默认项目ID，实际应从项目上下文获取
  status: "todo",
  dueDate: null,
  startDate: null,
  businessValue: null,
});

// Epic 表单
const epicForm = ref({
  title: "",
  description: "",
  priority: "medium",
  projectId: 1, // 默认项目ID，实际应从项目上下文获取
  status: "active",
  epicColor: "#6366F1", // 默认颜色
});

// 表单验证规则
const storySchema = z.object({
  epicId: z.string().min(1, "请选择 Epic"),
  title: z.string().min(1, "标题不能为空").max(200, "标题长度不能超过200字符"),
  description: z.string().optional(),
  priority: z.string().min(1, "请选择优先级"),
  storyPoints: z.number().optional(),
  assigneeId: z.string().optional(),
  labels: z.array(z.string()).optional(),
  acceptanceCriteria: z.string().optional(),
  projectId: z.number(),
  status: z.string().min(1, "请选择状态"),
  dueDate: z.date().optional().nullable(),
  startDate: z.date().optional().nullable(),
  businessValue: z.number().optional().nullable(),
  originalHours: z.number().optional().nullable(),
  remainingHours: z.number().optional().nullable(),
  taskType: z.string().default("user_story"),
});

const epicSchema = z.object({
  title: z.string().min(1, "标题不能为空").max(200, "标题长度不能超过200字符"),
  description: z.string().optional(),
  priority: z.string().min(1, "请选择优先级"),
  projectId: z.number(),
  status: z.string(),
  epicColor: z
    .string()
    .regex(/^#[0-9A-Fa-f]{6}$/, "颜色格式必须为 #RRGGBB")
    .optional(),
});

// Mock 数据
const epics = ref([
  {
    id: "EPIC-001",
    title: t("backlog.sampleEpic"),
    description: t("backlog.sampleEpicDescription"),
    status: "active",
    priority: "high",
    userStories: [
      {
        id: "US-001",
        title: t("backlog.sampleUserStory1"),
        description: t("backlog.sampleUserStoryDescription1"),
        priority: "high",
        storyPoints: 8,
        assignee: "John Doe",
        assigneeId: "1",
        status: "todo",
        labels: ["frontend", "authentication"],
        acceptanceCriteria:
          "用户能够成功登录系统\n用户信息正确显示\n错误信息友好提示",
      },
      {
        id: "US-002",
        title: t("backlog.sampleUserStory2"),
        description: t("backlog.sampleUserStoryDescription2"),
        priority: "medium",
        storyPoints: 5,
        assignee: null,
        assigneeId: null,
        status: "todo",
        labels: ["backend"],
        acceptanceCriteria: "API 接口正常响应\n数据验证完整\n错误处理健壮",
      },
    ],
  },
]);

// 团队成员选项 (Mock)
const teamMembers = ref([
  { id: "1", name: "John Doe", value: "1", label: "John Doe" },
  { id: "2", name: "Jane Smith", value: "2", label: "Jane Smith" },
  { id: "3", name: "Bob Johnson", value: "3", label: "Bob Johnson" },
]);

// 获取优先级颜色
const getPriorityColor = (priority) => {
  const colors = {
    critical: "red",
    urgent: "rose",
    high: "orange",
    medium: "yellow",
    low: "green",
    lowest: "blue",
  };
  return colors[priority] || "gray";
};

// 事件处理
const createUserStory = () => {
  console.log("createUserStory clicked");
  editingStory.value = null;
  resetStoryForm();
  showStoryModal.value = true;
};

const createEpic = () => {
  resetEpicForm();
  showEpicModal.value = true;
};

const editUserStory = (story) => {
  editingStory.value = story;
  populateStoryForm(story);
  showStoryModal.value = true;
};

const deleteUserStory = async (story) => {
  const confirmed = await $confirm({
    title: t("common.confirmDelete"),
    description: t("backlog.confirmDeleteStory", { title: story.title }),
  });

  if (confirmed) {
    // TODO: 调用 API 删除用户故事
    // 临时从列表中移除
    const epic = epics.value.find((e) =>
      e.userStories.some((s) => s.id === story.id)
    );
    if (epic) {
      const index = epic.userStories.findIndex((s) => s.id === story.id);
      if (index > -1) {
        epic.userStories.splice(index, 1);
      }
    }

    showToast({
      title: t("common.success"),
      description: t("backlog.storyDeleted"),
      color: "green",
    });
  }
};

// 表单操作
const resetStoryForm = () => {
  storyForm.value = {
    epicId: "",
    title: "",
    description: "",
    priority: "medium",
    storyPoints: null,
    assigneeId: null,
    labels: [],
    acceptanceCriteria: "",
    projectId: 1, // 默认项目ID，实际应从项目上下文获取
    status: "todo",
    dueDate: null,
    startDate: null,
    businessValue: null,
    originalHours: null,
    remainingHours: null,
    taskType: "user_story", // 默认任务类型
  };
  labelsInput.value = "";
};

const resetEpicForm = () => {
  epicForm.value = {
    title: "",
    description: "",
    priority: "medium",
    projectId: 1, // 默认项目ID，实际应从项目上下文获取
    status: "active",
    epicColor: "#6366F1", // 默认颜色
  };
};

const populateStoryForm = (story) => {
  const epic = epics.value.find((e) =>
    e.userStories.some((s) => s.id === story.id)
  );
  storyForm.value = {
    epicId: epic?.id || "",
    title: story.title || "",
    description: story.description || "",
    priority: story.priority || "medium",
    storyPoints: story.storyPoints || null,
    assigneeId: story.assigneeId || null,
    labels: [...(story.labels || [])],
    acceptanceCriteria: story.acceptanceCriteria || "",
    projectId: story.projectId || 1,
    status: story.status || "todo",
    dueDate: story.dueDate || null,
    startDate: story.startDate || null,
    businessValue: story.businessValue || null,
    originalHours: story.originalHours || null,
    remainingHours: story.remainingHours || null,
    taskType: story.taskType || "user_story",
  };
  labelsInput.value = "";
};

const closeStoryModal = () => {
  showStoryModal.value = false;
  editingStory.value = null;
  resetStoryForm();
};

const closeEpicModal = () => {
  showEpicModal.value = false;
  resetEpicForm();
};

const saveUserStory = async (data) => {
  saving.value = true;
  try {
    // TODO: 调用 API 保存用户故事

    if (editingStory.value) {
      // 更新现有故事
      const epic = epics.value.find((e) =>
        e.userStories.some((s) => s.id === editingStory.value.id)
      );
      if (epic) {
        const story = epic.userStories.find(
          (s) => s.id === editingStory.value.id
        );
        if (story) {
          Object.assign(story, {
            ...data,
            assignee: data.assigneeId
              ? teamMembers.value.find((m) => m.id === data.assigneeId)?.name
              : null,
          });
        }
      }

      showToast({
        title: t("common.success"),
        description: t("backlog.storyUpdated"),
        color: "green",
      });
    } else {
      // 创建新故事
      const epic = epics.value.find((e) => e.id === data.epicId);
      if (epic) {
        const newStory = {
          id: `US-${Date.now()}`,
          ...data,
          assignee: data.assigneeId
            ? teamMembers.value.find((m) => m.id === data.assigneeId)?.name
            : null,
        };
        epic.userStories.push(newStory);
      }

      showToast({
        title: t("common.success"),
        description: t("backlog.storyCreated"),
        color: "green",
      });
    }

    closeStoryModal();
  } catch (error) {
    showToast({
      title: t("common.error"),
      description: error.message || t("common.unexpectedError"),
      color: "red",
    });
  } finally {
    saving.value = false;
  }
};

const saveEpic = async (data) => {
  saving.value = true;
  try {
    // TODO: 调用 API 保存 Epic

    const newEpic = {
      id: `EPIC-${Date.now()}`,
      ...data,
      userStories: [],
    };

    epics.value.push(newEpic);

    showToast({
      title: t("common.success"),
      description: t("backlog.epicCreated"),
      color: "green",
    });

    closeEpicModal();
  } catch (error) {
    showToast({
      title: t("common.error"),
      description: error.message || t("common.unexpectedError"),
      color: "red",
    });
  } finally {
    saving.value = false;
  }
};

// 标签操作
const addLabel = () => {
  const label = labelsInput.value.trim();
  if (label && !storyForm.value.labels.includes(label)) {
    storyForm.value.labels.push(label);
    labelsInput.value = "";
  }
};

const removeLabel = (index) => {
  storyForm.value.labels.splice(index, 1);
};
</script>
