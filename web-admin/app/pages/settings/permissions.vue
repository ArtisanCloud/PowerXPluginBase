<template>
  <div class="space-y-6">
    <!-- 页面头部 -->
    <div class="border-b border-gray-200 dark:border-gray-700 pb-4">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
        {{ $t("navigation.permissions") }}
      </h1>
      <p class="text-gray-600 dark:text-gray-400 mt-1">
        {{ $t("settings.permissionsDescription") }}
      </p>
    </div>

    <!-- 角色管理 -->
    <UCard>
      <template #header>
        <div class="flex justify-between items-center">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ $t("settings.roleManagement") }}
          </h2>
          <UButton color="primary" @click="addNewRole">
            {{ $t("settings.addRole") }}
          </UButton>
        </div>
      </template>

      <div class="space-y-4">
        <div
          v-for="role in roles"
          :key="role.id"
          class="p-4 border border-gray-200 dark:border-gray-700 rounded-lg"
        >
          <div class="flex items-center justify-between mb-3">
            <div>
              <h3 class="font-medium text-gray-900 dark:text-white">
                {{ role.name }}
              </h3>
              <p class="text-sm text-gray-600 dark:text-gray-400">
                {{ role.description }}
              </p>
            </div>
            <div class="flex items-center space-x-2">
              <UBadge variant="soft">
                {{ role.userCount }} {{ $t("settings.users") }}
              </UBadge>
              <UButton size="xs" variant="ghost" @click="editRole(role)">
                {{ $t("common.edit") }}
              </UButton>
              <UButton
                size="xs"
                variant="ghost"
                color="red"
                :disabled="role.isSystem"
                @click="deleteRole(role)"
              >
                {{ $t("common.delete") }}
              </UButton>
            </div>
          </div>

          <!-- 权限列表 -->
          <div class="grid grid-cols-2 md:grid-cols-4 gap-2">
            <div
              v-for="permission in role.permissions"
              :key="permission"
              class="text-xs bg-gray-100 dark:bg-gray-800 px-2 py-1 rounded"
            >
              {{ $t(`permissions.${permission}`) }}
            </div>
          </div>
        </div>
      </div>
    </UCard>

    <!-- 用户权限分配 -->
    <UCard>
      <template #header>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ $t("settings.userPermissions") }}
        </h2>
      </template>

      <div class="space-y-4">
        <!-- 搜索过滤 -->
        <div class="flex space-x-4">
          <UInput
            v-model="userSearchQuery"
            :placeholder="$t('settings.searchUsers')"
            icon="i-heroicons-magnifying-glass"
            class="flex-1"
          />
          <USelectMenu
            v-model="selectedRoleFilter"
            :options="roleFilterOptions"
            :placeholder="$t('settings.filterByRole')"
            class="w-48"
          />
        </div>

        <!-- 用户列表 -->
        <div class="overflow-x-auto">
          <table
            class="min-w-full divide-y divide-gray-200 dark:divide-gray-700"
          >
            <thead class="bg-gray-50 dark:bg-gray-800">
              <tr>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                >
                  {{ $t("settings.user") }}
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                >
                  {{ $t("settings.role") }}
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                >
                  {{ $t("settings.lastActive") }}
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                >
                  {{ $t("common.actions") }}
                </th>
              </tr>
            </thead>
            <tbody
              class="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-700"
            >
              <tr v-for="user in filteredUsers" :key="user.id">
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <div class="flex-shrink-0 h-10 w-10">
                      <UAvatar :src="user.avatar" :alt="user.name" size="sm" />
                    </div>
                    <div class="ml-4">
                      <div
                        class="text-sm font-medium text-gray-900 dark:text-white"
                      >
                        {{ user.name }}
                      </div>
                      <div class="text-sm text-gray-500">
                        {{ user.email }}
                      </div>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <UBadge :color="getRoleColor(user.role)" variant="soft">
                    {{ user.role }}
                  </UBadge>
                </td>
                <td
                  class="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-400"
                >
                  {{ formatDate(user.lastActive) }}
                </td>
                <td
                  class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium"
                >
                  <UDropdownMenu :items="getUserActions(user)">
                    <UButton variant="ghost" size="xs">
                      {{ $t("common.actions") }}
                      <UIcon
                        name="i-heroicons-chevron-down"
                        class="w-4 h-4 ml-1"
                      />
                    </UButton>
                  </UDropdownMenu>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </UCard>

    <!-- 权限矩阵 -->
    <UCard>
      <template #header>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ $t("settings.permissionMatrix") }}
        </h2>
      </template>

      <div class="overflow-x-auto">
        <table class="min-w-full">
          <thead>
            <tr class="border-b border-gray-200 dark:border-gray-700">
              <th
                class="text-left py-3 px-4 font-medium text-gray-900 dark:text-white"
              >
                {{ $t("settings.permission") }}
              </th>
              <th
                v-for="role in roles"
                :key="role.id"
                class="text-center py-3 px-4 font-medium text-gray-900 dark:text-white"
              >
                {{ role.name }}
              </th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="permission in allPermissions"
              :key="permission.key"
              class="border-b border-gray-100 dark:border-gray-800"
            >
              <td class="py-3 px-4">
                <div>
                  <div class="font-medium text-gray-900 dark:text-white">
                    {{ permission.name }}
                  </div>
                  <div class="text-sm text-gray-500">
                    {{ permission.description }}
                  </div>
                </div>
              </td>
              <td
                v-for="role in roles"
                :key="role.id"
                class="text-center py-3 px-4"
              >
                <UCheckbox
                  :model-value="hasPermission(role, permission.key)"
                  :disabled="role.isSystem"
                  @update:model-value="
                    togglePermission(role, permission.key, $event)
                  "
                />
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>
  </div>
