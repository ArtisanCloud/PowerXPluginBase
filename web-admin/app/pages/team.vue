<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">
          {{ $t("navigation.team") }}
        </h1>
        <p class="mt-2 text-sm text-gray-600 dark:text-gray-400">
          {{ $t("team.description") }}
        </p>
      </div>
      <UButton
        color="primary"
        icon="i-heroicons-plus"
        @click="addMember"
      >
        {{ $t("team.addMember") }}
      </UButton>
    </div>

    <!-- 团队统计 -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <UCard>
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <UIcon name="i-heroicons-users" class="h-8 w-8 text-blue-500" />
          </div>
          <div class="ml-4">
            <div class="text-2xl font-bold text-gray-900 dark:text-white">
              {{ teamMembers.length }}
            </div>
            <div class="text-sm text-gray-500 dark:text-gray-400">
              {{ $t("team.totalMembers") }}
            </div>
          </div>
        </div>
      </UCard>

      <UCard>
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <UIcon name="i-heroicons-user-group" class="h-8 w-8 text-green-500" />
          </div>
          <div class="ml-4">
            <div class="text-2xl font-bold text-gray-900 dark:text-white">
              {{ activeMembers.length }}
            </div>
            <div class="text-sm text-gray-500 dark:text-gray-400">
              {{ $t("team.activeMembers") }}
            </div>
          </div>
        </div>
      </UCard>

      <UCard>
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <UIcon name="i-heroicons-chart-bar" class="h-8 w-8 text-purple-500" />
          </div>
          <div class="ml-4">
            <div class="text-2xl font-bold text-gray-900 dark:text-white">
              {{ averageCapacity }}%
            </div>
            <div class="text-sm text-gray-500 dark:text-gray-400">
              {{ $t("team.averageCapacity") }}
            </div>
          </div>
        </div>
      </UCard>

      <UCard>
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <UIcon name="i-heroicons-clock" class="h-8 w-8 text-orange-500" />
          </div>
          <div class="ml-4">
            <div class="text-2xl font-bold text-gray-900 dark:text-white">
              {{ currentSprintHours }}h
            </div>
            <div class="text-sm text-gray-500 dark:text-gray-400">
              {{ $t("team.currentSprintHours") }}
            </div>
          </div>
        </div>
      </UCard>
    </div>

    <!-- 搜索和过滤 -->
    <UCard>
      <div class="flex flex-col sm:flex-row gap-4">
        <div class="flex-1">
          <UInput
            v-model="searchQuery"
            :placeholder="$t('team.searchMembers')"
            icon="i-heroicons-magnifying-glass"
          />
        </div>
        <div class="flex gap-2">
          <USelect
            v-model="selectedRole"
            :options="roleOptions"
            :placeholder="$t('team.filterByRole')"
            size="sm"
          />
          <USelect
            v-model="selectedStatus"
            :options="statusOptions"
            :placeholder="$t('team.filterByStatus')"
            size="sm"
          />
        </div>
      </div>
    </UCard>

    <!-- 团队成员列表 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <UCard
        v-for="member in filteredMembers"
        :key="member.id"
        class="hover:shadow-lg transition-shadow"
      >
        <div class="flex items-start space-x-4">
          <UAvatar
            :src="member.avatar"
            :alt="member.name"
            size="lg"
            class="bg-primary-500"
          >
            <span class="text-lg font-medium text-white">
              {{ getInitials(member.name) }}
            </span>
          </UAvatar>

          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between">
              <h3 class="text-sm font-medium text-gray-900 dark:text-white truncate">
                {{ member.name }}
              </h3>
              <UBadge
                :color="member.status === 'active' ? 'green' : 'gray'"
                size="xs"
                variant="soft"
              >
                {{ $t(`team.status.${member.status}`) }}
              </UBadge>
            </div>

            <p class="text-sm text-gray-500 dark:text-gray-400 mb-2">
              {{ $t(`team.role.${member.role}`) }}
            </p>

            <div class="space-y-2">
              <!-- 容量进度条 -->
              <div>
                <div class="flex items-center justify-between text-xs mb-1">
                  <span class="text-gray-500 dark:text-gray-400">
                    {{ $t("team.capacity") }}
                  </span>
                  <span class="font-medium text-gray-900 dark:text-white">
                    {{ member.capacity }}%
                  </span>
                </div>
                <UProgress
                  :value="member.capacity"
                  :color="getCapacityColor(member.capacity)"
                  size="xs"
                />
              </div>

              <!-- 当前任务 -->
              <div class="text-xs text-gray-500 dark:text-gray-400">
                {{ $t("team.currentTasks") }}: {{ member.currentTasks }}
              </div>

              <!-- 联系信息 -->
              <div class="text-xs text-gray-500 dark:text-gray-400">
                {{ member.email }}
              </div>
            </div>

            <!-- 操作按钮 -->
            <div class="mt-4 flex space-x-2">
              <UButton
                size="xs"
                variant="outline"
                @click="editMember(member)"
              >
                {{ $t("common.edit") }}
              </UButton>
              <UButton
                size="xs"
                variant="outline"
                @click="viewProfile(member)"
              >
                {{ $t("team.viewProfile") }}
              </UButton>
            </div>
          </div>
        </div>
      </UCard>
    </div>

    <!-- 空状态 -->
    <UCard v-if="filteredMembers.length === 0">
      <div class="text-center py-12">
        <UIcon name="i-heroicons-user-group" class="mx-auto h-12 w-12 text-gray-400" />
        <h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white">
          {{ $t("team.noMembers") }}
        </h3>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
          {{ $t("team.noMembersDescription") }}
        </p>
        <div class="mt-6">
          <UButton color="primary" @click="addMember">
            {{ $t("team.addFirstMember") }}
          </UButton>
        </div>
      </div>
    </UCard>
  </div>
