<template>
  <aside
    class="w-64 min-w-64 max-w-64 bg-white dark:bg-gray-900 border-r border-gray-200 dark:border-gray-800 min-h-screen flex-shrink-0"
  >
    <nav class="p-4 space-y-6">
      <!-- 仪表盘 -->
      <div>
        <UButton
          to="/"
          variant="ghost"
          color="neutral"
          class="w-full justify-start"
          :class="{
            'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
              $route.path === '/',
          }"
        >
          <UIcon name="i-heroicons-home" class="w-4 h-4 mr-3" />
          {{ $t("navigation.dashboard") }}
        </UButton>
      </div>

      <!-- 笔记管理 -->
      <div>
        <div
          class="px-3 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider"
        >
          {{ $t("navigation.contentManagement") }}
        </div>
        <div class="space-y-1">
          <!-- 笔记列表 -->
          <UButton
            to="/notes/active"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/notes',
            }"
          >
            <UIcon name="i-heroicons-document-text" class="w-4 h-4 mr-3" />
            {{ $t("navigation.notesList") }}
          </UButton>
        </div>
      </div>

      <!-- 团队协作 -->
      <div>
        <div
          class="px-3 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider"
        >
          {{ $t("navigation.teamCollaboration") }}
        </div>
        <div class="space-y-1">
          <!-- 团队协作 -->
          <div>
            <UButton
              variant="ghost"
              color="neutral"
              class="w-full justify-start"
              @click="toggleTeamMenu"
            >
              <UIcon name="i-heroicons-users" class="w-4 h-4 mr-3" />
              {{ $t("navigation.teamCollaboration") }}
              <UIcon
                :name="
                  showTeamMenu
                    ? 'i-heroicons-chevron-down'
                    : 'i-heroicons-chevron-right'
                "
                class="w-4 h-4 ml-auto"
              />
            </UButton>

            <div v-show="showTeamMenu" class="ml-6 mt-1 space-y-1">
              <UButton
                to="/team/management"
                variant="ghost"
                color="neutral"
                size="sm"
                class="w-full justify-start text-sm"
                :class="{
                  'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                    $route.path === '/team/management',
                }"
              >
                <UIcon
                  name="i-heroicons-building-office"
                  class="w-3 h-3 mr-2"
                />
                {{ $t("navigation.teamManagement") }}
              </UButton>

              <UButton
                to="/team/members"
                variant="ghost"
                color="neutral"
                size="sm"
                class="w-full justify-start text-sm"
                :class="{
                  'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                    $route.path === '/team/members',
                }"
              >
                <UIcon name="i-heroicons-user-group" class="w-3 h-3 mr-2" />
                {{ $t("navigation.memberManagement") }}
              </UButton>
            </div>
          </div>
        </div>
      </div>

      <!-- 系统设置 -->
      <div>
        <div
          class="px-3 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider"
        >
          {{ $t("navigation.systemSettings") }}
        </div>
        <div class="space-y-1">
          <!-- 权限管理 -->
          <UButton
            to="/settings/permissions"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/settings/permissions',
            }"
          >
            <UIcon name="i-heroicons-key" class="w-4 h-4 mr-3" />
            {{ $t("navigation.permissions") }}
          </UButton>
        </div>
      </div>
    </nav>
  </aside>
</template>

<script setup lang="ts">
const { t } = useI18n();

// 控制团队协作子菜单的展开状态
const showTeamMenu = ref(false);

// 切换团队协作子菜单
const toggleTeamMenu = () => {
  showTeamMenu.value = !showTeamMenu.value;
};

// 监听路由变化，如果当前路由在团队模块下，自动展开子菜单
const route = useRoute();
watch(
  () => route.path,
  (newPath) => {
    if (newPath.startsWith("/team/")) {
      showTeamMenu.value = true;
    }
  },
  { immediate: true }
);
</script>
