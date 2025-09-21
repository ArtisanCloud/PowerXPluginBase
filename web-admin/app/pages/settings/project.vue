<template>
  <div class="space-y-6">
    <!-- 页面头部 -->
    <div class="border-b border-gray-200 dark:border-gray-700 pb-4">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
        {{ $t("navigation.projectSettings") }}
      </h1>
      <p class="text-gray-600 dark:text-gray-400 mt-1">
        {{ $t("settings.projectDescription") }}
      </p>
    </div>

    <!-- 基本信息 -->
    <UCard>
      <template #header>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ $t("settings.basicInfo") }}
        </h2>
      </template>

      <UForm
        :state="projectForm"
        class="space-y-4"
        @submit="saveProjectSettings"
      >
        <UFormField :label="$t('settings.projectName')" name="name">
          <UInput v-model="projectForm.name" placeholder="输入项目名称" />
        </UFormField>

        <UFormField
          :label="$t('settings.projectDescription')"
          name="description"
        >
          <UTextarea
            v-model="projectForm.description"
            :placeholder="$t('settings.projectDescriptionPlaceholder')"
            rows="3"
          />
        </UFormField>

        <UFormField :label="$t('settings.projectKey')" name="key">
          <UInput v-model="projectForm.key" placeholder="PROJECT_KEY" />
        </UFormField>

        <UFormField :label="$t('settings.projectLead')" name="lead">
          <USelectMenu
            v-model="projectForm.lead"
            :options="teamMembers"
            :placeholder="$t('settings.selectProjectLead')"
          />
        </UFormField>

        <div class="flex justify-end">
          <UButton type="submit" color="primary">
            {{ $t("common.save") }}
          </UButton>
        </div>
      </UForm>
    </UCard>

    <!-- 笔记设置 -->
    <UCard>
      <template #header>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ $t("settings.sprintSettings") }}
        </h2>
      </template>

      <UForm :state="sprintForm" class="space-y-4" @submit="saveSprintSettings">
        <UFormField
          :label="$t('settings.defaultSprintLength')"
          name="defaultLength"
        >
          <USelectMenu
            v-model="sprintForm.defaultLength"
            :options="sprintLengthOptions"
            :placeholder="$t('settings.selectSprintLength')"
          />
        </UFormField>

        <UFormField :label="$t('settings.sprintStartDay')" name="startDay">
          <USelectMenu
            v-model="sprintForm.startDay"
            :options="weekDayOptions"
            :placeholder="$t('settings.selectStartDay')"
          />
        </UFormField>

        <UFormField :label="$t('settings.autoCreateNext')" name="autoCreate">
          <UToggle v-model="sprintForm.autoCreate" />
          <span class="text-sm text-gray-600 dark:text-gray-400 ml-2">
            {{ $t("settings.autoCreateDescription") }}
          </span>
        </UFormField>

        <div class="flex justify-end">
          <UButton type="submit" color="primary">
            {{ $t("common.save") }}
          </UButton>
        </div>
      </UForm>
    </UCard>

    <!-- 工作流设置 -->
    <UCard>
      <template #header>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ $t("settings.workflowSettings") }}
        </h2>
      </template>

      <div class="space-y-4">
        <div>
          <h3 class="font-medium text-gray-900 dark:text-white mb-3">
            {{ $t("settings.noteStatuses") }}
          </h3>
          <div class="space-y-2">
            <div
              v-for="status in workflowStatuses"
              :key="status.id"
              class="flex items-center justify-between p-3 border border-gray-200 dark:border-gray-700 rounded-lg"
            >
              <div class="flex items-center space-x-3">
                <div class="flex items-center space-x-2">
                  <UIcon
                    name="i-heroicons-bars-3"
                    class="w-4 h-4 text-gray-400 cursor-move"
                  />
                  <UBadge :color="status.color" variant="soft">
                    {{ status.name }}
                  </UBadge>
                </div>
                <span class="text-sm text-gray-600 dark:text-gray-400">
                  {{ status.description }}
                </span>
              </div>
              <div class="flex items-center space-x-2">
                <UButton size="xs" variant="ghost" @click="editStatus(status)">
                  {{ $t("common.edit") }}
                </UButton>
                <UButton
                  size="xs"
                  variant="ghost"
                  color="red"
                  @click="deleteStatus(status)"
                >
                  {{ $t("common.delete") }}
                </UButton>
              </div>
            </div>
          </div>
          <UButton class="mt-3" variant="outline" @click="addNewStatus">
            {{ $t("settings.addStatus") }}
          </UButton>
        </div>
      </div>
    </UCard>

    <!-- 通知设置 -->
    <UCard>
      <template #header>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ $t("settings.notificationSettings") }}
        </h2>
      </template>

      <div class="space-y-4">
        <div class="flex items-center justify-between">
          <div>
            <h3 class="font-medium text-gray-900 dark:text-white">
              {{ $t("settings.emailNotifications") }}
            </h3>
            <p class="text-sm text-gray-600 dark:text-gray-400">
              {{ $t("settings.emailNotificationsDescription") }}
            </p>
          </div>
          <UToggle v-model="notificationSettings.email" />
        </div>

        <div class="flex items-center justify-between">
          <div>
            <h3 class="font-medium text-gray-900 dark:text-white">
              {{ $t("settings.dailyDigest") }}
            </h3>
            <p class="text-sm text-gray-600 dark:text-gray-400">
              {{ $t("settings.dailyDigestDescription") }}
            </p>
          </div>
          <UToggle v-model="notificationSettings.dailyDigest" />
        </div>

        <div class="flex items-center justify-between">
          <div>
            <h3 class="font-medium text-gray-900 dark:text-white">
              {{ $t("settings.sprintReminders") }}
            </h3>
            <p class="text-sm text-gray-600 dark:text-gray-400">
              {{ $t("settings.sprintRemindersDescription") }}
            </p>
          </div>
          <UToggle v-model="notificationSettings.sprintReminders" />
        </div>

        <div
          class="flex justify-end pt-4 border-t border-gray-200 dark:border-gray-700"
        >
          <UButton color="primary" @click="saveNotificationSettings">
            {{ $t("common.save") }}
          </UButton>
        </div>
      </div>
    </UCard>
  </div>