</template>

<script setup>
// 国际化
const { t } = useI18n();

// 页面元数据
useSeoMeta({
  title: () => `${t('common.appName')} - ${t('navigation.team')}`,
  description: () => t('team.description')
});

// 搜索和过滤
const searchQuery = ref('');
const selectedRole = ref('');
const selectedStatus = ref('');

// 角色选项
const roleOptions = computed(() => [
  { label: t('common.all'), value: '' },
  { label: t('team.role.scrumMaster'), value: 'scrumMaster' },
  { label: t('team.role.productOwner'), value: 'productOwner' },
  { label: t('team.role.developer'), value: 'developer' },
  { label: t('team.role.tester'), value: 'tester' },
  { label: t('team.role.designer'), value: 'designer' }
]);

// 状态选项
const statusOptions = computed(() => [
  { label: t('common.all'), value: '' },
  { label: t('team.status.active'), value: 'active' },
  { label: t('team.status.inactive'), value: 'inactive' },
  { label: t('team.status.vacation'), value: 'vacation' }
]);

// Mock 团队成员数据
const teamMembers = ref([
  {
    id: 1,
    name: 'John Doe',
    email: 'john.doe@example.com',
    role: 'scrumMaster',
    status: 'active',
    capacity: 85,
    currentTasks: 3,
    avatar: null
  },
  {
    id: 2,
    name: 'Alice Smith',
    email: 'alice.smith@example.com',
    role: 'productOwner',
    status: 'active',
    capacity: 75,
    currentTasks: 5,
    avatar: null
  },
  {
    id: 3,
    name: 'Bob Johnson',
    email: 'bob.johnson@example.com',
    role: 'developer',
    status: 'active',
    capacity: 90,
    currentTasks: 4,
    avatar: null
  },
  {
    id: 4,
    name: 'Carol Brown',
    email: 'carol.brown@example.com',
    role: 'developer',
    status: 'vacation',
    capacity: 0,
    currentTasks: 0,
    avatar: null
  },
  {
    id: 5,
    name: 'David Wilson',
    email: 'david.wilson@example.com',
    role: 'tester',
    status: 'active',
    capacity: 80,
    currentTasks: 2,
    avatar: null
  }
]);

// 过滤的成员列表
const filteredMembers = computed(() => {
  let filtered = teamMembers.value;

  // 按名称搜索
  if (searchQuery.value) {
    filtered = filtered.filter(member =>
      member.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      member.email.toLowerCase().includes(searchQuery.value.toLowerCase())
    );
  }

  // 按角色过滤
  if (selectedRole.value) {
    filtered = filtered.filter(member => member.role === selectedRole.value);
  }

  // 按状态过滤
  if (selectedStatus.value) {
    filtered = filtered.filter(member => member.status === selectedStatus.value);
  }

  return filtered;
});

// 计算统计信息
const activeMembers = computed(() => teamMembers.value.filter(m => m.status === 'active'));
const averageCapacity = computed(() => {
  const active = activeMembers.value;
  if (active.length === 0) return 0;
  return Math.round(active.reduce((sum, m) => sum + m.capacity, 0) / active.length);
});
const currentSprintHours = computed(() => {
  return activeMembers.value.reduce((sum, m) => sum + (m.capacity / 100 * 40), 0);
});

// 获取姓名首字母
const getInitials = (name) => {
  return name
    .split(' ')
    .map(word => word.charAt(0))
    .join('')
    .toUpperCase()
    .slice(0, 2);
};

// 获取容量颜色
const getCapacityColor = (capacity) => {
  if (capacity >= 80) return 'red';
  if (capacity >= 60) return 'yellow';
  return 'green';
};

// 事件处理
const addMember = () => {
  // TODO: 实现添加成员功能
  console.log('添加团队成员');
};

const editMember = (member) => {
  // TODO: 实现编辑成员功能
  console.log('编辑成员:', member);
};

const viewProfile = (member) => {
  // TODO: 实现查看成员详情功能
  console.log('查看成员详情:', member);
};
</script>