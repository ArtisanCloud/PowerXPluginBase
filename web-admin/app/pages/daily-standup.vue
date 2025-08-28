<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">
          {{ $t("navigation.dailyStandup") }}
        </h1>
        <p class="mt-2 text-sm text-gray-600 dark:text-gray-400">
          {{ $t("standup.description") }}
        </p>
      </div>
      <UButton
        color="primary"
        icon="i-heroicons-plus"
        @click="startStandup"
      >
        {{ $t("standup.startStandup") }}
      </UButton>
    </div>

    <!-- 站会状态 -->
    <UCard v-if="currentStandup">
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <div class="flex-shrink-0">
            <div class="w-3 h-3 bg-green-500 rounded-full animate-pulse"/>
          </div>
          <div>
            <h3 class="text-lg font-medium text-gray-900 dark:text-white">
              {{ $t("standup.standupInProgress") }}
            </h3>
            <p class="text-sm text-gray-500 dark:text-gray-400">
              {{ $t("standup.startedAt") }}: {{ formatTime(currentStandup.startTime) }}
            </p>
          </div>
        </div>
        <div class="flex items-center space-x-2">
          <UButton
            variant="outline"
            @click="pauseStandup"
          >
            {{ $t("standup.pause") }}
          </UButton>
          <UButton
            color="red"
            @click="endStandup"
          >
            {{ $t("standup.endStandup") }}
          </UButton>
        </div>
      </div>
    </UCard>

    <!-- 团队成员站会状态 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <UCard
        v-for="member in teamMembers"
        :key="member.id"
        :class="[
          'transition-all',
          member.hasSpoken ? 'border-green-200 bg-green-50 dark:bg-green-900/20' : '',
          member.isSpeaking ? 'border-blue-200 bg-blue-50 dark:bg-blue-900/20' : ''
        ]"
      >
        <div class="flex items-start space-x-4">
          <div class="relative">
            <UAvatar
              :alt="member.name"
              size="lg"
              :class="member.isSpeaking ? 'ring-2 ring-blue-500' : ''"
            >
              <span class="text-lg font-medium text-white">
                {{ getInitials(member.name) }}
              </span>
            </UAvatar>
            <div
              v-if="member.hasSpoken"
              class="absolute -bottom-1 -right-1 w-4 h-4 bg-green-500 rounded-full flex items-center justify-center"
            >
              <UIcon name="i-heroicons-check" class="w-3 h-3 text-white" />
            </div>
            <div
              v-else-if="member.isSpeaking"
              class="absolute -bottom-1 -right-1 w-4 h-4 bg-blue-500 rounded-full animate-pulse"
            />
          </div>

          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between mb-2">
              <h3 class="text-sm font-medium text-gray-900 dark:text-white truncate">
                {{ member.name }}
              </h3>
              <UBadge
                :color="member.hasSpoken ? 'green' : member.isSpeaking ? 'blue' : 'gray'"
                size="xs"
                variant="soft"
              >
                {{
                  member.hasSpoken
                    ? $t("standup.completed")
                    : member.isSpeaking
                    ? $t("standup.speaking")
                    : $t("standup.waiting")
                }}
              </UBadge>
            </div>

            <div class="text-xs text-gray-500 dark:text-gray-400 mb-3">
              {{ $t(`team.role.${member.role}`) }}
            </div>

            <!-- 站会内容 -->
            <div v-if="member.standupUpdate" class="space-y-2 text-sm">
              <div>
                <div class="font-medium text-gray-700 dark:text-gray-300">
                  {{ $t("standup.yesterday") }}:
                </div>
                <div class="text-gray-600 dark:text-gray-400">
                  {{ member.standupUpdate.yesterday || $t("standup.noUpdate") }}
                </div>
              </div>
              <div>
                <div class="font-medium text-gray-700 dark:text-gray-300">
                  {{ $t("standup.today") }}:
                </div>
                <div class="text-gray-600 dark:text-gray-400">
                  {{ member.standupUpdate.today || $t("standup.noUpdate") }}
                </div>
              </div>
              <div v-if="member.standupUpdate.blockers">
                <div class="font-medium text-red-600 dark:text-red-400">
                  {{ $t("standup.blockers") }}:
                </div>
                <div class="text-red-500 dark:text-red-400">
                  {{ member.standupUpdate.blockers }}
                </div>
              </div>
            </div>

            <!-- 操作按钮 -->
            <div class="mt-4 flex space-x-2">
              <UButton
                v-if="!member.hasSpoken && !member.isSpeaking"
                size="xs"
                @click="startMemberUpdate(member)"
              >
                {{ $t("standup.startUpdate") }}
              </UButton>
              <UButton
                v-if="member.isSpeaking"
                size="xs"
                color="green"
                @click="completeMemberUpdate(member)"
              >
                {{ $t("standup.completeUpdate") }}
              </UButton>
              <UButton
                v-if="member.hasSpoken"
                size="xs"
                variant="outline"
                @click="editUpdate(member)"
              >
                {{ $t("common.edit") }}
              </UButton>
            </div>
          </div>
        </div>
      </UCard>
    </div>

    <!-- 站会历史 -->
    <UCard>
      <template #header>
        <h3 class="text-lg font-medium text-gray-900 dark:text-white">
          {{ $t("standup.recentStandups") }}
        </h3>
      </template>

      <div class="space-y-4">
        <div
          v-for="standup in recentStandups"
          :key="standup.id"
          class="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-800 rounded-lg"
        >
          <div class="flex items-center space-x-4">
            <div class="flex-shrink-0">
              <UIcon name="i-heroicons-microphone" class="h-5 w-5 text-gray-400" />
            </div>
            <div>
              <div class="text-sm font-medium text-gray-900 dark:text-white">
                {{ formatDate(standup.date) }}
              </div>
              <div class="text-sm text-gray-500 dark:text-gray-400">
                {{ $t("standup.duration") }}: {{ standup.duration }} {{ $t("standup.minutes") }}
                • {{ standup.participants }} {{ $t("standup.participants") }}
              </div>
            </div>
          </div>
          <div class="flex items-center space-x-2">
            <UBadge
              :color="standup.status === 'completed' ? 'green' : 'gray'"
              variant="soft"
              size="sm"
            >
              {{ $t(`standup.status.${standup.status}`) }}
            </UBadge>
            <UButton
              icon="i-heroicons-eye"
              size="xs"
              variant="outline"
              @click="viewStandup(standup)"
            >
              {{ $t("standup.viewDetails") }}
            </UButton>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="recentStandups.length === 0" class="text-center py-8">
        <UIcon name="i-heroicons-microphone" class="mx-auto h-12 w-12 text-gray-400" />
        <h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white">
          {{ $t("standup.noStandups") }}
        </h3>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
          {{ $t("standup.noStandupsDescription") }}
        </p>
      </div>
    </UCard>
  </div>