</template>

<script setup>
// 页面元数据
definePageMeta({
  title: "permissions",
});

// 响应式数据
const userSearchQuery = ref("");
const selectedRoleFilter = ref("all");

// 角色数据
const roles = ref([
  {
    id: 1,
    name: "管理员",
    description: "系统管理员，拥有所有权限",
    userCount: 2,
    isSystem: true,
    permissions: ["viewNotes", "manageNotes", "manageTeam", "manageUsers", "viewReports", "manageSettings"],
  },
  {
    id: 2,
    name: "编辑者",
    description: "可以创建和编辑笔记内容",
    userCount: 5,
    isSystem: true,
    permissions: ["viewNotes", "manageNotes", "viewReports"],
  },
  {
    id: 3,
    name: "查看者",
    description: "只能查看笔记内容",
    userCount: 8,
    isSystem: true,
    permissions: ["viewNotes"],
  },
  {
    id: 4,
    name: "协作者",
    description: "可以查看和评论笔记",
    userCount: 3,
    isSystem: false,
    permissions: ["viewNotes", "commentNotes"],
  },
]);

// 用户数据
const users = ref([
  {
    id: 1,
    name: "张三",
    email: "zhangsan@example.com",
    role: "管理员",
    lastActive: "2024-02-15",
    avatar: null,
  },
  {
    id: 2,
    name: "李四",
    email: "lisi@example.com",
    role: "编辑者",
    lastActive: "2024-02-14",
    avatar: null,
  },
  {
    id: 3,
    name: "王五",
    email: "wangwu@example.com",
    role: "查看者",
    lastActive: "2024-02-15",
    avatar: null,
  },
  {
    id: 4,
    name: "赵六",
    email: "zhaoliu@example.com",
    role: "协作者",
    lastActive: "2024-02-13",
    avatar: null,
  },
]);

// 所有权限列表
const allPermissions = ref([
  {
    key: "viewNotes",
    name: "查看笔记",
    description: "可以查看笔记内容和列表",
  },
  {
    key: "manageNotes",
    name: "管理笔记",
    description: "可以创建、编辑和删除笔记",
  },
  {
    key: "commentNotes",
    name: "评论笔记",
    description: "可以对笔记进行评论和回复",
  },
  {
    key: "shareNotes",
    name: "分享笔记",
    description: "可以分享笔记给其他用户",
  },
  {
    key: "manageTeam",
    name: "管理团队",
    description: "可以添加、移除和管理团队成员",
  },
  {
    key: "manageUsers",
    name: "管理用户",
    description: "可以管理用户账户和权限",
  },
  {
    key: "viewReports",
    name: "查看报告",
    description: "可以查看使用统计和分析报告",
  },
  {
    key: "manageSettings",
    name: "管理设置",
    description: "可以修改系统设置和配置",
  },
]);

// 计算属性
const roleFilterOptions = computed(() => [
  { label: "全部角色", value: "all" },
  ...roles.value.map((role) => ({ label: role.name, value: role.name })),
]);

const filteredUsers = computed(() => {
  let filtered = users.value;

  if (userSearchQuery.value) {
    filtered = filtered.filter(
      (user) =>
        user.name.toLowerCase().includes(userSearchQuery.value.toLowerCase()) ||
        user.email.toLowerCase().includes(userSearchQuery.value.toLowerCase())
    );
  }

  if (selectedRoleFilter.value !== "all") {
    filtered = filtered.filter(
      (user) => user.role === selectedRoleFilter.value
    );
  }

  return filtered;
});

// 方法
const addNewRole = () => {
  console.log("添加新角色");
};

const editRole = (role) => {
  console.log("编辑角色", role);
};

const deleteRole = (role) => {
  console.log("删除角色", role);
};

const getRoleColor = (roleName) => {
  const colorMap = {
    管理员: "red",
    编辑者: "blue",
    查看者: "green",
    协作者: "purple",
  };
  return colorMap[roleName] || "gray";
};

const getUserActions = (user) => [
  [
    {
      label: "编辑角色",
      icon: "i-heroicons-pencil",
      click: () => editUserRole(user),
    },
  ],
  [
    {
      label: "重置密码",
      icon: "i-heroicons-key",
      click: () => resetPassword(user),
    },
  ],
  [
    {
      label: "禁用用户",
      icon: "i-heroicons-no-symbol",
      click: () => disableUser(user),
    },
  ],
];

const editUserRole = (user) => {
  console.log("编辑用户角色", user);
};

const resetPassword = (user) => {
  console.log("重置密码", user);
};

const disableUser = (user) => {
  console.log("禁用用户", user);
};

const hasPermission = (role, permissionKey) => {
  return role.permissions.includes(permissionKey);
};

const togglePermission = (role, permissionKey, hasPermission) => {
  if (role.isSystem) return;

  if (hasPermission) {
    role.permissions.push(permissionKey);
  } else {
    const index = role.permissions.indexOf(permissionKey);
    if (index > -1) {
      role.permissions.splice(index, 1);
    }
  }
};

const formatDate = (date) => {
  return new Date(date).toLocaleDateString();
};
</script>
