<template>
  <div class="space-y-6">
    <!-- 页面头部 -->
    <div class="border-b border-gray-200 dark:border-gray-700 pb-4">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
        {{ $t("navigation.teamManagement") }}
      </h1>
      <p class="text-gray-600 dark:text-gray-400 mt-1">
        管理团队信息、设置和配置
      </p>
    </div>

    <!-- 团队概览 -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <UCard>
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <UIcon name="i-heroicons-users" class="w-8 h-8 text-blue-500" />
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">
              团队成员
            </p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">
              {{ teamStats.totalMembers }}
            </p>
          </div>
        </div>
      </UCard>

      <UCard>
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <UIcon
              name="i-heroicons-document-text"
              class="w-8 h-8 text-green-500"
            />
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">
              团队笔记
            </p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">
              {{ teamStats.totalNotes }}
            </p>
          </div>
        </div>
      </UCard>

      <UCard>
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <UIcon
              name="i-heroicons-chart-bar"
              class="w-8 h-8 text-purple-500"
            />
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">
              活跃度
            </p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">
              {{ teamStats.activityRate }}%
            </p>
          </div>
        </div>
      </UCard>
    </div>

    <!-- 团队信息 -->
    <UCard>
      <template #header>
        <div class="flex justify-between items-center">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
            团队信息
          </h2>
          <UButton color="primary" @click="editTeamInfo"> 编辑信息 </UButton>
        </div>
      </template>

      <div class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label
              class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
            >
              团队名称
            </label>
            <UInput
              v-model="teamInfo.name"
              :disabled="!isEditing"
              placeholder="输入团队名称"
            />
          </div>
          <div>
            <label
              class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
            >
              团队类型
            </label>
            <USelectMenu
              v-model="teamInfo.type"
              :options="teamTypeOptions"
              :disabled="!isEditing"
              placeholder="选择团队类型"
            />
          </div>
        </div>

        <div>
          <label
            class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
          >
            团队描述
          </label>
          <UTextarea
            v-model="teamInfo.description"
            :disabled="!isEditing"
            placeholder="输入团队描述"
            :rows="3"
          />
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label
              class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
            >
              创建时间
            </label>
            <UInput :model-value="formatDate(teamInfo.createdAt)" disabled />
          </div>
          <div>
            <label
              class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
            >
              团队负责人
            </label>
            <USelectMenu
              v-model="teamInfo.leader"
              :options="memberOptions"
              :disabled="!isEditing"
              placeholder="选择团队负责人"
            />
          </div>
        </div>

        <div v-if="isEditing" class="flex justify-end space-x-3 pt-4">
          <UButton variant="ghost" @click="cancelEdit">
            {{ $t("common.cancel") }}
          </UButton>
          <UButton color="primary" @click="saveTeamInfo">
            {{ $t("common.save") }}
          </UButton>
        </div>
      </div>
    </UCard>

    <!-- 团队设置 -->
    <UCard>
      <template #header>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
          团队设置
        </h2>
      </template>

      <div class="space-y-6">
        <!-- 权限设置 -->
        <div>
          <h3 class="text-md font-medium text-gray-900 dark:text-white mb-4">
            权限设置
          </h3>
          <div class="space-y-3">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium text-gray-900 dark:text-white">
                  允许成员邀请新成员
                </p>
                <p class="text-xs text-gray-500">
                  团队成员可以邀请其他用户加入团队
                </p>
              </div>
              <UToggle v-model="teamSettings.allowMemberInvite" />
            </div>

            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium text-gray-900 dark:text-white">
                  公开团队笔记
                </p>
                <p class="text-xs text-gray-500">团队笔记对所有成员可见</p>
              </div>
              <UToggle v-model="teamSettings.publicNotes" />
            </div>

            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium text-gray-900 dark:text-white">
                  启用笔记评论
                </p>
                <p class="text-xs text-gray-500">
                  允许成员对笔记进行评论和讨论
                </p>
              </div>
              <UToggle v-model="teamSettings.enableComments" />
            </div>
          </div>
        </div>

        <!-- 通知设置 -->
        <div>
          <h3 class="text-md font-medium text-gray-900 dark:text-white mb-4">
            通知设置
          </h3>
          <div class="space-y-3">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium text-gray-900 dark:text-white">
                  新成员加入通知
                </p>
                <p class="text-xs text-gray-500">
                  有新成员加入时通知团队负责人
                </p>
              </div>
              <UToggle v-model="teamSettings.notifyNewMember" />
            </div>

            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium text-gray-900 dark:text-white">
                  笔记更新通知
                </p>
                <p class="text-xs text-gray-500">笔记有更新时通知相关成员</p>
              </div>
              <UToggle v-model="teamSettings.notifyNoteUpdate" />
            </div>
          </div>
        </div>

        <div class="flex justify-end pt-4">
          <UButton color="primary" @click="saveTeamSettings">
            保存设置
          </UButton>
        </div>
      </div>
    </UCard>
  </div>
</template>

<script setup>
// 页面元数据
definePageMeta({
  title: "teamManagement",
});

// 响应式数据
const isEditing = ref(false);

// 团队统计数据
const teamStats = ref({
  totalMembers: 12,
  totalNotes: 156,
  activityRate: 85,
});

// 团队信息
const teamInfo = ref({
  name: "产品开发团队",
  type: "development",
  description: "负责产品功能开发和维护的核心团队",
  createdAt: "2024-01-15",
  leader: "zhangsan",
});

// 团队设置
const teamSettings = ref({
  allowMemberInvite: true,
  publicNotes: true,
  enableComments: true,
  notifyNewMember: true,
  notifyNoteUpdate: false,
});

// 选项数据
const teamTypeOptions = [
  { label: "开发团队", value: "development" },
  { label: "设计团队", value: "design" },
  { label: "产品团队", value: "product" },
  { label: "运营团队", value: "operation" },
];

const memberOptions = [
  { label: "张三", value: "zhangsan" },
  { label: "李四", value: "lisi" },
  { label: "王五", value: "wangwu" },
];

// 方法
const editTeamInfo = () => {
  isEditing.value = true;
};

const cancelEdit = () => {
  isEditing.value = false;
  // 这里可以重置表单数据
};

const saveTeamInfo = () => {
  // 保存团队信息的逻辑
  console.log("保存团队信息", teamInfo.value);
  isEditing.value = false;
};

const saveTeamSettings = () => {
  // 保存团队设置的逻辑
  console.log("保存团队设置", teamSettings.value);
};

const formatDate = (date) => {
  return new Date(date).toLocaleDateString();
};
</script>