</template>

<script setup>
// 国际化
const { t } = useI18n();

// 页面元数据
useSeoMeta({
  title: () => `${t('common.appName')} - ${t('navigation.dailyStandup')}`,
  description: () => t('standup.description')
});

// 当前站会状态
const currentStandup = ref(null);

// Mock 团队成员数据
const teamMembers = ref([
  {
    id: 1,
    name: 'John Doe',
    role: 'scrumMaster',
    hasSpoken: true,
    isSpeaking: false,
    standupUpdate: {
      yesterday: '完成了用户认证模块的代码审查',
      today: '开始实现密码重置功能',
      blockers: null
    }
  },
  {
    id: 2,
    name: 'Alice Smith',
    role: 'developer',
    hasSpoken: false,
    isSpeaking: true,
    standupUpdate: {
      yesterday: '修复了登录页面的UI问题',
      today: '将开始用户个人资料页面的开发',
      blockers: '需要等待设计稿确认'
    }
  },
  {
    id: 3,
    name: 'Bob Johnson',
    role: 'developer',
    hasSpoken: false,
    isSpeaking: false,
    standupUpdate: null
  },
  {
    id: 4,
    name: 'Carol Brown',
    role: 'tester',
    hasSpoken: false,
    isSpeaking: false,
    standupUpdate: null
  }
]);

// Mock 最近站会数据
const recentStandups = ref([
  {
    id: 1,
    date: '2024-01-25',
    duration: 15,
    participants: 4,
    status: 'completed'
  },
  {
    id: 2,
    date: '2024-01-24',
    duration: 12,
    participants: 4,
    status: 'completed'
  },
  {
    id: 3,
    date: '2024-01-23',
    duration: 18,
    participants: 3,
    status: 'completed'
  }
]);

// 工具函数
const getInitials = (name) => {
  return name
    .split(' ')
    .map(word => word.charAt(0))
    .join('')
    .toUpperCase()
    .slice(0, 2);
};

const formatTime = (timeString) => {
  return new Date(timeString).toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit'
  });
};

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('zh-CN', {
    month: 'short',
    day: 'numeric'
  });
};

// 事件处理
const startStandup = () => {
  currentStandup.value = {
    id: Date.now(),
    startTime: new Date().toISOString(),
    participants: teamMembers.value.length
  };
  // 重置所有成员状态
  teamMembers.value.forEach(member => {
    member.hasSpoken = false;
    member.isSpeaking = false;
    member.standupUpdate = null;
  });
};

const pauseStandup = () => {
  // TODO: 实现暂停站会功能
  console.log('暂停站会');
};

const endStandup = () => {
  // TODO: 实现结束站会功能
  currentStandup.value = null;
  console.log('结束站会');
};

const startMemberUpdate = (member) => {
  // 停止其他人的发言
  teamMembers.value.forEach(m => {
    m.isSpeaking = false;
  });
  // 开始当前成员发言
  member.isSpeaking = true;
  member.standupUpdate = {
    yesterday: '',
    today: '',
    blockers: null
  };
};

const completeMemberUpdate = (member) => {
  member.isSpeaking = false;
  member.hasSpoken = true;
};

const editUpdate = (member) => {
  // TODO: 实现编辑更新功能
  console.log('编辑更新:', member);
};

const viewStandup = (standup) => {
  // TODO: 实现查看站会详情功能
  console.log('查看站会详情:', standup);
};
</script>