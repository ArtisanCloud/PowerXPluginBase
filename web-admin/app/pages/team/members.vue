<template>
  <div class="space-y-6">
    <!-- 页面头部 -->
    <div class="border-b border-gray-200 dark:border-gray-700 pb-4">
      <div class="flex justify-between items-center">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
            {{ $t("navigation.memberManagement") }}
          </h1>
          <p class="text-gray-600 dark:text-gray-400 mt-1">
            管理团队成员、角色分配和权限设置
          </p>
        </div>
        <UButton color="primary" @click="showInviteModal = true">
          <UIcon name="i-heroicons-plus" class="w-4 h-4 mr-2" />
          邀请成员
        </UButton>
      </div>
    </div>

    <!-- 成员统计 -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <UCard>
        <div class="text-center">
          <p class="text-2xl font-bold text-gray-900 dark:text-white">
            {{ memberStats.total }}
          </p>
          <p class="text-sm text-gray-600 dark:text-gray-400">总成员</p>
        </div>
      </UCard>
      <UCard>
        <div class="text-center">
          <p class="text-2xl font-bold text-green-600">
            {{ memberStats.active }}
          </p>
          <p class="text-sm text-gray-600 dark:text-gray-400">活跃成员</p>
        </div>
      </UCard>
      <UCard>
        <div class="text-center">
          <p class="text-2xl font-bold text-blue-600">
            {{ memberStats.admins }}
          </p>
          <p class="text-sm text-gray-600 dark:text-gray-400">管理员</p>
        </div>
      </UCard>
      <UCard>
        <div class="text-center">
          <p class="text-2xl font-bold text-purple-600">
            {{ memberStats.pending }}
          </p>
          <p class="text-sm text-gray-600 dark:text-gray-400">待确认</p>
        </div>
      </UCard>
    </div>

    <!-- 搜索和筛选 -->
    <UCard>
      <div class="flex flex-col sm:flex-row gap-4">
        <UInput
          v-model="searchQuery"
          placeholder="搜索成员姓名或邮箱"
          icon="i-heroicons-magnifying-glass"
          class="flex-1"
        />
        <USelectMenu
          v-model="selectedRole"
          :options="roleOptions"
          placeholder="筛选角色"
          class="w-full sm:w-48"
        />
        <USelectMenu
          v-model="selectedStatus"
          :options="statusOptions"
          placeholder="筛选状态"
          class="w-full sm:w-48"
        />
      </div>
    </UCard>

    <!-- 成员列表 -->
    <UCard>
      <template #header>
        <div class="flex justify-between items-center">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
            团队成员 ({{ filteredMembers.length }})
          </h2>
          <div class="flex items-center space-x-2">
            <UButton variant="ghost" size="sm" @click="exportMembers">
              <UIcon name="i-heroicons-arrow-down-tray" class="w-4 h-4 mr-1" />
              导出
            </UButton>
          </div>
        </div>
      </template>

      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
          <thead class="bg-gray-50 dark:bg-gray-800">
            <tr>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
              >
                成员
              </th>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
              >
                角色
              </th>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
              >
                状态
              </th>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
              >
                加入时间
              </th>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
              >
                最后活跃
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
              v-for="member in filteredMembers"
              :key="member.id"
              class="hover:bg-gray-50 dark:hover:bg-gray-800"
            >
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <div class="flex-shrink-0 h-10 w-10">
                    <UAvatar
                      :src="member.avatar"
                      :alt="member.name"
                      size="sm"
                    />
                  </div>
                  <div class="ml-4">
                    <div
                      class="text-sm font-medium text-gray-900 dark:text-white"
                    >
                      {{ member.name }}
                    </div>
                    <div class="text-sm text-gray-500">
                      {{ member.email }}
                    </div>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <UBadge :color="getRoleColor(member.role)" variant="soft">
                  {{ member.role }}
                </UBadge>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <UBadge :color="getStatusColor(member.status)" variant="soft">
                  {{ getStatusText(member.status) }}
                </UBadge>
              </td>
              <td
                class="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-400"
              >
                {{ formatDate(member.joinedAt) }}
              </td>
              <td
                class="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-400"
              >
                {{ formatDate(member.lastActive) }}
              </td>
              <td
                class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium"
              >
                <UDropdownMenu :items="getMemberActions(member)">
                  <UButton variant="ghost" size="xs">
                    操作
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

      <!-- 空状态 -->
      <div v-if="filteredMembers.length === 0" class="text-center py-12">
        <UIcon
          name="i-heroicons-users"
          class="w-12 h-12 text-gray-400 mx-auto mb-4"
        />
        <p class="text-gray-500 dark:text-gray-400">没有找到匹配的成员</p>
      </div>
    </UCard>

    <!-- 邀请成员模态框 -->
    <UModal v-model="showInviteModal">
      <UCard>
        <template #header>
          <div class="flex justify-between items-center">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
              邀请新成员
            </h3>
            <UButton variant="ghost" size="sm" @click="showInviteModal = false">
              <UIcon name="i-heroicons-x-mark" class="w-4 h-4" />
            </UButton>
          </div>
        </template>

        <div class="space-y-4">
          <div>
            <label
              class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
            >
              邮箱地址
            </label>
            <UInput
              v-model="inviteForm.email"
              type="email"
              placeholder="输入邮箱地址"
            />
          </div>

          <div>
            <label
              class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
            >
              分配角色
            </label>
            <USelectMenu
              v-model="inviteForm.role"
              :options="roleOptions.filter((r) => r.value !== 'all')"
              placeholder="选择角色"
            />
          </div>

          <div>
            <label
              class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
            >
              邀请消息 (可选)
            </label>
            <UTextarea
              v-model="inviteForm.message"
              placeholder="输入邀请消息"
              :rows="3"
            />
          </div>
        </div>

        <template #footer>
          <div class="flex justify-end space-x-3">
            <UButton variant="ghost" @click="showInviteModal = false">
              取消
            </UButton>
            <UButton color="primary" @click="sendInvite"> 发送邀请 </UButton>
          </div>
        </template>
      </UCard>
    </UModal>

    <!-- 编辑成员模态框 -->
    <UModal v-model="showEditModal">
      <UCard>
        <template #header>
          <div class="flex justify-between items-center">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
              编辑成员信息
            </h3>
            <UButton variant="ghost" size="sm" @click="showEditModal = false">
              <UIcon name="i-heroicons-x-mark" class="w-4 h-4" />
            </UButton>
          </div>
        </template>

        <div v-if="editingMember" class="space-y-4">
          <div class="flex items-center space-x-4">
            <UAvatar
              :src="editingMember.avatar"
              :alt="editingMember.name"
              size="lg"
            />
            <div>
              <p class="font-medium text-gray-900 dark:text-white">
                {{ editingMember.name }}
              </p>
              <p class="text-sm text-gray-500">
                {{ editingMember.email }}
              </p>
            </div>
          </div>

          <div>
            <label
              class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
            >
              角色
            </label>
            <USelectMenu
              v-model="editForm.role"
              :options="roleOptions.filter((r) => r.value !== 'all')"
              placeholder="选择角色"
            />
          </div>

          <div>
            <label
              class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
            >
              状态
            </label>
            <USelectMenu
              v-model="editForm.status"
              :options="statusOptions.filter((s) => s.value !== 'all')"
              placeholder="选择状态"
            />
          </div>
        </div>

        <template #footer>
          <div class="flex justify-end space-x-3">
            <UButton variant="ghost" @click="showEditModal = false">
              取消
            </UButton>
            <UButton color="primary" @click="saveMemberEdit">
              保存更改
            </UButton>
          </div>
        </template>
      </UCard>
    </UModal>
  </div>