</template>

<script setup>
// 页面元数据
definePageMeta({
  title: "projectSettings",
});

// 响应式数据
const projectForm = ref({
  name: "PowerX Note项目",
  description: "Note项目管理系统",
  key: "PSM",
  lead: null,
});

const sprintForm = ref({
  defaultLength: 2,
  startDay: 1,
  autoCreate: true,
});

const notificationSettings = ref({
  email: true,
  dailyDigest: false,
  sprintReminders: true,
});

// 选项数据
const teamMembers = ref([
  { label: "张三", value: "zhangsan" },
  { label: "李四", value: "lisi" },
  { label: "王五", value: "wangwu" },
]);

const sprintLengthOptions = ref([
  { label: "1周", value: 1 },
  { label: "2周", value: 2 },
  { label: "3周", value: 3 },
  { label: "4周", value: 4 },
]);

const weekDayOptions = ref([
  { label: "周一", value: 1 },
  { label: "周二", value: 2 },
  { label: "周三", value: 3 },
  { label: "周四", value: 4 },
  { label: "周五", value: 5 },
]);

const workflowStatuses = ref([
  {
    id: 1,
    name: "待办",
    description: "待开始的任务",
    color: "gray",
  },
  {
    id: 2,
    name: "进行中",
    description: "正在开发的任务",
    color: "blue",
  },
  {
    id: 3,
    name: "待审核",
    description: "等待代码审查的任务",
    color: "yellow",
  },
  {
    id: 4,
    name: "已完成",
    description: "已完成的任务",
    color: "green",
  },
]);

// 方法
const saveProjectSettings = () => {
  console.log("保存项目设置", projectForm.value);
  // 这里添加保存逻辑
};

const saveSprintSettings = () => {
  console.log("保存笔记设置", sprintForm.value);
  // 这里添加保存逻辑
};

const saveNotificationSettings = () => {
  console.log("保存通知设置", notificationSettings.value);
  // 这里添加保存逻辑
};

const editStatus = (status) => {
  console.log("编辑状态", status);
  // 这里添加编辑逻辑
};

const deleteStatus = (status) => {
  console.log("删除状态", status);
  // 这里添加删除逻辑
};

const addNewStatus = () => {
  console.log("添加新状态");
  // 这里添加新增逻辑
};

// 初始化
onMounted(() => {
  projectForm.value.lead = teamMembers.value[0].value;
});
</script>