</template>

<script setup>
// 页面元数据
definePageMeta({
  title: "memberManagement",
});

// 响应式数据
const searchQuery = ref("");
const selectedRole = ref("all");
const selectedStatus = ref("all");
const showInviteModal = ref(false);
const showEditModal = ref(false);
const editingMember = ref(null);

// 成员统计
const memberStats = ref({
  total: 12,
  active: 10,
  admins: 3,
  pending: 2,
});

// 邀请表单
const inviteForm = ref({
  email: "",
  role: "",
  message: "",
});

// 编辑表单
const editForm = ref({
  role: "",
  status: "",
});

// 成员数据
const members = ref([
  {
    id: 1,
    name: "张三",
    email: "zhangsan@example.com",
    role: "管理员",
    status: "active",
    joinedAt: "2024-01-15",
    lastActive: "2024-02-15",
    avatar: null,
  },
  {
    id: 2,
    name: "李四",
    email: "lisi@example.com",
    role: "编辑者",
    status: "active",
    joinedAt: "2024-01-20",
    lastActive: "2024-02-14",
    avatar: null,
  },
  {
    id: 3,
    name: "王五",
    email: "wangwu@example.com",
    role: "查看者",
    status: "active",
    joinedAt: "2024-02-01",
    lastActive: "2024-02-15",
    avatar: null,
  },
  {
    id: 4,
    name: "赵六",
    email: "zhaoliu@example.com",
    role: "协作者",
    status: "pending",
    joinedAt: "2024-02-10",
    lastActive: "2024-02-13",
    avatar: null,
  },
]);

// 选项数据
const roleOptions = [
  { label: "全部角色", value: "all" },
  { label: "管理员", value: "管理员" },
  { label: "编辑者", value: "编辑者" },
  { label: "查看者", value: "查看者" },
  { label: "协作者", value: "协作者" },
];

const statusOptions = [
  { label: "全部状态", value: "all" },
  { label: "活跃", value: "active" },
  { label: "待确认", value: "pending" },
  { label: "已禁用", value: "disabled" },
];

// 计算属性
const filteredMembers = computed(() => {
  let filtered = members.value;

  if (searchQuery.value) {
    filtered = filtered.filter(
      (member) =>
        member.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        member.email.toLowerCase().includes(searchQuery.value.toLowerCase())
    );
  }

  if (selectedRole.value !== "all") {
    filtered = filtered.filter((member) => member.role === selectedRole.value);
  }

  if (selectedStatus.value !== "all") {
    filtered = filtered.filter(
      (member) => member.status === selectedStatus.value
    );
  }

  return filtered;
});

// 方法
const getRoleColor = (role) => {
  const colorMap = {
    管理员: "red",
    编辑者: "blue",
    查看者: "green",
    协作者: "purple",
  };
  return colorMap[role] || "gray";
};

const getStatusColor = (status) => {
  const colorMap = {
    active: "green",
    pending: "yellow",
    disabled: "red",
  };
  return colorMap[status] || "gray";
};

const getStatusText = (status) => {
  const textMap = {
    active: "活跃",
    pending: "待确认",
    disabled: "已禁用",
  };
  return textMap[status] || status;
};

const getMemberActions = (member) => [
  [
    {
      label: "编辑",
      icon: "i-heroicons-pencil",
      click: () => editMember(member),
    },
    {
      label: "查看详情",
      icon: "i-heroicons-eye",
      click: () => viewMemberDetails(member),
    },
  ],
  [
    {
      label: member.status === "active" ? "禁用" : "启用",
      icon:
        member.status === "active"
          ? "i-heroicons-no-symbol"
          : "i-heroicons-check",
      click: () => toggleMemberStatus(member),
    },
  ],
  [
    {
      label: "移除成员",
      icon: "i-heroicons-trash",
      click: () => removeMember(member),
    },
  ],
];

const editMember = (member) => {
  editingMember.value = member;
  editForm.value = {
    role: member.role,
    status: member.status,
  };
  showEditModal.value = true;
};

const viewMemberDetails = (member) => {
  console.log("查看成员详情", member);
};

const toggleMemberStatus = (member) => {
  console.log("切换成员状态", member);
};

const removeMember = (member) => {
  console.log("移除成员", member);
};

const sendInvite = () => {
  console.log("发送邀请", inviteForm.value);
  showInviteModal.value = false;
  // 重置表单
  inviteForm.value = {
    email: "",
    role: "",
    message: "",
  };
};

const saveMemberEdit = () => {
  console.log("保存成员编辑", editForm.value);
  showEditModal.value = false;
};

const exportMembers = () => {
  console.log("导出成员列表");
};

const formatDate = (date) => {
  return new Date(date).toLocaleDateString();
};
</script>
